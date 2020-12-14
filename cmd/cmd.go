// Package cmd provides functions to help build command-line applications.
// Each application (*Cmd) can contain nested applications.
package cmd

import (
	"bytes"
	"io"
	"text/template"
)

// A Cmd is a basic command-line application, consisting of a name, a version,
// a short description, a long description, options, sub-commands and a
// default action.
type Cmd struct {
	name, version              string
	summary, description, help *template.Template
	aliases                    []string
	action                     func(*Ctx)
	opts                       []*Opt
	cmds                       []*Cmd
	stdin                      io.Reader
	stdout, stderr             io.Writer
	getenv                     func(string) string
	exit                       func(int)
	cmd                        *Cmd
}

// Execute runs the Cmd in the given Ctx
func (c *Cmd) Execute(ctx *Ctx) {
	action := c.Parse(ctx)
	action(ctx)
}

// Parse gets a function which will run the command
func (c *Cmd) Parse(ctx *Ctx) func(*Ctx) *Ctx {
	return func(c *Ctx) *Ctx {
		return c
	}
}

// Name gets the name of the Cmd
func (c *Cmd) Name() string {
	return c.name
}

// Version gets the version of the Cmd
func (c *Cmd) Version() string {
	return c.version
}

// Summary gets the summary of the Cmd in the given context
func (c *Cmd) Summary() string {
	return c.render(c.summary)
}

// Description gets the description of the Cmd in the given context
func (c *Cmd) Description() string {
	return c.render(c.description)
}

// Options gets the attached options
func (c *Cmd) Options() []*Opt {
	return c.opts
}

// Commands gets the attached sub-commands
func (c *Cmd) Commands() []*Cmd {
	return c.cmds
}

// Command gets a Cmd with the given name, version, summary and description
func Command(name, version, summary, description string) *Cmd {
	return New(
		Name(name),
		Version(version),
		Summary(summary),
		Description(description),
	)
}

// New builds a new Cmd from Directives
func New(directives ...Modifier) *Cmd {
	c := new(Cmd)
	for _, directive := range directives {
		c = directive(c)
	}
	return c
}

// A Modifier is a command-modifying function, such as Name (which sets the
// command's name).
type Modifier func(*Cmd) *Cmd

func (c *Cmd) setName(name string) *Cmd {
	c.name = name
	return c
}

func (c *Cmd) setVersion(version string) *Cmd {
	c.version = version
	return c
}

func (c *Cmd) setSummary(summary string) *Cmd {
	c.summary = template.Must(template.New("summary").Parse(summary))
	return c
}

func (c *Cmd) setDescription(description string) *Cmd {
	c.description = template.Must(template.New("description").Parse(description))
	return c
}

func (c *Cmd) setAction(action func(*Ctx)) *Cmd {
	c.action = action
	return c
}

func (c *Cmd) setStdout(w io.Writer) *Cmd {
	c.stdout = w
	return c
}

func (c *Cmd) setStderr(w io.Writer) *Cmd {
	c.stderr = w
	return c
}

func (c *Cmd) setStdin(r io.Reader) *Cmd {
	c.stdin = r
	return c
}

func (c *Cmd) setExit(f func(int)) *Cmd {
	c.exit = f
	return c
}

func (c *Cmd) setGetenv(f func(string) string) *Cmd {
	c.getenv = f
	return c
}

func (c *Cmd) addOptions(opts ...*Opt) *Cmd {
	for _, o := range opts {
		c.opts = append(c.opts, o.on(c))
	}
	return c
}

func (c *Cmd) addCommands(cmds ...*Cmd) *Cmd {
	for _, sc := range cmds {
		c.cmds = append(c.cmds, sc.on(c))
	}
	return c
}

// Name gets a directive to set a command's name
func Name(name string) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setName(name)
	})
}

// Version gets a directive to set a command's version
func Version(version string) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setVersion(version)
	})
}

// Summary gets a directive to set a command's summary
func Summary(summary string) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setSummary(summary)
	})
}

// Description gets a directive to set a command's description
func Description(description string) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setDescription(description)
	})
}

// Action gets a directive to add an action function to a command
func Action(action func(*Ctx)) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setAction(action)
	})
}

// Stdout gets a modifier to set the standard output stream of the command
func Stdout(w io.Writer) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setStdout(w)
	})
}

// Stderr gets a modifier to set the standard error stream of the command
func Stderr(w io.Writer) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setStderr(w)
	})
}

// Stdin gets a modifier to set the standard input stream of the command
func Stdin(r io.Reader) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setStdin(r)
	})
}

// Exit gets a modifier to set the exit function of the command
func Exit(f func(int)) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setExit(f)
	})
}

// Getenv gets a modifier to set the exit function of the command
func Getenv(f func(string) string) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.setGetenv(f)
	})
}

// Options gets a directive to add Opts to a Cmd
func Options(opts ...*Opt) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.addOptions(opts...)
	})
}

// Commands gets a directive to add sub-Cmds to a Cmd
func Commands(cmds ...*Cmd) Modifier {
	return Modifier(func(c *Cmd) *Cmd {
		return c.addCommands(cmds...)
	})
}

func (c *Cmd) on(parent *Cmd) *Cmd {
	c.cmd = parent
	return c
}

func (c *Cmd) render(t *template.Template) string {
	b := new(bytes.Buffer)
	if err := t.Execute(b, c); err != nil {
		fatal(c.stderr, c.exit, ErrnoRenderFailed, err.Error())
	}
	return b.String()
}
