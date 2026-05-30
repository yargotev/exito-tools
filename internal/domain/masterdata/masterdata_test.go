package masterdata_test

import (
	"context"
	"errors"
	"testing"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/domain/masterdata"
)

type recordingClient struct {
	getInput    masterdata.GetDocumentInput
	searchInput masterdata.SearchDocumentsInput
	scrollInput masterdata.ScrollDocumentsInput
}

func (c *recordingClient) GetDocument(_ context.Context, input masterdata.GetDocumentInput) (masterdata.DocumentResult, error) {
	c.getInput = input
	return masterdata.DocumentResult{Brand: input.Brand, Entity: input.Entity, DocumentID: input.DocumentID, Data: map[string]any{"email": "customer@example.test"}}, nil
}

func (c *recordingClient) SearchDocuments(_ context.Context, input masterdata.SearchDocumentsInput) (masterdata.DocumentsPage, error) {
	c.searchInput = input
	return masterdata.DocumentsPage{Brand: input.Brand, Entity: input.Entity, Documents: []map[string]any{{"id": "doc-1"}}}, nil
}

func (c *recordingClient) ScrollDocuments(_ context.Context, input masterdata.ScrollDocumentsInput) (masterdata.DocumentsPage, error) {
	c.scrollInput = input
	return masterdata.DocumentsPage{Brand: input.Brand, Entity: input.Entity, Documents: []map[string]any{{"id": "doc-1"}}, NextToken: "scroll-token"}, nil
}

func (c *recordingClient) ListSchemas(context.Context, masterdata.EntityInput) (masterdata.SchemasResult, error) {
	return masterdata.SchemasResult{Schemas: []string{"client-v1"}}, nil
}

func (c *recordingClient) GetSchema(context.Context, masterdata.GetSchemaInput) (masterdata.SchemaResult, error) {
	return masterdata.SchemaResult{Schema: "client-v1", Definition: map[string]any{"type": "object"}}, nil
}

func (c *recordingClient) ListIndices(context.Context, masterdata.EntityInput) (masterdata.IndicesResult, error) {
	return masterdata.IndicesResult{Indices: []string{"email-index"}}, nil
}

func TestGetDocumentUseCaseNormalizesInput(t *testing.T) {
	t.Parallel()

	client := &recordingClient{}
	got, err := masterdata.NewGetDocumentUseCase(client).Execute(context.Background(), masterdata.GetDocumentInput{Entity: " CL ", DocumentID: " doc-1 ", Fields: []string{" email ", " firstName "}})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if client.getInput.Brand != "exito" || client.getInput.Entity != "CL" || client.getInput.DocumentID != "doc-1" {
		t.Fatalf("input = %#v, want normalized brand/entity/document", client.getInput)
	}
	if len(client.getInput.Fields) != 2 || client.getInput.Fields[0] != "email" || client.getInput.Fields[1] != "firstName" {
		t.Fatalf("fields = %#v, want trimmed fields", client.getInput.Fields)
	}
	if got.Document.Brand != "exito" || got.Document.Data["email"] != "customer@example.test" {
		t.Fatalf("result = %#v, want wrapped document data", got)
	}
}

func TestSearchDocumentsRejectsRangeAboveVTEXLimit(t *testing.T) {
	t.Parallel()

	client := &recordingClient{}
	_, err := masterdata.NewSearchDocumentsUseCase(client).Execute(context.Background(), masterdata.SearchDocumentsInput{Brand: "exito", Entity: "CL", RangeFrom: 0, RangeTo: 100})
	if err == nil {
		t.Fatalf("Execute() error = nil, want range validation error")
	}
	var structured capability.StructuredError
	if !errors.As(err, &structured) || structured.Code != masterdata.ErrorMasterDataInvalidInput {
		t.Fatalf("error = %#v, want %s", err, masterdata.ErrorMasterDataInvalidInput)
	}
	if client.searchInput.Entity != "" {
		t.Fatalf("client was called with %#v", client.searchInput)
	}
}

func TestScrollDocumentsReturnsPaginationAndEnforcesSize(t *testing.T) {
	t.Parallel()

	client := &recordingClient{}
	got, err := masterdata.NewScrollDocumentsCapability(client).Handler(context.Background(), capability.ExecutionRequest{Input: capability.Input{"entity": "CL", "size": 1000}})
	if err != nil {
		t.Fatalf("Handler() error = %v", err)
	}
	if got.Pagination == nil || got.Pagination.NextCursor != "scroll-token" || !got.Pagination.HasMore {
		t.Fatalf("pagination = %#v, want scroll token cursor", got.Pagination)
	}

	_, err = masterdata.NewScrollDocumentsUseCase(client).Execute(context.Background(), masterdata.ScrollDocumentsInput{Brand: "exito", Entity: "CL", Size: 1001})
	if err == nil {
		t.Fatalf("Execute() error = nil, want invalid size")
	}
}

func TestSearchDocumentsCapabilityWarnsWithoutSort(t *testing.T) {
	t.Parallel()

	client := &recordingClient{}
	got, err := masterdata.NewSearchDocumentsCapability(client).Handler(context.Background(), capability.ExecutionRequest{Input: capability.Input{"entity": "CL", "fields": []any{"email"}, "rangeFrom": 0, "rangeTo": 99}})
	if err != nil {
		t.Fatalf("Handler() error = %v", err)
	}
	if len(got.Warnings) != 1 || got.Warnings[0].Code != masterdata.WarningMasterDataSearchWithoutSort {
		t.Fatalf("warnings = %#v, want missing sort warning", got.Warnings)
	}
}

func TestMasterDataDefinitionsAreReadOnly(t *testing.T) {
	t.Parallel()

	definitions := []capability.Definition{
		masterdata.GetDocumentDefinition(),
		masterdata.SearchDocumentsDefinition(),
		masterdata.ScrollDocumentsDefinition(),
		masterdata.ListSchemasDefinition(),
		masterdata.GetSchemaDefinition(),
		masterdata.ListIndicesDefinition(),
	}
	for _, definition := range definitions {
		if definition.Domain != masterdata.DomainName || definition.Risk != capability.RiskReadOnly || definition.RequiresConfirmation {
			t.Fatalf("definition = %#v, want read-only Master Data capability", definition)
		}
	}
}

func TestBrandClientRoutesBrandsAndFailsWhenUnavailable(t *testing.T) {
	t.Parallel()

	exitoClient := &recordingClient{}
	carullaClient := &recordingClient{}
	client := masterdata.NewBrandClient(exitoClient, carullaClient)

	if _, err := client.GetDocument(context.Background(), masterdata.GetDocumentInput{Brand: "carulla", Entity: "EX", DocumentID: "doc-1"}); err != nil {
		t.Fatalf("GetDocument(carulla) error = %v", err)
	}
	if carullaClient.getInput.Brand != "carulla" || exitoClient.getInput.Brand != "" {
		t.Fatalf("routing = exito:%#v carulla:%#v, want carulla client only", exitoClient.getInput, carullaClient.getInput)
	}

	missing := masterdata.NewBrandClient(nil, nil)
	_, err := missing.ListSchemas(context.Background(), masterdata.EntityInput{Brand: "exito", Entity: "EX"})
	if err == nil {
		t.Fatalf("ListSchemas() error = nil, want not configured")
	}
	var structured capability.StructuredError
	if !errors.As(err, &structured) || structured.Code != masterdata.ErrorMasterDataNotConfigured {
		t.Fatalf("error = %#v, want %s", err, masterdata.ErrorMasterDataNotConfigured)
	}
}
