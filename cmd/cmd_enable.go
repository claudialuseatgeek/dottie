package main

import (
	"context"

	"dotfedi/pkg"

	"github.com/urfave/cli/v3"
)

var enableCommand = &cli.Command{
	Name:      "enable",
	Usage:     "Uncomment/enable a key if it exists",
	Before:    setup,
	ArgsUsage: "KEY",
	Action: func(_ context.Context, cmd *cli.Command) error {
		key := cmd.Args().Get(0)
		existing := env.Get(key)
		existing.Commented = false

		return pkg.Save(cmd.String("file"), env)
	},
}
