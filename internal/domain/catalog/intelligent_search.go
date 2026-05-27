package catalog

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
)

const CapabilityIntelligentSearchProductsID = "catalog.intelligent-search-products"

// IntelligentSearchProductsInput is the schema-shaped input accepted by catalog.intelligent-search-products.
type IntelligentSearchProductsInput struct {
	Brand              string
	TradePolicy        string
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
	Cookies            []string
}

// Facet describes an Intelligent Search path facet.
type Facet struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RequestDiagnostics captures non-secret request details useful for troubleshooting.
type RequestDiagnostics struct {
	RequestPath      string            `json:"requestPath,omitempty"`
	RequestQuery     map[string]string `json:"requestQuery,omitempty"`
	CookieNames      []string          `json:"cookieNames,omitempty"`
	ProviderPayload  map[string]any    `json:"providerPayload,omitempty"`
	ProviderProducts []map[string]any  `json:"providerProducts,omitempty"`
}

// IntelligentSearchProductsResult is the stable use-case result shape for VTEX Intelligent Search.
type IntelligentSearchProductsResult struct {
	Brand       string             `json:"brand"`
	Query       string             `json:"query,omitempty"`
	Facets      []Facet            `json:"facets"`
	Page        int                `json:"page"`
	Count       int                `json:"count"`
	Sort        string             `json:"sort,omitempty"`
	Products    []Product          `json:"products"`
	Total       *int               `json:"total,omitempty"`
	Diagnostics RequestDiagnostics `json:"diagnostics,omitempty"`
}

// IntelligentSearchProductsSearcher retrieves VTEX Intelligent Search products using domain-owned models.
type IntelligentSearchProductsSearcher interface {
	IntelligentSearchProducts(ctx context.Context, input IntelligentSearchProductsInput) (IntelligentSearchProductsResult, error)
}

// IntelligentSearchProductsUseCase runs catalog.intelligent-search-products without surface dependencies.
type IntelligentSearchProductsUseCase struct {
	searcher IntelligentSearchProductsSearcher
}

// NewIntelligentSearchProductsUseCase creates the Intelligent Search use case.
func NewIntelligentSearchProductsUseCase(searcher IntelligentSearchProductsSearcher) IntelligentSearchProductsUseCase {
	return IntelligentSearchProductsUseCase{searcher: searcher}
}

// Execute searches VTEX Intelligent Search products.
func (u IntelligentSearchProductsUseCase) Execute(ctx context.Context, input IntelligentSearchProductsInput) (IntelligentSearchProductsResult, error) {
	if u.searcher == nil {
		return IntelligentSearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Intelligent Search client is not configured."}
	}
	input.Brand = normalizedBrand(input.Brand)
	input.By = normalizedIntelligentBy(input.By)
	if input.Page == 0 {
		input.Page = 1
	}
	if input.Count == 0 {
		input.Count = 24
	}
	if err := validateIntelligentSearchInput(input); err != nil {
		return IntelligentSearchProductsResult{}, err
	}
	return u.searcher.IntelligentSearchProducts(ctx, input)
}

// NewIntelligentSearchProductsCapability adapts the use case into a neutral executable Capability.
func NewIntelligentSearchProductsCapability(searcher IntelligentSearchProductsSearcher) capability.Executable {
	useCase := NewIntelligentSearchProductsUseCase(searcher)
	return capability.Executable{
		Definition: IntelligentSearchProductsDefinition(),
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			input, err := intelligentInputFromCapability(request.Input)
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

// IntelligentSearchProductsDefinition returns the neutral discovery contract.
func IntelligentSearchProductsDefinition() capability.Definition {
	return capability.Definition{
		ID:          CapabilityIntelligentSearchProductsID,
		Domain:      DomainName,
		Version:     "1.0.0",
		Title:       "Search VTEX Intelligent Search products",
		Description: "Searches VTEX Intelligent Search products by text, typed lookup, raw query, or path facets.",
		Risk:        capability.RiskReadOnly,
		Audiences:   []capability.Audience{capability.AudienceAgents, capability.AudiencePeople},
		Visibility:  []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."},
			{Name: "tradePolicy", Type: capability.InputTypeString, Required: true, Description: "VTEX trade policy/sales channel encoded as a required path facet."},
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
			{Name: "cookie", Type: capability.InputTypeArray, Required: false, Description: "Optional VTEX cookie strings; values are redacted from diagnostics."},
		}},
	}
}

func intelligentInputFromCapability(input capability.Input) (IntelligentSearchProductsInput, error) {
	out := IntelligentSearchProductsInput{Page: 1, Count: 24}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["tradePolicy"].(string); ok {
		out.TradePolicy = value
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
		parsed, err := stringSlice(values)
		if err != nil {
			return IntelligentSearchProductsInput{}, capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "value must be a string array."}
		}
		out.Values = parsed
	}
	if values, ok := input["facet"]; ok {
		parsed, err := stringSlice(values)
		if err != nil {
			return IntelligentSearchProductsInput{}, capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "facet must be a string array."}
		}
		out.Facets = parsed
	}
	if values, ok := input["cookie"]; ok {
		parsed, err := stringSlice(values)
		if err != nil {
			return IntelligentSearchProductsInput{}, capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "cookie must be a string array."}
		}
		out.Cookies = parsed
	}
	if value, ok := input["page"]; ok {
		parsed, err := intValue(value, "page")
		if err != nil {
			return IntelligentSearchProductsInput{}, err
		}
		out.Page = parsed
	}
	if value, ok := input["count"]; ok {
		parsed, err := intValue(value, "count")
		if err != nil {
			return IntelligentSearchProductsInput{}, err
		}
		out.Count = parsed
	}
	return out, nil
}

func validateIntelligentSearchInput(input IntelligentSearchProductsInput) error {
	if strings.TrimSpace(input.TradePolicy) == "" {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "tradePolicy is required."}
	}
	if input.Page < 1 || input.Page > 50 || input.Count < 1 || input.Count > 50 {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "Intelligent Search pagination must satisfy 1 <= page <= 50 and 1 <= count <= 50."}
	}
	if input.By != "" && len(nonBlank(input.Values)) == 0 {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "value is required when by is provided."}
	}
	if input.By == "" && len(nonBlank(input.Values)) > 0 {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "by is required when value is provided."}
	}
	modes := 0
	for _, present := range []bool{strings.TrimSpace(input.Text) != "", strings.TrimSpace(input.Query) != "", input.By != ""} {
		if present {
			modes++
		}
	}
	if modes > 1 {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "Provide only one query mode: text, query, or by/value."}
	}
	if modes == 0 && len(nonBlank(input.Facets)) == 0 {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "Provide text, query, by/value, or at least one facet."}
	}
	for _, facet := range nonBlank(input.Facets) {
		key, value, ok := strings.Cut(facet, "=")
		if !ok || strings.TrimSpace(key) == "" || strings.TrimSpace(value) == "" {
			return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "facet must use key=value format."}
		}
	}
	if input.SimulationBehavior != "" {
		switch input.SimulationBehavior {
		case "default", "skip", "only1P":
		default:
			return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "simulationBehavior must be default, skip, or only1P."}
		}
	}
	if input.By != "" && intelligentQueryPrefix(input.By) == "" {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: fmt.Sprintf("Unsupported Intelligent Search lookup mode %q.", input.By)}
	}
	return nil
}

func normalizedIntelligentBy(by string) string {
	switch strings.ToLower(strings.TrimSpace(by)) {
	case "product", "product-id", "product.id":
		return "product-id"
	case "sku", "sku-id", "sku.id":
		return "sku-id"
	case "sku-ref", "sku-reference", "sku.reference":
		return "sku-reference"
	case "ean", "sku.ean":
		return "ean"
	case "slug", "product-link", "product.link":
		return "slug"
	case "id":
		return "id"
	default:
		return strings.TrimSpace(by)
	}
}

func intelligentQueryPrefix(by string) string {
	switch by {
	case "product-id":
		return "product.id:"
	case "sku-id":
		return "sku.id:"
	case "ean":
		return "sku.ean:"
	case "sku-reference":
		return "sku.reference:"
	case "slug":
		return "product.link:"
	case "id":
		return "id:"
	default:
		return ""
	}
}

func intelligentQuery(input IntelligentSearchProductsInput) string {
	if strings.TrimSpace(input.Query) != "" {
		return strings.TrimSpace(input.Query)
	}
	if strings.TrimSpace(input.Text) != "" {
		return strings.TrimSpace(input.Text)
	}
	if input.By != "" {
		return intelligentQueryPrefix(input.By) + strings.Join(nonBlank(input.Values), ";")
	}
	return ""
}

func boolString(value bool) string {
	return strconv.FormatBool(value)
}
