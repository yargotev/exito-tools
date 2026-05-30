package masterdata

import (
	"context"
	"fmt"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
)

const (
	CapabilityGetDocumentID     = "masterdata.get-document"
	CapabilitySearchDocumentsID = "masterdata.search-documents"
	CapabilityScrollDocumentsID = "masterdata.scroll-documents"
	CapabilityListSchemasID     = "masterdata.list-schemas"
	CapabilityGetSchemaID       = "masterdata.get-schema"
	CapabilityListIndicesID     = "masterdata.list-indices"
	DomainName                  = "masterdata"

	ErrorMasterDataNotConfigured           = "MASTERDATA_NOT_CONFIGURED"
	ErrorMasterDataProviderUnavailable     = "MASTERDATA_PROVIDER_UNAVAILABLE"
	ErrorMasterDataProviderInvalidResponse = "MASTERDATA_PROVIDER_INVALID_RESPONSE"
	ErrorMasterDataInvalidInput            = "MASTERDATA_INVALID_INPUT"

	WarningMasterDataSearchWithoutSort = "MASTERDATA_SEARCH_WITHOUT_SORT"
)

type GetDocumentInput struct {
	Brand      string
	Entity     string
	DocumentID string
	Schema     string
	Fields     []string
}

type SearchDocumentsInput struct {
	Brand     string
	Entity    string
	Fields    []string
	Where     string
	Schema    string
	Sort      string
	RangeFrom int
	RangeTo   int
}

type ScrollDocumentsInput struct {
	Brand  string
	Entity string
	Fields []string
	Where  string
	Schema string
	Size   int
	Token  string
}

type EntityInput struct {
	Brand  string
	Entity string
}

type GetSchemaInput struct {
	Brand  string
	Entity string
	Schema string
}

type Diagnostics struct {
	RequestPath    string `json:"requestPath,omitempty"`
	ProviderStatus int    `json:"providerStatus,omitempty"`
}

type DocumentResult struct {
	Brand       string         `json:"brand"`
	Entity      string         `json:"entity"`
	DocumentID  string         `json:"documentId,omitempty"`
	Fields      []string       `json:"fields,omitempty"`
	Data        map[string]any `json:"data"`
	Diagnostics Diagnostics    `json:"diagnostics,omitempty"`
}

type DocumentsPage struct {
	Brand       string           `json:"brand"`
	Entity      string           `json:"entity"`
	Fields      []string         `json:"fields,omitempty"`
	Documents   []map[string]any `json:"documents"`
	RangeFrom   int              `json:"rangeFrom,omitempty"`
	RangeTo     int              `json:"rangeTo,omitempty"`
	Total       *int             `json:"total,omitempty"`
	NextToken   string           `json:"nextToken,omitempty"`
	Diagnostics Diagnostics      `json:"diagnostics,omitempty"`
}

type SchemasResult struct {
	Brand       string      `json:"brand,omitempty"`
	Entity      string      `json:"entity,omitempty"`
	Schemas     []string    `json:"schemas"`
	Diagnostics Diagnostics `json:"diagnostics,omitempty"`
}

type SchemaResult struct {
	Brand       string         `json:"brand,omitempty"`
	Entity      string         `json:"entity,omitempty"`
	Schema      string         `json:"schema"`
	Definition  map[string]any `json:"definition"`
	Diagnostics Diagnostics    `json:"diagnostics,omitempty"`
}

type IndicesResult struct {
	Brand       string      `json:"brand,omitempty"`
	Entity      string      `json:"entity,omitempty"`
	Indices     []string    `json:"indices"`
	Diagnostics Diagnostics `json:"diagnostics,omitempty"`
}

type GetDocumentResult struct {
	Document DocumentResult `json:"document"`
}

type SearchDocumentsResult struct {
	Page DocumentsPage `json:"page"`
}

type ScrollDocumentsResult struct {
	Page DocumentsPage `json:"page"`
}

type Client interface {
	GetDocument(context.Context, GetDocumentInput) (DocumentResult, error)
	SearchDocuments(context.Context, SearchDocumentsInput) (DocumentsPage, error)
	ScrollDocuments(context.Context, ScrollDocumentsInput) (DocumentsPage, error)
	ListSchemas(context.Context, EntityInput) (SchemasResult, error)
	GetSchema(context.Context, GetSchemaInput) (SchemaResult, error)
	ListIndices(context.Context, EntityInput) (IndicesResult, error)
}

type (
	GetDocumentUseCase     struct{ client Client }
	SearchDocumentsUseCase struct{ client Client }
	ScrollDocumentsUseCase struct{ client Client }
	ListSchemasUseCase     struct{ client Client }
	GetSchemaUseCase       struct{ client Client }
	ListIndicesUseCase     struct{ client Client }
)

func NewGetDocumentUseCase(client Client) GetDocumentUseCase {
	return GetDocumentUseCase{client: client}
}

func NewSearchDocumentsUseCase(client Client) SearchDocumentsUseCase {
	return SearchDocumentsUseCase{client: client}
}

func NewScrollDocumentsUseCase(client Client) ScrollDocumentsUseCase {
	return ScrollDocumentsUseCase{client: client}
}

func NewListSchemasUseCase(client Client) ListSchemasUseCase {
	return ListSchemasUseCase{client: client}
}
func NewGetSchemaUseCase(client Client) GetSchemaUseCase { return GetSchemaUseCase{client: client} }
func NewListIndicesUseCase(client Client) ListIndicesUseCase {
	return ListIndicesUseCase{client: client}
}

func (u GetDocumentUseCase) Execute(ctx context.Context, input GetDocumentInput) (GetDocumentResult, error) {
	if u.client == nil {
		return GetDocumentResult{}, notConfiguredError()
	}
	input = normalizeGetDocumentInput(input)
	if err := validateGetDocumentInput(input); err != nil {
		return GetDocumentResult{}, err
	}
	document, err := u.client.GetDocument(ctx, input)
	if err != nil {
		return GetDocumentResult{}, err
	}
	return GetDocumentResult{Document: document}, nil
}

func (u SearchDocumentsUseCase) Execute(ctx context.Context, input SearchDocumentsInput) (SearchDocumentsResult, error) {
	if u.client == nil {
		return SearchDocumentsResult{}, notConfiguredError()
	}
	input = normalizeSearchDocumentsInput(input)
	if err := validateSearchDocumentsInput(input); err != nil {
		return SearchDocumentsResult{}, err
	}
	page, err := u.client.SearchDocuments(ctx, input)
	if err != nil {
		return SearchDocumentsResult{}, err
	}
	return SearchDocumentsResult{Page: page}, nil
}

func (u ScrollDocumentsUseCase) Execute(ctx context.Context, input ScrollDocumentsInput) (ScrollDocumentsResult, error) {
	if u.client == nil {
		return ScrollDocumentsResult{}, notConfiguredError()
	}
	input = normalizeScrollDocumentsInput(input)
	if err := validateScrollDocumentsInput(input); err != nil {
		return ScrollDocumentsResult{}, err
	}
	page, err := u.client.ScrollDocuments(ctx, input)
	if err != nil {
		return ScrollDocumentsResult{}, err
	}
	return ScrollDocumentsResult{Page: page}, nil
}

func (u ListSchemasUseCase) Execute(ctx context.Context, input EntityInput) (SchemasResult, error) {
	if u.client == nil {
		return SchemasResult{}, notConfiguredError()
	}
	input = normalizeEntityInput(input)
	if err := validateEntityInput(input); err != nil {
		return SchemasResult{}, err
	}
	return u.client.ListSchemas(ctx, input)
}

func (u GetSchemaUseCase) Execute(ctx context.Context, input GetSchemaInput) (SchemaResult, error) {
	if u.client == nil {
		return SchemaResult{}, notConfiguredError()
	}
	input = normalizeGetSchemaInput(input)
	if err := validateGetSchemaInput(input); err != nil {
		return SchemaResult{}, err
	}
	return u.client.GetSchema(ctx, input)
}

func (u ListIndicesUseCase) Execute(ctx context.Context, input EntityInput) (IndicesResult, error) {
	if u.client == nil {
		return IndicesResult{}, notConfiguredError()
	}
	input = normalizeEntityInput(input)
	if err := validateEntityInput(input); err != nil {
		return IndicesResult{}, err
	}
	return u.client.ListIndices(ctx, input)
}

func NewGetDocumentCapability(client Client) capability.Executable {
	useCase := NewGetDocumentUseCase(client)
	return capability.Executable{Definition: GetDocumentDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, getDocumentInputFromCapability(request.Input))
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		return capability.ExecutionResult{Data: result}, nil
	}}
}

func NewSearchDocumentsCapability(client Client) capability.Executable {
	useCase := NewSearchDocumentsUseCase(client)
	return capability.Executable{Definition: SearchDocumentsDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		input := searchDocumentsInputFromCapability(request.Input)
		result, err := useCase.Execute(ctx, input)
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		warnings := []capability.StructuredWarning(nil)
		if strings.TrimSpace(input.Sort) == "" {
			warnings = append(warnings, capability.StructuredWarning{Code: WarningMasterDataSearchWithoutSort, Message: "Master Data search pagination can be unstable without sort."})
		}
		return capability.ExecutionResult{Data: result, Warnings: warnings}, nil
	}}
}

func NewScrollDocumentsCapability(client Client) capability.Executable {
	useCase := NewScrollDocumentsUseCase(client)
	return capability.Executable{Definition: ScrollDocumentsDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, scrollDocumentsInputFromCapability(request.Input))
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		pagination := &capability.PaginationMeta{NextCursor: result.Page.NextToken, HasMore: result.Page.NextToken != ""}
		return capability.ExecutionResult{Data: result, Pagination: pagination}, nil
	}}
}

func NewListSchemasCapability(client Client) capability.Executable {
	useCase := NewListSchemasUseCase(client)
	return capability.Executable{Definition: ListSchemasDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, entityInputFromCapability(request.Input))
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		return capability.ExecutionResult{Data: result}, nil
	}}
}

func NewGetSchemaCapability(client Client) capability.Executable {
	useCase := NewGetSchemaUseCase(client)
	return capability.Executable{Definition: GetSchemaDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, getSchemaInputFromCapability(request.Input))
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		return capability.ExecutionResult{Data: result}, nil
	}}
}

func NewListIndicesCapability(client Client) capability.Executable {
	useCase := NewListIndicesUseCase(client)
	return capability.Executable{Definition: ListIndicesDefinition(), Handler: func(ctx context.Context, request capability.ExecutionRequest) (capability.ExecutionResult, error) {
		result, err := useCase.Execute(ctx, entityInputFromCapability(request.Input))
		if err != nil {
			return capability.ExecutionResult{}, err
		}
		return capability.ExecutionResult{Data: result}, nil
	}}
}

func GetDocumentDefinition() capability.Definition {
	return baseDefinition(CapabilityGetDocumentID, "Get Master Data document", "Gets one VTEX Master Data document by ID.", []capability.InputField{brandField(), entityField(), {Name: "documentId", Type: capability.InputTypeString, Required: true, Description: "Master Data document identifier."}, schemaField(false), fieldsField(false)})
}

func SearchDocumentsDefinition() capability.Definition {
	return baseDefinition(CapabilitySearchDocumentsID, "Search Master Data documents", "Searches VTEX Master Data documents with bounded pagination.", []capability.InputField{brandField(), entityField(), fieldsField(false), {Name: "where", Type: capability.InputTypeString, Description: "Optional VTEX _where predicate."}, schemaField(false), {Name: "sort", Type: capability.InputTypeString, Description: "Optional VTEX _sort value. Recommended for pagination stability."}, {Name: "rangeFrom", Type: capability.InputTypeNumber, Description: "First result index. Defaults to 0."}, {Name: "rangeTo", Type: capability.InputTypeNumber, Description: "Last result index. Defaults to 99 and cannot exceed 100 documents per request."}})
}

func ScrollDocumentsDefinition() capability.Definition {
	return baseDefinition(CapabilityScrollDocumentsID, "Scroll Master Data documents", "Starts or continues a VTEX Master Data scroll read.", []capability.InputField{brandField(), entityField(), fieldsField(false), {Name: "where", Type: capability.InputTypeString, Description: "Optional VTEX _where predicate."}, schemaField(false), {Name: "size", Type: capability.InputTypeNumber, Description: "Scroll page size. Defaults to 100 and cannot exceed 1000."}, {Name: "token", Type: capability.InputTypeString, Description: "Previous X-VTEX-MD-TOKEN value."}})
}

func ListSchemasDefinition() capability.Definition {
	return baseDefinition(CapabilityListSchemasID, "List Master Data schemas", "Lists VTEX Master Data v2 schemas for an entity.", []capability.InputField{brandField(), entityField()})
}

func GetSchemaDefinition() capability.Definition {
	return baseDefinition(CapabilityGetSchemaID, "Get Master Data schema", "Gets a VTEX Master Data v2 schema definition.", []capability.InputField{brandField(), entityField(), schemaField(true)})
}

func ListIndicesDefinition() capability.Definition {
	return baseDefinition(CapabilityListIndicesID, "List Master Data indices", "Lists VTEX Master Data v2 indices for an entity.", []capability.InputField{brandField(), entityField()})
}

func baseDefinition(id, title, description string, fields []capability.InputField) capability.Definition {
	return capability.Definition{ID: id, Domain: DomainName, Version: "1.0.0", Title: title, Description: description, Risk: capability.RiskReadOnly, Audiences: []capability.Audience{capability.AudienceAgents}, Visibility: []capability.Visibility{capability.VisibilityCLI, capability.VisibilityCommandPalette}, InputSchema: &capability.InputSchema{Fields: fields}}
}

func brandField() capability.InputField {
	return capability.InputField{Name: "brand", Type: capability.InputTypeString, Description: "VTEX brand account to query: exito or carulla. Defaults to exito."}
}

func entityField() capability.InputField {
	return capability.InputField{Name: "entity", Type: capability.InputTypeString, Required: true, Description: "Master Data entity acronym/name."}
}

func schemaField(required bool) capability.InputField {
	return capability.InputField{Name: "schema", Type: capability.InputTypeString, Required: required, Description: "Optional v2 schema selector."}
}

func fieldsField(required bool) capability.InputField {
	return capability.InputField{Name: "fields", Type: capability.InputTypeArray, Required: required, Description: "Field names to request from VTEX Master Data."}
}

func getDocumentInputFromCapability(input capability.Input) GetDocumentInput {
	return GetDocumentInput{Brand: stringInput(input, "brand"), Entity: stringInput(input, "entity"), DocumentID: stringInput(input, "documentId"), Schema: stringInput(input, "schema"), Fields: stringSliceInput(input["fields"])}
}

func searchDocumentsInputFromCapability(input capability.Input) SearchDocumentsInput {
	return SearchDocumentsInput{Brand: stringInput(input, "brand"), Entity: stringInput(input, "entity"), Fields: stringSliceInput(input["fields"]), Where: stringInput(input, "where"), Schema: stringInput(input, "schema"), Sort: stringInput(input, "sort"), RangeFrom: intInput(input["rangeFrom"]), RangeTo: intInput(input["rangeTo"])}
}

func scrollDocumentsInputFromCapability(input capability.Input) ScrollDocumentsInput {
	return ScrollDocumentsInput{Brand: stringInput(input, "brand"), Entity: stringInput(input, "entity"), Fields: stringSliceInput(input["fields"]), Where: stringInput(input, "where"), Schema: stringInput(input, "schema"), Size: intInput(input["size"]), Token: stringInput(input, "token")}
}

func entityInputFromCapability(input capability.Input) EntityInput {
	return EntityInput{Brand: stringInput(input, "brand"), Entity: stringInput(input, "entity")}
}

func getSchemaInputFromCapability(input capability.Input) GetSchemaInput {
	return GetSchemaInput{Brand: stringInput(input, "brand"), Entity: stringInput(input, "entity"), Schema: stringInput(input, "schema")}
}

func normalizeGetDocumentInput(input GetDocumentInput) GetDocumentInput {
	input.Brand = normalizedBrand(input.Brand)
	input.Entity = strings.TrimSpace(input.Entity)
	input.DocumentID = strings.TrimSpace(input.DocumentID)
	input.Schema = strings.TrimSpace(input.Schema)
	input.Fields = nonBlank(input.Fields)
	return input
}

func normalizeSearchDocumentsInput(input SearchDocumentsInput) SearchDocumentsInput {
	input.Brand = normalizedBrand(input.Brand)
	input.Entity = strings.TrimSpace(input.Entity)
	input.Schema = strings.TrimSpace(input.Schema)
	input.Where = strings.TrimSpace(input.Where)
	input.Sort = strings.TrimSpace(input.Sort)
	input.Fields = nonBlank(input.Fields)
	if input.RangeTo == 0 {
		input.RangeTo = 99
	}
	return input
}

func normalizeScrollDocumentsInput(input ScrollDocumentsInput) ScrollDocumentsInput {
	input.Brand = normalizedBrand(input.Brand)
	input.Entity = strings.TrimSpace(input.Entity)
	input.Schema = strings.TrimSpace(input.Schema)
	input.Where = strings.TrimSpace(input.Where)
	input.Token = strings.TrimSpace(input.Token)
	input.Fields = nonBlank(input.Fields)
	if input.Size == 0 {
		input.Size = 100
	}
	return input
}

func normalizeEntityInput(input EntityInput) EntityInput {
	input.Brand = normalizedBrand(input.Brand)
	input.Entity = strings.TrimSpace(input.Entity)
	return input
}

func normalizeGetSchemaInput(input GetSchemaInput) GetSchemaInput {
	input.Brand = normalizedBrand(input.Brand)
	input.Entity = strings.TrimSpace(input.Entity)
	input.Schema = strings.TrimSpace(input.Schema)
	return input
}

func validateGetDocumentInput(input GetDocumentInput) error {
	if err := validateEntityInput(EntityInput{Brand: input.Brand, Entity: input.Entity}); err != nil {
		return err
	}
	if input.DocumentID == "" {
		return capability.StructuredError{Code: ErrorMasterDataInvalidInput, Message: "documentId is required."}
	}
	return nil
}

func validateSearchDocumentsInput(input SearchDocumentsInput) error {
	if err := validateEntityInput(EntityInput{Brand: input.Brand, Entity: input.Entity}); err != nil {
		return err
	}
	if input.RangeFrom < 0 || input.RangeTo < input.RangeFrom || input.RangeTo-input.RangeFrom+1 > 100 {
		return capability.StructuredError{Code: ErrorMasterDataInvalidInput, Message: "Master Data search range must request between 1 and 100 documents."}
	}
	return nil
}

func validateScrollDocumentsInput(input ScrollDocumentsInput) error {
	if err := validateEntityInput(EntityInput{Brand: input.Brand, Entity: input.Entity}); err != nil {
		return err
	}
	if input.Size < 1 || input.Size > 1000 {
		return capability.StructuredError{Code: ErrorMasterDataInvalidInput, Message: "Master Data scroll size must be between 1 and 1000."}
	}
	return nil
}

func validateEntityInput(input EntityInput) error {
	if input.Entity == "" {
		return capability.StructuredError{Code: ErrorMasterDataInvalidInput, Message: "entity is required."}
	}
	return validateBrand(input.Brand)
}

func validateGetSchemaInput(input GetSchemaInput) error {
	if err := validateEntityInput(EntityInput{Brand: input.Brand, Entity: input.Entity}); err != nil {
		return err
	}
	if input.Schema == "" {
		return capability.StructuredError{Code: ErrorMasterDataInvalidInput, Message: "schema is required."}
	}
	return nil
}

func normalizedBrand(brand string) string {
	trimmed := strings.ToLower(strings.TrimSpace(brand))
	if trimmed == "" {
		return "exito"
	}
	return trimmed
}

func validateBrand(brand string) error {
	if brand != "exito" && brand != "carulla" {
		return capability.StructuredError{Code: ErrorMasterDataInvalidInput, Message: fmt.Sprintf("Unsupported VTEX brand %q.", brand)}
	}
	return nil
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

func stringInput(input capability.Input, key string) string {
	if value, ok := input[key].(string); ok {
		return value
	}
	return ""
}

func stringSliceInput(value any) []string {
	switch values := value.(type) {
	case []string:
		return append([]string(nil), values...)
	case []any:
		out := make([]string, 0, len(values))
		for _, value := range values {
			if text, ok := value.(string); ok {
				out = append(out, text)
			}
		}
		return out
	default:
		return nil
	}
}

func intInput(value any) int {
	switch number := value.(type) {
	case int:
		return number
	case int64:
		return int(number)
	case float64:
		if number == float64(int(number)) {
			return int(number)
		}
	}
	return 0
}

func notConfiguredError() error {
	return capability.StructuredError{Code: ErrorMasterDataNotConfigured, Message: "VTEX Master Data client is not configured."}
}
