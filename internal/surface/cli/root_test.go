package cli_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/config"
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
			} {
				if !strings.Contains(rendered, fragment) {
					t.Fatalf("help output missing %q\n%s", fragment, rendered)
				}
			}

			for _, forbidden := range []string{"orders", "geo", "\"ok\"", "\"data\"", "\"error\""} {
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

func TestRunCommandExecutesBootstrappedOrdersGetCapability(t *testing.T) {
	t.Parallel()

	root := clisurface.NewRoot(app.New)
	var stdout bytes.Buffer
	root.SetOut(&stdout)
	root.SetErr(&stdout)
	root.SetArgs([]string{"run", orders.CapabilityGetID, "--input-json", `{"id":"A123"}`})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
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
		t.Fatalf("ok = true, want false until Orders client is configured")
	}
	if got.Error.Code != orders.ErrorOrdersNotConfigured {
		t.Fatalf("error.code = %q, want %s", got.Error.Code, orders.ErrorOrdersNotConfigured)
	}
	if got.Meta.CapabilityID != orders.CapabilityGetID {
		t.Fatalf("meta.capabilityId = %q, want %s", got.Meta.CapabilityID, orders.CapabilityGetID)
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

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
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

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
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
	if got.Error.Code != "CAPABILITY_NOT_FOUND" {
		t.Fatalf("error.code = %q, want CAPABILITY_NOT_FOUND", got.Error.Code)
	}
	if got.Meta.CapabilityID != "missing.example" {
		t.Fatalf("meta.capabilityId = %q, want missing.example", got.Meta.CapabilityID)
	}
}
