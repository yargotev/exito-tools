package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/app"
	clisurface "github.com/yargotev/exito-tools/internal/surface/cli"
)

func TestRootHelpPaths(t *testing.T) {
	t.Parallel()

	application, err := app.New()
	if err != nil {
		t.Fatalf("app.New() error = %v", err)
	}

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

			root := clisurface.NewRoot(application)
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
