package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/capability"
	"github.com/yargotev/exito-tools/internal/config"
	"github.com/yargotev/exito-tools/internal/execution"
	"github.com/yargotev/exito-tools/internal/presenter"
)

// Bootstrapper builds the application after Cobra has parsed CLI boot flags.
type Bootstrapper func(app.Options) (*app.Application, error)

type rootOptions struct {
	configPath    string
	profile       string
	correlationID string
}

type capabilitiesData struct {
	Capabilities []capability.Definition `json:"capabilities"`
}

// NewRoot builds the minimal English-only CLI root surface.
func NewRoot(bootstrap Bootstrapper) *cobra.Command {
	if bootstrap == nil {
		bootstrap = app.New
	}

	options := rootOptions{}
	command := &cobra.Command{
		Use:   "exito",
		Short: "Exito Tools command-line interface",
		Long:  rootLong(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			application, err := bootstrap(appOptions(options))
			if err != nil {
				return err
			}

			cmd.Long = rootLong(len(application.Registry.All()))
			return cmd.Help()
		},
		SilenceErrors:     true,
		SilenceUsage:      true,
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}

	command.PersistentFlags().StringVar(&options.configPath, "config", "", "Path to the Exito Tools configuration file")
	command.PersistentFlags().StringVar(&options.profile, "profile", "", "Configuration profile to use")
	command.PersistentFlags().StringVar(&options.correlationID, "correlation-id", "", "Correlation ID to include in JSON command metadata")
	command.SetHelpCommand(&cobra.Command{Hidden: true})
	command.AddCommand(newCapabilitiesCommand(bootstrap, &options))
	command.AddCommand(newRunCommand(bootstrap, &options))
	return command
}

func newCapabilitiesCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "capabilities",
		Short: "Print the machine-readable capability inventory",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			startedAt := time.Now()
			requestID, err := execution.NewRequestID()
			if err != nil {
				return err
			}

			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			data := capabilitiesData{Capabilities: application.Registry.All()}
			metadata := execution.NewMetadata(requestID, options.correlationID, startedAt, time.Now())
			envelope := capability.Envelope[capabilitiesData]{
				OK:   true,
				Data: &data,
				Meta: metadata.EnvelopeMeta(application.Config.Profile, ""),
			}

			return presenter.WriteJSON(cmd.OutOrStdout(), envelope)
		},
	}
}

func newRunCommand(bootstrap Bootstrapper, options *rootOptions) *cobra.Command {
	var inputJSON string
	var inputFile string

	command := &cobra.Command{
		Use:   "run <capability-id>",
		Short: "Run a capability by its stable ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input, err := parseRunInput(cmd, inputJSON, inputFile)
			if err != nil {
				return err
			}

			application, err := bootstrap(appOptions(*options))
			if err != nil {
				return err
			}

			pipeline := execution.NewPipeline(application.Registry)
			envelope, err := pipeline.Execute(cmd.Context(), execution.ExecuteRequest{
				CapabilityID:  args[0],
				Input:         input,
				Profile:       application.Config.Profile,
				CorrelationID: options.correlationID,
			})
			if err != nil {
				return err
			}

			return presenter.WriteJSON(cmd.OutOrStdout(), envelope)
		},
	}

	command.Flags().StringVar(&inputJSON, "input-json", "", "Complete capability input object as inline JSON")
	command.Flags().StringVar(&inputFile, "input-file", "", "Path to a JSON file containing the complete capability input object")
	return command
}

func parseRunInput(cmd *cobra.Command, inputJSON string, inputFile string) (capability.Input, error) {
	sources := 0
	if inputJSON != "" {
		sources++
	}
	if inputFile != "" {
		sources++
	}
	stdinAvailable := runStdinAvailable(cmd)
	if stdinAvailable {
		sources++
	}
	if sources > 1 {
		return nil, fmt.Errorf("run input must be provided by only one source")
	}

	switch {
	case inputJSON != "":
		return decodeRunInput([]byte(inputJSON))
	case inputFile != "":
		content, err := os.ReadFile(inputFile) // #nosec G304 -- users explicitly choose the generic run input file
		if err != nil {
			return nil, err
		}
		return decodeRunInput(content)
	case stdinAvailable:
		content, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return nil, err
		}
		return decodeRunInput(content)
	default:
		return capability.Input{}, nil
	}
}

func decodeRunInput(content []byte) (capability.Input, error) {
	var input capability.Input
	if err := json.Unmarshal(content, &input); err != nil {
		return nil, err
	}
	if input == nil {
		return nil, fmt.Errorf("run input must be a JSON object")
	}
	return input, nil
}

func runStdinAvailable(cmd *cobra.Command) bool {
	input := cmd.InOrStdin()
	if input == nil {
		return false
	}
	if input != os.Stdin {
		return true
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return stat.Mode()&os.ModeCharDevice == 0
}

func appOptions(options rootOptions) app.Options {
	return app.Options{
		Config: config.Options{
			ConfigPath: options.configPath,
			Profile:    options.profile,
		},
	}
}

func rootLong(registeredEntries int) string {
	return fmt.Sprintf(
		"Exito Tools command-line interface\n\nExito Tools is the machine-first CLI surface for the application.\n\nRegistered capabilities: %d\n\nUse an implemented subcommand for machine-readable JSON output.",
		registeredEntries,
	)
}
