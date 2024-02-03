package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/jippi/dottie/pkg"
	"github.com/jippi/dottie/pkg/ast"

	"github.com/hashicorp/go-getter"
	"github.com/urfave/cli/v3"
)

var updateCommand = &cli.Command{
	Name:   "update",
	Usage:  "Update the .env file from a source",
	Before: setup,
	Action: func(_ context.Context, _ *cli.Command) error {
		fmt.Print("Finding source")
		source, err := env.GetConfig("dottie/source")
		if err != nil {
			return err
		}
		fmt.Println(" ✅")

		fmt.Print("Grabbing .env.docker from [", source, "]")

		if _, err := os.Stat(".env.source"); errors.Is(err, os.ErrNotExist) {
			tmp, err := os.OpenFile(".env.source", os.O_RDWR|os.O_CREATE, 0o666)
			if err != nil {
				return err
			}

			// Grab source file

			client := getter.Client{
				DisableSymlinks: true,
				Mode:            getter.ClientModeFile,
				Src:             source,
				Dst:             tmp.Name(),
			}

			if err := client.Get(); err != nil {
				return err
			}
		}
		fmt.Println(" ✅")

		// Copy source to "new"
		fmt.Print("Copying .env.source into .env.merged")
		if err := Copy(".env.merged", ".env.source"); err != nil {
			return err
		}
		fmt.Println(" ✅")

		// Load the soon-to-be-merged file
		fmt.Print("Loading and parsing .env.merged")
		mergedEnv, err := pkg.Load(".env.merged")
		if err != nil {
			return err
		}
		fmt.Println(" ✅")

		// Take current assignments and set them in the new doc
		fmt.Println("Updating .env.merged with key/value pairs from .env")
		for _, stmt := range env.Assignments() {
			if stmt.Active {
				continue
			}

			changed, err := mergedEnv.Set(&ast.Assignment{Name: stmt.Name, Literal: stmt.Literal}, ast.SetOptions{SkipIfSame: true})
			if err != nil {
				return err
			}

			if changed {
				fmt.Println("  ✅", stmt.Name)
			}
		}

		fmt.Println()
		fmt.Println("Saving .env.merged")

		return pkg.Save(".env.merged", mergedEnv)
	},
}

func Copy(dst, src string) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	info, err := srcF.Stat()
	if err != nil {
		return err
	}

	dstF, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstF.Close()

	if _, err := io.Copy(dstF, srcF); err != nil {
		return err
	}

	return nil
}
