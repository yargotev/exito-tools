package config

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	// DefaultProfile is the profile used when no explicit or saved profile exists.
	DefaultProfile = "staging"

	envConfigPath = "EXITO_CONFIG"
	envProfile    = "EXITO_PROFILE"
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
}

// ConfigCandidate describes a configuration file candidate considered by the resolver.
type ConfigCandidate struct {
	Source Source
	Path   string
	Exists bool
}

// CredentialLayer describes a possible source of sensitive values without reading them.
type CredentialLayer struct {
	Source Source
	Name   string
	Path   string
}

// Resolve applies Exito Tools configuration precedence rules without parsing YAML or secrets.
func Resolve(options Options) (Effective, error) {
	resolvedOptions, err := normalizeOptions(options)
	if err != nil {
		return Effective{}, err
	}

	profile, profileSource := resolveProfile(resolvedOptions)
	configPath, configSource, candidates := resolveConfigPath(resolvedOptions)

	return Effective{
		Profile:          profile,
		ProfileSource:    profileSource,
		ConfigPath:       configPath,
		ConfigSource:     configSource,
		ConfigCandidates: candidates,
		CredentialLayers: credentialLayers(resolvedOptions.WorkDir, profile),
	}, nil
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
