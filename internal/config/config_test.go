package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yargotev/exito-tools/internal/config"
)

func TestResolveProfilePrecedence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		options    config.Options
		want       string
		wantSource config.Source
	}{
		{
			name: "explicit profile wins over environment and saved default",
			options: config.Options{
				Profile:             "prod",
				Env:                 map[string]string{"EXITO_PROFILE": "qa"},
				SavedDefaultProfile: "staging-team",
			},
			want:       "prod",
			wantSource: config.SourceExplicit,
		},
		{
			name: "environment profile wins over saved default",
			options: config.Options{
				Env:                 map[string]string{"EXITO_PROFILE": "qa"},
				SavedDefaultProfile: "staging-team",
			},
			want:       "qa",
			wantSource: config.SourceEnvironment,
		},
		{
			name: "saved default wins over staging fallback",
			options: config.Options{
				Env:                 map[string]string{},
				SavedDefaultProfile: "dev",
			},
			want:       "dev",
			wantSource: config.SourceSavedDefault,
		},
		{
			name: "staging fallback is used without other profile sources",
			options: config.Options{
				Env: map[string]string{},
			},
			want:       "staging",
			wantSource: config.SourceDefault,
		},
		{
			name: "blank values are ignored",
			options: config.Options{
				Profile:             "  ",
				Env:                 map[string]string{"EXITO_PROFILE": "  "},
				SavedDefaultProfile: "  ",
			},
			want:       "staging",
			wantSource: config.SourceDefault,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resolved := resolveForTest(t, tc.options)

			if resolved.Profile != tc.want {
				t.Fatalf("Profile = %q, want %q", resolved.Profile, tc.want)
			}
			if resolved.ProfileSource != tc.wantSource {
				t.Fatalf("ProfileSource = %q, want %q", resolved.ProfileSource, tc.wantSource)
			}
		})
	}
}

func TestYAMLDefaultProfileResolution(t *testing.T) {
	t.Parallel()

	t.Run("local YAML default profile is used as saved default", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), "# Exito Tools\ndefaultProfile: prod\n")

		resolved := resolveForTest(t, config.Options{
			Env:     map[string]string{},
			WorkDir: workDir,
		})

		if resolved.Profile != "prod" {
			t.Fatalf("Profile = %q, want prod", resolved.Profile)
		}
		if resolved.ProfileSource != config.SourceSavedDefault {
			t.Fatalf("ProfileSource = %q, want %q", resolved.ProfileSource, config.SourceSavedDefault)
		}
	})

	t.Run("explicit profile overrides YAML default profile", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), "defaultProfile: prod\n")

		resolved := resolveForTest(t, config.Options{
			Profile: "qa",
			Env:     map[string]string{},
			WorkDir: workDir,
		})

		if resolved.Profile != "qa" {
			t.Fatalf("Profile = %q, want qa", resolved.Profile)
		}
		if resolved.ProfileSource != config.SourceExplicit {
			t.Fatalf("ProfileSource = %q, want %q", resolved.ProfileSource, config.SourceExplicit)
		}
	})

	t.Run("environment profile overrides YAML default profile", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), "defaultProfile: prod\n")

		resolved := resolveForTest(t, config.Options{
			Env:     map[string]string{"EXITO_PROFILE": "dev"},
			WorkDir: workDir,
		})

		if resolved.Profile != "dev" {
			t.Fatalf("Profile = %q, want dev", resolved.Profile)
		}
		if resolved.ProfileSource != config.SourceEnvironment {
			t.Fatalf("ProfileSource = %q, want %q", resolved.ProfileSource, config.SourceEnvironment)
		}
	})

	t.Run("quoted YAML default profile allows inline comments", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), "defaultProfile: 'prod' # team default\n")

		resolved := resolveForTest(t, config.Options{
			Env:     map[string]string{},
			WorkDir: workDir,
		})

		if resolved.Profile != "prod" {
			t.Fatalf("Profile = %q, want prod", resolved.Profile)
		}
	})

	t.Run("blank YAML default profile falls back to staging", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), "defaultProfile:   \n")

		resolved := resolveForTest(t, config.Options{
			Env:     map[string]string{},
			WorkDir: workDir,
		})

		if resolved.Profile != config.DefaultProfile {
			t.Fatalf("Profile = %q, want %q", resolved.Profile, config.DefaultProfile)
		}
		if resolved.ProfileSource != config.SourceDefault {
			t.Fatalf("ProfileSource = %q, want %q", resolved.ProfileSource, config.SourceDefault)
		}
	})
}

func TestResolveConfigPathPrecedence(t *testing.T) {
	t.Parallel()

	t.Run("explicit config path wins over environment and files", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		homeDir := t.TempDir()
		writeFile(t, filepath.Join(workDir, "exito.yaml"))
		writeFile(t, filepath.Join(homeDir, ".config", "exito-tools", "config.yaml"))

		resolved := resolveForTest(t, config.Options{
			ConfigPath: "explicit.yaml",
			Env:        map[string]string{"EXITO_CONFIG": filepath.Join(workDir, "env.yaml")},
			WorkDir:    workDir,
			HomeDir:    homeDir,
		})

		want := filepath.Join(workDir, "explicit.yaml")
		if resolved.ConfigPath != want {
			t.Fatalf("ConfigPath = %q, want %q", resolved.ConfigPath, want)
		}
		if resolved.ConfigSource != config.SourceExplicit {
			t.Fatalf("ConfigSource = %q, want %q", resolved.ConfigSource, config.SourceExplicit)
		}
		if len(resolved.ConfigCandidates) != 0 {
			t.Fatalf("ConfigCandidates length = %d, want 0 for explicit path", len(resolved.ConfigCandidates))
		}
	})

	t.Run("environment config path wins over local and user files", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		homeDir := t.TempDir()
		writeFile(t, filepath.Join(workDir, "exito.yaml"))
		writeFile(t, filepath.Join(homeDir, ".config", "exito-tools", "config.yaml"))

		resolved := resolveForTest(t, config.Options{
			Env:     map[string]string{"EXITO_CONFIG": "env.yaml"},
			WorkDir: workDir,
			HomeDir: homeDir,
		})

		want := filepath.Join(workDir, "env.yaml")
		if resolved.ConfigPath != want {
			t.Fatalf("ConfigPath = %q, want %q", resolved.ConfigPath, want)
		}
		if resolved.ConfigSource != config.SourceEnvironment {
			t.Fatalf("ConfigSource = %q, want %q", resolved.ConfigSource, config.SourceEnvironment)
		}
	})

	t.Run("local project config wins over user config", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		homeDir := t.TempDir()
		localPath := filepath.Join(workDir, "exito.yaml")
		writeFile(t, localPath)
		writeFile(t, filepath.Join(homeDir, ".config", "exito-tools", "config.yaml"))

		resolved := resolveForTest(t, config.Options{
			Env:     map[string]string{},
			WorkDir: workDir,
			HomeDir: homeDir,
		})

		if resolved.ConfigPath != localPath {
			t.Fatalf("ConfigPath = %q, want %q", resolved.ConfigPath, localPath)
		}
		if resolved.ConfigSource != config.SourceLocalProject {
			t.Fatalf("ConfigSource = %q, want %q", resolved.ConfigSource, config.SourceLocalProject)
		}
		assertCandidates(t, resolved.ConfigCandidates, []candidateWant{
			{source: config.SourceLocalProject, path: localPath, exists: true},
			{source: config.SourceUserConfig, path: filepath.Join(homeDir, ".config", "exito-tools", "config.yaml"), exists: true},
		})
	})

	t.Run("user config is selected when local config is absent", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		homeDir := t.TempDir()
		userPath := filepath.Join(homeDir, ".config", "exito-tools", "config.yaml")
		writeFile(t, userPath)

		resolved := resolveForTest(t, config.Options{
			Env:     map[string]string{},
			WorkDir: workDir,
			HomeDir: homeDir,
		})

		if resolved.ConfigPath != userPath {
			t.Fatalf("ConfigPath = %q, want %q", resolved.ConfigPath, userPath)
		}
		if resolved.ConfigSource != config.SourceUserConfig {
			t.Fatalf("ConfigSource = %q, want %q", resolved.ConfigSource, config.SourceUserConfig)
		}
	})

	t.Run("defaults are selected when no config file exists", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		homeDir := t.TempDir()

		resolved := resolveForTest(t, config.Options{
			Env:     map[string]string{},
			WorkDir: workDir,
			HomeDir: homeDir,
		})

		if resolved.ConfigPath != "" {
			t.Fatalf("ConfigPath = %q, want empty", resolved.ConfigPath)
		}
		if resolved.ConfigSource != config.SourceDefault {
			t.Fatalf("ConfigSource = %q, want %q", resolved.ConfigSource, config.SourceDefault)
		}
		assertCandidates(t, resolved.ConfigCandidates, []candidateWant{
			{source: config.SourceLocalProject, path: filepath.Join(workDir, "exito.yaml"), exists: false},
			{source: config.SourceUserConfig, path: filepath.Join(homeDir, ".config", "exito-tools", "config.yaml"), exists: false},
		})
	})
}

func TestResolveCredentialLayers(t *testing.T) {
	t.Parallel()

	workDir := t.TempDir()
	resolved := resolveForTest(t, config.Options{
		Profile: "staging",
		Env:     map[string]string{},
		WorkDir: workDir,
		HomeDir: t.TempDir(),
	})

	want := []config.CredentialLayer{
		{Source: config.SourceEnvironment, Name: "process environment"},
		{Source: config.SourceDotenv, Name: ".env.staging", Path: filepath.Join(workDir, ".env.staging")},
		{Source: config.SourceDotenv, Name: ".env", Path: filepath.Join(workDir, ".env")},
	}

	if len(resolved.CredentialLayers) != len(want) {
		t.Fatalf("CredentialLayers length = %d, want %d", len(resolved.CredentialLayers), len(want))
	}
	for i := range want {
		if resolved.CredentialLayers[i] != want[i] {
			t.Fatalf("CredentialLayers[%d] = %#v, want %#v", i, resolved.CredentialLayers[i], want[i])
		}
	}
}

func TestResolveGeoProviderConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("environment values configure geo provider", func(t *testing.T) {
		t.Parallel()

		resolved := resolveForTest(t, config.Options{
			Env: map[string]string{ // #nosec G101 -- test-only fake tokens exercise environment-source precedence.
				"EXITO_GEO_BASE_URL": " https://geo.example.test ",
				"EXITO_GEO_TOKEN":    " secret-token ",
			},
		})

		if !resolved.GeoProvider.Configured {
			t.Fatalf("GeoProvider.Configured = false, want true")
		}
		if resolved.GeoProvider.BaseURL != "https://geo.example.test" {
			t.Fatalf("GeoProvider.BaseURL = %q, want trimmed base URL", resolved.GeoProvider.BaseURL)
		}
		if resolved.GeoProvider.BaseURLSource != config.SourceEnvironment {
			t.Fatalf("GeoProvider.BaseURLSource = %q, want %q", resolved.GeoProvider.BaseURLSource, config.SourceEnvironment)
		}
		if resolved.GeoProvider.Token != "secret-token" {
			t.Fatalf("GeoProvider.Token = %q, want resolved token", resolved.GeoProvider.Token)
		}
		if resolved.GeoProvider.TokenSource != config.SourceEnvironment {
			t.Fatalf("GeoProvider.TokenSource = %q, want %q", resolved.GeoProvider.TokenSource, config.SourceEnvironment)
		}
		if !resolved.GeoProvider.TokenSet {
			t.Fatalf("GeoProvider.TokenSet = false, want true")
		}
	})

	t.Run("process environment wins over profile dotenv and general dotenv", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, ".env.staging"), "EXITO_GEO_BASE_URL=https://profile.example.test\nEXITO_GEO_TOKEN=profile-token\n")
		writeTextFile(t, filepath.Join(workDir, ".env"), "EXITO_GEO_BASE_URL=https://general.example.test\nEXITO_GEO_TOKEN=general-token\n")

		resolved := resolveForTest(t, config.Options{
			Profile: "staging",
			Env: map[string]string{ // #nosec G101 -- test-only fake token exercises environment-source precedence.
				"EXITO_GEO_TOKEN": "env-token",
			},
			WorkDir: workDir,
		})

		if resolved.GeoProvider.BaseURL != "https://profile.example.test" {
			t.Fatalf("GeoProvider.BaseURL = %q, want profile dotenv value", resolved.GeoProvider.BaseURL)
		}
		if resolved.GeoProvider.BaseURLSource != config.SourceDotenv {
			t.Fatalf("GeoProvider.BaseURLSource = %q, want %q", resolved.GeoProvider.BaseURLSource, config.SourceDotenv)
		}
		if resolved.GeoProvider.Token != "env-token" {
			t.Fatalf("GeoProvider.Token = %q, want environment value", resolved.GeoProvider.Token)
		}
		if resolved.GeoProvider.TokenSource != config.SourceEnvironment {
			t.Fatalf("GeoProvider.TokenSource = %q, want %q", resolved.GeoProvider.TokenSource, config.SourceEnvironment)
		}
		if !resolved.GeoProvider.Configured {
			t.Fatalf("GeoProvider.Configured = false, want true")
		}
	})

	t.Run("general dotenv is used after profile dotenv", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, ".env"), "export EXITO_GEO_BASE_URL='https://general.example.test'\nEXITO_GEO_TOKEN=general-token\n")

		resolved := resolveForTest(t, config.Options{
			Profile: "staging",
			Env:     map[string]string{},
			WorkDir: workDir,
		})

		if resolved.GeoProvider.BaseURL != "https://general.example.test" {
			t.Fatalf("GeoProvider.BaseURL = %q, want general dotenv value", resolved.GeoProvider.BaseURL)
		}
		if resolved.GeoProvider.Token != "general-token" {
			t.Fatalf("GeoProvider.Token = %q, want general dotenv token", resolved.GeoProvider.Token)
		}
		if !resolved.GeoProvider.Configured {
			t.Fatalf("GeoProvider.Configured = false, want true")
		}
	})

	t.Run("missing token remains unconfigured without exposing a token", func(t *testing.T) {
		t.Parallel()

		resolved := resolveForTest(t, config.Options{
			Env: map[string]string{"EXITO_GEO_BASE_URL": "https://geo.example.test"},
		})

		if resolved.GeoProvider.Configured {
			t.Fatalf("GeoProvider.Configured = true, want false")
		}
		if resolved.GeoProvider.TokenSet {
			t.Fatalf("GeoProvider.TokenSet = true, want false")
		}
		if resolved.GeoProvider.Token != "" {
			t.Fatalf("GeoProvider.Token = %q, want empty", resolved.GeoProvider.Token)
		}
		if resolved.GeoProvider.TokenSource != config.SourceDefault {
			t.Fatalf("GeoProvider.TokenSource = %q, want %q", resolved.GeoProvider.TokenSource, config.SourceDefault)
		}
	})
}

func TestResolveProviderBaseURLsFromYAMLProfiles(t *testing.T) {
	t.Parallel()

	t.Run("YAML profile base URLs configure providers with environment tokens", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `
defaultProfile: staging
profiles:
  staging:
    geo:
      baseUrl: https://geo-yaml.example.test
    orders:
      baseUrl: https://orders-yaml.example.test
`)

		resolved := resolveForTest(t, config.Options{
			Env: map[string]string{ // #nosec G101 -- test-only fake tokens exercise environment-source precedence.
				"EXITO_GEO_TOKEN":    "geo-env-token",
				"EXITO_ORDERS_TOKEN": "orders-env-token",
			},
			WorkDir: workDir,
		})

		if !resolved.GeoProvider.Configured || !resolved.OrdersProvider.Configured {
			t.Fatalf("providers configured = geo:%v orders:%v, want both true", resolved.GeoProvider.Configured, resolved.OrdersProvider.Configured)
		}
		if resolved.GeoProvider.BaseURL != "https://geo-yaml.example.test" || resolved.GeoProvider.BaseURLSource != config.SourceConfigFile {
			t.Fatalf("GeoProvider base = (%q,%q), want YAML config-file", resolved.GeoProvider.BaseURL, resolved.GeoProvider.BaseURLSource)
		}
		if resolved.OrdersProvider.BaseURL != "https://orders-yaml.example.test" || resolved.OrdersProvider.BaseURLSource != config.SourceConfigFile {
			t.Fatalf("OrdersProvider base = (%q,%q), want YAML config-file", resolved.OrdersProvider.BaseURL, resolved.OrdersProvider.BaseURLSource)
		}
		if resolved.GeoProvider.Token != "geo-env-token" || resolved.OrdersProvider.Token != "orders-env-token" {
			t.Fatalf("tokens should come from environment, got geo=%q orders=%q", resolved.GeoProvider.Token, resolved.OrdersProvider.Token)
		}
	})

	t.Run("environment base URL overrides YAML profile base URL", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `
profiles:
  staging:
    geo:
      baseUrl: https://geo-yaml.example.test
`)

		resolved := resolveForTest(t, config.Options{
			Env: map[string]string{ // #nosec G101 -- test-only fake token exercises environment-source precedence.
				"EXITO_GEO_BASE_URL": "https://geo-env.example.test",
				"EXITO_GEO_TOKEN":    "geo-env-token",
			},
			WorkDir: workDir,
		})

		if resolved.GeoProvider.BaseURL != "https://geo-env.example.test" {
			t.Fatalf("GeoProvider.BaseURL = %q, want environment value", resolved.GeoProvider.BaseURL)
		}
		if resolved.GeoProvider.BaseURLSource != config.SourceEnvironment {
			t.Fatalf("GeoProvider.BaseURLSource = %q, want %q", resolved.GeoProvider.BaseURLSource, config.SourceEnvironment)
		}
	})

	t.Run("effective profile selects matching YAML profile", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `
profiles:
  staging:
    geo:
      baseUrl: https://geo-staging.example.test
  prod:
    geo:
      baseUrl: https://geo-prod.example.test
    orders:
      baseUrl: https://orders-prod.example.test
`)

		resolved := resolveForTest(t, config.Options{
			Profile: "prod",
			Env: map[string]string{ // #nosec G101 -- test-only fake tokens exercise environment-source precedence.
				"EXITO_GEO_TOKEN":    "geo-env-token",
				"EXITO_ORDERS_TOKEN": "orders-env-token",
			},
			WorkDir: workDir,
		})

		if resolved.GeoProvider.BaseURL != "https://geo-prod.example.test" {
			t.Fatalf("GeoProvider.BaseURL = %q, want prod YAML profile", resolved.GeoProvider.BaseURL)
		}
		if resolved.OrdersProvider.BaseURL != "https://orders-prod.example.test" {
			t.Fatalf("OrdersProvider.BaseURL = %q, want prod YAML profile", resolved.OrdersProvider.BaseURL)
		}
	})

	t.Run("YAML token-like keys are ignored", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `
profiles:
  staging:
    geo:
      baseUrl: https://geo-yaml.example.test
      token: must-not-be-read
    orders:
      baseUrl: https://orders-yaml.example.test
      token: must-not-be-read
`)

		resolved := resolveForTest(t, config.Options{
			Env:     map[string]string{},
			WorkDir: workDir,
		})

		if resolved.GeoProvider.Configured || resolved.OrdersProvider.Configured {
			t.Fatalf("providers should remain unconfigured without environment/dotenv tokens")
		}
		if resolved.GeoProvider.Token != "" || resolved.OrdersProvider.Token != "" {
			t.Fatalf("YAML token-like values must be ignored, got geo=%q orders=%q", resolved.GeoProvider.Token, resolved.OrdersProvider.Token)
		}
	})
}

func TestGeoProviderTokenIsOmittedFromEffectiveJSON(t *testing.T) {
	t.Parallel()

	resolved := resolveForTest(t, config.Options{
		Env: map[string]string{
			"EXITO_GEO_BASE_URL": "https://geo.example.test",
			"EXITO_GEO_TOKEN":    "super-secret-token",
		},
	})

	encoded, err := json.Marshal(resolved)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if strings.Contains(string(encoded), "super-secret-token") || strings.Contains(string(encoded), `"Token":`) {
		t.Fatalf("encoded effective config exposed token: %s", string(encoded))
	}
	if !strings.Contains(string(encoded), `"tokenSet":true`) {
		t.Fatalf("encoded effective config should expose only token presence: %s", string(encoded))
	}
}

func TestResolveOrdersProviderConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("environment values configure orders provider", func(t *testing.T) {
		t.Parallel()

		resolved := resolveForTest(t, config.Options{
			Env: map[string]string{
				"EXITO_ORDERS_BASE_URL": " https://orders.example.test ",
				"EXITO_ORDERS_TOKEN":    " provider-value ",
			},
		})

		if !resolved.OrdersProvider.Configured {
			t.Fatalf("OrdersProvider.Configured = false, want true")
		}
		if resolved.OrdersProvider.BaseURL != "https://orders.example.test" {
			t.Fatalf("OrdersProvider.BaseURL = %q, want trimmed base URL", resolved.OrdersProvider.BaseURL)
		}
		if resolved.OrdersProvider.BaseURLSource != config.SourceEnvironment {
			t.Fatalf("OrdersProvider.BaseURLSource = %q, want %q", resolved.OrdersProvider.BaseURLSource, config.SourceEnvironment)
		}
		if resolved.OrdersProvider.Token != "provider-value" {
			t.Fatalf("OrdersProvider.Token = %q, want resolved token", resolved.OrdersProvider.Token)
		}
		if resolved.OrdersProvider.TokenSource != config.SourceEnvironment {
			t.Fatalf("OrdersProvider.TokenSource = %q, want %q", resolved.OrdersProvider.TokenSource, config.SourceEnvironment)
		}
		if !resolved.OrdersProvider.TokenSet {
			t.Fatalf("OrdersProvider.TokenSet = false, want true")
		}
	})

	t.Run("geoms credentials configure client credentials", func(t *testing.T) {
		t.Parallel()

		resolved := resolveForTest(t, config.Options{
			Profile: "staging",
			Env: map[string]string{ // #nosec G101 -- test-only fake GEOMS credential bundle.
				"EXITO_ORDERS_BASE_URL": "https://geoms.example.test/apioms/api/v1/scope/geoms",
				"GEOMS_CREDENTIALS_QA":  `{ 'client_id': 'client-id', 'client_secret': 'secret-value', 'grant_type': 'client_credentials', 'scope': 'scope-value' }`,
			},
		})

		if !resolved.OrdersProvider.Configured {
			t.Fatalf("OrdersProvider.Configured = false, want true")
		}
		if resolved.OrdersProvider.ClientID != "client-id" || resolved.OrdersProvider.ClientSecret != "secret-value" || resolved.OrdersProvider.Scope != "scope-value" {
			t.Fatalf("GEOMS credentials = (%q,%q,%q), want parsed values", resolved.OrdersProvider.ClientID, resolved.OrdersProvider.ClientSecret, resolved.OrdersProvider.Scope)
		}
		if !resolved.OrdersProvider.ClientIDSet || !resolved.OrdersProvider.SecretSet || !resolved.OrdersProvider.ScopeSet {
			t.Fatalf("credential presence flags should be true: %#v", resolved.OrdersProvider)
		}
	})

	t.Run("prod profile reads PDN GEOMS credential bundle", func(t *testing.T) {
		t.Parallel()

		resolved := resolveForTest(t, config.Options{
			Profile: "prod",
			Env: map[string]string{ // #nosec G101 -- test-only fake GEOMS credential bundle.
				"EXITO_ORDERS_BASE_URL": "https://geoms.example.test/apioms/api/v1/scope/geoms",
				"GEOMS_CREDENTIALS_QA":  `{ 'client_id': 'qa-client-id', 'client_secret': 'qa-secret', 'scope': 'qa-scope' }`,
				"GEOMS_CREDENTIALS_PDN": `{ 'client_id': 'pdn-client-id', 'client_secret': 'pdn-secret', 'scope': 'pdn-scope' }`,
			},
		})

		if !resolved.OrdersProvider.Configured {
			t.Fatalf("OrdersProvider.Configured = false, want true")
		}
		if resolved.OrdersProvider.ClientID != "pdn-client-id" || resolved.OrdersProvider.ClientSecret != "pdn-secret" || resolved.OrdersProvider.Scope != "pdn-scope" {
			t.Fatalf("GEOMS PDN credentials = (%q,%q,%q), want PDN values", resolved.OrdersProvider.ClientID, resolved.OrdersProvider.ClientSecret, resolved.OrdersProvider.Scope)
		}
	})

	t.Run("process environment wins over profile dotenv and general dotenv", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, ".env.staging"), "EXITO_ORDERS_BASE_URL=https://profile-orders.example.test\nEXITO_ORDERS_TOKEN=profile-value\n")
		writeTextFile(t, filepath.Join(workDir, ".env"), "EXITO_ORDERS_BASE_URL=https://general-orders.example.test\nEXITO_ORDERS_TOKEN=general-value\n")

		resolved := resolveForTest(t, config.Options{
			Profile: "staging",
			Env: map[string]string{
				"EXITO_ORDERS_TOKEN": "env-value",
			},
			WorkDir: workDir,
		})

		if resolved.OrdersProvider.BaseURL != "https://profile-orders.example.test" {
			t.Fatalf("OrdersProvider.BaseURL = %q, want profile dotenv value", resolved.OrdersProvider.BaseURL)
		}
		if resolved.OrdersProvider.BaseURLSource != config.SourceDotenv {
			t.Fatalf("OrdersProvider.BaseURLSource = %q, want %q", resolved.OrdersProvider.BaseURLSource, config.SourceDotenv)
		}
		if resolved.OrdersProvider.Token != "env-value" {
			t.Fatalf("OrdersProvider.Token = %q, want environment value", resolved.OrdersProvider.Token)
		}
		if resolved.OrdersProvider.TokenSource != config.SourceEnvironment {
			t.Fatalf("OrdersProvider.TokenSource = %q, want %q", resolved.OrdersProvider.TokenSource, config.SourceEnvironment)
		}
		if !resolved.OrdersProvider.Configured {
			t.Fatalf("OrdersProvider.Configured = false, want true")
		}
	})

	t.Run("missing token remains unconfigured without exposing a token", func(t *testing.T) {
		t.Parallel()

		resolved := resolveForTest(t, config.Options{
			Env: map[string]string{"EXITO_ORDERS_BASE_URL": "https://orders.example.test"},
		})

		if resolved.OrdersProvider.Configured {
			t.Fatalf("OrdersProvider.Configured = true, want false")
		}
		if resolved.OrdersProvider.TokenSet {
			t.Fatalf("OrdersProvider.TokenSet = true, want false")
		}
		if resolved.OrdersProvider.Token != "" {
			t.Fatalf("OrdersProvider.Token = %q, want empty", resolved.OrdersProvider.Token)
		}
		if resolved.OrdersProvider.TokenSource != config.SourceDefault {
			t.Fatalf("OrdersProvider.TokenSource = %q, want %q", resolved.OrdersProvider.TokenSource, config.SourceDefault)
		}
	})
}

func TestOrdersProviderTokenIsOmittedFromEffectiveJSON(t *testing.T) {
	t.Parallel()

	resolved := resolveForTest(t, config.Options{
		Env: map[string]string{
			"EXITO_ORDERS_BASE_URL": "https://orders.example.test",
			"EXITO_ORDERS_TOKEN":    "super-secret-provider-value",
		},
	})

	encoded, err := json.Marshal(resolved)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if strings.Contains(string(encoded), "super-secret-provider-value") || strings.Contains(string(encoded), `"Token":`) {
		t.Fatalf("encoded effective config exposed token: %s", string(encoded))
	}
	if !strings.Contains(string(encoded), `"OrdersProvider"`) && !strings.Contains(string(encoded), `"ordersProvider"`) {
		t.Fatalf("encoded effective config should include orders provider metadata: %s", string(encoded))
	}
}

type candidateWant struct {
	source config.Source
	path   string
	exists bool
}

func assertCandidates(t *testing.T, got []config.ConfigCandidate, want []candidateWant) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("ConfigCandidates length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Source != want[i].source || got[i].Path != want[i].path || got[i].Exists != want[i].exists {
			t.Fatalf("ConfigCandidates[%d] = %#v, want source=%q path=%q exists=%v", i, got[i], want[i].source, want[i].path, want[i].exists)
		}
	}
}

func resolveForTest(t *testing.T, options config.Options) config.Effective {
	t.Helper()

	if options.WorkDir == "" {
		options.WorkDir = t.TempDir()
	}
	if options.HomeDir == "" {
		options.HomeDir = t.TempDir()
	}

	resolved, err := config.Resolve(options)
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	return resolved
}

func writeFile(t *testing.T, path string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte("placeholder: true\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func writeTextFile(t *testing.T, path string, content string) {
	t.Helper()

	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		t.Fatalf("MkdirAll(%q) error = %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%q) error = %v", path, err)
	}
}

func TestSetDefaultProfile(t *testing.T) {
	t.Parallel()

	t.Run("updates existing selected local YAML default profile", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		path := filepath.Join(workDir, "exito.yaml")
		writeTextFile(t, path, "# team config\ndefaultProfile: staging\n")

		result, err := config.SetDefaultProfile(config.Options{Env: map[string]string{}, WorkDir: workDir, HomeDir: t.TempDir()}, "prod")
		if err != nil {
			t.Fatalf("SetDefaultProfile() error = %v", err)
		}

		if result.Profile != "prod" || result.ConfigPath != path || result.ConfigSource != config.SourceLocalProject {
			t.Fatalf("result = %#v, want prod local path", result)
		}
		content := readTextFile(t, path)
		if !strings.Contains(content, "defaultProfile: prod") {
			t.Fatalf("updated file missing default profile:\n%s", content)
		}
		if strings.Contains(content, "EXITO_GEO_TOKEN") || strings.Contains(content, "EXITO_ORDERS_TOKEN") {
			t.Fatalf("updated file should not write credential keys:\n%s", content)
		}
	})

	t.Run("appends missing default profile", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		path := filepath.Join(workDir, "exito.yaml")
		writeTextFile(t, path, "profiles:\n  staging: {}\n")

		_, err := config.SetDefaultProfile(config.Options{Env: map[string]string{}, WorkDir: workDir, HomeDir: t.TempDir()}, "qa")
		if err != nil {
			t.Fatalf("SetDefaultProfile() error = %v", err)
		}

		content := readTextFile(t, path)
		if !strings.Contains(content, "profiles:\n  staging: {}\ndefaultProfile: qa\n") {
			t.Fatalf("updated file did not append default profile as expected:\n%s", content)
		}
	})

	t.Run("creates local config when no file exists", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		path := filepath.Join(workDir, "exito.yaml")

		result, err := config.SetDefaultProfile(config.Options{Env: map[string]string{}, WorkDir: workDir, HomeDir: t.TempDir()}, "dev")
		if err != nil {
			t.Fatalf("SetDefaultProfile() error = %v", err)
		}

		if result.ConfigPath != path || result.ConfigSource != config.SourceLocalProject {
			t.Fatalf("result = %#v, want created local config", result)
		}
		if content := readTextFile(t, path); content != "defaultProfile: dev\n" {
			t.Fatalf("created file = %q, want default profile only", content)
		}
	})

	t.Run("blank profile is rejected before writing", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		_, err := config.SetDefaultProfile(config.Options{Env: map[string]string{}, WorkDir: workDir, HomeDir: t.TempDir()}, "   ")
		if err == nil {
			t.Fatalf("SetDefaultProfile() error = nil, want validation error")
		}
		if _, statErr := os.Stat(filepath.Join(workDir, "exito.yaml")); !os.IsNotExist(statErr) {
			t.Fatalf("exito.yaml should not be written, stat error = %v", statErr)
		}
	})
}

func readTextFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path) // #nosec G304 -- tests read their own temporary files.
	if err != nil {
		t.Fatalf("ReadFile(%q) error = %v", path, err)
	}
	return string(content)
}

func TestResolveVTEXOMSProviderConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("staging resolves Exito QA credentials and YAML base URL", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `defaultProfile: staging
profiles:
  staging:
    vtexOms:
      exito:
        baseUrl: https://master--exito.myvtex.com
`)

		resolved, err := config.Resolve(config.Options{
			WorkDir: workDir,
			HomeDir: t.TempDir(),
			Env: map[string]string{ // #nosec G101 -- test-only fake VTEX credentials.
				"EXITO_APP_KEY_QA":   "qa-key",
				"EXITO_APP_TOKEN_QA": "qa-token",
			},
		})
		if err != nil {
			t.Fatalf("Resolve() error = %v", err)
		}

		provider := resolved.VTEXOMSProvider.Exito
		if !provider.Configured {
			t.Fatalf("Exito VTEX OMS provider should be configured: %#v", provider)
		}
		if provider.BaseURL != "https://master--exito.myvtex.com" || provider.BaseURLSource != config.SourceConfigFile {
			t.Fatalf("base URL = (%q,%q), want YAML config-file", provider.BaseURL, provider.BaseURLSource)
		}
		if provider.AppKey != "qa-key" || provider.AppToken != "qa-token" {
			t.Fatalf("credentials = (%q,%q), want env values", provider.AppKey, provider.AppToken)
		}
	})

	t.Run("prod resolves Carulla production credential names", func(t *testing.T) {
		t.Parallel()

		resolved, err := config.Resolve(config.Options{
			Profile: "prod",
			WorkDir: t.TempDir(),
			HomeDir: t.TempDir(),
			Env: map[string]string{ // #nosec G101 -- test-only fake VTEX credentials.
				"CARULLA_VTEX_OMS_BASE_URL_PROD": "https://carulla.myvtex.com",
				"CARULLA_APP_KEY":                "prod-key",
				"CARULLA_APP_TOKEN":              "prod-token",
			},
		})
		if err != nil {
			t.Fatalf("Resolve() error = %v", err)
		}

		provider := resolved.VTEXOMSProvider.Carulla
		if !provider.Configured || provider.BaseURL != "https://carulla.myvtex.com" || provider.AppKey != "prod-key" || provider.AppToken != "prod-token" {
			t.Fatalf("Carulla VTEX OMS provider = %#v, want configured prod values", provider)
		}
	})
}

func TestVTEXOMSCredentialsAreOmittedFromEffectiveJSON(t *testing.T) {
	t.Parallel()

	resolved, err := config.Resolve(config.Options{
		WorkDir: t.TempDir(),
		HomeDir: t.TempDir(),
		Env: map[string]string{ // #nosec G101 -- test-only fake VTEX credentials.
			"EXITO_VTEX_OMS_BASE_URL_QA": "https://master--exito.myvtex.com",
			"EXITO_APP_KEY_QA":           "secret-key",
			"EXITO_APP_TOKEN_QA":         "secret-token",
		},
	})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	encoded, err := json.Marshal(resolved)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if strings.Contains(string(encoded), "secret-key") || strings.Contains(string(encoded), "secret-token") {
		t.Fatalf("effective config JSON leaked VTEX credentials: %s", string(encoded))
	}
	if !strings.Contains(string(encoded), `"appKeySet":true`) || !strings.Contains(string(encoded), `"appTokenSet":true`) {
		t.Fatalf("effective config should expose only VTEX credential presence metadata: %s", string(encoded))
	}
}

func TestResolveVTEXIntelligentSearchProviderConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("staging resolves Exito intelligent search YAML base URL", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `defaultProfile: staging
profiles:
  staging:
    vtexIntelligentSearch:
      exito:
        baseUrl: https://exito.vtexcommercestable.com.br
`)

		resolved := resolveForTest(t, config.Options{WorkDir: workDir, Env: map[string]string{}})
		provider := resolved.VTEXIntelligentSearchProvider.Exito
		if !provider.Configured {
			t.Fatalf("Exito VTEX Intelligent Search provider should be configured: %#v", provider)
		}
		if provider.BaseURL != "https://exito.vtexcommercestable.com.br" || provider.BaseURLSource != config.SourceConfigFile {
			t.Fatalf("base URL = (%q,%q), want YAML config-file", provider.BaseURL, provider.BaseURLSource)
		}
	})

	t.Run("prod resolves Carulla intelligent search environment base URL", func(t *testing.T) {
		t.Parallel()

		resolved := resolveForTest(t, config.Options{
			Profile: "prod",
			Env: map[string]string{
				"CARULLA_VTEX_INTELLIGENT_SEARCH_BASE_URL_PROD": "https://carulla.vtexcommercestable.com.br",
			},
		})
		provider := resolved.VTEXIntelligentSearchProvider.Carulla
		if !provider.Configured || provider.BaseURL != "https://carulla.vtexcommercestable.com.br" || provider.BaseURLSource != config.SourceEnvironment {
			t.Fatalf("Carulla VTEX Intelligent Search provider = %#v, want configured prod account/environment value", provider)
		}
	})
}

func TestResolveVTEXCatalogProviderConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("staging resolves Exito catalog YAML base URL", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `defaultProfile: staging
profiles:
  staging:
    vtexCatalog:
      exito:
        baseUrl: https://exito.vtexcommercestable.com.br
`)

		resolved := resolveForTest(t, config.Options{WorkDir: workDir, Env: map[string]string{}})
		provider := resolved.VTEXCatalogProvider.Exito
		if !provider.Configured {
			t.Fatalf("Exito VTEX Catalog provider should be configured: %#v", provider)
		}
		if provider.BaseURL != "https://exito.vtexcommercestable.com.br" || provider.BaseURLSource != config.SourceConfigFile {
			t.Fatalf("base URL = (%q,%q), want YAML config-file", provider.BaseURL, provider.BaseURLSource)
		}
	})

	t.Run("prod resolves Carulla catalog environment base URL", func(t *testing.T) {
		t.Parallel()

		resolved := resolveForTest(t, config.Options{
			Profile: "prod",
			Env: map[string]string{
				"CARULLA_VTEX_CATALOG_BASE_URL_PROD": "https://www.carulla.com",
			},
		})
		provider := resolved.VTEXCatalogProvider.Carulla
		if !provider.Configured || provider.BaseURL != "https://www.carulla.com" || provider.BaseURLSource != config.SourceEnvironment {
			t.Fatalf("Carulla VTEX Catalog provider = %#v, want configured prod environment value", provider)
		}
	})
}

func TestResolveVTEXCheckoutProvider(t *testing.T) {
	t.Parallel()

	t.Run("YAML profile configures checkout brand base URLs", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `defaultProfile: staging
profiles:
  staging:
    vtexCheckout:
      exito:
        baseUrl: https://checkout-exito.example.test
      carulla:
        baseUrl: https://checkout-carulla.example.test
`)

		resolved := resolveForTest(t, config.Options{Env: map[string]string{}, WorkDir: workDir})
		if !resolved.VTEXCheckoutProvider.Exito.Configured || resolved.VTEXCheckoutProvider.Exito.BaseURL != "https://checkout-exito.example.test" {
			t.Fatalf("exito checkout provider = %#v, want YAML configured", resolved.VTEXCheckoutProvider.Exito)
		}
		if resolved.VTEXCheckoutProvider.Carulla.BaseURL != "https://checkout-carulla.example.test" {
			t.Fatalf("carulla checkout provider = %#v, want YAML base URL", resolved.VTEXCheckoutProvider.Carulla)
		}
	})

	t.Run("environment overrides YAML checkout base URL by profile", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `profiles:
  staging:
    vtexCheckout:
      exito:
        baseUrl: https://yaml.example.test
`)

		resolved := resolveForTest(t, config.Options{Env: map[string]string{"EXITO_VTEX_CHECKOUT_BASE_URL_QA": "https://env.example.test"}, WorkDir: workDir})
		if resolved.VTEXCheckoutProvider.Exito.BaseURL != "https://env.example.test" {
			t.Fatalf("checkout base URL = %q, want environment override", resolved.VTEXCheckoutProvider.Exito.BaseURL)
		}
		if resolved.VTEXCheckoutProvider.Exito.BaseURLSource != config.SourceEnvironment {
			t.Fatalf("checkout source = %q, want environment", resolved.VTEXCheckoutProvider.Exito.BaseURLSource)
		}
	})
}

func TestResolveVTEXMasterDataProviderConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("staging resolves Exito YAML base URL with QA credentials", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `defaultProfile: staging
profiles:
  staging:
    vtexMasterData:
      exito:
        baseUrl: https://exito.vtexcommercestable.com.br
`)

		resolved := resolveForTest(t, config.Options{
			WorkDir: workDir,
			Env: map[string]string{ // #nosec G101 -- test-only fake VTEX credentials.
				"EXITO_APP_KEY_QA":   "md-qa-key",
				"EXITO_APP_TOKEN_QA": "md-qa-token",
			},
		})

		provider := resolved.VTEXMasterDataProvider.Exito
		if !provider.Configured {
			t.Fatalf("Exito Master Data provider should be configured: %#v", provider)
		}
		if provider.BaseURL != "https://exito.vtexcommercestable.com.br" || provider.BaseURLSource != config.SourceConfigFile {
			t.Fatalf("base URL = (%q,%q), want YAML config-file", provider.BaseURL, provider.BaseURLSource)
		}
		if provider.AppKey != "md-qa-key" || provider.AppToken != "md-qa-token" {
			t.Fatalf("credentials = (%q,%q), want QA env credentials", provider.AppKey, provider.AppToken)
		}
	})

	t.Run("prod environment overrides YAML Carulla base URL", func(t *testing.T) {
		t.Parallel()

		workDir := t.TempDir()
		writeTextFile(t, filepath.Join(workDir, "exito.yaml"), `profiles:
  prod:
    vtexMasterData:
      carulla:
        baseUrl: https://yaml-carulla.example.test
`)

		resolved := resolveForTest(t, config.Options{
			Profile: "prod",
			WorkDir: workDir,
			Env: map[string]string{ // #nosec G101 -- test-only fake VTEX credentials.
				"CARULLA_VTEX_MASTERDATA_BASE_URL_PROD": "https://env-carulla.example.test",
				"CARULLA_APP_KEY":                       "md-prod-key",
				"CARULLA_APP_TOKEN":                     "md-prod-token",
			},
		})

		provider := resolved.VTEXMasterDataProvider.Carulla
		if !provider.Configured || provider.BaseURL != "https://env-carulla.example.test" || provider.BaseURLSource != config.SourceEnvironment {
			t.Fatalf("Carulla Master Data provider = %#v, want env-configured provider", provider)
		}
	})

	t.Run("missing credentials leave brand unconfigured", func(t *testing.T) {
		t.Parallel()

		resolved := resolveForTest(t, config.Options{Env: map[string]string{"EXITO_VTEX_MASTERDATA_BASE_URL_QA": "https://master.example.test"}})
		provider := resolved.VTEXMasterDataProvider.Exito
		if provider.Configured || provider.BaseURL == "" || provider.AppKeySet || provider.AppTokenSet {
			t.Fatalf("Exito Master Data provider = %#v, want base URL present but unconfigured without credentials", provider)
		}
	})
}

func TestVTEXMasterDataCredentialsAreOmittedFromEffectiveJSON(t *testing.T) {
	t.Parallel()

	resolved := resolveForTest(t, config.Options{Env: map[string]string{ // #nosec G101 -- test-only fake VTEX credentials.
		"EXITO_VTEX_MASTERDATA_BASE_URL_QA": "https://exito.vtexcommercestable.com.br",
		"EXITO_APP_KEY_QA":                  "md-secret-key",
		"EXITO_APP_TOKEN_QA":                "md-secret-token",
	}})

	encoded, err := json.Marshal(resolved)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if strings.Contains(string(encoded), "md-secret-key") || strings.Contains(string(encoded), "md-secret-token") {
		t.Fatalf("effective config JSON leaked Master Data credentials: %s", string(encoded))
	}
	if !strings.Contains(string(encoded), `"vtexMasterDataProvider"`) || !strings.Contains(string(encoded), `"appKeySet":true`) || !strings.Contains(string(encoded), `"appTokenSet":true`) {
		t.Fatalf("effective config should expose Master Data presence metadata: %s", string(encoded))
	}
}
