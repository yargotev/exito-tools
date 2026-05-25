package cli

import (
	"fmt"
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
