package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/raystack/salt/cmdx"
	"github.com/sairahul1526/morphic/cmd"
)

const (
	exitOK    = 0
	exitError = 1
)

// @title Morphic API
// @version 0.0.1
// @schemes http
// @host localhost:8060

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	execute(ctx)
}

func execute(ctx context.Context) {
	root := cmd.New()
	command, err := root.ExecuteContextC(ctx)

	if err == nil {
		return
	}

	if cmdx.IsCmdErr(err) {
		if !strings.HasSuffix(err.Error(), "\n") {
			fmt.Println()
		}
		fmt.Println(command.UsageString())
		os.Exit(exitOK)
	}

	fmt.Println(err)
	os.Exit(exitError)
}
