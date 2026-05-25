package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yargotev/exito-tools/internal/app"
	"github.com/yargotev/exito-tools/internal/config"
)

// Bootstrapper builds the application after Cobra has parsed CLI boot flags.
type Bootstrapper func(app.Options) (*app.Application, error)

type rootOptions struct {
	configPath string
	profile    string
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
			application, err := bootstrap(app.Options{
				Config: config.Options{
					ConfigPath: options.configPath,
					Profile:    options.profile,
				},
			})
			if err != nil {
				return err
			}

			cmd.Long = rootLong(len(application.Registry.All()))
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	command.PersistentFlags().StringVar(&options.configPath, "config", "", "Path to the Exito Tools configuration file")
	command.PersistentFlags().StringVar(&options.profile, "profile", "", "Configuration profile to use")
	command.SetHelpCommand(&cobra.Command{Hidden: true})
	return command
}

func rootLong(registeredEntries int) string {
	return fmt.Sprintf(
		"Exito Tools command-line interface\n\nExito Tools is the machine-first CLI surface for the application.\n\nRegistered foundation entries in this scaffold: %d\n\nThis foundation slice only provides bootstrap and help.",
		registeredEntries,
	)
}
