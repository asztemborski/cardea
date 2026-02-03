package cmd

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

var rootCommand = &cli.Command{
	Name:                   "cardea",
	UseShortOptionHandling: true,
	EnableShellCompletion:  true,
	Commands: []*cli.Command{
		runCommand,
	},
}

func Execute(ctx context.Context, args []string) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if err := rootCommand.Run(ctx, args); err != nil {
		log.Fatal().Err(err).Msg("failed to execute command")
	}
}
