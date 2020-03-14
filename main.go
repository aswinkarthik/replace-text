package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aswinkarthik/replace-text/fs"
	"github.com/aswinkarthik/replace-text/replacer"
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
		Name:            AppName,
		Usage:           "Find & Replace multiple texts in files",
		ArgsUsage:       "[PATH ...]",
		Action:          run(fs),
		Writer:          fs.DevNull(),
		HideHelpCommand: true,
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

func run(fs fs.Fs) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		patternsFileName := ctx.String(flagPatternsFile)
		patternsFile, err := fs.Open(patternsFileName)
		if err != nil {
			return fmt.Errorf("error opening patterns-file: %v", err)
		}

		var patterns map[string]string
		if err := json.NewDecoder(patternsFile).Decode(&patterns); err != nil {
			return fmt.Errorf("error decoding patterns file: %v", err)
		}

		r, err := replacer.NewReplacer(patterns)
		if err != nil {
			return fmt.Errorf("error creating replacer for given patterns: %v", err)
		}

		for _, inputFile := range ctx.Args().Slice() {
			file, err := fs.Open(inputFile)
			if err != nil {
				return fmt.Errorf("error opening input file: %v", err)
			}

			if err := r.Replace(file, os.Stdout); err != nil {
				return fmt.Errorf("error finding and replacing content in input file: %v", err)
			}
		}

		return nil
	}
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

		}

		if ctx.NArg() > 0 {
			for _, arg := range ctx.Args().Slice() {
				if !fs.IsFile(arg) {
					ctx.App.Metadata[metadataValidationErrorsKey] = true
					return cli.Exit(
						fmt.Sprintf(`%s: file "%s" does not exist`, AppName, arg),
						ExitCodeValidationError,
					)
				}
			}
		}

		if ctx.NArg() == 0 {
			ctx.App.Metadata[metadataValidationErrorsKey] = true
			return cli.Exit(
				fmt.Sprintf("%s: reading input from stdin is not supported yet", AppName),
				ExitCodeValidationError,
			)
		}

		return nil
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
