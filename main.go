package main

import (
	"fmt"
	"io"
	"os"

	"github.com/aswinkarthik/replace-text/fs"
	cli "github.com/urfave/cli/v2"
)

const (
	// AppName is the name of the CLI
	AppName = "replace-text"
	// ExitCodeValidationError is returned whenever there are validation errors
	// in the input
	ExitCodeValidationError = 2

	flagPatternsFile            = "patterns-file"
	metadataValidationErrorsKey = "validation-errors"
)

func main() {
	fs := fs.NewOsFs()
	cli.HelpPrinter = overrideDefaultPrinter(cli.HelpPrinter)
	app := &cli.App{
		Name:   AppName,
		Usage:  "Find & Replace multiple texts in files",
		Action: run,
		Writer: fs.DevNull(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagPatternsFile,
				Aliases: []string{"p"},
				Usage:   "Load find & replace patterns from a JSON file",
				EnvVars: []string{"PATTERNS_FILE", "REPLACE_TEXT_PATTERNS_FILE"},
			},
		},
		Before: parseInput(fs),
	}
	if err := app.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", AppName, err)
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	return fmt.Errorf("not implemented")
}

func parseInput(fs fs.Fs) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		flagPreset := ctx.IsSet(flagPatternsFile)

		if flagPreset {
			filename := ctx.String(flagPatternsFile)
			if !fs.IsFile(filename) {
				ctx.App.Metadata[metadataValidationErrorsKey] = true
				return cli.Exit(
					fmt.Sprintf(`%s: file "%s" does not exist`, AppName, filename),
					ExitCodeValidationError,
				)
			}
			return nil
		}

		return fmt.Errorf("not implemented")
	}
}

func overrideDefaultPrinter(defaultPrinter func(w io.Writer, templ string, data interface{})) func(w io.Writer, templ string, data interface{}) {
	return func(w io.Writer, templ string, data interface{}) {
		app, ok := data.(*cli.App)
		if ok {
			// Do not show help for validation errors
			if _, present := app.Metadata[metadataValidationErrorsKey]; present {
				return
			}
		}

		defaultPrinter(os.Stdout, templ, data)
	}
}
