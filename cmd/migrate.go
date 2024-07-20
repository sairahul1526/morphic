package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/sairahul1526/morphic/config"
	"github.com/sairahul1526/morphic/store"
	"github.com/spf13/cobra"
)

func migrateCommand() *cobra.Command {
	var configFile string
	cmd := &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"m"},
		Short:   "Migrations management",
		Long:    "Migration to apply changes to database",
		Example: heredoc.Doc(`
			$ morphic migrate
		`),
		Annotations: map[string]string{
			"group": "database",
		},
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configFile)
		if err != nil {
			return err
		}

		return store.Migrate(cfg.Database)
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}

func rollbackMigrations() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:     "rollback",
		Aliases: []string{"r"},
		Short:   "Rollback management",
		Long:    "Rollback to previously migrated changes from database",
		Example: heredoc.Doc(`
			$ morphic rollback
		`),
		Annotations: map[string]string{
			"group": "database",
		},
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configFile)
		if err != nil {
			return err
		}

		return store.Rollback(cfg.Database)
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}
