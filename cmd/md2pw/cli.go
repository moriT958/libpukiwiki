package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/moriT958/libpukiwiki/converter"
)

type command struct {
	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
}

func newCommand(inStream io.Reader, outStream, errStream io.Writer) *command {
	return &command{
		inStream:  inStream,
		outStream: outStream,
		errStream: errStream,
	}
}

// isStdinPiped checks if stdin is a pipe (not a terminal)
func isStdinPiped() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) == 0
}

func (c *command) run(args []string) int {
	var outputFile string

	flags := flag.NewFlagSet("md2pw", flag.ContinueOnError)
	flags.SetOutput(c.errStream)
	flags.StringVar(&outputFile, "o", "", "output file path (default: stdout)")

	flags.Usage = func() {
		_, _ = fmt.Fprintf(c.errStream, "Usage: md2pw [options] [<file.md>|-]\n\n")
		_, _ = fmt.Fprintf(c.errStream, "Options:\n")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		return 1
	}

	var content []byte
	var err error

	if flags.NArg() >= 1 {
		filename := flags.Arg(0)
		if filename == "-" {
			// Explicit stdin with "-"
			content, err = io.ReadAll(c.inStream)
		} else {
			// File argument
			content, err = os.ReadFile(filename)
		}
	} else if isStdinPiped() {
		// No argument but stdin is piped
		content, err = io.ReadAll(c.inStream)
	} else {
		// No argument and no pipe
		_, _ = fmt.Fprintln(c.errStream, "Error: input file required")
		flags.Usage()
		return 1
	}

	if err != nil {
		_, _ = fmt.Fprintf(c.errStream, "Error reading input: %v\n", err)
		return 1
	}

	result, err := converter.ConvertToPukiwiki(content)
	if err != nil {
		_, _ = fmt.Fprintf(c.errStream, "Error converting: %v\n", err)
		return 1
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(result), 0644); err != nil {
			_, _ = fmt.Fprintf(c.errStream, "Error writing to file %s: %v\n", outputFile, err)
			return 1
		}
	} else {
		_, _ = fmt.Fprint(c.outStream, result)
	}

	return 0
}
