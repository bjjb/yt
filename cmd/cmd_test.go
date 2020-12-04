package cmd

import (
	"testing"
)

func TestNew(t *testing.T) {
	var a, b *Cmd
	a = New()
	if a == nil {
		t.Errorf("expected New() to make a new *Cmd")
	}
	b = New(Directive(func(c *Cmd) *Cmd { return a }))
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
	if New(Description("foo")).Description() != "foo" {
		t.Errorf("expected description to be set")
	}
}

func TestAction(t *testing.T) {
	var ok bool
	New(Action(func(*Ctx) { ok = true })).action(nil)
	if !ok {
		t.Errorf("expected action to be set")
	}
}

func TestOpts(t *testing.T) {
	f := func() {}
	c := New(Options(Option("x", "xxx", "xxxx", f)))
	if len(c.Options()) == 0 {
		t.Errorf("expected an option")
	}
}

func TestCmds(t *testing.T) {
	c := New(Commands(Command("a", "b", "c", "d")))
	if len(c.Commands()) == 0 {
		t.Fatal("expected a sub command")
	}
	c = c.cmds[0]
	if c.Name() != "a" {
		t.Errorf("expected a name")
	}
	if c.Version() != "b" {
		t.Errorf("expected a version")
	}
	if c.Summary() != "c" {
		t.Errorf("expected a summary")
	}
	if c.Description() != "d" {
		t.Errorf("expected a description")
	}
}
