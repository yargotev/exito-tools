package config_test

import (
	"os"
	"path/filepath"
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
