package main

import (
	"context"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/render"

	"github.com/urfave/cli/v3"
)

func setup(_ context.Context, cmd *cli.Command) error {
	var err error

	env, err = pkg.Load(cmd.String("file"))
	if err != nil {
		return err
	}

	settings = render.NewSettings(
		render.WithBlankLines(cmd.Root().Bool("with-blank-lines")),
		render.WithComments(cmd.Root().Bool("with-comments")),
		render.WithGroupBanners(cmd.Root().Bool("with-groups")),
		render.WithFilterGroup(cmd.Root().String("group")),
		render.WithFilterKeyPrefix(cmd.Root().String("key-prefix")),
		render.WithIncludeDisabled(cmd.Root().Bool("include-commented")),
		render.WithFormattedOutput(cmd.Root().Bool("pretty")),
		render.WithColors(cmd.Root().Bool("colors")),
	)

	return nil
}
