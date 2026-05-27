package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/config"
	"github.com/yargotev/exito-tools/internal/domain/catalog"
	"github.com/yargotev/exito-tools/internal/domain/geo"
	"github.com/yargotev/exito-tools/internal/domain/orders"
	"github.com/yargotev/exito-tools/internal/registry"
	clisurface "github.com/yargotev/exito-tools/internal/surface/cli"
)

func TestRootHelpPaths(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{name: "bare root shows help", args: nil},
		{name: "explicit help flag shows help", args: []string{"--help"}},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			root := clisurface.NewRoot(app.New)
			var output bytes.Buffer
			root.SetOut(&output)
			root.SetErr(&output)
			root.SetArgs(tc.args)

			if err := root.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			rendered := output.String()
			for _, fragment := range []string{
				"Exito Tools command-line interface",
				"Use an implemented subcommand for machine-readable JSON output.",
				"Usage:",
				"capabilities",
				"config",
				"orders",
				"geo",
				"catalog",
				"tui",
			} {
				if !strings.Contains(rendered, fragment) {
					t.Fatalf("help output missing %q\n%s", fragment, rendered)
				}
			}

			for _, forbidden := range []string{"\"ok\"", "\"data\"", "\"error\""} {
				if strings.Contains(strings.ToLower(rendered), forbidden) {
					t.Fatalf("help output unexpectedly contains %q\n%s", forbidden, rendered)
				}
			}
		})
	}
}

func TestRootWiresBootFlagsIntoApplicationOptions(t *testing.T) {
	t.Parallel()

	var got app.Options
	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		got = options
		return &app.Application{Registry: registry.NewBuilder().Finalize()}, nil
	})

	var output bytes.Buffer
	root.SetOut(&output)
	root.SetErr(&output)
	root.SetArgs([]string{"--config", "./team.yaml", "--profile", "prod"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got.Config.ConfigPath != "./team.yaml" {
		t.Fatalf("ConfigPath = %q, want ./team.yaml", got.Config.ConfigPath)
	}
	if got.Config.Profile != "prod" {
		t.Fatalf("Profile = %q, want prod", got.Config.Profile)
	}
}

func TestRootHelpAdvertisesBootFlags(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(app.New)
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetErr(&output)
	root.SetArgs([]string{"--help"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	rendered := output.String()
	for _, fragment := range []string{"--config", "--profile", "--correlation-id"} {
		if !strings.Contains(rendered, fragment) {
			t.Fatalf("help output missing %q\n%s", fragment, rendered)
		}
	}
}

func TestCapabilitiesCommandEmitsInventoryEnvelope(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.Register(capability.Definition{
		ID:          "foundation.example",
		Domain:      "foundation",
		Version:     "1.0.0",
		Title:       "Foundation Example",
		Description: "Registered during application boot.",
		Risk:        capability.RiskReadOnly,
		Audiences:   []capability.Audience{capability.AudienceAgents, capability.AudiencePeople},
		Visibility:  []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "id", Type: capability.InputTypeString, Required: true, Description: "Identifier."},
		}},
	}); err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		return &app.Application{
			Config:   config.Effective{Profile: "staging"},
			Registry: builder.Finalize(),
		}, nil
	})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"capabilities"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("capabilities output should not render help\n%s", stdout.String())
	}

	var got struct {
		OK   bool `json:"ok"`
		Data struct {
			Capabilities []capability.Definition `json:"capabilities"`
		} `json:"data"`
		Meta struct {
			RequestID  string `json:"requestId"`
			Profile    string `json:"profile"`
			DurationMS int64  `json:"durationMs"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("capabilities output is not JSON: %v\n%s", err, stdout.String())
	}

	if !got.OK {
		t.Fatalf("ok = false, want true")
	}
	if got.Meta.RequestID == "" {
		t.Fatalf("meta.requestId is empty")
	}
	if got.Meta.Profile != "staging" {
		t.Fatalf("meta.profile = %q, want staging", got.Meta.Profile)
	}
	if got.Meta.DurationMS < 0 {
		t.Fatalf("meta.durationMs = %d, want non-negative", got.Meta.DurationMS)
	}
	if len(got.Data.Capabilities) != 1 {
		t.Fatalf("capabilities length = %d, want 1", len(got.Data.Capabilities))
	}
	capabilityGot := got.Data.Capabilities[0]
	if capabilityGot.ID != "foundation.example" {
		t.Fatalf("capability ID = %q, want foundation.example", capabilityGot.ID)
	}
	if capabilityGot.Domain != "foundation" || capabilityGot.Version != "1.0.0" || capabilityGot.Risk != capability.RiskReadOnly {
		t.Fatalf("capability metadata = %#v, want domain, version, and risk", capabilityGot)
	}
	if len(capabilityGot.Audiences) != 2 || capabilityGot.Audiences[0] != capability.AudienceAgents || capabilityGot.Audiences[1] != capability.AudiencePeople {
		t.Fatalf("capability audiences = %#v, want agents and people", capabilityGot.Audiences)
	}
	if len(capabilityGot.Visibility) != 2 || capabilityGot.Visibility[0] != capability.VisibilityCLI || capabilityGot.Visibility[1] != capability.VisibilityCommandPalette {
		t.Fatalf("capability visibility = %#v, want cli and command-palette", capabilityGot.Visibility)
	}
	if capabilityGot.InputSchema == nil || len(capabilityGot.InputSchema.Fields) != 1 {
		t.Fatalf("capability input schema = %#v, want one field", capabilityGot.InputSchema)
	}
	field := capabilityGot.InputSchema.Fields[0]
	if field.Name != "id" || field.Type != capability.InputTypeString || !field.Required || field.Description != "Identifier." {
		t.Fatalf("capability input field = %#v, want string required id", field)
	}
}

func TestCapabilitiesCommandPropagatesCorrelationID(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		return &app.Application{
			Config:   config.Effective{Profile: "staging"},
			Registry: registry.NewBuilder().Finalize(),
		}, nil
	})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"--correlation-id", "corr-123", "capabilities"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var got struct {
		Meta struct {
			CorrelationID string `json:"correlationId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("capabilities output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.Meta.CorrelationID != "corr-123" {
		t.Fatalf("meta.correlationId = %q, want corr-123", got.Meta.CorrelationID)
	}
}

func TestCapabilitiesCommandWiresBootFlagsIntoApplicationOptions(t *testing.T) {
	t.Parallel()

	var got app.Options
	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		got = options
		return &app.Application{
			Config:   config.Effective{Profile: options.Config.Profile},
			Registry: registry.NewBuilder().Finalize(),
		}, nil
	})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"--config", "./team.yaml", "--profile", "prod", "capabilities"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if got.Config.ConfigPath != "./team.yaml" {
		t.Fatalf("ConfigPath = %q, want ./team.yaml", got.Config.ConfigPath)
	}
	if got.Config.Profile != "prod" {
		t.Fatalf("Profile = %q, want prod", got.Config.Profile)
	}
	if !strings.Contains(stdout.String(), "\"profile\":\"prod\"") {
		t.Fatalf("output should include selected profile\n%s", stdout.String())
	}
}

func TestRunCommandExecutesCapabilityWithInlineJSON(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	var gotRequest capability.ExecutionRequest
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{ID: "orders.get"},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			gotRequest = request
			return capability.ExecutionResult{Data: map[string]any{"orderId": request.Input["id"]}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		return &app.Application{
			Config:   config.Effective{Profile: options.Config.Profile},
			Registry: builder.Finalize(),
		}, nil
	})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"--profile", "prod", "--correlation-id", "corr-123", "run", "orders.get", "--input-json", `{"id":"A123"}`})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if gotRequest.Input["id"] != "A123" {
		t.Fatalf("handler input id = %#v, want A123", gotRequest.Input["id"])
	}
	if gotRequest.Context.Profile != "prod" {
		t.Fatalf("handler profile = %q, want prod", gotRequest.Context.Profile)
	}
	if gotRequest.Context.CorrelationID != "corr-123" {
		t.Fatalf("handler correlation ID = %q, want corr-123", gotRequest.Context.CorrelationID)
	}

	var got struct {
		OK   bool `json:"ok"`
		Data struct {
			OrderID string `json:"orderId"`
		} `json:"data"`
		Meta struct {
			RequestID     string `json:"requestId"`
			CorrelationID string `json:"correlationId"`
			Profile       string `json:"profile"`
			CapabilityID  string `json:"capabilityId"`
			DurationMS    int64  `json:"durationMs"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("run output is not JSON: %v\n%s", err, stdout.String())
	}
	if !got.OK {
		t.Fatalf("ok = false, want true\n%s", stdout.String())
	}
	if got.Data.OrderID != "A123" {
		t.Fatalf("data.orderId = %q, want A123", got.Data.OrderID)
	}
	if got.Meta.RequestID == "" || got.Meta.CorrelationID != "corr-123" || got.Meta.Profile != "prod" || got.Meta.CapabilityID != "orders.get" || got.Meta.DurationMS < 0 {
		t.Fatalf("unexpected metadata: %#v", got.Meta)
	}
}

func TestRunCommandAcceptsInputFileAndStdin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args func(t *testing.T) []string
		in   string
		want string
	}{
		{
			name: "input file",
			args: func(t *testing.T) []string {
				t.Helper()
				path := t.TempDir() + "/input.json"
				if err := os.WriteFile(path, []byte(`{"id":"file-123"}`), 0o600); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
				return []string{"run", "orders.get", "--input-file", path}
			},
			want: "file-123",
		},
		{
			name: "stdin",
			args: func(t *testing.T) []string { return []string{"run", "orders.get"} },
			in:   `{"id":"stdin-123"}`,
			want: "stdin-123",
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			builder := registry.NewBuilder()
			if err := builder.RegisterExecutable(capability.Executable{
				Definition: capability.Definition{ID: "orders.get"},
				Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
					return capability.ExecutionResult{Data: map[string]any{"orderId": request.Input["id"]}}, nil
				},
			}); err != nil {
				t.Fatalf("RegisterExecutable() error = %v", err)
			}

			root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
				return &app.Application{Config: config.Effective{Profile: "staging"}, Registry: builder.Finalize()}, nil
			})
			var stdout bytes.Buffer
			root.SetOut(&stdout)
			root.SetErr(&stdout)
			if tc.in != "" {
				root.SetIn(strings.NewReader(tc.in))
			}
			root.SetArgs(tc.args(t))

			if err := root.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}
			if !strings.Contains(stdout.String(), tc.want) {
				t.Fatalf("run output missing %q\n%s", tc.want, stdout.String())
			}
		})
	}
}

func TestOrdersGetCommandRunsOrdersGetCapability(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(app.New)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"--profile", "prod", "--correlation-id", "corr-123", "orders", "get", "--id", "A123"})

	assertFailureExitCode(t, root.Execute())

	var got struct {
		OK    bool `json:"ok"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
		Meta struct {
			CorrelationID string `json:"correlationId"`
			Profile       string `json:"profile"`
			CapabilityID  string `json:"capabilityId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("orders get output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.OK {
		t.Fatalf("ok = true, want false until Orders client is configured")
	}
	if got.Error.Code != orders.ErrorOrdersNotConfigured {
		t.Fatalf("error.code = %q, want %s", got.Error.Code, orders.ErrorOrdersNotConfigured)
	}
	if got.Meta.CorrelationID != "corr-123" || got.Meta.Profile != "prod" || got.Meta.CapabilityID != orders.CapabilityGetID {
		t.Fatalf("unexpected metadata: %#v", got.Meta)
	}
}

func TestOrdersGetCommandPassesOrderTypeFlag(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	var gotRequest capability.ExecutionRequest
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID: orders.CapabilityGetID,
			InputSchema: &capability.InputSchema{Fields: []capability.InputField{
				{Name: "id", Type: capability.InputTypeString, Required: true},
				{Name: "orderType", Type: capability.InputTypeString, Required: false},
			}},
		},
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			gotRequest = request
			return capability.ExecutionResult{Data: map[string]any{"orderId": request.Input["id"]}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		return &app.Application{Config: config.Effective{Profile: options.Config.Profile}, Registry: builder.Finalize()}, nil
	})
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"--profile", "prod", "orders", "get", "--id", "A123", "--order-type", "CarullaEcomm"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v\n%s", err, stdout.String())
	}

	if gotRequest.Input["id"] != "A123" || gotRequest.Input["orderType"] != "CarullaEcomm" {
		t.Fatalf("handler input = %#v, want id A123 and orderType CarullaEcomm", gotRequest.Input)
	}
}

func TestOrdersGetCommandRequiresID(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(app.New)
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetErr(&output)
	root.SetArgs([]string{"orders", "get"})

	if err := root.Execute(); err == nil {
		t.Fatalf("Execute() error = nil, want missing required flag error")
	}
	if strings.Contains(output.String(), "{\"ok\"") {
		t.Fatalf("missing flag error should not emit JSON envelope\n%s", output.String())
	}
}

func TestRunCommandExecutesBootstrappedOrdersGetCapability(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(app.New)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"run", orders.CapabilityGetID, "--input-json", `{"id":"A123"}`})

	assertFailureExitCode(t, root.Execute())

	var got struct {
		OK    bool `json:"ok"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
		Meta struct {
			CapabilityID string `json:"capabilityId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("run output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.OK {
		t.Fatalf("ok = true, want false until Orders client is configured")
	}
	if got.Error.Code != orders.ErrorOrdersNotConfigured {
		t.Fatalf("error.code = %q, want %s", got.Error.Code, orders.ErrorOrdersNotConfigured)
	}
	if got.Meta.CapabilityID != orders.CapabilityGetID {
		t.Fatalf("meta.capabilityId = %q, want %s", got.Meta.CapabilityID, orders.CapabilityGetID)
	}
}

func TestGeoGeocodeAddressCommandRunsGeoCapability(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(app.New)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"--profile", "prod", "--correlation-id", "corr-geo", "geo", "geocode-address", "--city", "Bogota", "--address", "CL 57 H SUR # 68 D - 75"})

	assertFailureExitCode(t, root.Execute())

	var got struct {
		OK    bool `json:"ok"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
		Meta struct {
			CorrelationID string `json:"correlationId"`
			Profile       string `json:"profile"`
			CapabilityID  string `json:"capabilityId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("geo geocode-address output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.OK {
		t.Fatalf("ok = true, want false until Geo client is configured")
	}
	if got.Error.Code != geo.ErrorGeoNotConfigured {
		t.Fatalf("error.code = %q, want %s", got.Error.Code, geo.ErrorGeoNotConfigured)
	}
	if got.Meta.CorrelationID != "corr-geo" || got.Meta.Profile != "prod" || got.Meta.CapabilityID != geo.CapabilityGeocodeAddressID {
		t.Fatalf("unexpected metadata: %#v", got.Meta)
	}
}

func TestGeoGeocodeAddressCommandRequiresCityAndAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{name: "missing city", args: []string{"geo", "geocode-address", "--address", "CL 57 H SUR # 68 D - 75"}},
		{name: "missing address", args: []string{"geo", "geocode-address", "--city", "Bogota"}},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			root := clisurface.NewRoot(app.New)
			var output bytes.Buffer
			root.SetOut(&output)
			root.SetErr(&output)
			root.SetArgs(tc.args)

			if err := root.Execute(); err == nil {
				t.Fatalf("Execute() error = nil, want missing required flag error")
			}
			if strings.Contains(output.String(), "{\"ok\"") {
				t.Fatalf("missing flag error should not emit JSON envelope\n%s", output.String())
			}
		})
	}
}

func TestRunCommandExecutesBootstrappedGeoGeocodeAddressCapability(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(app.New)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"run", geo.CapabilityGeocodeAddressID, "--input-json", `{"city":"Bogota","address":"CL 57 H SUR # 68 D - 75"}`})

	assertFailureExitCode(t, root.Execute())

	var got struct {
		OK    bool `json:"ok"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
		Meta struct {
			CapabilityID string `json:"capabilityId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("run output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.OK {
		t.Fatalf("ok = true, want false until Geo client is configured")
	}
	if got.Error.Code != geo.ErrorGeoNotConfigured {
		t.Fatalf("error.code = %q, want %s", got.Error.Code, geo.ErrorGeoNotConfigured)
	}
	if got.Meta.CapabilityID != geo.CapabilityGeocodeAddressID {
		t.Fatalf("meta.capabilityId = %q, want %s", got.Meta.CapabilityID, geo.CapabilityGeocodeAddressID)
	}
}

func TestRunCommandReturnsInvalidInputEnvelope(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{
			ID: "orders.get",
			InputSchema: &capability.InputSchema{Fields: []capability.InputField{
				{Name: "id", Type: capability.InputTypeString, Required: true},
			}},
		},
		Handler: func(context.Context, capability.ExecutionRequest) (capability.ExecutionResult, error) {
			return capability.ExecutionResult{Data: map[string]any{"unreachable": true}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		return &app.Application{Config: config.Effective{Profile: "staging"}, Registry: builder.Finalize()}, nil
	})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"run", "orders.get", "--input-json", `{}`})

	assertFailureExitCode(t, root.Execute())

	var got struct {
		OK    bool `json:"ok"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
		Meta struct {
			CapabilityID string `json:"capabilityId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("run output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.OK {
		t.Fatalf("ok = true, want false")
	}
	if got.Error.Code != "INVALID_INPUT" {
		t.Fatalf("error.code = %q, want INVALID_INPUT", got.Error.Code)
	}
	if got.Meta.CapabilityID != "orders.get" {
		t.Fatalf("meta.capabilityId = %q, want orders.get", got.Meta.CapabilityID)
	}
}

func TestRunCommandReturnsCapabilityNotFoundEnvelope(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		return &app.Application{Config: config.Effective{Profile: "staging"}, Registry: registry.NewBuilder().Finalize()}, nil
	})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"run", "missing.example"})

	assertFailureExitCode(t, root.Execute())

	var got struct {
		OK    bool `json:"ok"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
		Meta struct {
			CapabilityID string `json:"capabilityId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("run output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.OK {
		t.Fatalf("ok = true, want false")
	}
	if got.Error.Code != "CAPABILITY_NOT_FOUND" {
		t.Fatalf("error.code = %q, want CAPABILITY_NOT_FOUND", got.Error.Code)
	}
	if got.Meta.CapabilityID != "missing.example" {
		t.Fatalf("meta.capabilityId = %q, want missing.example", got.Meta.CapabilityID)
	}
}

func TestRunCommandRequiresExplicitConfirmationForRiskyCapability(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	handlerCalled := false
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{ID: "orders.cancel", RequiresConfirmation: true},
		Handler: func(context.Context, capability.ExecutionRequest) (capability.ExecutionResult, error) {
			handlerCalled = true
			return capability.ExecutionResult{Data: map[string]any{"cancelled": true}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		return &app.Application{Config: config.Effective{Profile: "staging"}, Registry: builder.Finalize()}, nil
	})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"run", "orders.cancel"})

	assertFailureExitCode(t, root.Execute())

	if handlerCalled {
		t.Fatalf("handler was called without --confirm")
	}
	var got struct {
		OK    bool `json:"ok"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
		Meta struct {
			CapabilityID string `json:"capabilityId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("run output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.OK {
		t.Fatalf("ok = true, want false")
	}
	if got.Error.Code != "CONFIRMATION_REQUIRED" {
		t.Fatalf("error.code = %q, want CONFIRMATION_REQUIRED", got.Error.Code)
	}
	if got.Meta.CapabilityID != "orders.cancel" {
		t.Fatalf("meta.capabilityId = %q, want orders.cancel", got.Meta.CapabilityID)
	}
}

func TestRunCommandPassesExplicitConfirmation(t *testing.T) {
	t.Parallel()

	builder := registry.NewBuilder()
	handlerCalled := false
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: capability.Definition{ID: "orders.cancel", RequiresConfirmation: true},
		Handler: func(context.Context, capability.ExecutionRequest) (capability.ExecutionResult, error) {
			handlerCalled = true
			return capability.ExecutionResult{Data: map[string]any{"cancelled": true}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		return &app.Application{Config: config.Effective{Profile: "staging"}, Registry: builder.Finalize()}, nil
	})

	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"run", "orders.cancel", "--confirm"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !handlerCalled {
		t.Fatalf("handler was not called with --confirm")
	}
	var got struct {
		OK   bool `json:"ok"`
		Data struct {
			Cancelled bool `json:"cancelled"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("run output is not JSON: %v\n%s", err, stdout.String())
	}
	if !got.OK || !got.Data.Cancelled {
		t.Fatalf("output = %#v, want successful cancellation result", got)
	}
}

func assertFailureExitCode(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Fatalf("Execute() error = nil, want failure exit code")
	}
	if got := clisurface.ExitCode(err); got != clisurface.ExitCodeFailure {
		t.Fatalf("ExitCode(error) = %d, want %d (err: %v)", got, clisurface.ExitCodeFailure, err)
	}
}

func TestConfigSetDefaultProfileCommandWritesJSONEnvelopeAndConfig(t *testing.T) {
	workDir := t.TempDir()
	configPath := filepath.Join(workDir, "team.yaml")
	if err := os.WriteFile(configPath, []byte("defaultProfile: staging\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	root := clisurface.NewRoot(app.New)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"--config", configPath, "--correlation-id", "corr-profile", "config", "set-default-profile", "prod"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	content, err := os.ReadFile(configPath) // #nosec G304 -- test reads its own temporary config file.
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "defaultProfile: prod\n" {
		t.Fatalf("config file = %q, want updated default profile", string(content))
	}

	var got struct {
		OK   bool `json:"ok"`
		Data struct {
			Profile      string        `json:"profile"`
			ConfigPath   string        `json:"configPath"`
			ConfigSource config.Source `json:"configSource"`
		} `json:"data"`
		Meta struct {
			RequestID     string `json:"requestId"`
			CorrelationID string `json:"correlationId"`
			Profile       string `json:"profile"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("config output is not JSON: %v\n%s", err, stdout.String())
	}
	if !got.OK {
		t.Fatalf("ok = false, want true")
	}
	if got.Data.Profile != "prod" || got.Data.ConfigPath != configPath || got.Data.ConfigSource != config.SourceExplicit {
		t.Fatalf("data = %#v, want persisted profile data", got.Data)
	}
	if got.Meta.RequestID == "" || got.Meta.CorrelationID != "corr-profile" || got.Meta.Profile != "prod" {
		t.Fatalf("meta = %#v, want request, correlation, and profile", got.Meta)
	}
}

func TestConfigSetDefaultProfileCommandCreatesLocalConfigByDefault(t *testing.T) {
	workDir := t.TempDir()
	t.Chdir(workDir)

	root := clisurface.NewRoot(app.New)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"config", "set-default-profile", "qa"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	configPath := filepath.Join(workDir, "exito.yaml")
	content, err := os.ReadFile(configPath) // #nosec G304 -- test reads its own temporary config file.
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "defaultProfile: qa\n" {
		t.Fatalf("created config file = %q, want default profile", string(content))
	}

	var got struct {
		Data struct {
			ConfigPath   string        `json:"configPath"`
			ConfigSource config.Source `json:"configSource"`
		} `json:"data"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("config output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.Data.ConfigPath != configPath || got.Data.ConfigSource != config.SourceLocalProject {
		t.Fatalf("data = %#v, want local config target", got.Data)
	}
}

func TestConfigSetDefaultProfileCommandRejectsBlankProfile(t *testing.T) {
	workDir := t.TempDir()
	t.Chdir(workDir)

	root := clisurface.NewRoot(app.New)
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetErr(&output)
	root.SetArgs([]string{"config", "set-default-profile", "   "})

	if err := root.Execute(); err == nil {
		t.Fatalf("Execute() error = nil, want validation error")
	}
	if _, err := os.Stat(filepath.Join(workDir, "exito.yaml")); !os.IsNotExist(err) {
		t.Fatalf("blank profile should not create config file, stat error = %v", err)
	}
}

func TestOrdersGetVTEXCommandRunsOrdersGetVTEXCapability(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(app.New)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"--profile", "prod", "--correlation-id", "corr-vtex", "orders", "get-vtex", "--id", "1611511090420-01", "--brand", "carulla"})

	assertFailureExitCode(t, root.Execute())

	var got struct {
		OK    bool `json:"ok"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
		Meta struct {
			CorrelationID string `json:"correlationId"`
			Profile       string `json:"profile"`
			CapabilityID  string `json:"capabilityId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("orders get-vtex output is not JSON: %v\n%s", err, stdout.String())
	}
	if got.OK {
		t.Fatalf("ok = true, want false until VTEX OMS client is configured")
	}
	if got.Error.Code != orders.ErrorOrdersNotConfigured {
		t.Fatalf("error.code = %q, want %s", got.Error.Code, orders.ErrorOrdersNotConfigured)
	}
	if got.Meta.CorrelationID != "corr-vtex" || got.Meta.Profile != "prod" || got.Meta.CapabilityID != orders.CapabilityGetVTEXID {
		t.Fatalf("unexpected metadata: %#v", got.Meta)
	}
}

func TestOrdersGetVTEXCommandRequiresID(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(app.New)
	var output bytes.Buffer
	root.SetOut(&output)
	root.SetErr(&output)
	root.SetArgs([]string{"orders", "get-vtex"})

	if err := root.Execute(); err == nil {
		t.Fatalf("Execute() error = nil, want missing required flag error")
	}
	if strings.Contains(output.String(), "{\"ok\"") {
		t.Fatalf("missing flag error should not emit JSON envelope\n%s", output.String())
	}
}

func TestCatalogSearchProductsCommandRunsCapability(t *testing.T) {
	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		builder := registry.NewBuilder()
		if err := builder.RegisterExecutable(capability.Executable{
			Definition: catalog.Definition(),
			Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
				if request.Input["brand"] != "carulla" || request.Input["by"] != "sku-id" || request.Input["value"] != "912350" {
					t.Fatalf("input = %#v, want catalog search flags", request.Input)
				}
				if request.Input["from"] != 0 || request.Input["to"] != 0 {
					t.Fatalf("pagination input = %#v", request.Input)
				}
				return capability.ExecutionResult{Data: catalog.SearchProductsResult{Brand: "carulla", Count: 1, Products: []catalog.Product{{ProductID: "534690"}}}}, nil
			},
		}); err != nil {
			t.Fatalf("RegisterExecutable() error = %v", err)
		}
		return &app.Application{Config: config.Effective{Profile: "staging"}, ConfigOptions: options.Config, Registry: builder.Finalize()}, nil
	})
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"--correlation-id", "corr-catalog", "catalog", "search-products", "--brand", "carulla", "--by", "sku-id", "--value", "912350", "--from", "0", "--to", "0"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	var got struct {
		OK   bool `json:"ok"`
		Data struct {
			Brand string `json:"brand"`
			Count int    `json:"count"`
		} `json:"data"`
		Meta struct {
			CorrelationID string `json:"correlationId"`
			CapabilityID  string `json:"capabilityId"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("catalog search output is not JSON: %v\n%s", err, stdout.String())
	}
	if !got.OK || got.Data.Brand != "carulla" || got.Data.Count != 1 {
		t.Fatalf("output = %#v, want successful catalog result", got)
	}
	if got.Meta.CorrelationID != "corr-catalog" || got.Meta.CapabilityID != catalog.CapabilitySearchProductsID {
		t.Fatalf("metadata = %#v, want catalog capability", got.Meta)
	}
}

func TestCatalogSearchProductsCommandPassesAdvancedFilters(t *testing.T) {
	root := clisurface.NewRoot(func(options app.Options) (*app.Application, error) {
		builder := registry.NewBuilder()
		if err := builder.RegisterExecutable(capability.Executable{
			Definition: catalog.Definition(),
			Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
				filters, ok := request.Input["fq"].([]string)
				if !ok || len(filters) != 2 || filters[0] != "sellerId:VMIABBA" || filters[1] != "skuId:912350" {
					t.Fatalf("fq = %#v, want repeated filters", request.Input["fq"])
				}
				if request.Input["ft"] != "minibar" || request.Input["order"] != "OrderByPriceASC" {
					t.Fatalf("advanced input = %#v", request.Input)
				}
				return capability.ExecutionResult{Data: catalog.SearchProductsResult{Brand: "exito", Products: []catalog.Product{}}}, nil
			},
		}); err != nil {
			t.Fatalf("RegisterExecutable() error = %v", err)
		}
		return &app.Application{Config: config.Effective{Profile: "staging"}, ConfigOptions: options.Config, Registry: builder.Finalize()}, nil
	})
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"catalog", "search-products", "--fq", "sellerId:VMIABBA", "--fq", "skuId:912350", "--ft", "minibar", "--order", "OrderByPriceASC"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v\nstdout: %s", err, stdout.String())
	}
}
