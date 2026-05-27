package orders

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/platform/httpclient"
)

const defaultGEOMSTokenURL = "https://login.microsoftonline.com/40f94963-1b34-45ce-a5fb-6f1fde2f1a27/oauth2/v2.0/token" // #nosec G101 -- token endpoint URL, not a credential value.

// HTTPGetterConfig contains the provider settings required by HTTPGetter.
type HTTPGetterConfig struct {
	BaseURL      string
	Token        string
	TokenURL     string
	ClientID     string
	ClientSecret string
	Scope        string
}

// HTTPGetter calls GEOMS findOrders and maps its DTO to domain output.
type HTTPGetter struct {
	baseURL string
	client  httpclient.Client
	tokens  *geomsTokenSource
}

// NewHTTPGetter creates an Orders provider-backed getter.
func NewHTTPGetter(config HTTPGetterConfig, client *http.Client) HTTPGetter {
	baseClient := httpclient.New(httpclient.Config{
		BaseURL: config.BaseURL,
		Token:   strings.TrimSpace(config.Token),
		Client:  client,
	})
	return HTTPGetter{
		baseURL: strings.TrimSpace(config.BaseURL),
		client:  baseClient,
		tokens: newGEOMSTokenSource(geomsTokenConfig{
			Token:        config.Token,
			TokenURL:     config.TokenURL,
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			Scope:        config.Scope,
		}, client),
	}
}

// Get calls GEOMS findOrders endpoint with the order number filter and maps the first matching order.
func (g HTTPGetter) Get(ctx context.Context, input GetInput) (Order, error) {
	if strings.TrimSpace(g.baseURL) == "" || g.tokens == nil || !g.tokens.configured() {
		return Order{}, capability.StructuredError{Code: ErrorOrdersNotConfigured, Message: "Orders client is not configured."}
	}

	accessToken, err := g.tokens.token(ctx)
	if err != nil {
		return Order{}, err
	}

	var providerResponse geomsFindOrdersResponse
	if err := g.postGEOMS(ctx, accessToken, "/findOrders", geomsFindOrdersRequest(input), &providerResponse); err != nil {
		return Order{}, err
	}

	order, ok := providerResponse.firstOrder()
	if !ok {
		return Order{}, capability.StructuredError{Code: ErrorOrderNotFound, Message: "Order not found."}
	}

	details, err := g.getOrderDetails(ctx, accessToken, order.ID)
	if err != nil {
		return Order{}, err
	}
	items, err := g.getOrderItems(ctx, accessToken, order.ID)
	if err != nil {
		return Order{}, err
	}
	order.Details = details
	order.Items = &items
	return order, nil
}

func (g HTTPGetter) getOrderDetails(ctx context.Context, accessToken string, orderNumber string) (map[string]any, error) {
	var response struct {
		Data map[string]any `json:"data"`
	}
	if err := g.postGEOMS(ctx, accessToken, "/getOrder", geomsGetOrderRequest(orderNumber), &response); err != nil {
		return nil, err
	}
	if response.Data == nil {
		return map[string]any{}, nil
	}
	return response.Data, nil
}

func (g HTTPGetter) getOrderItems(ctx context.Context, accessToken string, orderNumber string) (OrderItems, error) {
	food, err := g.findItemsByOrder(ctx, accessToken, orderNumber, false)
	if err != nil {
		return OrderItems{}, err
	}
	notFood, err := g.findItemsByOrder(ctx, accessToken, orderNumber, true)
	if err != nil {
		return OrderItems{}, err
	}
	return OrderItems{Food: food, NotFood: notFood}, nil
}

func (g HTTPGetter) findItemsByOrder(ctx context.Context, accessToken string, orderNumber string, notFood bool) ([]map[string]any, error) {
	var response struct {
		Data []map[string]any `json:"data"`
	}
	if err := g.postGEOMS(ctx, accessToken, "/findItemsByOrder", geomsFindItemsByOrderRequest(orderNumber, notFood), &response); err != nil {
		if structuredCode(err) == ErrorOrderNotFound {
			return []map[string]any{}, nil
		}
		return nil, err
	}
	if response.Data == nil {
		return []map[string]any{}, nil
	}
	return response.Data, nil
}

func (g HTTPGetter) postGEOMS(ctx context.Context, accessToken string, path string, payload any, target any) error {
	request, err := httpclient.New(httpclient.Config{BaseURL: g.baseURL, Token: accessToken}).NewJSONRequest(ctx, http.MethodPost, path, payload)
	if err != nil {
		return capability.StructuredError{Code: ErrorOrdersNotConfigured, Message: "Orders provider base URL is invalid."}
	}
	// Preserve the test/server client selected at construction time.
	response, err := g.client.Do(request)
	if err != nil {
		return capability.StructuredError{Code: ErrorOrdersProviderUnavailable, Message: "Orders provider request failed."}
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == http.StatusNotFound {
		return capability.StructuredError{Code: ErrorOrderNotFound, Message: "Order not found."}
	}
	if !httpclient.Successful(response) {
		return capability.StructuredError{Code: ErrorOrdersProviderUnavailable, Message: "Orders provider returned an unsuccessful response."}
	}
	if err := httpclient.DecodeJSONResponse(response, target); err != nil {
		return capability.StructuredError{Code: ErrorOrdersProviderInvalidResponse, Message: "Orders provider returned an invalid response."}
	}
	return nil
}

func structuredCode(err error) string {
	var structured capability.StructuredError
	if errors.As(err, &structured) {
		return structured.Code
	}
	return ""
}

type geomsFindOrdersProviderRequest struct {
	TransactionID   string           `json:"transactionId"`
	Hostname        string           `json:"hostname"`
	User            string           `json:"user"`
	TransactionDate string           `json:"transactionDate"`
	Filters         geomsOrderFilter `json:"filters"`
}

type geomsGetOrderProviderRequest struct {
	TransactionID   string `json:"transactionId"`
	Hostname        string `json:"hostname"`
	User            string `json:"user"`
	TransactionDate string `json:"transactionDate"`
	Order           string `json:"order"`
}

type geomsFindItemsByOrderProviderRequest struct {
	TransactionID   string `json:"transactionId"`
	Hostname        string `json:"hostname"`
	User            string `json:"user"`
	TransactionDate string `json:"transactionDate"`
	Order           string `json:"order"`
	PerPageItem     int    `json:"perPageItem"`
	PageNumberItem  int    `json:"pageNumberItem"`
	NotFood         bool   `json:"notFood"`
}

type geomsOrderFilter struct {
	PerPage         int     `json:"perPage"`
	PageNumber      int     `json:"pageNumber"`
	OrderNumber     *string `json:"orderNumber"`
	Email           *string `json:"email"`
	CreatedDateFrom *string `json:"createdDateFrom"`
	CreatedDateTo   *string `json:"createdDateTo"`
	OrderType       string  `json:"orderType"`
}

func geomsGetOrderRequest(orderNumber string) geomsGetOrderProviderRequest {
	return geomsGetOrderProviderRequest{
		TransactionID:   requestTransactionID(),
		Hostname:        "vtex",
		User:            "web",
		TransactionDate: transactionDate(),
		Order:           cleanOrderNumber(orderNumber),
	}
}

func geomsFindItemsByOrderRequest(orderNumber string, notFood bool) geomsFindItemsByOrderProviderRequest {
	return geomsFindItemsByOrderProviderRequest{
		TransactionID:   requestTransactionID(),
		Hostname:        "vtex",
		User:            "web",
		TransactionDate: transactionDate(),
		Order:           cleanOrderNumber(orderNumber),
		PerPageItem:     999,
		PageNumberItem:  1,
		NotFood:         notFood,
	}
}

func geomsFindOrdersRequest(input GetInput) geomsFindOrdersProviderRequest {
	id := cleanOrderNumber(input.ID)
	return geomsFindOrdersProviderRequest{
		TransactionID:   requestTransactionID(),
		Hostname:        "vtex",
		User:            "web",
		TransactionDate: transactionDate(),
		Filters: geomsOrderFilter{
			PerPage:     10,
			PageNumber:  1,
			OrderNumber: &id,
			Email:       nil,
			OrderType:   geomsOrderType(input.OrderType),
		},
	}
}

func cleanOrderNumber(orderNumber string) string {
	trimmed := strings.TrimSpace(orderNumber)
	if value, _, ok := strings.Cut(trimmed, "-"); ok {
		return strings.TrimSpace(value)
	}
	return trimmed
}

func transactionDate() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

func geomsOrderType(orderType string) string {
	if trimmed := strings.TrimSpace(orderType); trimmed != "" {
		return trimmed
	}
	return "ExitoEcomm"
}

func requestTransactionID() string {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return fmt.Sprintf("%d", time.Now().UTC().UnixNano())
	}
	value[6] = (value[6] & 0x0f) | 0x40
	value[8] = (value[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", value[0:4], value[4:6], value[6:8], value[8:10], value[10:16])
}

type geomsFindOrdersResponse map[string]any

func (r geomsFindOrdersResponse) firstOrder() (Order, bool) {
	for _, key := range []string{"orders", "Orders", "data", "items", "results"} {
		if raw, exists := r[key]; exists {
			if order, ok := firstOrderFromAny(raw); ok {
				return order, true
			}
		}
	}
	if raw, exists := r["order"]; exists {
		return orderFromAny(raw)
	}
	return orderFromAny(map[string]any(r))
}

func firstOrderFromAny(raw any) (Order, bool) {
	switch value := raw.(type) {
	case []any:
		if len(value) == 0 {
			return Order{}, false
		}
		return orderFromAny(value[0])
	case map[string]any:
		if order, ok := orderFromAny(value); ok {
			return order, true
		}
		for _, key := range []string{"orders", "items", "results", "data"} {
			if nested, exists := value[key]; exists {
				if order, ok := firstOrderFromAny(nested); ok {
					return order, true
				}
			}
		}
	}
	return Order{}, false
}

func orderFromAny(raw any) (Order, bool) {
	fields, ok := raw.(map[string]any)
	if !ok {
		return Order{}, false
	}
	order := Order{
		ID:             firstString(fields, "id", "orderId", "orderID", "orderNumber", "orderNo"),
		Status:         firstString(fields, "status", "orderStatus", "state", "statusOrderMax", "statusOrderMin"),
		CreatedAt:      firstString(fields, "createdAt", "creationDate", "createdDate", "dateCreated"),
		CustomerName:   firstString(fields, "customerName"),
		Email:          firstString(fields, "email"),
		OrderTotal:     firstFloat(fields, "orderTotal"),
		StatusOrderMax: firstString(fields, "statusOrderMax"),
		StatusOrderMin: firstString(fields, "statusOrderMin"),
	}
	return order, order.ID != ""
}

func firstFloat(fields map[string]any, names ...string) float64 {
	for _, name := range names {
		if value, exists := fields[name]; exists {
			switch typed := value.(type) {
			case float64:
				return typed
			case int:
				return float64(typed)
			case json.Number:
				parsed, _ := typed.Float64()
				return parsed
			}
		}
	}
	return 0
}

func firstString(fields map[string]any, names ...string) string {
	for _, name := range names {
		if value, exists := fields[name]; exists {
			switch typed := value.(type) {
			case string:
				return typed
			case json.Number:
				return typed.String()
			case float64:
				return fmt.Sprintf("%.0f", typed)
			}
		}
	}
	return ""
}

type geomsTokenConfig struct {
	Token        string
	TokenURL     string
	ClientID     string
	ClientSecret string
	Scope        string
}

type geomsTokenSource struct {
	config      geomsTokenConfig
	client      *http.Client
	mu          sync.Mutex
	cachedToken string
	expiry      time.Time
}

func newGEOMSTokenSource(config geomsTokenConfig, client *http.Client) *geomsTokenSource {
	if client == nil {
		client = &http.Client{Timeout: httpclient.DefaultTimeout}
	}
	config.Token = strings.TrimSpace(config.Token)
	config.TokenURL = strings.TrimSpace(config.TokenURL)
	config.ClientID = strings.TrimSpace(config.ClientID)
	config.ClientSecret = strings.TrimSpace(config.ClientSecret)
	config.Scope = strings.TrimSpace(config.Scope)
	if config.TokenURL == "" {
		config.TokenURL = defaultGEOMSTokenURL
	}
	return &geomsTokenSource{config: config, client: client}
}

func (s *geomsTokenSource) configured() bool {
	return s.config.Token != "" || (s.config.ClientID != "" && s.config.ClientSecret != "" && s.config.Scope != "")
}

func (s *geomsTokenSource) token(ctx context.Context) (string, error) {
	if s.config.Token != "" {
		return s.config.Token, nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cachedToken != "" && time.Until(s.expiry) > time.Minute {
		return s.cachedToken, nil
	}

	form := url.Values{}
	form.Set("client_id", s.config.ClientID)
	form.Set("client_secret", s.config.ClientSecret)
	form.Set("grant_type", "client_credentials")
	form.Set("scope", "api://"+s.config.Scope+"/.default")

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", capability.StructuredError{Code: ErrorOrdersNotConfigured, Message: "Orders token URL is invalid."}
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	response, err := s.client.Do(request) // #nosec G107,G704 -- URL is resolved from non-sensitive configuration.
	if err != nil {
		return "", capability.StructuredError{Code: ErrorOrdersProviderUnavailable, Message: "Orders token request failed."}
	}
	defer func() { _ = response.Body.Close() }()
	if !httpclient.Successful(response) {
		return "", capability.StructuredError{Code: ErrorOrdersProviderUnavailable, Message: "Orders token endpoint returned an unsuccessful response."}
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	if err := httpclient.DecodeJSONResponse(response, &tokenResponse); err != nil || strings.TrimSpace(tokenResponse.AccessToken) == "" {
		return "", capability.StructuredError{Code: ErrorOrdersProviderInvalidResponse, Message: "Orders token endpoint returned an invalid response."}
	}

	expiresIn := time.Duration(tokenResponse.ExpiresIn) * time.Second
	if expiresIn <= 0 {
		expiresIn = time.Hour
	}
	s.cachedToken = strings.TrimSpace(tokenResponse.AccessToken)
	s.expiry = time.Now().Add(expiresIn)
	return s.cachedToken, nil
}
