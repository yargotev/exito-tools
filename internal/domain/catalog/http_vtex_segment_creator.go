package catalog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

const (
	vtexSessionsPath               = "/io/api/sessions"
	maxVTEXSegmentResponseBodySize = 4 << 20
	redactedValue                  = "[REDACTED]"
)

// HTTPVTEXSegmentCreatorConfig contains the provider settings required by HTTPVTEXSegmentCreator.
type HTTPVTEXSegmentCreatorConfig struct {
	BaseURL string
}

// HTTPVTEXSegmentCreator calls VTEX Sessions and maps its DTO to domain output.
type HTTPVTEXSegmentCreator struct {
	baseURL string
	client  httpclient.Client
}

// NewHTTPVTEXSegmentCreator creates a VTEX Sessions-backed segment creator.
func NewHTTPVTEXSegmentCreator(config HTTPVTEXSegmentCreatorConfig, client *http.Client) HTTPVTEXSegmentCreator {
	return HTTPVTEXSegmentCreator{
		baseURL: strings.TrimSpace(config.BaseURL),
		client:  httpclient.New(httpclient.Config{BaseURL: config.BaseURL, Client: client}),
	}
}

// CreateVTEXSegment creates a VTEX segment token through Sessions.
func (c HTTPVTEXSegmentCreator) CreateVTEXSegment(ctx context.Context, input CreateVTEXSegmentInput) (CreateVTEXSegmentResult, error) {
	if strings.TrimSpace(c.baseURL) == "" {
		return CreateVTEXSegmentResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Sessions client is not configured."}
	}

	payload := vtexSegmentRequestPayload(input)
	body, err := json.Marshal(payload)
	if err != nil {
		return CreateVTEXSegmentResult{}, capability.StructuredError{Code: ErrorCatalogInvalidInput, Message: "VTEX segment request payload is invalid."}
	}
	request, err := c.client.NewRequest(ctx, http.MethodPost, vtexSessionsPath, bytes.NewReader(body))
	if err != nil {
		return CreateVTEXSegmentResult{}, capability.StructuredError{Code: ErrorCatalogNotConfigured, Message: "VTEX Sessions provider base URL is invalid."}
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := c.client.Do(request)
	if err != nil {
		return CreateVTEXSegmentResult{}, capability.StructuredError{Code: ErrorCatalogProviderUnavailable, Message: "VTEX Sessions provider request failed."}
	}
	defer func() { _ = response.Body.Close() }()

	var providerResponse map[string]any
	limited := io.LimitReader(response.Body, maxVTEXSegmentResponseBodySize)
	if err := json.NewDecoder(limited).Decode(&providerResponse); err != nil {
		return CreateVTEXSegmentResult{}, capability.StructuredError{Code: ErrorCatalogProviderInvalidResponse, Message: "VTEX Sessions provider returned an invalid response."}
	}

	if !httpclient.Successful(response) {
		return CreateVTEXSegmentResult{}, capability.StructuredError{Code: ErrorCatalogProviderUnavailable, Message: "VTEX Sessions provider returned an unsuccessful response."}
	}

	token := segmentToken(providerResponse)
	result := CreateVTEXSegmentResult{
		Brand:        normalizedBrand(input.Brand),
		RegionID:     input.RegionID,
		SalesChannel: input.SalesChannel,
		TokenSet:     token != "",
		TokenLength:  len(token),
		Diagnostics: SegmentDiagnostics{
			RequestPath:     vtexSessionsPath,
			RequestPayload:  payload,
			ProviderPayload: redactSegmentProviderPayload(providerResponse),
		},
	}
	if input.IncludeCookie && token != "" {
		result.Cookie = "vtex_segment=" + token
	}
	return result, nil
}

func vtexSegmentRequestPayload(input CreateVTEXSegmentInput) map[string]any {
	return map[string]any{
		"public": map[string]any{
			"regionId": map[string]any{"value": input.RegionID},
			"sc":       map[string]any{"value": input.SalesChannel},
		},
	}
}

func segmentToken(payload map[string]any) string {
	for _, key := range []string{"segmentToken", "vtexSegment", "token"} {
		if value, ok := payload[key].(string); ok && strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	if namespace, ok := payload["namespaces"].(map[string]any); ok {
		if segment, ok := namespace["segment"].(map[string]any); ok {
			if value, ok := segment["token"].(string); ok && strings.TrimSpace(value) != "" {
				return strings.TrimSpace(value)
			}
		}
	}
	return ""
}

func redactSegmentProviderPayload(payload map[string]any) map[string]any {
	redacted := make(map[string]any, len(payload))
	for key, value := range payload {
		redacted[key] = redactSegmentValue(key, value)
	}
	return redacted
}

func redactSegmentValue(key string, value any) any {
	lowerKey := strings.ToLower(key)
	if strings.Contains(lowerKey, "token") || strings.Contains(lowerKey, "cookie") {
		if _, ok := value.(string); ok {
			return redactedValue
		}
	}
	switch typed := value.(type) {
	case map[string]any:
		return redactSegmentProviderPayload(typed)
	case []any:
		out := make([]any, 0, len(typed))
		for _, item := range typed {
			out = append(out, redactSegmentValue(key, item))
		}
		return out
	default:
		return value
	}
}
