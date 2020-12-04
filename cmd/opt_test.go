package cmd

import (
	"bytes"
	"testing"
)

func TestOption(t *testing.T) {
	var actionCalled bool
	f := func() { actionCalled = true }
	o := Option("a", "b", "c", f)
	if o.Short() != "a" {
		t.Errorf("wrong short name")
	}
	if o.Long() != "b" {
		t.Errorf("wrong long name")
	}
	if o.Description() != "c" {
		t.Errorf("wrong description")
	}
	o.Action()()
	if !actionCalled {
		t.Errorf("action didn't work")
	}
	// test rendering with a command
	o.cmd = New()
	if o.Description() != "c" {
		t.Errorf("wrong description")
	}
	// tests failed rendering
	b := new(bytes.Buffer)
	var errno int
	o = Option("", "", "{{.X}}", nil)
	o.cmd = New()
	o.cmd.stderr = b
	o.cmd.exit = func(n int) { errno = n }
	if o.Description() != "" {
		t.Errorf("expected description rendering to fail")
	}
	if b.String() != `template: description:1:2: executing "description" at <.X>: can't evaluate field X in type *cmd.Opt` {
		t.Errorf("expected a different error message (got: %s)", b.String())
	}
	if errno != ErrnoRenderFailed {
		t.Errorf("expected exit to be called properly")
	}
}
