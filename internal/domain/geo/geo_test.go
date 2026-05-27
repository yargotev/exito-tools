package geo_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/domain/geo"
	"github.com/yargotev/exito-tools/internal/execution"
	"github.com/yargotev/exito-tools/internal/registry"
)

func TestDefinition(t *testing.T) {
	t.Parallel()

	definition := geo.Definition()
	if definition.ID != geo.CapabilityGeocodeAddressID {
		t.Fatalf("ID = %q, want %q", definition.ID, geo.CapabilityGeocodeAddressID)
	}
	if definition.Domain != geo.DomainName {
		t.Fatalf("Domain = %q, want %q", definition.Domain, geo.DomainName)
	}
	if definition.Risk != capability.RiskReadOnly {
		t.Fatalf("Risk = %q, want %q", definition.Risk, capability.RiskReadOnly)
	}
	if len(definition.Audiences) != 2 || definition.Audiences[0] != capability.AudienceAgents || definition.Audiences[1] != capability.AudiencePeople {
		t.Fatalf("Audiences = %#v, want agents and people", definition.Audiences)
	}
	if len(definition.Visibility) != 3 || definition.Visibility[0] != capability.VisibilityCLI || definition.Visibility[1] != capability.VisibilityTUI || definition.Visibility[2] != capability.VisibilityCommandPalette {
		t.Fatalf("Visibility = %#v, want cli, tui, command-palette", definition.Visibility)
	}
	if definition.InputSchema == nil || len(definition.InputSchema.Fields) != 2 {
		t.Fatalf("InputSchema = %#v, want city and address fields", definition.InputSchema)
	}
	if field := definition.InputSchema.Fields[0]; field.Name != "city" || field.Type != capability.InputTypeString || !field.Required {
		t.Fatalf("first input field = %#v, want required string city", field)
	}
	if field := definition.InputSchema.Fields[1]; field.Name != "address" || field.Type != capability.InputTypeString || !field.Required {
		t.Fatalf("second input field = %#v, want required string address", field)
	}
}

func TestGeocodeAddressCapabilityExecutesUseCase(t *testing.T) {
	t.Parallel()

	result := geo.GeocodeAddressResult{
		Message:           "Geocoding successful.",
		Success:           true,
		Location:          geo.Location{Latitude: "4.598090587", Longitude: "-74.160580822"},
		Status:            "M",
		NormalizedAddress: "CL 57 H SUR # 68 D - 75",
		Neighborhood:      "VILLA DEL RIO",
		DANECode:          "110010001",
	}
	geocoder := &fakeGeocoder{result: result}
	envelope, err := pipelineWithGeocoder(t, geocoder).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID: geo.CapabilityGeocodeAddressID,
		Input:        capability.Input{"city": "Bogota", "address": "CL 57 H SUR # 68 D - 75"},
		Profile:      "staging",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if !envelope.OK {
		t.Fatalf("OK = false, want true: %#v", envelope.Error)
	}
	if geocoder.got.City != "Bogota" || geocoder.got.Address != "CL 57 H SUR # 68 D - 75" {
		t.Fatalf("geocoder input = %#v, want city and address", geocoder.got)
	}
	gotResult, ok := (*envelope.Data).(geo.GeocodeAddressResult)
	if !ok {
		t.Fatalf("Data = %T, want geo.GeocodeAddressResult", *envelope.Data)
	}
	if gotResult.Location.Latitude != "4.598090587" || gotResult.DANECode != "110010001" {
		t.Fatalf("Result = %#v, want fake geocode result", gotResult)
	}
}

func TestGeocodeAddressCapabilityPropagatesStructuredDomainError(t *testing.T) {
	t.Parallel()

	envelope, err := pipelineWithGeocoder(t, &fakeGeocoder{err: capability.StructuredError{Code: geo.ErrorGeoNotConfigured, Message: "Geo client is not configured."}}).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID: geo.CapabilityGeocodeAddressID,
		Input:        capability.Input{"city": "Bogota", "address": "missing"},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	if envelope.OK {
		t.Fatalf("OK = true, want false")
	}
	if envelope.Error == nil || envelope.Error.Code != geo.ErrorGeoNotConfigured {
		t.Fatalf("Error = %#v, want %s", envelope.Error, geo.ErrorGeoNotConfigured)
	}
}

func TestUnavailableGeocoderReturnsStructuredError(t *testing.T) {
	t.Parallel()

	_, err := geo.UnavailableGeocoder{}.GeocodeAddress(context.Background(), geo.GeocodeAddressInput{City: "Bogota", Address: "Avenida Siempre Viva"})
	var structured capability.StructuredError
	if !errors.As(err, &structured) {
		t.Fatalf("GeocodeAddress() error = %T, want StructuredError", err)
	}
	if structured.Code != geo.ErrorGeoNotConfigured {
		t.Fatalf("StructuredError.Code = %q, want %q", structured.Code, geo.ErrorGeoNotConfigured)
	}
}

func pipelineWithGeocoder(t *testing.T, geocoder geo.Geocoder) execution.Pipeline {
	t.Helper()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(geo.NewGeocodeAddressCapability(geocoder)); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}
	return execution.NewPipeline(
		builder.Finalize(),
		execution.WithRequestIDGenerator(func() (string, error) { return "req_test", nil }),
	)
}

type fakeGeocoder struct {
	result geo.GeocodeAddressResult
	err    error
	got    geo.GeocodeAddressInput
}

func (g *fakeGeocoder) GeocodeAddress(_ context.Context, input geo.GeocodeAddressInput) (geo.GeocodeAddressResult, error) {
	g.got = input
	if g.err != nil {
		return geo.GeocodeAddressResult{}, g.err
	}
	return g.result, nil
}

func TestResolveVTEXRegionCapabilityExecutesUseCase(t *testing.T) {
	t.Parallel()

	resolver := &fakeVTEXRegionResolver{result: geo.ResolveVTEXRegionResult{Brand: "carulla", HasCoverage: true}}
	envelope, err := pipelineWithVTEXRegionResolver(t, resolver).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID: geo.CapabilityResolveVTEXRegionID,
		Input: capability.Input{
			"brand":        "carulla",
			"country":      "COL",
			"salesChannel": "1",
			"longitude":    "-74.160580822",
			"latitude":     "4.598090587",
		},
		Profile: "staging",
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !envelope.OK {
		t.Fatalf("OK = false, want true: %#v", envelope.Error)
	}
	if resolver.got.Brand != "carulla" || resolver.got.Country != "COL" || resolver.got.SalesChannel != "1" || resolver.got.Longitude != "-74.160580822" || resolver.got.Latitude != "4.598090587" {
		t.Fatalf("resolver input = %#v", resolver.got)
	}
}

func TestHTTPVTEXRegionResolverBuildsCoordinatesAndCoverage(t *testing.T) {
	t.Parallel()

	var gotPath string
	var gotQuery string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotQuery = r.URL.RawQuery
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"REGION-1","sellers":[{"id":"exito","name":"Account"},{"id":"seller-2","name":"Marketplace"}]}]`))
	}))
	defer server.Close()

	resolver := geo.NewHTTPVTEXRegionResolver(geo.HTTPVTEXRegionResolverConfig{BaseURL: server.URL}, server.Client())
	result, err := resolver.ResolveVTEXRegion(context.Background(), geo.ResolveVTEXRegionInput{
		Brand:        "exito",
		Country:      "COL",
		SalesChannel: "1",
		Longitude:    "-74.160580822",
		Latitude:     "4.598090587",
	})
	if err != nil {
		t.Fatalf("ResolveVTEXRegion() error = %v", err)
	}
	if gotPath != "/api/checkout/pub/regions" {
		t.Fatalf("path = %q", gotPath)
	}
	for _, want := range []string{"country=COL", "sc=1", "geoCoordinates=-74.160580822%3B4.598090587"} {
		if !strings.Contains(gotQuery, want) {
			t.Fatalf("query = %q, missing %s", gotQuery, want)
		}
	}
	if !result.HasCoverage || len(result.Sellers) != 2 || result.Sellers[1].ID != "seller-2" {
		t.Fatalf("result = %#v, want coverage from non-account seller", result)
	}
	if result.Diagnostics.RequestQuery["geoCoordinates"] != "-74.160580822;4.598090587" {
		t.Fatalf("diagnostic query = %#v", result.Diagnostics.RequestQuery)
	}
}

func TestHTTPVTEXRegionResolverCoverageTrueForAccountSeller(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"sellers":[{"id":"exito","name":"Account"}]}]`))
	}))
	defer server.Close()

	resolver := geo.NewHTTPVTEXRegionResolver(geo.HTTPVTEXRegionResolverConfig{BaseURL: server.URL}, server.Client())
	result, err := resolver.ResolveVTEXRegion(context.Background(), geo.ResolveVTEXRegionInput{Brand: "exito", Country: "COL", SalesChannel: "1", Longitude: "-74", Latitude: "4"})
	if err != nil {
		t.Fatalf("ResolveVTEXRegion() error = %v", err)
	}
	if !result.HasCoverage {
		t.Fatalf("HasCoverage = false, want true for account seller returned by Regions")
	}
}

func TestHTTPVTEXRegionResolverCoverageFalseForNoSellers(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"id":"REGION-EMPTY","sellers":[]}]`))
	}))
	defer server.Close()

	resolver := geo.NewHTTPVTEXRegionResolver(geo.HTTPVTEXRegionResolverConfig{BaseURL: server.URL}, server.Client())
	result, err := resolver.ResolveVTEXRegion(context.Background(), geo.ResolveVTEXRegionInput{Brand: "exito", Country: "COL", SalesChannel: "1", Longitude: "-74", Latitude: "4"})
	if err != nil {
		t.Fatalf("ResolveVTEXRegion() error = %v", err)
	}
	if result.HasCoverage {
		t.Fatalf("HasCoverage = true, want false when Regions returns no sellers")
	}
}

func pipelineWithVTEXRegionResolver(t *testing.T, resolver geo.VTEXRegionResolver) execution.Pipeline {
	t.Helper()

	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(geo.NewResolveVTEXRegionCapability(resolver)); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}
	return execution.NewPipeline(
		builder.Finalize(),
		execution.WithRequestIDGenerator(func() (string, error) { return "req_test", nil }),
	)
}

type fakeVTEXRegionResolver struct {
	result geo.ResolveVTEXRegionResult
	err    error
	got    geo.ResolveVTEXRegionInput
}

func (r *fakeVTEXRegionResolver) ResolveVTEXRegion(_ context.Context, input geo.ResolveVTEXRegionInput) (geo.ResolveVTEXRegionResult, error) {
	r.got = input
	if r.err != nil {
		return geo.ResolveVTEXRegionResult{}, r.err
	}
	return r.result, nil
}
