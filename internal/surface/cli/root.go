package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yargotev/exito-tools/internal/app"
)

// NewRoot builds the minimal English-only CLI root surface.
func NewRoot(application *app.Application) *cobra.Command {
	command := &cobra.Command{
		Use:   "exito",
		Short: "Exito Tools command-line interface",
		Long: fmt.Sprintf(
			"Exito Tools command-line interface\n\nExito Tools is the machine-first CLI surface for the application.\n\nRegistered foundation entries in this scaffold: %d\n\nThis foundation slice only provides bootstrap and help.",
			len(application.Registry.All()),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	command.SetHelpCommand(&cobra.Command{Hidden: true})
	return command
}
