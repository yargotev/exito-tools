package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/app"
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
				"This foundation slice only provides bootstrap and help.",
				"Usage:",
			} {
				if !strings.Contains(rendered, fragment) {
					t.Fatalf("help output missing %q\n%s", fragment, rendered)
				}
			}

			for _, forbidden := range []string{"orders", "geo", "capabilities", " run ", "\"ok\"", "\"data\"", "\"error\""} {
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
