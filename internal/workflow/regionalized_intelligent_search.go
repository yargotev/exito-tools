package workflow

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/domain/catalog"
	"github.com/yargotev/exito-tools/internal/domain/geo"
)

const (
	// CapabilityRegionalizedIntelligentSearchProductsID identifies the Phase 4 regionalized product workflow.
	CapabilityRegionalizedIntelligentSearchProductsID = "catalog.regionalized-intelligent-search-products"

	ErrorRegionalizedSearchInvalidInput       = "REGIONALIZED_SEARCH_INVALID_INPUT"
	ErrorRegionalizedSearchNoRegion           = "REGIONALIZED_SEARCH_NO_REGION"
	ErrorRegionalizedSearchSegmentUnavailable = "REGIONALIZED_SEARCH_SEGMENT_UNAVAILABLE"
)

// RegionalizedIntelligentSearchProductsInput is the schema-shaped input accepted by the workflow.
type RegionalizedIntelligentSearchProductsInput struct {
	Brand              string
	Country            string
	TradePolicy        string
	Longitude          string
	Latitude           string
	Text               string
	By                 string
	Values             []string
	Query              string
	Facets             []string
	Page               int
	Count              int
	Sort               string
	Locale             string
	HideUnavailable    *bool
	SimulationBehavior string
}

// SelectedRegion records the VTEX region selected for the segment step.
type SelectedRegion struct {
	ID      string             `json:"id"`
	Sellers []geo.RegionSeller `json:"sellers,omitempty"`
}

// SegmentSummary exposes non-secret VTEX segment metadata.
type SegmentSummary struct {
	RegionID     string `json:"regionId"`
	SalesChannel string `json:"salesChannel"`
	TokenSet     bool   `json:"tokenSet"`
	TokenLength  int    `json:"tokenLength,omitempty"`
}

// RegionalizedIntelligentSearchProductsResult is the stable workflow result.
type RegionalizedIntelligentSearchProductsResult struct {
	Brand        string                                  `json:"brand"`
	Country      string                                  `json:"country"`
	TradePolicy  string                                  `json:"tradePolicy"`
	Coordinates  geo.Coordinates                         `json:"coordinates"`
	HasCoverage  bool                                    `json:"hasCoverage"`
	Region       SelectedRegion                          `json:"region"`
	Segment      SegmentSummary                          `json:"segment"`
	Search       catalog.IntelligentSearchProductsResult `json:"search"`
	RegionResult geo.ResolveVTEXRegionResult             `json:"regionResult"`
}

// RegionalizedIntelligentSearchProductsUseCase orchestrates region resolution, segment creation, and search.
type RegionalizedIntelligentSearchProductsUseCase struct {
	regionResolver geo.VTEXRegionResolver
	segmentCreator catalog.VTEXSegmentCreator
	searcher       catalog.IntelligentSearchProductsSearcher
}

// NewRegionalizedIntelligentSearchProductsUseCase creates the workflow use case.
func NewRegionalizedIntelligentSearchProductsUseCase(regionResolver geo.VTEXRegionResolver, segmentCreator catalog.VTEXSegmentCreator, searcher catalog.IntelligentSearchProductsSearcher) RegionalizedIntelligentSearchProductsUseCase {
	return RegionalizedIntelligentSearchProductsUseCase{regionResolver: regionResolver, segmentCreator: segmentCreator, searcher: searcher}
}

// Execute runs the Phase 4 regionalized Intelligent Search workflow.
func (u RegionalizedIntelligentSearchProductsUseCase) Execute(ctx context.Context, input RegionalizedIntelligentSearchProductsInput) (RegionalizedIntelligentSearchProductsResult, error) {
	input = normalizeRegionalizedInput(input)
	if err := validateRegionalizedInput(input); err != nil {
		return RegionalizedIntelligentSearchProductsResult{}, err
	}

	regionUseCase := geo.NewResolveVTEXRegionUseCase(u.regionResolver)
	regionResult, err := regionUseCase.Execute(ctx, geo.ResolveVTEXRegionInput{
		Brand:        input.Brand,
		Country:      input.Country,
		SalesChannel: input.TradePolicy,
		Longitude:    input.Longitude,
		Latitude:     input.Latitude,
	})
	if err != nil {
		return RegionalizedIntelligentSearchProductsResult{}, err
	}

	selected, ok := firstRegionWithID(regionResult.Regions)
	if !ok {
		return RegionalizedIntelligentSearchProductsResult{}, capability.StructuredError{Code: ErrorRegionalizedSearchNoRegion, Message: "VTEX region resolution did not return a region ID."}
	}

	segmentUseCase := catalog.NewCreateVTEXSegmentUseCase(u.segmentCreator)
	segmentResult, err := segmentUseCase.Execute(ctx, catalog.CreateVTEXSegmentInput{
		Brand:         input.Brand,
		RegionID:      selected.ID,
		SalesChannel:  input.TradePolicy,
		IncludeCookie: true,
	})
	if err != nil {
		return RegionalizedIntelligentSearchProductsResult{}, err
	}
	if !segmentResult.TokenSet || strings.TrimSpace(segmentResult.Cookie) == "" {
		return RegionalizedIntelligentSearchProductsResult{}, capability.StructuredError{Code: ErrorRegionalizedSearchSegmentUnavailable, Message: "VTEX segment creation did not return a usable segment cookie."}
	}

	searchUseCase := catalog.NewIntelligentSearchProductsUseCase(u.searcher)
	searchResult, err := searchUseCase.Execute(ctx, catalog.IntelligentSearchProductsInput{
		Brand:              input.Brand,
		TradePolicy:        input.TradePolicy,
		Text:               input.Text,
		By:                 input.By,
		Values:             input.Values,
		Query:              input.Query,
		Facets:             input.Facets,
		Page:               input.Page,
		Count:              input.Count,
		Sort:               input.Sort,
		Locale:             input.Locale,
		HideUnavailable:    input.HideUnavailable,
		SimulationBehavior: input.SimulationBehavior,
		Cookies:            []string{segmentResult.Cookie},
	})
	if err != nil {
		return RegionalizedIntelligentSearchProductsResult{}, err
	}

	return RegionalizedIntelligentSearchProductsResult{
		Brand:       input.Brand,
		Country:     input.Country,
		TradePolicy: input.TradePolicy,
		Coordinates: geo.Coordinates{Longitude: input.Longitude, Latitude: input.Latitude},
		HasCoverage: regionResult.HasCoverage,
		Region:      SelectedRegion{ID: selected.ID, Sellers: selected.Sellers},
		Segment: SegmentSummary{
			RegionID:     segmentResult.RegionID,
			SalesChannel: segmentResult.SalesChannel,
			TokenSet:     segmentResult.TokenSet,
			TokenLength:  segmentResult.TokenLength,
		},
		Search:       searchResult,
		RegionResult: regionResult,
	}, nil
}

// NewRegionalizedIntelligentSearchProductsCapability adapts the workflow into an executable Capability.
func NewRegionalizedIntelligentSearchProductsCapability(regionResolver geo.VTEXRegionResolver, segmentCreator catalog.VTEXSegmentCreator, searcher catalog.IntelligentSearchProductsSearcher) capability.Executable {
	useCase := NewRegionalizedIntelligentSearchProductsUseCase(regionResolver, segmentCreator, searcher)
	return capability.Executable{
		Definition: RegionalizedIntelligentSearchProductsDefinition(),
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			input, err := regionalizedInputFromCapability(request.Input)
			if err != nil {
				return capability.ExecutionResult{}, err
			}
			result, err := useCase.Execute(ctx, input)
			if err != nil {
				return capability.ExecutionResult{}, err
			}
			return capability.ExecutionResult{Data: result}, nil
		},
	}
}

// RegionalizedIntelligentSearchProductsDefinition returns the discovery contract.
func RegionalizedIntelligentSearchProductsDefinition() capability.Definition {
	return capability.Definition{
		ID:                   CapabilityRegionalizedIntelligentSearchProductsID,
		Domain:               catalog.DomainName,
		Version:              "1.0.0",
		Title:                "Run regionalized VTEX Intelligent Search products",
		Description:          "Resolves a VTEX region, creates a transient segment, and runs Intelligent Search with that segment.",
		Risk:                 capability.RiskSafeWrite,
		RequiresConfirmation: true,
		Audiences:            []capability.Audience{capability.AudienceAgents},
		Visibility:           []capability.Visibility{capability.VisibilityCLI},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."},
			{Name: "country", Type: capability.InputTypeString, Required: false, Description: "Country code for Checkout Regions. Defaults to COL."},
			{Name: "tradePolicy", Type: capability.InputTypeString, Required: true, Description: "VTEX trade policy/sales channel used for region, segment, and search."},
			{Name: "longitude", Type: capability.InputTypeString, Required: true, Description: "Longitude in VTEX geoCoordinates order."},
			{Name: "latitude", Type: capability.InputTypeString, Required: true, Description: "Latitude in VTEX geoCoordinates order."},
			{Name: "text", Type: capability.InputTypeString, Required: false, Description: "Natural-language search text."},
			{Name: "by", Type: capability.InputTypeString, Required: false, Description: "Typed lookup mode: product-id, sku-id, ean, sku-reference, slug, or id."},
			{Name: "value", Type: capability.InputTypeArray, Required: false, Description: "Lookup values used with by; repeat for same-type multi-ID lookup."},
			{Name: "query", Type: capability.InputTypeString, Required: false, Description: "Raw Intelligent Search query expression for diagnostics."},
			{Name: "facet", Type: capability.InputTypeArray, Required: false, Description: "Additional path facet in key=value form; repeat for multiple facets."},
			{Name: "page", Type: capability.InputTypeNumber, Required: false, Description: "Result page. Defaults to 1."},
			{Name: "count", Type: capability.InputTypeNumber, Required: false, Description: "Products per page. Defaults to 24."},
			{Name: "sort", Type: capability.InputTypeString, Required: false, Description: "Sort value such as price:asc or orders:desc."},
			{Name: "locale", Type: capability.InputTypeString, Required: false, Description: "BCP 47 locale."},
			{Name: "hideUnavailable", Type: capability.InputTypeBoolean, Required: false, Description: "Whether unavailable products should be hidden."},
			{Name: "simulationBehavior", Type: capability.InputTypeString, Required: false, Description: "Simulation behavior: default, skip, or only1P."},
		}},
	}
}

func regionalizedInputFromCapability(input capability.Input) (RegionalizedIntelligentSearchProductsInput, error) {
	out := RegionalizedIntelligentSearchProductsInput{Page: 1, Count: 24}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["country"].(string); ok {
		out.Country = value
	}
	if value, ok := input["tradePolicy"].(string); ok {
		out.TradePolicy = value
	}
	if value, ok := input["longitude"].(string); ok {
		out.Longitude = value
	}
	if value, ok := input["latitude"].(string); ok {
		out.Latitude = value
	}
	if value, ok := input["text"].(string); ok {
		out.Text = value
	}
	if value, ok := input["by"].(string); ok {
		out.By = value
	}
	if value, ok := input["query"].(string); ok {
		out.Query = value
	}
	if value, ok := input["sort"].(string); ok {
		out.Sort = value
	}
	if value, ok := input["locale"].(string); ok {
		out.Locale = value
	}
	if value, ok := input["simulationBehavior"].(string); ok {
		out.SimulationBehavior = value
	}
	if value, ok := input["hideUnavailable"].(bool); ok {
		out.HideUnavailable = &value
	}
	if values, ok := input["value"]; ok {
		parsed, err := stringSlice(values, "value")
		if err != nil {
			return RegionalizedIntelligentSearchProductsInput{}, err
		}
		out.Values = parsed
	}
	if values, ok := input["facet"]; ok {
		parsed, err := stringSlice(values, "facet")
		if err != nil {
			return RegionalizedIntelligentSearchProductsInput{}, err
		}
		out.Facets = parsed
	}
	if value, ok := input["page"]; ok {
		parsed, err := intValue(value, "page")
		if err != nil {
			return RegionalizedIntelligentSearchProductsInput{}, err
		}
		out.Page = parsed
	}
	if value, ok := input["count"]; ok {
		parsed, err := intValue(value, "count")
		if err != nil {
			return RegionalizedIntelligentSearchProductsInput{}, err
		}
		out.Count = parsed
	}
	return out, nil
}

func normalizeRegionalizedInput(input RegionalizedIntelligentSearchProductsInput) RegionalizedIntelligentSearchProductsInput {
	input.Brand = strings.ToLower(strings.TrimSpace(input.Brand))
	if input.Brand == "" {
		input.Brand = "exito"
	}
	input.Country = strings.TrimSpace(input.Country)
	if input.Country == "" {
		input.Country = "COL"
	}
	input.TradePolicy = strings.TrimSpace(input.TradePolicy)
	input.Longitude = strings.TrimSpace(input.Longitude)
	input.Latitude = strings.TrimSpace(input.Latitude)
	input.Text = strings.TrimSpace(input.Text)
	input.By = strings.TrimSpace(input.By)
	input.Query = strings.TrimSpace(input.Query)
	input.Sort = strings.TrimSpace(input.Sort)
	input.Locale = strings.TrimSpace(input.Locale)
	input.SimulationBehavior = strings.TrimSpace(input.SimulationBehavior)
	if input.Page == 0 {
		input.Page = 1
	}
	if input.Count == 0 {
		input.Count = 24
	}
	return input
}

func validateRegionalizedInput(input RegionalizedIntelligentSearchProductsInput) error {
	if input.Brand != "exito" && input.Brand != "carulla" {
		return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: fmt.Sprintf("Unsupported VTEX brand %q.", input.Brand)}
	}
	for name, value := range map[string]string{"tradePolicy": input.TradePolicy, "longitude": input.Longitude, "latitude": input.Latitude} {
		if strings.TrimSpace(value) == "" {
			return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: name + " is required."}
		}
	}
	modes := 0
	for _, present := range []bool{input.Text != "", input.Query != "", input.By != ""} {
		if present {
			modes++
		}
	}
	if modes > 1 {
		return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: "Provide only one query mode: text, query, or by/value."}
	}
	if modes == 0 && len(nonBlank(input.Facets)) == 0 {
		return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: "Provide text, query, by/value, or at least one facet."}
	}
	if input.By != "" && len(nonBlank(input.Values)) == 0 {
		return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: "value is required when by is provided."}
	}
	if input.By == "" && len(nonBlank(input.Values)) > 0 {
		return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: "by is required when value is provided."}
	}
	if input.By != "" && !supportedIntelligentBy(input.By) {
		return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: fmt.Sprintf("Unsupported Intelligent Search lookup mode %q.", input.By)}
	}
	for _, facet := range nonBlank(input.Facets) {
		key, value, ok := strings.Cut(facet, "=")
		if !ok || strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: "facet must use key=value format."}
		}
	}
	if input.Page < 1 || input.Page > 50 || input.Count < 1 || input.Count > 50 {
		return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: "Intelligent Search pagination must satisfy 1 <= page <= 50 and 1 <= count <= 50."}
	}
	if input.SimulationBehavior != "" {
		switch input.SimulationBehavior {
		case "default", "skip", "only1P":
		default:
			return capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: "simulationBehavior must be default, skip, or only1P."}
		}
	}
	return nil
}

func supportedIntelligentBy(by string) bool {
	switch strings.ToLower(strings.TrimSpace(by)) {
	case "product", "product-id", "product.id",
		"sku", "sku-id", "sku.id",
		"sku-ref", "sku-reference", "sku.reference",
		"ean", "slug", "id":
		return true
	default:
		return false
	}
}

func firstRegionWithID(regions []geo.Region) (geo.Region, bool) {
	for _, region := range regions {
		if strings.TrimSpace(region.ID) != "" {
			return region, true
		}
	}
	return geo.Region{}, false
}

func stringSlice(value any, field string) ([]string, error) {
	switch typed := value.(type) {
	case []string:
		return typed, nil
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if !ok {
				return nil, capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: field + " must be a string array."}
			}
			out = append(out, text)
		}
		return out, nil
	default:
		return nil, capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: field + " must be a string array."}
	}
}

func intValue(value any, field string) (int, error) {
	switch typed := value.(type) {
	case int:
		return typed, nil
	case int8:
		return int(typed), nil
	case int16:
		return int(typed), nil
	case int32:
		return int(typed), nil
	case int64:
		return int(typed), nil
	case uint:
		return int(typed), nil
	case uint8:
		return int(typed), nil
	case uint16:
		return int(typed), nil
	case uint32:
		return int(typed), nil
	case uint64:
		if typed > uint64(math.MaxInt) {
			break
		}
		return int(typed), nil
	case float64:
		if typed == math.Trunc(typed) {
			return int(typed), nil
		}
	case float32:
		if typed == float32(math.Trunc(float64(typed))) {
			return int(typed), nil
		}
	}
	return 0, capability.StructuredError{Code: ErrorRegionalizedSearchInvalidInput, Message: field + " must be an integer."}
}

func nonBlank(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
