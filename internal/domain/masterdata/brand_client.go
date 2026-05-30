package masterdata

import "context"

type BrandClient struct {
	exito   Client
	carulla Client
}

func NewBrandClient(exito Client, carulla Client) BrandClient {
	return BrandClient{exito: exito, carulla: carulla}
}

func (c BrandClient) GetDocument(ctx context.Context, input GetDocumentInput) (DocumentResult, error) {
	return c.client(input.Brand).GetDocument(ctx, input)
}

func (c BrandClient) SearchDocuments(ctx context.Context, input SearchDocumentsInput) (DocumentsPage, error) {
	return c.client(input.Brand).SearchDocuments(ctx, input)
}

func (c BrandClient) ScrollDocuments(ctx context.Context, input ScrollDocumentsInput) (DocumentsPage, error) {
	return c.client(input.Brand).ScrollDocuments(ctx, input)
}

func (c BrandClient) ListSchemas(ctx context.Context, input EntityInput) (SchemasResult, error) {
	return c.client(input.Brand).ListSchemas(ctx, input)
}

func (c BrandClient) GetSchema(ctx context.Context, input GetSchemaInput) (SchemaResult, error) {
	return c.client(input.Brand).GetSchema(ctx, input)
}

func (c BrandClient) ListIndices(ctx context.Context, input EntityInput) (IndicesResult, error) {
	return c.client(input.Brand).ListIndices(ctx, input)
}

func (c BrandClient) client(brand string) Client {
	if normalizedBrand(brand) == "carulla" {
		if c.carulla != nil {
			return c.carulla
		}
		return UnavailableClient{}
	}
	if c.exito != nil {
		return c.exito
	}
	return UnavailableClient{}
}

type UnavailableClient struct{}

func (UnavailableClient) GetDocument(context.Context, GetDocumentInput) (DocumentResult, error) {
	return DocumentResult{}, notConfiguredError()
}

func (UnavailableClient) SearchDocuments(context.Context, SearchDocumentsInput) (DocumentsPage, error) {
	return DocumentsPage{}, notConfiguredError()
}

func (UnavailableClient) ScrollDocuments(context.Context, ScrollDocumentsInput) (DocumentsPage, error) {
	return DocumentsPage{}, notConfiguredError()
}

func (UnavailableClient) ListSchemas(context.Context, EntityInput) (SchemasResult, error) {
	return SchemasResult{}, notConfiguredError()
}

func (UnavailableClient) GetSchema(context.Context, GetSchemaInput) (SchemaResult, error) {
	return SchemaResult{}, notConfiguredError()
}

func (UnavailableClient) ListIndices(context.Context, EntityInput) (IndicesResult, error) {
	return IndicesResult{}, notConfiguredError()
}
