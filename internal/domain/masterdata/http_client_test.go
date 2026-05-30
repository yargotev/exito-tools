package masterdata_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/domain/masterdata"
)

func TestHTTPClientGetDocumentUsesAuthHeadersAndMapsData(t *testing.T) {
	t.Parallel()

	var gotAppKey string
	var gotAppToken string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAppKey = r.Header.Get("X-VTEX-API-AppKey")
		gotAppToken = r.Header.Get("X-VTEX-API-AppToken")
		if r.Method != http.MethodGet || r.URL.Path != "/api/dataentities/CL/documents/doc-1" {
			t.Fatalf("request = %s %s, want GET document", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("_fields") != "email,firstName" || r.URL.Query().Get("_schema") != "client-v1" {
			t.Fatalf("query = %q, want fields and schema", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"doc-1","email":"customer@example.test","firstName":"Jane"}`))
	}))
	defer server.Close()

	client := masterdata.NewHTTPClient(masterdata.HTTPClientConfig{BaseURL: server.URL, AppKey: "app-key", AppToken: "app-token"}, server.Client())
	got, err := client.GetDocument(context.Background(), masterdata.GetDocumentInput{Brand: "exito", Entity: "CL", DocumentID: "doc-1", Schema: "client-v1", Fields: []string{"email", "firstName"}})
	if err != nil {
		t.Fatalf("GetDocument() error = %v", err)
	}
	if gotAppKey != "app-key" || gotAppToken != "app-token" {
		t.Fatalf("auth headers = (%q,%q), want app credentials", gotAppKey, gotAppToken)
	}
	if got.DocumentID != "doc-1" || got.Data["email"] != "customer@example.test" || got.Diagnostics.ProviderStatus != http.StatusOK {
		t.Fatalf("result = %#v, want mapped document and diagnostics", got)
	}
}

func TestHTTPClientSearchAndScrollPagination(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/dataentities/CL/search":
			if r.Header.Get("REST-Range") != "resources=10-19" {
				t.Fatalf("REST-Range = %q, want resources=10-19", r.Header.Get("REST-Range"))
			}
			if r.URL.Query().Get("_where") != "email is not null" || r.URL.Query().Get("_sort") != "email ASC" {
				t.Fatalf("query = %q, want where/sort", r.URL.RawQuery)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("resources", "10-19/25")
			_, _ = w.Write([]byte(`[{"id":"doc-10"}]`))
		case "/api/dataentities/CL/scroll":
			if r.URL.Query().Get("_size") != "1000" || r.Header.Get("X-VTEX-MD-TOKEN") != "prev-token" {
				t.Fatalf("scroll query/token = %q/%q, want size and previous token", r.URL.RawQuery, r.Header.Get("X-VTEX-MD-TOKEN"))
			}
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-VTEX-MD-TOKEN", "next-token")
			_, _ = w.Write([]byte(`[{"id":"doc-20"}]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := masterdata.NewHTTPClient(masterdata.HTTPClientConfig{BaseURL: server.URL, AppKey: "app-key", AppToken: "app-token"}, server.Client())
	search, err := client.SearchDocuments(context.Background(), masterdata.SearchDocumentsInput{Brand: "exito", Entity: "CL", Fields: []string{"id"}, Where: "email is not null", Sort: "email ASC", RangeFrom: 10, RangeTo: 19})
	if err != nil {
		t.Fatalf("SearchDocuments() error = %v", err)
	}
	if search.Total == nil || *search.Total != 25 || search.RangeFrom != 10 || search.RangeTo != 19 || len(search.Documents) != 1 {
		t.Fatalf("search = %#v, want parsed resources pagination", search)
	}

	scroll, err := client.ScrollDocuments(context.Background(), masterdata.ScrollDocumentsInput{Brand: "exito", Entity: "CL", Fields: []string{"id"}, Size: 1000, Token: "prev-token"})
	if err != nil {
		t.Fatalf("ScrollDocuments() error = %v", err)
	}
	if scroll.NextToken != "next-token" || len(scroll.Documents) != 1 {
		t.Fatalf("scroll = %#v, want next token and documents", scroll)
	}
}

func TestHTTPClientSchemaAndIndexReads(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/dataentities/CL/schemas":
			_, _ = w.Write([]byte(`["client-v1","client-v2"]`))
		case "/api/dataentities/CL/schemas/client-v1":
			_ = json.NewEncoder(w).Encode(map[string]any{"type": "object", "v-indexed": []string{"email"}})
		case "/api/dataentities/CL/indices":
			_, _ = w.Write([]byte(`["email-index"]`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := masterdata.NewHTTPClient(masterdata.HTTPClientConfig{BaseURL: server.URL, AppKey: "app-key", AppToken: "app-token"}, server.Client())
	schemas, err := client.ListSchemas(context.Background(), masterdata.EntityInput{Brand: "exito", Entity: "CL"})
	if err != nil {
		t.Fatalf("ListSchemas() error = %v", err)
	}
	if len(schemas.Schemas) != 2 || schemas.Schemas[0] != "client-v1" {
		t.Fatalf("schemas = %#v, want list", schemas)
	}
	schema, err := client.GetSchema(context.Background(), masterdata.GetSchemaInput{Brand: "exito", Entity: "CL", Schema: "client-v1"})
	if err != nil {
		t.Fatalf("GetSchema() error = %v", err)
	}
	if schema.Definition["type"] != "object" || schema.Diagnostics.RequestPath != "/api/dataentities/CL/schemas/client-v1" {
		t.Fatalf("schema = %#v, want definition and diagnostics", schema)
	}
	indices, err := client.ListIndices(context.Background(), masterdata.EntityInput{Brand: "exito", Entity: "CL"})
	if err != nil {
		t.Fatalf("ListIndices() error = %v", err)
	}
	if len(indices.Indices) != 1 || indices.Indices[0] != "email-index" {
		t.Fatalf("indices = %#v, want list", indices)
	}
}

func TestHTTPClientDiagnosticsOmitCredentialsAndHeaders(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-VTEX-API-AppKey") != "diagnostic-key" || r.Header.Get("X-VTEX-API-AppToken") != "diagnostic-token" {
			t.Fatalf("missing auth headers")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Set-Cookie", "sensitive-cookie=value")
		_, _ = w.Write([]byte(`{"id":"doc-safe"}`))
	}))
	defer server.Close()

	client := masterdata.NewHTTPClient(masterdata.HTTPClientConfig{BaseURL: server.URL, AppKey: "diagnostic-key", AppToken: "diagnostic-token"}, server.Client())
	got, err := client.GetDocument(context.Background(), masterdata.GetDocumentInput{Brand: "exito", Entity: "EX", DocumentID: "doc-safe"})
	if err != nil {
		t.Fatalf("GetDocument() error = %v", err)
	}

	encoded, err := json.Marshal(got.Diagnostics)
	if err != nil {
		t.Fatalf("Marshal diagnostics error = %v", err)
	}
	diagnostics := string(encoded)
	for _, forbidden := range []string{"diagnostic-key", "diagnostic-token", "X-VTEX-API-AppKey", "X-VTEX-API-AppToken", "sensitive-cookie", "Cookie"} {
		if strings.Contains(diagnostics, forbidden) {
			t.Fatalf("diagnostics leaked %q: %s", forbidden, diagnostics)
		}
	}
}
