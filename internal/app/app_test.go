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
	"github.com/yargotev/exito-tools/internal/domain/orders"
	"github.com/yargotev/exito-tools/internal/execution"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
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

	var gotRequestID string
	var gotCorrelationID string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestID = r.Header.Get(httpclient.HeaderRequestID)
		gotCorrelationID = r.Header.Get(httpclient.HeaderCorrelationID)
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

	envelope, err := execution.NewPipeline(
		application.Registry,
		execution.WithRequestIDGenerator(func() (string, error) { return "req_app_geo", nil }),
	).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID:  geo.CapabilityGeocodeAddressID,
		Input:         capability.Input{"city": "Bogota", "address": "Avenida Siempre Viva"},
		Profile:       application.Config.Profile,
		CorrelationID: "corr-app",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !envelope.OK {
		t.Fatalf("OK = false, want true: %#v", envelope.Error)
	}
	if gotRequestID != "req_app_geo" {
		t.Fatalf("%s = %q, want req_app_geo", httpclient.HeaderRequestID, gotRequestID)
	}
	if gotCorrelationID != "corr-app" {
		t.Fatalf("%s = %q, want corr-app", httpclient.HeaderCorrelationID, gotCorrelationID)
	}
	got, ok := (*envelope.Data).(geo.GeocodeAddressResult)
	if !ok {
		t.Fatalf("Data = %T, want geo.GeocodeAddressResult", *envelope.Data)
	}
	if got.NormalizedAddress != "NORMALIZED" || got.DANECode != "11001" {
		t.Fatalf("result = %#v, want mapped provider data", got)
	}
}

func TestNewWiresConfiguredOrdersHTTPGetter(t *testing.T) {
	t.Parallel()

	var gotRequestID string
	var gotCorrelationID string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotRequestID = r.Header.Get(httpclient.HeaderRequestID)
		gotCorrelationID = r.Header.Get(httpclient.HeaderCorrelationID)
		if r.URL.Path != "/orders/get" {
			t.Fatalf("request path = %q, want /orders/get", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer orders-token" {
			t.Fatalf("Authorization = %q, want bearer token", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"order":{"id":"A123","status":"created","createdAt":"2026-05-26T00:00:00Z"}}`))
	}))
	defer server.Close()

	application, err := app.New(app.Options{Config: config.Options{Env: map[string]string{
		"EXITO_ORDERS_BASE_URL": server.URL,
		"EXITO_ORDERS_TOKEN":    "orders-token",
	}}})
	if err != nil {
		t.Fatalf("app.New() error = %v", err)
	}

	envelope, err := execution.NewPipeline(
		application.Registry,
		execution.WithRequestIDGenerator(func() (string, error) { return "req_app_orders", nil }),
	).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID:  orders.CapabilityGetID,
		Input:         capability.Input{"id": "A123"},
		Profile:       application.Config.Profile,
		CorrelationID: "corr-orders-app",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !envelope.OK {
		t.Fatalf("OK = false, want true: %#v", envelope.Error)
	}
	if gotRequestID != "req_app_orders" {
		t.Fatalf("%s = %q, want req_app_orders", httpclient.HeaderRequestID, gotRequestID)
	}
	if gotCorrelationID != "corr-orders-app" {
		t.Fatalf("%s = %q, want corr-orders-app", httpclient.HeaderCorrelationID, gotCorrelationID)
	}
	got, ok := (*envelope.Data).(orders.GetResult)
	if !ok {
		t.Fatalf("Data = %T, want orders.GetResult", *envelope.Data)
	}
	if got.Order.ID != "A123" || got.Order.Status != "created" {
		t.Fatalf("result = %#v, want mapped provider order", got)
	}
}
