package main

import (
	"os"
)

func main() {
	c := newCommand(os.Stdin, os.Stdout, os.Stderr)
	os.Exit(c.run(os.Args))
}
