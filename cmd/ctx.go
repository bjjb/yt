package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
)

// A Ctx provides the tools needed for a command-line application. A zero Ctx
// will use values from the environment.
type Ctx struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	Args           []string
	Getenv         func(string) string
	Exit           func(int)
}

// DefaultContext is the context used when none is supplied
var DefaultContext = &Ctx{
	Stdin:  os.Stdin,
	Stdout: os.Stdout,
	Stderr: os.Stderr,
	Args:   os.Args,
	Getenv: os.Getenv,
	Exit:   os.Exit,
}

// ErrnoRenderFailed indicates some template rendering failure
const ErrnoRenderFailed = 2

// Render renders the given template into a string, using itself as the view;
// if an error occurs, it will use the *Ctx's Stderr to print the error, and
// Exit(ErrnoRenderFailed).
func (c *Ctx) Render(t *template.Template) string {
	b := new(bytes.Buffer)
	if err := t.Execute(b, c); err != nil {
		fmt.Fprintf(c.Stderr, "couldn't render %s: %s", t.Name(), err)
		c.Exit(1)
	}
	return b.String()
}
