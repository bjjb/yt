package cmd

import (
	"bytes"
	"io"
	"text/template"
)

// An Opt is a command-line option, with a value and optional action
type Opt struct {
	short, long string
	description *template.Template
	action      func()
	cmd         *Cmd
}

// Option builds an Opt
func Option(short, long, description string, action func()) *Opt {
	desc := template.Must(template.New("description").Parse(description))
	return &Opt{short, long, desc, action, nil}
}

// Short gets the option's short name (without the dash)
func (o *Opt) Short() string {
	return o.short
}

// Long gets the option's long name (without the dashes)
func (o *Opt) Long() string {
	return o.long
}

// Description gets the description of o in the given context
func (o *Opt) Description() string {
	return o.render(o.description)
}

// Action gets the action which will be performed when the option is parsed
func (o *Opt) Action() func() {
	return o.action
}

func (o *Opt) on(c *Cmd) *Opt {
	o.cmd = c
	return o
}

func (o *Opt) render(t *template.Template) string {
	b := new(bytes.Buffer)
	var stderr io.Writer
	var exit func(int)
	if o.cmd != nil {
		stderr = o.cmd.stderr
		exit = o.cmd.exit
	}
	if err := t.Execute(b, o); err != nil {
		fatal(stderr, exit, ErrnoRenderFailed, err.Error())
	}
	return b.String()
}
