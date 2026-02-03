package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/asztemborski/cardea/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

const (
	configRoot = "config"
	varsDir    = "vars"
)

var runCommandFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    configRoot,
		Value:   configRoot,
		Aliases: []string{"c"},
	},
	&cli.StringFlag{
		Name:    varsDir,
		Value:   "_vars",
		Aliases: []string{"v"},
	},
}

var runCommand = &cli.Command{
	Name:   "run",
	Usage:  "runs cardea instance",
	Action: runCommandAction,
	Flags:  runCommandFlags,
}

func runCommandAction(ctx context.Context, cmd *cli.Command) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		select {
		case sig := <-sigs:
			log.Info().Stringer("sig", sig).Msg("signal intercepted")
			cancel()
		case <-ctx.Done():
		}
	}()

	loader := config.NewLoader(
		cmd.String(configRoot),
		config.WithVarsDir(cmd.String(varsDir)),
	)

	cfg, err := loader.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	fmt.Printf("%#v\n", cfg)
	return nil
}
