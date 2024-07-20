package cmd

import (
	"github.com/raystack/salt/cmdx"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "morphic <command> <subcommand> [flags]",
		Short:         "API for Morphic",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(
		serverCommand(),
		migrateCommand(),
		rollbackMigrations(),
	)

	cmdx.SetHelp(cmd)

	return cmd
}
