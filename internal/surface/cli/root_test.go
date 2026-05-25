package cli_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/config"
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

			for _, forbidden := range []string{"orders", "geo", " run ", "\"ok\"", "\"data\"", "\"error\""} {
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
	for _, fragment := range []string{"--config", "--profile"} {
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
		Title:       "Foundation Example",
		Description: "Registered during application boot.",
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
			Profile string `json:"profile"`
		} `json:"meta"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("capabilities output is not JSON: %v\n%s", err, stdout.String())
	}

	if !got.OK {
		t.Fatalf("ok = false, want true")
	}
	if got.Meta.Profile != "staging" {
		t.Fatalf("meta.profile = %q, want staging", got.Meta.Profile)
	}
	if len(got.Data.Capabilities) != 1 {
		t.Fatalf("capabilities length = %d, want 1", len(got.Data.Capabilities))
	}
	if got.Data.Capabilities[0].ID != "foundation.example" {
		t.Fatalf("capability ID = %q, want foundation.example", got.Data.Capabilities[0].ID)
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
