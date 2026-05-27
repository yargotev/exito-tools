package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	// DefaultProfile is the profile used when no explicit or saved profile exists.
	DefaultProfile = "staging"

	envConfigPath = "EXITO_CONFIG"
	envProfile    = "EXITO_PROFILE"

	envGeoBaseURL = "EXITO_GEO_BASE_URL"
	envGeoToken   = "EXITO_GEO_TOKEN" // #nosec G101 -- environment variable name, not a credential value.

	envOrdersBaseURL       = "EXITO_ORDERS_BASE_URL"
	envOrdersToken         = "EXITO_ORDERS_TOKEN" // #nosec G101 -- environment variable name, not a credential value.
	envOrdersClientID      = "EXITO_ORDERS_CLIENT_ID"
	envOrdersClientSecret  = "EXITO_ORDERS_CLIENT_SECRET" // #nosec G101 -- environment variable name, not a credential value.
	envOrdersScope         = "EXITO_ORDERS_SCOPE"
	envOrdersTokenURL      = "EXITO_ORDERS_TOKEN_URL"
	envGEOMSCredentialsQA  = "GEOMS_CREDENTIALS_QA"  // #nosec G101 -- environment variable name, not a credential value.
	envGEOMSCredentialsPDN = "GEOMS_CREDENTIALS_PDN" // #nosec G101 -- environment variable name, not a credential value.

	envExitoVTEXOMSBaseURLQA      = "EXITO_VTEX_OMS_BASE_URL_QA"
	envExitoVTEXOMSBaseURLProd    = "EXITO_VTEX_OMS_BASE_URL_PROD"
	envExitoVTEXOMSAppKeyQA       = "EXITO_APP_KEY_QA"     // #nosec G101 -- environment variable name, not a credential value.
	envExitoVTEXOMSAppTokenQA     = "EXITO_APP_TOKEN_QA"   // #nosec G101 -- environment variable name, not a credential value.
	envExitoVTEXOMSAppKeyProd     = "EXITO_APP_KEY_PROD"   // #nosec G101 -- environment variable name, not a credential value.
	envExitoVTEXOMSAppTokenProd   = "EXITO_APP_TOKEN_PROD" // #nosec G101 -- environment variable name, not a credential value.
	envCarullaVTEXOMSBaseURLQA    = "CARULLA_VTEX_OMS_BASE_URL_QA"
	envCarullaVTEXOMSBaseURLProd  = "CARULLA_VTEX_OMS_BASE_URL_PROD"
	envCarullaVTEXOMSAppKeyQA     = "CARULLA_APP_KEY_QA"   // #nosec G101 -- environment variable name, not a credential value.
	envCarullaVTEXOMSAppTokenQA   = "CARULLA_APP_TOKEN_QA" // #nosec G101 -- environment variable name, not a credential value.
	envCarullaVTEXOMSAppKeyProd   = "CARULLA_APP_KEY"      // #nosec G101 -- environment variable name, not a credential value.
	envCarullaVTEXOMSAppTokenProd = "CARULLA_APP_TOKEN"    // #nosec G101 -- environment variable name, not a credential value.
)

// Source identifies where a resolved value came from.
type Source string

const (
	SourceExplicit     Source = "explicit"
	SourceEnvironment  Source = "environment"
	SourceLocalProject Source = "local-project"
	SourceUserConfig   Source = "user-config"
	SourceConfigFile   Source = "config-file"
	SourceSavedDefault Source = "saved-default"
	SourceDefault      Source = "default"
	SourceDotenv       Source = "dotenv"
)

// Options contains all inputs needed to resolve Application Configuration.
type Options struct {
	// ConfigPath is the explicit configuration file path, typically from --config.
	ConfigPath string
	// Profile is the explicit profile, typically from --profile.
	Profile string
	// SavedDefaultProfile is the configured default profile loaded by a future parser.
	SavedDefaultProfile string

	// Env provides environment values. When nil, os.Environ is used.
	Env map[string]string
	// WorkDir is the directory used for local project configuration and dotenv files.
	WorkDir string
	// HomeDir is the user home directory used for user-level configuration.
	HomeDir string
}

// Effective is the deterministic output of configuration resolution.
type Effective struct {
	Profile       string
	ProfileSource Source

	ConfigPath       string
	ConfigSource     Source
	ConfigCandidates []ConfigCandidate

	CredentialLayers []CredentialLayer
	GeoProvider      GeoProvider
	OrdersProvider   OrdersProvider
	VTEXOMSProvider  VTEXOMSProvider
}

// GeoProvider contains the resolved Geo provider configuration.
type GeoProvider struct {
	BaseURL       string `json:"baseUrl,omitempty"`
	BaseURLSource Source `json:"baseUrlSource,omitempty"`

	// Token is intentionally omitted from JSON so secrets are not exposed by
	// accidental envelope/debug serialization of Effective configuration.
	Token       string `json:"-"`
	TokenSource Source `json:"tokenSource,omitempty"`
	TokenSet    bool   `json:"tokenSet"`

	Configured bool `json:"configured"`
}

// VTEXOMSProvider contains the resolved VTEX OMS provider configuration.
type VTEXOMSProvider struct {
	Exito   VTEXOMSBrandProvider `json:"exito"`
	Carulla VTEXOMSBrandProvider `json:"carulla"`
}

// VTEXOMSBrandProvider contains one brand's resolved VTEX OMS provider configuration.
type VTEXOMSBrandProvider struct {
	BaseURL       string `json:"baseUrl,omitempty"`
	BaseURLSource Source `json:"baseUrlSource,omitempty"`

	// AppKey and AppToken are intentionally omitted from JSON so VTEX server-side
	// credentials cannot leak through effective configuration serialization.
	AppKey         string `json:"-"`
	AppKeySource   Source `json:"appKeySource,omitempty"`
	AppKeySet      bool   `json:"appKeySet"`
	AppToken       string `json:"-"`
	AppTokenSource Source `json:"appTokenSource,omitempty"`
	AppTokenSet    bool   `json:"appTokenSet"`
	Configured     bool   `json:"configured"`
}

// OrdersProvider contains the resolved Orders provider configuration.
type OrdersProvider struct {
	BaseURL       string `json:"baseUrl,omitempty"`
	BaseURLSource Source `json:"baseUrlSource,omitempty"`

	// Token is intentionally omitted from JSON so secrets are not exposed by
	// accidental envelope/debug serialization of Effective configuration. A
	// pre-fetched token is still supported for tests/manual fallback, but GEOMS
	// normally uses client credentials and the Orders HTTP client obtains tokens.
	Token       string `json:"-"`
	TokenSource Source `json:"tokenSource,omitempty"`
	TokenSet    bool   `json:"tokenSet"`

	TokenURL       string `json:"tokenUrl,omitempty"`
	TokenURLSource Source `json:"tokenUrlSource,omitempty"`
	ClientID       string `json:"-"`
	ClientIDSource Source `json:"clientIdSource,omitempty"`
	ClientIDSet    bool   `json:"clientIdSet"`
	ClientSecret   string `json:"-"`
	SecretSource   Source `json:"secretSource,omitempty"`
	SecretSet      bool   `json:"secretSet"`
	Scope          string `json:"-"`
	ScopeSource    Source `json:"scopeSource,omitempty"`
	ScopeSet       bool   `json:"scopeSet"`

	Configured bool `json:"configured"`
}

// ConfigCandidate describes a configuration file candidate considered by the resolver.
type ConfigCandidate struct {
	Source Source
	Path   string
	Exists bool
}

// CredentialLayer describes a possible source of sensitive values.
type CredentialLayer struct {
	Source Source
	Name   string
	Path   string
}

// Resolve applies Exito Tools configuration precedence rules.
func Resolve(options Options) (Effective, error) {
	resolvedOptions, err := normalizeOptions(options)
	if err != nil {
		return Effective{}, err
	}

	configPath, configSource, candidates := resolveConfigPath(resolvedOptions)
	savedDefaultProfile, err := savedDefaultProfile(resolvedOptions, configPath)
	if err != nil {
		return Effective{}, err
	}
	resolvedOptions.SavedDefaultProfile = savedDefaultProfile

	profile, profileSource := resolveProfile(resolvedOptions)
	credentialLayers := credentialLayers(resolvedOptions.WorkDir, profile)
	yamlProviders, err := readYAMLProfileProviders(configPath, profile)
	if err != nil {
		return Effective{}, err
	}
	geoProvider, err := resolveGeoProvider(resolvedOptions.Env, credentialLayers, yamlProviders.GeoBaseURL)
	if err != nil {
		return Effective{}, err
	}
	ordersProvider, err := resolveOrdersProvider(resolvedOptions.Env, credentialLayers, yamlProviders.OrdersBaseURL, profile)
	if err != nil {
		return Effective{}, err
	}
	vtexOMSProvider, err := resolveVTEXOMSProvider(resolvedOptions.Env, credentialLayers, yamlProviders.VTEXOMSBaseURLs, profile)
	if err != nil {
		return Effective{}, err
	}

	return Effective{
		Profile:          profile,
		ProfileSource:    profileSource,
		ConfigPath:       configPath,
		ConfigSource:     configSource,
		ConfigCandidates: candidates,
		CredentialLayers: credentialLayers,
		GeoProvider:      geoProvider,
		OrdersProvider:   ordersProvider,
		VTEXOMSProvider:  vtexOMSProvider,
	}, nil
}

func savedDefaultProfile(options Options, configPath string) (string, error) {
	if profile := strings.TrimSpace(options.SavedDefaultProfile); profile != "" {
		return profile, nil
	}
	if strings.TrimSpace(configPath) == "" || !fileExists(configPath) {
		return "", nil
	}

	profile, err := readDefaultProfile(configPath)
	if err != nil {
		return "", err
	}
	return profile, nil
}

func readDefaultProfile(path string) (string, error) {
	file, err := os.Open(path) // #nosec G304 -- path comes from deterministic configuration path resolution.
	if err != nil {
		return "", fmt.Errorf("read configuration file %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if key, value, ok := strings.Cut(line, ":"); ok && strings.TrimSpace(key) == "defaultProfile" {
			return strings.TrimSpace(unquoteDotenvValue(stripInlineYAMLComment(value))), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("scan configuration file %q: %w", path, err)
	}

	return "", nil
}

func stripInlineYAMLComment(value string) string {
	inSingleQuote := false
	inDoubleQuote := false
	for index, r := range value {
		switch r {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case '#':
			if !inSingleQuote && !inDoubleQuote {
				return strings.TrimSpace(value[:index])
			}
		}
	}
	return strings.TrimSpace(value)
}

func resolveGeoProvider(env map[string]string, layers []CredentialLayer, yamlBaseURL string) (GeoProvider, error) {
	provider, err := resolveProvider(env, layers, envGeoBaseURL, envGeoToken, yamlBaseURL)
	if err != nil {
		return GeoProvider{}, err
	}

	return GeoProvider(provider), nil
}

func resolveOrdersProvider(env map[string]string, layers []CredentialLayer, yamlBaseURL string, profile string) (OrdersProvider, error) {
	provider, err := resolveProvider(env, layers, envOrdersBaseURL, envOrdersToken, yamlBaseURL)
	if err != nil {
		return OrdersProvider{}, err
	}

	values := map[string]resolvedValue{}
	for _, layer := range layers {
		switch layer.Source {
		case SourceEnvironment:
			setLayerValue(values, envOrdersClientID, env[envOrdersClientID], SourceEnvironment)
			setLayerValue(values, envOrdersClientSecret, env[envOrdersClientSecret], SourceEnvironment)
			setLayerValue(values, envOrdersScope, env[envOrdersScope], SourceEnvironment)
			setLayerValue(values, envOrdersTokenURL, env[envOrdersTokenURL], SourceEnvironment)
			setOrdersCredentialsValue(values, env[geomsCredentialsKey(profile)], SourceEnvironment)
		case SourceDotenv:
			dotenv, err := readDotenvFile(layer.Path)
			if err != nil {
				return OrdersProvider{}, err
			}
			setLayerValue(values, envOrdersClientID, dotenv[envOrdersClientID], SourceDotenv)
			setLayerValue(values, envOrdersClientSecret, dotenv[envOrdersClientSecret], SourceDotenv)
			setLayerValue(values, envOrdersScope, dotenv[envOrdersScope], SourceDotenv)
			setLayerValue(values, envOrdersTokenURL, dotenv[envOrdersTokenURL], SourceDotenv)
			setOrdersCredentialsValue(values, dotenv[geomsCredentialsKey(profile)], SourceDotenv)
		}
	}

	clientID := values[envOrdersClientID]
	clientSecret := values[envOrdersClientSecret]
	scope := values[envOrdersScope]
	tokenURL := values[envOrdersTokenURL]
	credentialsConfigured := clientID.Value != "" && clientSecret.Value != "" && scope.Value != ""

	return OrdersProvider{
		BaseURL:        provider.BaseURL,
		BaseURLSource:  provider.BaseURLSource,
		Token:          provider.Token,
		TokenSource:    provider.TokenSource,
		TokenSet:       provider.TokenSet,
		TokenURL:       tokenURL.Value,
		TokenURLSource: sourceOrDefault(tokenURL.Source),
		ClientID:       clientID.Value,
		ClientIDSource: sourceOrDefault(clientID.Source),
		ClientIDSet:    clientID.Value != "",
		ClientSecret:   clientSecret.Value,
		SecretSource:   sourceOrDefault(clientSecret.Source),
		SecretSet:      clientSecret.Value != "",
		Scope:          scope.Value,
		ScopeSource:    sourceOrDefault(scope.Source),
		ScopeSet:       scope.Value != "",
		Configured:     provider.BaseURL != "" && (provider.TokenSet || credentialsConfigured),
	}, nil
}

func geomsCredentialsKey(profile string) string {
	switch strings.ToLower(strings.TrimSpace(profile)) {
	case "prod", "production", "pdn":
		return envGEOMSCredentialsPDN
	default:
		return envGEOMSCredentialsQA
	}
}

func resolveVTEXOMSProvider(env map[string]string, layers []CredentialLayer, yamlBaseURLs map[string]string, profile string) (VTEXOMSProvider, error) {
	exito, err := resolveVTEXOMSBrandProvider(env, layers, yamlBaseURLs["exito"], profile, vtexOMSEnvKeys{
		BaseURLQA: envExitoVTEXOMSBaseURLQA, BaseURLProd: envExitoVTEXOMSBaseURLProd,
		AppKeyQA: envExitoVTEXOMSAppKeyQA, AppTokenQA: envExitoVTEXOMSAppTokenQA,
		AppKeyProd: envExitoVTEXOMSAppKeyProd, AppTokenProd: envExitoVTEXOMSAppTokenProd,
	})
	if err != nil {
		return VTEXOMSProvider{}, err
	}
	carulla, err := resolveVTEXOMSBrandProvider(env, layers, yamlBaseURLs["carulla"], profile, vtexOMSEnvKeys{
		BaseURLQA: envCarullaVTEXOMSBaseURLQA, BaseURLProd: envCarullaVTEXOMSBaseURLProd,
		AppKeyQA: envCarullaVTEXOMSAppKeyQA, AppTokenQA: envCarullaVTEXOMSAppTokenQA,
		AppKeyProd: envCarullaVTEXOMSAppKeyProd, AppTokenProd: envCarullaVTEXOMSAppTokenProd,
	})
	if err != nil {
		return VTEXOMSProvider{}, err
	}
	return VTEXOMSProvider{Exito: exito, Carulla: carulla}, nil
}

type vtexOMSEnvKeys struct {
	BaseURLQA    string
	BaseURLProd  string
	AppKeyQA     string
	AppTokenQA   string
	AppKeyProd   string
	AppTokenProd string
}

func resolveVTEXOMSBrandProvider(env map[string]string, layers []CredentialLayer, yamlBaseURL string, profile string, keys vtexOMSEnvKeys) (VTEXOMSBrandProvider, error) {
	baseURLKey, appKeyKey, appTokenKey := keys.forProfile(profile)
	values := map[string]resolvedValue{}
	for _, layer := range layers {
		switch layer.Source {
		case SourceEnvironment:
			setLayerValue(values, baseURLKey, env[baseURLKey], SourceEnvironment)
			setLayerValue(values, appKeyKey, env[appKeyKey], SourceEnvironment)
			setLayerValue(values, appTokenKey, env[appTokenKey], SourceEnvironment)
		case SourceDotenv:
			dotenv, err := readDotenvFile(layer.Path)
			if err != nil {
				return VTEXOMSBrandProvider{}, err
			}
			setLayerValue(values, baseURLKey, dotenv[baseURLKey], SourceDotenv)
			setLayerValue(values, appKeyKey, dotenv[appKeyKey], SourceDotenv)
			setLayerValue(values, appTokenKey, dotenv[appTokenKey], SourceDotenv)
		}
	}
	setLayerValue(values, baseURLKey, yamlBaseURL, SourceConfigFile)

	baseURL := values[baseURLKey]
	appKey := values[appKeyKey]
	appToken := values[appTokenKey]
	return VTEXOMSBrandProvider{
		BaseURL: baseURL.Value, BaseURLSource: sourceOrDefault(baseURL.Source),
		AppKey: appKey.Value, AppKeySource: sourceOrDefault(appKey.Source), AppKeySet: appKey.Value != "",
		AppToken: appToken.Value, AppTokenSource: sourceOrDefault(appToken.Source), AppTokenSet: appToken.Value != "",
		Configured: baseURL.Value != "" && appKey.Value != "" && appToken.Value != "",
	}, nil
}

func (k vtexOMSEnvKeys) forProfile(profile string) (baseURL string, appKey string, appToken string) {
	switch strings.ToLower(strings.TrimSpace(profile)) {
	case "prod", "production", "pdn":
		return k.BaseURLProd, k.AppKeyProd, k.AppTokenProd
	default:
		return k.BaseURLQA, k.AppKeyQA, k.AppTokenQA
	}
}

var credentialsPairPattern = regexp.MustCompile(`["']([^"']+)["']\s*:\s*["']([^"']*)["']`)

func setOrdersCredentialsValue(values map[string]resolvedValue, raw string, source Source) {
	credentials := parseCredentialsMap(raw)
	setLayerValue(values, envOrdersClientID, credentials["client_id"], source)
	setLayerValue(values, envOrdersClientSecret, credentials["client_secret"], source)
	setLayerValue(values, envOrdersScope, credentials["scope"], source)
}

func parseCredentialsMap(raw string) map[string]string {
	parsed := map[string]string{}
	for _, match := range credentialsPairPattern.FindAllStringSubmatch(raw, -1) {
		if len(match) == 3 {
			parsed[strings.TrimSpace(match[1])] = strings.TrimSpace(match[2])
		}
	}
	return parsed
}

type provider struct {
	BaseURL       string `json:"baseUrl,omitempty"`
	BaseURLSource Source `json:"baseUrlSource,omitempty"`
	Token         string `json:"-"`
	TokenSource   Source `json:"tokenSource,omitempty"`
	TokenSet      bool   `json:"tokenSet"`
	Configured    bool   `json:"configured"`
}

func resolveProvider(env map[string]string, layers []CredentialLayer, baseURLKey string, tokenKey string, yamlBaseURL string) (provider, error) {
	values := map[string]resolvedValue{}

	for _, layer := range layers {
		switch layer.Source {
		case SourceEnvironment:
			setLayerValue(values, baseURLKey, env[baseURLKey], SourceEnvironment)
			setLayerValue(values, tokenKey, env[tokenKey], SourceEnvironment)
		case SourceDotenv:
			dotenv, err := readDotenvFile(layer.Path)
			if err != nil {
				return provider{}, err
			}
			setLayerValue(values, baseURLKey, dotenv[baseURLKey], SourceDotenv)
			setLayerValue(values, tokenKey, dotenv[tokenKey], SourceDotenv)
		}
	}
	setLayerValue(values, baseURLKey, yamlBaseURL, SourceConfigFile)

	baseURL := values[baseURLKey]
	token := values[tokenKey]

	return provider{
		BaseURL:       baseURL.Value,
		BaseURLSource: sourceOrDefault(baseURL.Source),
		Token:         token.Value,
		TokenSource:   sourceOrDefault(token.Source),
		TokenSet:      token.Value != "",
		Configured:    baseURL.Value != "" && token.Value != "",
	}, nil
}

type yamlProfileProviders struct {
	GeoBaseURL      string
	OrdersBaseURL   string
	VTEXOMSBaseURLs map[string]string
}

func readYAMLProfileProviders(path string, profile string) (yamlProfileProviders, error) {
	if strings.TrimSpace(path) == "" || strings.TrimSpace(profile) == "" || !fileExists(path) {
		return yamlProfileProviders{}, nil
	}

	file, err := os.Open(path) // #nosec G304 -- path comes from deterministic configuration path resolution.
	if err != nil {
		return yamlProfileProviders{}, fmt.Errorf("read configuration file %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	var providers yamlProfileProviders
	providers.VTEXOMSBaseURLs = map[string]string{}
	inProfiles := false
	inSelectedProfile := false
	currentProvider := ""
	currentVTEXOMSBrand := ""

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		indent := leadingSpaces(raw)
		key, value, ok := yamlKeyValue(trimmed)
		if !ok {
			continue
		}

		switch indent {
		case 0:
			inProfiles = key == "profiles" && value == ""
			inSelectedProfile = false
			currentProvider = ""
			currentVTEXOMSBrand = ""
		case 2:
			if !inProfiles {
				continue
			}
			inSelectedProfile = key == profile && value == ""
			currentProvider = ""
			currentVTEXOMSBrand = ""
		case 4:
			if !inProfiles || !inSelectedProfile {
				continue
			}
			currentVTEXOMSBrand = ""
			if (key == "geo" || key == "orders" || key == "vtexOms") && value == "" {
				currentProvider = key
			} else {
				currentProvider = ""
			}
		case 6:
			if !inProfiles || !inSelectedProfile || currentProvider == "" {
				continue
			}
			if currentProvider == "vtexOms" {
				if (key == "exito" || key == "carulla") && value == "" {
					currentVTEXOMSBrand = key
				}
				continue
			}
			if key != "baseUrl" && key != "baseURL" {
				continue
			}
			switch currentProvider {
			case "geo":
				providers.GeoBaseURL = value
			case "orders":
				providers.OrdersBaseURL = value
			}
		case 8:
			if !inProfiles || !inSelectedProfile || currentProvider != "vtexOms" || currentVTEXOMSBrand == "" {
				continue
			}
			if key == "baseUrl" || key == "baseURL" {
				providers.VTEXOMSBaseURLs[currentVTEXOMSBrand] = value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return yamlProfileProviders{}, fmt.Errorf("scan configuration file %q: %w", path, err)
	}

	return providers, nil
}

func leadingSpaces(value string) int {
	return len(value) - len(strings.TrimLeft(value, " "))
}

func yamlKeyValue(line string) (string, string, bool) {
	key, value, ok := strings.Cut(line, ":")
	if !ok {
		return "", "", false
	}
	return strings.TrimSpace(unquoteDotenvValue(key)), strings.TrimSpace(unquoteDotenvValue(stripInlineYAMLComment(value))), true
}

type resolvedValue struct {
	Value  string
	Source Source
}

func setLayerValue(values map[string]resolvedValue, key string, value string, source Source) {
	if _, exists := values[key]; exists {
		return
	}
	if trimmed := strings.TrimSpace(value); trimmed != "" {
		values[key] = resolvedValue{Value: trimmed, Source: source}
	}
}

func sourceOrDefault(source Source) Source {
	if source == "" {
		return SourceDefault
	}
	return source
}

func readDotenvFile(path string) (map[string]string, error) {
	values := map[string]string{}
	if strings.TrimSpace(path) == "" {
		return values, nil
	}

	file, err := os.Open(path) // #nosec G304 -- path comes from deterministic credential layer locations.
	if err != nil {
		if os.IsNotExist(err) {
			return values, nil
		}
		return nil, fmt.Errorf("read dotenv file %q: %w", path, err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		values[key] = unquoteDotenvValue(strings.TrimSpace(value))
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan dotenv file %q: %w", path, err)
	}

	return values, nil
}

func unquoteDotenvValue(value string) string {
	if len(value) < 2 {
		return value
	}
	if (strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`)) || (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
		return value[1 : len(value)-1]
	}
	return value
}

func normalizeOptions(options Options) (Options, error) {
	if options.Env == nil {
		options.Env = envMap(os.Environ())
	}

	if strings.TrimSpace(options.WorkDir) == "" {
		workDir, err := os.Getwd()
		if err != nil {
			return Options{}, err
		}
		options.WorkDir = workDir
	}

	if strings.TrimSpace(options.HomeDir) == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return Options{}, err
		}
		options.HomeDir = homeDir
	}

	return options, nil
}

func resolveProfile(options Options) (string, Source) {
	if profile := strings.TrimSpace(options.Profile); profile != "" {
		return profile, SourceExplicit
	}

	if profile := strings.TrimSpace(options.Env[envProfile]); profile != "" {
		return profile, SourceEnvironment
	}

	if profile := strings.TrimSpace(options.SavedDefaultProfile); profile != "" {
		return profile, SourceSavedDefault
	}

	return DefaultProfile, SourceDefault
}

func resolveConfigPath(options Options) (string, Source, []ConfigCandidate) {
	if configPath := strings.TrimSpace(options.ConfigPath); configPath != "" {
		return cleanPath(options.WorkDir, configPath), SourceExplicit, nil
	}

	if configPath := strings.TrimSpace(options.Env[envConfigPath]); configPath != "" {
		return cleanPath(options.WorkDir, configPath), SourceEnvironment, nil
	}

	candidates := []ConfigCandidate{
		{Source: SourceLocalProject, Path: filepath.Join(options.WorkDir, "exito.yaml")},
		{Source: SourceUserConfig, Path: filepath.Join(options.HomeDir, ".config", "exito-tools", "config.yaml")},
	}

	selectedIndex := -1
	for i := range candidates {
		candidates[i].Exists = fileExists(candidates[i].Path)
		if candidates[i].Exists && selectedIndex == -1 {
			selectedIndex = i
		}
	}

	if selectedIndex >= 0 {
		selected := candidates[selectedIndex]
		return selected.Path, selected.Source, candidates
	}

	return "", SourceDefault, candidates
}

func credentialLayers(workDir string, profile string) []CredentialLayer {
	return []CredentialLayer{
		{Source: SourceEnvironment, Name: "process environment"},
		{Source: SourceDotenv, Name: ".env." + profile, Path: filepath.Join(workDir, ".env."+profile)},
		{Source: SourceDotenv, Name: ".env", Path: filepath.Join(workDir, ".env")},
	}
}

func cleanPath(workDir string, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}

	return filepath.Join(workDir, path)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func envMap(environ []string) map[string]string {
	mapped := make(map[string]string, len(environ))
	for _, entry := range environ {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			mapped[key] = value
		}
	}
	return mapped
}
