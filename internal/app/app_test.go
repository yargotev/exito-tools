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

func TestNewWiresOrdersGetCapability(t *testing.T) {
	t.Parallel()

	application, err := app.New(app.Options{Config: config.Options{Env: map[string]string{}}})
	if err != nil {
		t.Fatalf("app.New() error = %v", err)
	}

	entry, ok := application.Registry.Find("orders.get")
	if !ok {
		t.Fatalf("orders.get capability was not registered")
	}
	if entry.Definition.Domain != "orders" {
		t.Fatalf("Domain = %q, want orders", entry.Definition.Domain)
	}
	if entry.Handler == nil {
		t.Fatalf("orders.get handler is nil")
	}
}
