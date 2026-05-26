package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// DefaultProfile is the profile used when no explicit or saved profile exists.
	DefaultProfile = "staging"

	envConfigPath = "EXITO_CONFIG"
	envProfile    = "EXITO_PROFILE"

	envGeoBaseURL = "EXITO_GEO_BASE_URL"
	envGeoToken   = "EXITO_GEO_TOKEN" // #nosec G101 -- environment variable name, not a credential value.

	envOrdersBaseURL = "EXITO_ORDERS_BASE_URL"
	envOrdersToken   = "EXITO_ORDERS_TOKEN" // #nosec G101 -- environment variable name, not a credential value.
)

// Source identifies where a resolved value came from.
type Source string

const (
	SourceExplicit     Source = "explicit"
	SourceEnvironment  Source = "environment"
	SourceLocalProject Source = "local-project"
	SourceUserConfig   Source = "user-config"
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

// OrdersProvider contains the resolved Orders provider configuration.
type OrdersProvider struct {
	BaseURL       string `json:"baseUrl,omitempty"`
	BaseURLSource Source `json:"baseUrlSource,omitempty"`

	// Token is intentionally omitted from JSON so secrets are not exposed by
	// accidental envelope/debug serialization of Effective configuration.
	Token       string `json:"-"`
	TokenSource Source `json:"tokenSource,omitempty"`
	TokenSet    bool   `json:"tokenSet"`

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

// Resolve applies Exito Tools configuration precedence rules without parsing YAML configuration files.
func Resolve(options Options) (Effective, error) {
	resolvedOptions, err := normalizeOptions(options)
	if err != nil {
		return Effective{}, err
	}

	profile, profileSource := resolveProfile(resolvedOptions)
	configPath, configSource, candidates := resolveConfigPath(resolvedOptions)
	credentialLayers := credentialLayers(resolvedOptions.WorkDir, profile)
	geoProvider, err := resolveGeoProvider(resolvedOptions.Env, credentialLayers)
	if err != nil {
		return Effective{}, err
	}
	ordersProvider, err := resolveOrdersProvider(resolvedOptions.Env, credentialLayers)
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
	}, nil
}

func resolveGeoProvider(env map[string]string, layers []CredentialLayer) (GeoProvider, error) {
	provider, err := resolveProvider(env, layers, envGeoBaseURL, envGeoToken)
	if err != nil {
		return GeoProvider{}, err
	}

	return GeoProvider(provider), nil
}

func resolveOrdersProvider(env map[string]string, layers []CredentialLayer) (OrdersProvider, error) {
	provider, err := resolveProvider(env, layers, envOrdersBaseURL, envOrdersToken)
	if err != nil {
		return OrdersProvider{}, err
	}

	return OrdersProvider(provider), nil
}

type provider struct {
	BaseURL       string `json:"baseUrl,omitempty"`
	BaseURLSource Source `json:"baseUrlSource,omitempty"`
	Token         string `json:"-"`
	TokenSource   Source `json:"tokenSource,omitempty"`
	TokenSet      bool   `json:"tokenSet"`
	Configured    bool   `json:"configured"`
}

func resolveProvider(env map[string]string, layers []CredentialLayer, baseURLKey string, tokenKey string) (provider, error) {
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
