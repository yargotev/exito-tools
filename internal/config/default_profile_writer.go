package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DefaultProfileWriteResult describes the saved Default Profile persistence result.
type DefaultProfileWriteResult struct {
	Profile      string `json:"profile"`
	ConfigPath   string `json:"configPath"`
	ConfigSource Source `json:"configSource"`
}

// SetDefaultProfile persists the saved Default Profile to the selected YAML Configuration File.
func SetDefaultProfile(options Options, profile string) (DefaultProfileWriteResult, error) {
	resolvedOptions, err := normalizeOptions(options)
	if err != nil {
		return DefaultProfileWriteResult{}, err
	}

	profile = strings.TrimSpace(profile)
	if profile == "" {
		return DefaultProfileWriteResult{}, fmt.Errorf("profile must not be blank")
	}
	if strings.ContainsAny(profile, "\r\n") {
		return DefaultProfileWriteResult{}, fmt.Errorf("profile must be a single-line value")
	}

	configPath, configSource, _ := resolveConfigPath(resolvedOptions)
	if configPath == "" {
		configPath = filepath.Join(resolvedOptions.WorkDir, "exito.yaml")
		configSource = SourceLocalProject
	}

	if err := writeDefaultProfile(configPath, profile); err != nil {
		return DefaultProfileWriteResult{}, err
	}

	return DefaultProfileWriteResult{Profile: profile, ConfigPath: configPath, ConfigSource: configSource}, nil
}

func writeDefaultProfile(path string, profile string) error {
	content, err := os.ReadFile(path) // #nosec G304 -- path comes from deterministic configuration path resolution.
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("read configuration file %q: %w", path, err)
	}

	updated := setDefaultProfileLine(string(content), profile)
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("create configuration directory %q: %w", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(updated), 0o600); err != nil { // #nosec G306,G703 -- path comes from deterministic configuration path resolution and is written with owner-only permissions.
		return fmt.Errorf("write configuration file %q: %w", path, err)
	}
	return nil
}

func setDefaultProfileLine(content string, profile string) string {
	line := "defaultProfile: " + profile
	if content == "" {
		return line + "\n"
	}

	endsWithNewline := strings.HasSuffix(content, "\n")
	lines := strings.Split(strings.TrimSuffix(content, "\n"), "\n")
	for index, existing := range lines {
		trimmed := strings.TrimSpace(existing)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		key, _, ok := strings.Cut(trimmed, ":")
		if ok && strings.TrimSpace(key) == "defaultProfile" {
			lines[index] = line
			result := strings.Join(lines, "\n")
			if endsWithNewline {
				result += "\n"
			}
			return result
		}
	}

	result := strings.Join(lines, "\n")
	if result != "" && !strings.HasSuffix(result, "\n") {
		result += "\n"
	}
	result += line + "\n"
	return result
}
