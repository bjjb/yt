package cmd

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var a, b *Cmd
	a = New()
	if a == nil {
		t.Errorf("expected New() to make a new *Cmd")
	}
	b = New(
		Stdin(new(bytes.Buffer)),
		Stdout(new(bytes.Buffer)),
		Stderr(new(bytes.Buffer)),
		Getenv(func(_ string) string { return "value" }),
		Exit(func(_ int) {}),
		Modifier(func(c *Cmd) *Cmd { return a }),
	)
	if b != a {
		t.Errorf("expected b to == a")
	}
}

func TestName(t *testing.T) {
	if New(Name("foo")).Name() != "foo" {
		t.Errorf("expected name to be set")
	}
}

func TestVersion(t *testing.T) {
	if New(Version("foo")).Version() != "foo" {
		t.Errorf("expected version to be set")
	}
}

func TestSummary(t *testing.T) {
	if New(Summary("foo")).Summary() != "foo" {
		t.Errorf("expected summary to be set")
	}
}

func TestDescription(t *testing.T) {
	x := -1
	b := new(bytes.Buffer)
	e := func(i int) { x = i }
	c := New(Description("foo"))
	if c.Description() != "foo" {
		t.Errorf("expected description to be set")
	}
	// Tests description with a broken template
	c = New(Description("{{.X}}"), Stderr(b), Exit(e))
	if c.Description() != "" {
		t.Errorf("expected no description")
	}
	if x != ErrnoRenderFailed {
		t.Errorf("expected exit ErrnoRenderFailed")
	}
	if b.String() != `template: description:1:2: executing "description" at <.X>: can't evaluate field X in type *cmd.Cmd` {
		t.Errorf("expected a different error message; got:\n%s", b.String())
	}
}

func TestAction(t *testing.T) {
	var ok bool
	New(Action(func(*Ctx) { ok = true })).action(nil)
	if !ok {
		t.Errorf("expected action to be set")
	}
}

func TestOptions(t *testing.T) {
	f := func() {}
	c := New(Options(Option("x", "xxx", "xxxx", f)))
	if len(c.Options()) == 0 {
		t.Errorf("expected an option")
	}
}

func TestCommands(t *testing.T) {
	c := New(Commands(Command("a", "b", "c", "d")))
	if len(c.Commands()) == 0 {
		t.Fatal("expected a sub command")
	}
}
