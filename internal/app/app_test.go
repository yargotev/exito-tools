package app_test

import (
	"testing"

	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/config"
)

func TestNewResolvesConfigurationAtBoot(t *testing.T) {
	t.Parallel()

	application, err := app.New(app.Options{
		Config: config.Options{
			ConfigPath: "./exito.yaml",
			Profile:    "prod",
			Env:        map[string]string{"EXITO_PROFILE": "qa", "EXITO_CONFIG": "env.yaml"},
			WorkDir:    "/workspace/project",
			HomeDir:    "/home/tester",
		},
	})
	if err != nil {
		t.Fatalf("app.New() error = %v", err)
	}

	if application.Config.Profile != "prod" {
		t.Fatalf("Profile = %q, want prod", application.Config.Profile)
	}
	if application.Config.ProfileSource != config.SourceExplicit {
		t.Fatalf("ProfileSource = %q, want %q", application.Config.ProfileSource, config.SourceExplicit)
	}
	if application.Config.ConfigPath != "/workspace/project/exito.yaml" {
		t.Fatalf("ConfigPath = %q, want /workspace/project/exito.yaml", application.Config.ConfigPath)
	}
	if application.Config.ConfigSource != config.SourceExplicit {
		t.Fatalf("ConfigSource = %q, want %q", application.Config.ConfigSource, config.SourceExplicit)
	}
}

func TestNewWiresBootCapabilities(t *testing.T) {
	t.Parallel()

	application, err := app.New(app.Options{Config: config.Options{Env: map[string]string{}}})
	if err != nil {
		t.Fatalf("app.New() error = %v", err)
	}

	tests := []struct {
		id     string
		domain string
	}{
		{id: "orders.get", domain: "orders"},
		{id: "geo.geocode-address", domain: "geo"},
	}

	for _, tt := range tests {
		entry, ok := application.Registry.Find(tt.id)
		if !ok {
			t.Fatalf("%s capability was not registered", tt.id)
		}
		if entry.Definition.Domain != tt.domain {
			t.Fatalf("%s domain = %q, want %s", tt.id, entry.Definition.Domain, tt.domain)
		}
		if entry.Handler == nil {
			t.Fatalf("%s handler is nil", tt.id)
		}
	}
}
