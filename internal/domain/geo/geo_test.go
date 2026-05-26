package geo_test

import (
	"context"
	"errors"
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
