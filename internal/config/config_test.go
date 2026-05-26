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
			Env: map[string]string{
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
			Env: map[string]string{
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
