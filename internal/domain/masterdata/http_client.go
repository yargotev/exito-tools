package masterdata

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

type HTTPClientConfig struct {
	BaseURL  string
	AppKey   string
	AppToken string
}

type HTTPClient struct {
	baseURL  string
	appKey   string
	appToken string
	client   httpclient.Client
}

func NewHTTPClient(config HTTPClientConfig, client *http.Client) HTTPClient {
	return HTTPClient{baseURL: strings.TrimSpace(config.BaseURL), appKey: strings.TrimSpace(config.AppKey), appToken: strings.TrimSpace(config.AppToken), client: httpclient.New(httpclient.Config{BaseURL: config.BaseURL, Client: client})}
}

func (c HTTPClient) GetDocument(ctx context.Context, input GetDocumentInput) (DocumentResult, error) {
	path := "/api/dataentities/" + url.PathEscape(input.Entity) + "/documents/" + url.PathEscape(input.DocumentID)
	request, err := c.newGET(ctx, path)
	if err != nil {
		return DocumentResult{}, err
	}
	query := request.URL.Query()
	addQuery(query, "_schema", input.Schema)
	addQuery(query, "_fields", strings.Join(input.Fields, ","))
	request.URL.RawQuery = query.Encode()

	var payload map[string]any
	status, err := c.doJSON(request, &payload)
	if err != nil {
		return DocumentResult{}, err
	}
	return DocumentResult{Brand: normalizedBrand(input.Brand), Entity: input.Entity, DocumentID: input.DocumentID, Fields: append([]string(nil), input.Fields...), Data: payload, Diagnostics: Diagnostics{RequestPath: pathWithQuery(path, request.URL.RawQuery), ProviderStatus: status}}, nil
}

func (c HTTPClient) SearchDocuments(ctx context.Context, input SearchDocumentsInput) (DocumentsPage, error) {
	path := "/api/dataentities/" + url.PathEscape(input.Entity) + "/search"
	request, err := c.newGET(ctx, path)
	if err != nil {
		return DocumentsPage{}, err
	}
	query := request.URL.Query()
	addQuery(query, "_fields", strings.Join(input.Fields, ","))
	addQuery(query, "_where", input.Where)
	addQuery(query, "_schema", input.Schema)
	addQuery(query, "_sort", input.Sort)
	request.URL.RawQuery = query.Encode()
	request.Header.Set("REST-Range", "resources="+strconv.Itoa(input.RangeFrom)+"-"+strconv.Itoa(input.RangeTo))

	var payload []map[string]any
	status, headers, err := c.doDocuments(request, &payload)
	if err != nil {
		return DocumentsPage{}, err
	}
	total := parseResourcesTotal(headers.Get("resources"))
	return c.searchResultFromResponse(input, pathWithQuery(path, request.URL.RawQuery), payload, status, total), nil
}

func (c HTTPClient) ScrollDocuments(ctx context.Context, input ScrollDocumentsInput) (DocumentsPage, error) {
	path := "/api/dataentities/" + url.PathEscape(input.Entity) + "/scroll"
	request, err := c.newGET(ctx, path)
	if err != nil {
		return DocumentsPage{}, err
	}
	query := request.URL.Query()
	addQuery(query, "_fields", strings.Join(input.Fields, ","))
	addQuery(query, "_where", input.Where)
	addQuery(query, "_schema", input.Schema)
	query.Set("_size", strconv.Itoa(input.Size))
	request.URL.RawQuery = query.Encode()
	if input.Token != "" {
		request.Header.Set("X-VTEX-MD-TOKEN", input.Token)
	}

	var payload []map[string]any
	status, headers, err := c.doDocuments(request, &payload)
	nextToken := headers.Get("X-VTEX-MD-TOKEN")
	if err != nil {
		return DocumentsPage{}, err
	}
	return DocumentsPage{Brand: normalizedBrand(input.Brand), Entity: input.Entity, Fields: append([]string(nil), input.Fields...), Documents: payload, NextToken: nextToken, Diagnostics: Diagnostics{RequestPath: pathWithQuery(path, request.URL.RawQuery), ProviderStatus: status}}, nil
}

func (c HTTPClient) ListSchemas(ctx context.Context, input EntityInput) (SchemasResult, error) {
	path := "/api/dataentities/" + url.PathEscape(input.Entity) + "/schemas"
	request, err := c.newGET(ctx, path)
	if err != nil {
		return SchemasResult{}, err
	}
	var payload []string
	status, err := c.doJSON(request, &payload)
	if err != nil {
		return SchemasResult{}, err
	}
	return SchemasResult{Brand: normalizedBrand(input.Brand), Entity: input.Entity, Schemas: payload, Diagnostics: Diagnostics{RequestPath: path, ProviderStatus: status}}, nil
}

func (c HTTPClient) GetSchema(ctx context.Context, input GetSchemaInput) (SchemaResult, error) {
	path := "/api/dataentities/" + url.PathEscape(input.Entity) + "/schemas/" + url.PathEscape(input.Schema)
	request, err := c.newGET(ctx, path)
	if err != nil {
		return SchemaResult{}, err
	}
	var payload map[string]any
	status, err := c.doJSON(request, &payload)
	if err != nil {
		return SchemaResult{}, err
	}
	return SchemaResult{Brand: normalizedBrand(input.Brand), Entity: input.Entity, Schema: input.Schema, Definition: payload, Diagnostics: Diagnostics{RequestPath: path, ProviderStatus: status}}, nil
}

func (c HTTPClient) ListIndices(ctx context.Context, input EntityInput) (IndicesResult, error) {
	path := "/api/dataentities/" + url.PathEscape(input.Entity) + "/indices"
	request, err := c.newGET(ctx, path)
	if err != nil {
		return IndicesResult{}, err
	}
	var payload []string
	status, err := c.doJSON(request, &payload)
	if err != nil {
		return IndicesResult{}, err
	}
	return IndicesResult{Brand: normalizedBrand(input.Brand), Entity: input.Entity, Indices: payload, Diagnostics: Diagnostics{RequestPath: path, ProviderStatus: status}}, nil
}

func (c HTTPClient) newGET(ctx context.Context, path string) (*http.Request, error) {
	if strings.TrimSpace(c.baseURL) == "" || strings.TrimSpace(c.appKey) == "" || strings.TrimSpace(c.appToken) == "" {
		return nil, notConfiguredError()
	}
	request, err := c.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, capability.StructuredError{Code: ErrorMasterDataNotConfigured, Message: "VTEX Master Data provider base URL is invalid."}
	}
	request.Header.Set("X-VTEX-API-AppKey", c.appKey)
	request.Header.Set("X-VTEX-API-AppToken", c.appToken)
	return request, nil
}

func (c HTTPClient) doJSON(request *http.Request, target any) (int, error) {
	response, err := c.client.Do(request)
	if err != nil {
		return 0, capability.StructuredError{Code: ErrorMasterDataProviderUnavailable, Message: "VTEX Master Data provider request failed."}
	}
	defer func() { _ = response.Body.Close() }()
	if !httpclient.Successful(response) {
		return response.StatusCode, capability.StructuredError{Code: ErrorMasterDataProviderUnavailable, Message: "VTEX Master Data provider returned an unsuccessful response."}
	}
	if err := httpclient.DecodeJSONResponse(response, target); err != nil {
		return response.StatusCode, capability.StructuredError{Code: ErrorMasterDataProviderInvalidResponse, Message: "VTEX Master Data provider returned an invalid response."}
	}
	return response.StatusCode, nil
}

func (c HTTPClient) doDocuments(request *http.Request, target any) (status int, headers http.Header, err error) {
	response, err := c.client.Do(request)
	if err != nil {
		return 0, nil, capability.StructuredError{Code: ErrorMasterDataProviderUnavailable, Message: "VTEX Master Data provider request failed."}
	}
	defer func() { _ = response.Body.Close() }()
	if !httpclient.Successful(response) {
		return response.StatusCode, response.Header, capability.StructuredError{Code: ErrorMasterDataProviderUnavailable, Message: "VTEX Master Data provider returned an unsuccessful response."}
	}
	if err := httpclient.DecodeJSONResponse(response, target); err != nil {
		return response.StatusCode, response.Header, capability.StructuredError{Code: ErrorMasterDataProviderInvalidResponse, Message: "VTEX Master Data provider returned an invalid response."}
	}
	return response.StatusCode, response.Header, nil
}

func (c HTTPClient) searchResultFromResponse(input SearchDocumentsInput, requestPath string, payload []map[string]any, status int, total *int) DocumentsPage {
	return DocumentsPage{Brand: normalizedBrand(input.Brand), Entity: input.Entity, Fields: append([]string(nil), input.Fields...), Documents: payload, RangeFrom: input.RangeFrom, RangeTo: input.RangeTo, Total: total, Diagnostics: Diagnostics{RequestPath: requestPath, ProviderStatus: status}}
}

func addQuery(query url.Values, key string, value string) {
	if strings.TrimSpace(value) != "" {
		query.Set(key, strings.TrimSpace(value))
	}
}

func pathWithQuery(path string, rawQuery string) string {
	if rawQuery == "" {
		return path
	}
	return path + "?" + rawQuery
}

func parseResourcesTotal(raw string) *int {
	_, totalText, ok := strings.Cut(raw, "/")
	if !ok {
		return nil
	}
	total, err := strconv.Atoi(strings.TrimSpace(totalText))
	if err != nil {
		return nil
	}
	return &total
}
