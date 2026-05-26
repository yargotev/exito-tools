package geo_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/domain/geo"
)

func TestHTTPGeocoderPostsRequestAndMapsProviderResponse(t *testing.T) {
	t.Parallel()

	var gotPath string
	var gotAuth string
	var gotBody struct {
		City    string `json:"city"`
		Address string `json:"address"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotAuth = r.Header.Get("Authorization")
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Fatalf("Decode(request body) error = %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"message":"Geocoding successful.",
			"success":true,
			"data":{
				"latitude":"4.598090587",
				"longitude":"-74.160580822",
				"estado":"M",
				"dirtrad":"CL 57 H SUR # 68 D - 75",
				"barrio":"VILLA DEL RIO",
				"coddane":"110010001"
			}
		}`))
	}))
	defer server.Close()

	geocoder := geo.NewHTTPGeocoder(geo.HTTPGeocoderConfig{BaseURL: server.URL, Token: "token-123"}, server.Client())
	got, err := geocoder.GeocodeAddress(context.Background(), geo.GeocodeAddressInput{City: "Bogota", Address: "CL 57 H SUR # 68 D - 75"})
	if err != nil {
		t.Fatalf("GeocodeAddress() error = %v", err)
	}

	if gotPath != "/geocode-address" {
		t.Fatalf("request path = %q, want /geocode-address", gotPath)
	}
	if gotAuth != "Bearer token-123" {
		t.Fatalf("Authorization = %q, want bearer token", gotAuth)
	}
	if gotBody.City != "Bogota" || gotBody.Address != "CL 57 H SUR # 68 D - 75" {
		t.Fatalf("request body = %#v, want city/address", gotBody)
	}
	if got.Message != "Geocoding successful." || !got.Success || got.Location.Latitude != "4.598090587" || got.Location.Longitude != "-74.160580822" || got.Status != "M" || got.NormalizedAddress != "CL 57 H SUR # 68 D - 75" || got.Neighborhood != "VILLA DEL RIO" || got.DANECode != "110010001" {
		t.Fatalf("mapped result = %#v, want provider fields mapped", got)
	}
}

func TestHTTPGeocoderReturnsStructuredErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		geocoder geo.HTTPGeocoder
		wantCode string
	}{
		{
			name:     "missing config",
			geocoder: geo.NewHTTPGeocoder(geo.HTTPGeocoderConfig{}, nil),
			wantCode: geo.ErrorGeoNotConfigured,
		},
		{
			name:     "invalid base URL",
			geocoder: geo.NewHTTPGeocoder(geo.HTTPGeocoderConfig{BaseURL: "://bad", Token: "token"}, nil),
			wantCode: geo.ErrorGeoNotConfigured,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := tc.geocoder.GeocodeAddress(context.Background(), geo.GeocodeAddressInput{City: "Bogota", Address: "Avenida Siempre Viva"})
			assertStructuredCode(t, err, tc.wantCode)
		})
	}
}

func TestHTTPGeocoderHandlesProviderFailures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		handler  http.HandlerFunc
		wantCode string
	}{
		{
			name: "non success status",
			handler: func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "provider down", http.StatusBadGateway)
			},
			wantCode: geo.ErrorGeoProviderUnavailable,
		},
		{
			name: "invalid JSON",
			handler: func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`not-json`))
			},
			wantCode: geo.ErrorGeoProviderInvalidResponse,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(tc.handler)
			defer server.Close()

			geocoder := geo.NewHTTPGeocoder(geo.HTTPGeocoderConfig{BaseURL: server.URL, Token: "token"}, server.Client())
			_, err := geocoder.GeocodeAddress(context.Background(), geo.GeocodeAddressInput{City: "Bogota", Address: "Avenida Siempre Viva"})
			assertStructuredCode(t, err, tc.wantCode)
		})
	}
}

func assertStructuredCode(t *testing.T, err error, want string) {
	t.Helper()

	var structured capability.StructuredError
	if !errors.As(err, &structured) {
		t.Fatalf("error = %T, want StructuredError", err)
	}
	if structured.Code != want {
		t.Fatalf("StructuredError.Code = %q, want %q", structured.Code, want)
	}
}
