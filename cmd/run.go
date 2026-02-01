package cmd

import (
	"context"
	"fmt"

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
	cfg, err := config.NewLoader(cmd.String(configRoot), config.WithVarsDir(cmd.String(varsDir))).Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	fmt.Printf("%#v\n", cfg)
	return nil
}
