package app_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/config"
	"github.com/yargotev/exito-tools/internal/domain/geo"
	"github.com/yargotev/exito-tools/internal/execution"
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

func TestNewWiresConfiguredGeoHTTPGeocoder(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/geocode-address" {
			t.Fatalf("request path = %q, want /geocode-address", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer token-123" {
			t.Fatalf("Authorization = %q, want bearer token", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"message":"ok","success":true,"data":{"latitude":"4.1","longitude":"-74.1","estado":"M","dirtrad":"NORMALIZED","barrio":"BARRIO","coddane":"11001"}}`))
	}))
	defer server.Close()

	application, err := app.New(app.Options{Config: config.Options{Env: map[string]string{
		"EXITO_GEO_BASE_URL": server.URL,
		"EXITO_GEO_TOKEN":    "token-123",
	}}})
	if err != nil {
		t.Fatalf("app.New() error = %v", err)
	}

	envelope, err := execution.NewPipeline(application.Registry).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID: geo.CapabilityGeocodeAddressID,
		Input:        capability.Input{"city": "Bogota", "address": "Avenida Siempre Viva"},
		Profile:      application.Config.Profile,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !envelope.OK {
		t.Fatalf("OK = false, want true: %#v", envelope.Error)
	}
	got, ok := (*envelope.Data).(geo.GeocodeAddressResult)
	if !ok {
		t.Fatalf("Data = %T, want geo.GeocodeAddressResult", *envelope.Data)
	}
	if got.NormalizedAddress != "NORMALIZED" || got.DANECode != "11001" {
		t.Fatalf("result = %#v, want mapped provider data", got)
	}
}
