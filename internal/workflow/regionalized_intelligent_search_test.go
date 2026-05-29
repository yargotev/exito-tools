package workflow_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/domain/catalog"
	"github.com/yargotev/exito-tools/internal/domain/geo"
	"github.com/yargotev/exito-tools/internal/execution"
	"github.com/yargotev/exito-tools/internal/registry"
	"github.com/yargotev/exito-tools/internal/workflow"
)

func TestRegionalizedIntelligentSearchExecutesRegionSegmentAndSearch(t *testing.T) {
	resolver := &fakeRegionResolver{result: geo.ResolveVTEXRegionResult{
		Brand:        "exito",
		Country:      "COL",
		SalesChannel: "1",
		Coordinates:  geo.Coordinates{Longitude: "-74", Latitude: "4"},
		HasCoverage:  true,
		Regions:      []geo.Region{{ID: "REGION-1", Sellers: []geo.RegionSeller{{ID: "seller-1"}}}},
		Sellers:      []geo.RegionSeller{{ID: "seller-1"}},
	}}
	segmenter := &fakeSegmentCreator{result: catalog.CreateVTEXSegmentResult{Brand: "exito", RegionID: "REGION-1", SalesChannel: "1", TokenSet: true, TokenLength: len("secret-token"), Cookie: "vtex_segment=secret-token"}}
	searcher := &fakeIntelligentSearcher{result: catalog.IntelligentSearchProductsResult{Brand: "exito", Query: "arroz", Products: []catalog.Product{{ProductID: "P1"}}}}
	useCase := workflow.NewRegionalizedIntelligentSearchProductsUseCase(resolver, segmenter, searcher)

	result, err := useCase.Execute(context.Background(), workflow.RegionalizedIntelligentSearchProductsInput{
		Brand:       "exito",
		Country:     "COL",
		TradePolicy: "1",
		Longitude:   "-74",
		Latitude:    "4",
		Text:        "arroz",
		Count:       3,
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if resolver.got.SalesChannel != "1" || resolver.got.Longitude != "-74" || resolver.got.Latitude != "4" {
		t.Fatalf("region input = %#v", resolver.got)
	}
	if segmenter.got.RegionID != "REGION-1" || !segmenter.got.IncludeCookie {
		t.Fatalf("segment input = %#v, want selected region with internal cookie", segmenter.got)
	}
	if searcher.got.Text != "arroz" || len(searcher.got.Cookies) != 1 || searcher.got.Cookies[0] != "vtex_segment=secret-token" {
		t.Fatalf("search input = %#v, want generated segment cookie", searcher.got)
	}
	if result.Region.ID != "REGION-1" || !result.Segment.TokenSet || len(result.Search.Products) != 1 {
		t.Fatalf("result = %#v", result)
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if strings.Contains(string(encoded), "secret-token") || strings.Contains(string(encoded), "vtex_segment=secret-token") {
		t.Fatalf("workflow result leaked segment token: %s", string(encoded))
	}
}

func TestRegionalizedIntelligentSearchDefinitionIsPeopleFacingPaletteAction(t *testing.T) {
	definition := workflow.RegionalizedIntelligentSearchProductsDefinition()

	if definition.ID != workflow.CapabilityRegionalizedIntelligentSearchProductsID {
		t.Fatalf("ID = %q, want %q", definition.ID, workflow.CapabilityRegionalizedIntelligentSearchProductsID)
	}
	if !containsAudience(definition.Audiences, capability.AudiencePeople) || !containsAudience(definition.Audiences, capability.AudienceAgents) {
		t.Fatalf("Audiences = %#v, want agents and people", definition.Audiences)
	}
	if !containsVisibility(definition.Visibility, capability.VisibilityCLI) || !containsVisibility(definition.Visibility, capability.VisibilityCommandPalette) {
		t.Fatalf("Visibility = %#v, want cli and command-palette", definition.Visibility)
	}
	if containsVisibility(definition.Visibility, capability.VisibilityTUI) {
		t.Fatalf("Visibility = %#v, should not promote regionalized workflow as a primary TUI action", definition.Visibility)
	}
	if !definition.RequiresConfirmation || definition.Risk != capability.RiskSafeWrite {
		t.Fatalf("confirmation/risk = %v/%q, want confirmation-required safe-write", definition.RequiresConfirmation, definition.Risk)
	}
}

func TestRegionalizedIntelligentSearchRequiresConfirmationInPipeline(t *testing.T) {
	called := false
	builder := registry.NewBuilder()
	if err := builder.RegisterExecutable(capability.Executable{
		Definition: workflow.RegionalizedIntelligentSearchProductsDefinition(),
		Handler: func(context.Context, capability.ExecutionRequest) (capability.ExecutionResult, error) {
			called = true
			return capability.ExecutionResult{Data: map[string]bool{"called": true}}, nil
		},
	}); err != nil {
		t.Fatalf("RegisterExecutable() error = %v", err)
	}

	envelope, err := execution.NewPipeline(builder.Finalize()).Execute(context.Background(), execution.ExecuteRequest{
		CapabilityID: workflow.CapabilityRegionalizedIntelligentSearchProductsID,
		Input: capability.Input{
			"tradePolicy": "1",
			"longitude":   "-74",
			"latitude":    "4",
			"text":        "arroz",
		},
	})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if called {
		t.Fatalf("handler was called without confirmation")
	}
	if envelope.OK || envelope.Error == nil || envelope.Error.Code != execution.ErrorConfirmationRequired {
		t.Fatalf("envelope = %#v, want confirmation required", envelope)
	}
}

func TestRegionalizedIntelligentSearchStopsWhenNoRegionID(t *testing.T) {
	resolver := &fakeRegionResolver{result: geo.ResolveVTEXRegionResult{HasCoverage: false}}
	segmenter := &fakeSegmentCreator{result: catalog.CreateVTEXSegmentResult{Cookie: "vtex_segment=secret"}}
	searcher := &fakeIntelligentSearcher{}
	useCase := workflow.NewRegionalizedIntelligentSearchProductsUseCase(resolver, segmenter, searcher)

	_, err := useCase.Execute(context.Background(), workflow.RegionalizedIntelligentSearchProductsInput{TradePolicy: "1", Longitude: "-74", Latitude: "4", Text: "arroz"})
	if err == nil {
		t.Fatalf("Execute() error = nil, want no region error")
	}
	var structured capability.StructuredError
	if !strings.Contains(err.Error(), "region ID") {
		t.Fatalf("error = %v, want no region ID", err)
	}
	if segmenter.called || searcher.called {
		t.Fatalf("segment/search should not be called when no region is resolved")
	}
	_ = structured
}

func containsAudience(values []capability.Audience, target capability.Audience) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func containsVisibility(values []capability.Visibility, target capability.Visibility) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

type fakeRegionResolver struct {
	result geo.ResolveVTEXRegionResult
	err    error
	got    geo.ResolveVTEXRegionInput
}

func (r *fakeRegionResolver) ResolveVTEXRegion(_ context.Context, input geo.ResolveVTEXRegionInput) (geo.ResolveVTEXRegionResult, error) {
	r.got = input
	if r.err != nil {
		return geo.ResolveVTEXRegionResult{}, r.err
	}
	return r.result, nil
}

type fakeSegmentCreator struct {
	result catalog.CreateVTEXSegmentResult
	err    error
	got    catalog.CreateVTEXSegmentInput
	called bool
}

func (c *fakeSegmentCreator) CreateVTEXSegment(_ context.Context, input catalog.CreateVTEXSegmentInput) (catalog.CreateVTEXSegmentResult, error) {
	c.called = true
	c.got = input
	if c.err != nil {
		return catalog.CreateVTEXSegmentResult{}, c.err
	}
	return c.result, nil
}

type fakeIntelligentSearcher struct {
	result catalog.IntelligentSearchProductsResult
	err    error
	got    catalog.IntelligentSearchProductsInput
	called bool
}

func (s *fakeIntelligentSearcher) IntelligentSearchProducts(_ context.Context, input catalog.IntelligentSearchProductsInput) (catalog.IntelligentSearchProductsResult, error) {
	s.called = true
	s.got = input
	if s.err != nil {
		return catalog.IntelligentSearchProductsResult{}, s.err
	}
	return s.result, nil
}
