package cmd

import (
	"bytes"
	"fmt"
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
	return c.render(c.summary, nil)
}

// Description gets the description of the Cmd in the given context
func (c *Cmd) Description() string {
	return c.render(c.description, nil)
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
func New(directives ...Directive) *Cmd {
	c := new(Cmd)
	for _, directive := range directives {
		c = directive(c)
	}
	return c
}

// A Directive is a command-modifying function, such as Name (which sets the
// name).
type Directive func(*Cmd) *Cmd

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

func (c *Cmd) render(t *template.Template, ctx interface{}) string {
	b := new(bytes.Buffer)
	if ctx == nil {
		ctx = c
	}
	if err := t.Execute(b, ctx); err != nil {
		fmt.Fprintf(c.stderr, "render failed - %s (%e)", t.Name(), err)
		c.exit(ErrnoRenderFailed)
	}
	return b.String()
}

// Name gets a directive to set a command's name
func Name(name string) Directive {
	return Directive(func(c *Cmd) *Cmd {
		return c.setName(name)
	})
}

// Version gets a directive to set a command's version
func Version(version string) Directive {
	return Directive(func(c *Cmd) *Cmd {
		return c.setVersion(version)
	})
}

// Summary gets a directive to set a command's summary
func Summary(summary string) Directive {
	return Directive(func(c *Cmd) *Cmd {
		return c.setSummary(summary)
	})
}

// Description gets a directive to set a command's description
func Description(description string) Directive {
	return Directive(func(c *Cmd) *Cmd {
		return c.setDescription(description)
	})
}

// Action gets a directive to add an action function to a command
func Action(action func(*Ctx)) Directive {
	return Directive(func(c *Cmd) *Cmd {
		return c.setAction(action)
	})
}

// Options gets a directive to add Opts to a Cmd
func Options(opts ...*Opt) Directive {
	return Directive(func(c *Cmd) *Cmd {
		return c.addOptions(opts...)
	})
}

// Commands gets a directive to add sub-Cmds to a Cmd
func Commands(cmds ...*Cmd) Directive {
	return Directive(func(c *Cmd) *Cmd {
		return c.addCommands(cmds...)
	})
}

func (c *Cmd) on(parent *Cmd) *Cmd {
	c.cmd = parent
	return c
}
