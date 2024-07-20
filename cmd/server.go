package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/MakeNowJust/heredoc"
	"github.com/sairahul1526/morphic/config"
	"github.com/sairahul1526/morphic/logger"
	"github.com/sairahul1526/morphic/server"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func serverCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server <command>",
		Aliases: []string{"s"},
		Short:   "Server management",
		Long:    "Server management commands.",
		Example: heredoc.Doc(`
			$ morphic server start
			$ morphic server start -c ./config.yaml
		`),
		Annotations: map[string]string{
			"group": "server",
		},
	}

	cmd.AddCommand(startCommand())
	return cmd
}

func startCommand() *cobra.Command {
	var configFile string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the server",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configFile)
		if err != nil {
			return err
		}

		logger.SetLevel(cfg.Log.Level)

		return runServer(cmd.Context(), cfg)
	}

	cmd.Flags().StringVarP(&configFile, "config", "c", "./config.yaml", "Config file path")
	return cmd
}

func runServer(baseCtx context.Context, cfg config.Config) error {
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()

	// Handle OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-quit
		logger.Info("receive interrupt signal", zap.Any("signal", sig))
		cancel() // Cancel the context when signal is received
	}()

	// start api server
	return server.Serve(ctx, cfg.Service.Addr(), cfg)
}
