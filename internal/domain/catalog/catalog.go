package catalog

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
)

const (
	CapabilitySearchProductsID = "catalog.search-products"
	DomainName                 = "catalog"

	ErrorCatalogNotConfigured           = "CATALOG_NOT_CONFIGURED"
	ErrorCatalogProviderUnavailable     = "CATALOG_PROVIDER_UNAVAILABLE"
	ErrorCatalogProviderInvalidResponse = "CATALOG_PROVIDER_INVALID_RESPONSE"
	ErrorCatalogInvalidInput            = "CATALOG_INVALID_INPUT"
)

// SearchProductsInput is the schema-shaped input accepted by catalog.search-products.
type SearchProductsInput struct {
	Brand string
	By    string
	Value string
	FQ    []string
	FT    string
	Order string
	From  int
	To    int
}

// Product is the domain-owned summary for a VTEX catalog product.
type Product struct {
	ProductID        string         `json:"productId"`
	ProductName      string         `json:"productName"`
	Brand            string         `json:"brand,omitempty"`
	BrandID          int            `json:"brandId,omitempty"`
	LinkText         string         `json:"linkText,omitempty"`
	ProductReference string         `json:"productReference,omitempty"`
	CategoryID       string         `json:"categoryId,omitempty"`
	Link             string         `json:"link,omitempty"`
	Items            []SKU          `json:"items,omitempty"`
	Details          map[string]any `json:"details,omitempty"`
}

// SKU is the domain-owned summary for a VTEX catalog SKU.
type SKU struct {
	ItemID       string   `json:"itemId"`
	Name         string   `json:"name"`
	NameComplete string   `json:"nameComplete,omitempty"`
	EAN          string   `json:"ean,omitempty"`
	ReferenceIDs []RefID  `json:"referenceIds,omitempty"`
	Sellers      []Seller `json:"sellers,omitempty"`
}

// RefID describes an SKU reference identifier returned by VTEX.
type RefID struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Seller describes a VTEX seller offer summary.
type Seller struct {
	SellerID      string  `json:"sellerId"`
	SellerName    string  `json:"sellerName,omitempty"`
	SellerDefault bool    `json:"sellerDefault,omitempty"`
	Price         float64 `json:"price,omitempty"`
	ListPrice     float64 `json:"listPrice,omitempty"`
	Available     bool    `json:"available"`
	Quantity      int     `json:"quantity,omitempty"`
}

// SearchProductsResult is the stable use-case result shape for catalog.search-products.
type SearchProductsResult struct {
	Brand      string    `json:"brand"`
	Products   []Product `json:"products"`
	Count      int       `json:"count"`
	RangeStart *int      `json:"rangeStart,omitempty"`
	RangeEnd   *int      `json:"rangeEnd,omitempty"`
	Total      *int      `json:"total,omitempty"`
	Resources  string    `json:"resources,omitempty"`
}

// Searcher retrieves catalog products using domain-owned models.
type Searcher interface {
	SearchProducts(ctx context.Context, input SearchProductsInput) (SearchProductsResult, error)
}

// SearchProductsUseCase runs catalog.search-products without surface dependencies.
type SearchProductsUseCase struct {
	searcher Searcher
}

// NewSearchProductsUseCase creates the catalog.search-products use case.
func NewSearchProductsUseCase(searcher Searcher) SearchProductsUseCase {
	return SearchProductsUseCase{searcher: searcher}
}

// Execute searches VTEX catalog products.
func (u SearchProductsUseCase) Execute(ctx context.Context, input SearchProductsInput) (SearchProductsResult, error) {
	if u.searcher == nil {
		return SearchProductsResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Catalog client is not configured."}
	}
	input.Brand = normalizedBrand(input.Brand)
	input.By = normalizedBy(input.By)
	if input.To == 0 && input.From == 0 {
		input.To = 9
	}
	if err := validateInput(input); err != nil {
		return SearchProductsResult{}, err
	}
	return u.searcher.SearchProducts(ctx, input)
}

// NewSearchProductsCapability adapts the use case into a neutral executable Capability.
func NewSearchProductsCapability(searcher Searcher) capability.Executable {
	useCase := NewSearchProductsUseCase(searcher)
	return capability.Executable{
		Definition: Definition(),
		Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
			input, err := inputFromCapability(request.Input)
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

// Definition returns the neutral catalog.search-products discovery contract.
func Definition() capability.Definition {
	return capability.Definition{
		ID:          CapabilitySearchProductsID,
		Domain:      DomainName,
		Version:     "1.0.0",
		Title:       "Search VTEX catalog products",
		Description: "Searches VTEX public catalog products by friendly lookup mode or raw search filters.",
		Risk:        capability.RiskReadOnly,
		Audiences:   []capability.Audience{capability.AudienceAgents, capability.AudiencePeople},
		Visibility:  []capability.Visibility{capability.VisibilityCLI, capability.VisibilityTUI, capability.VisibilityCommandPalette},
		InputSchema: &capability.InputSchema{Fields: []capability.InputField{
			{Name: "brand", Type: capability.InputTypeString, Required: false, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."},
			{Name: "by", Type: capability.InputTypeString, Required: false, Description: "Friendly lookup mode: sku-id, product-id, ref-id, ean, seller-id, category, brand-id, collection-id, text, or slug."},
			{Name: "value", Type: capability.InputTypeString, Required: false, Description: "Lookup value used with by."},
			{Name: "fq", Type: capability.InputTypeArray, Required: false, Description: "Raw VTEX fq filters. May include multiple values."},
			{Name: "ft", Type: capability.InputTypeString, Required: false, Description: "VTEX full-text search term."},
			{Name: "order", Type: capability.InputTypeString, Required: false, Description: "VTEX O sorting value, such as OrderByPriceASC."},
			{Name: "from", Type: capability.InputTypeNumber, Required: false, Description: "Initial result index. Defaults to 0."},
			{Name: "to", Type: capability.InputTypeNumber, Required: false, Description: "Final result index. Defaults to 9; VTEX allows at most 50 results per request."},
		}},
	}
}

func inputFromCapability(input capability.Input) (SearchProductsInput, error) {
	out := SearchProductsInput{From: 0, To: 9}
	if value, ok := input["brand"].(string); ok {
		out.Brand = value
	}
	if value, ok := input["by"].(string); ok {
		out.By = value
	}
	if value, ok := input["value"].(string); ok {
		out.Value = value
	}
	if value, ok := input["ft"].(string); ok {
		out.FT = value
	}
	if value, ok := input["order"].(string); ok {
		out.Order = value
	}
	if values, ok := input["fq"]; ok {
		fq, err := stringSlice(values)
		if err != nil {
			return SearchProductsInput{}, err
		}
		out.FQ = fq
	}
	if value, ok := input["from"]; ok {
		parsed, err := intValue(value, "from")
		if err != nil {
			return SearchProductsInput{}, err
		}
		out.From = parsed
	}
	if value, ok := input["to"]; ok {
		parsed, err := intValue(value, "to")
		if err != nil {
			return SearchProductsInput{}, err
		}
		out.To = parsed
	}
	return out, nil
}

func stringSlice(value any) ([]string, error) {
	switch typed := value.(type) {
	case []string:
		return typed, nil
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			text, ok := item.(string)
			if !ok {
				return nil, capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "fq must contain only strings."}
			}
			out = append(out, text)
		}
		return out, nil
	case string:
		if strings.TrimSpace(typed) == "" {
			return nil, nil
		}
		return []string{typed}, nil
	default:
		return nil, capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "fq must be a string array."}
	}
}

func intValue(value any, field string) (int, error) {
	switch typed := value.(type) {
	case int:
		return typed, nil
	case float64:
		return int(typed), nil
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err != nil {
			return 0, capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: fmt.Sprintf("%s must be a number.", field)}
		}
		return parsed, nil
	default:
		return 0, capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: fmt.Sprintf("%s must be a number.", field)}
	}
}

func validateInput(input SearchProductsInput) error {
	if input.From < 0 || input.To < 0 || input.To < input.From || input.To-input.From >= 50 || input.From > 2500 {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "Catalog pagination must satisfy 0 <= from <= to, to-from < 50, and from <= 2500."}
	}
	if input.By != "" && strings.TrimSpace(input.Value) == "" {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "value is required when by is provided."}
	}
	if input.By == "" && strings.TrimSpace(input.Value) != "" {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "by is required when value is provided."}
	}
	if input.By == "" && strings.TrimSpace(input.FT) == "" && len(nonBlank(input.FQ)) == 0 {
		return capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "Provide by/value, ft, or at least one fq filter."}
	}
	return nil
}

func normalizedBrand(brand string) string {
	switch strings.ToLower(strings.TrimSpace(brand)) {
	case "carulla":
		return "carulla"
	default:
		return "exito"
	}
}

func normalizedBy(by string) string {
	return strings.ToLower(strings.TrimSpace(by))
}

func nonBlank(values []string) []string {
	out := values[:0]
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, strings.TrimSpace(value))
		}
	}
	return out
}
