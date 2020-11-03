package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

type appFunc func()

func (t appFunc) main() { t() }

func Test_init(t *testing.T) {
	t.Run("ui", func(t *testing.T) {
		_, ok := ui.(app)
		t.Run("is an app", assert(ok))
		cli, ok := ui.(*cli)
		t.Run("is a *cli", assert(ok))
		t.Run("uses os.Args", assertSame(os.Args, cli.args))
		t.Run("uses os.Stdin", assertSame(os.Stdin, cli.in))
		t.Run("uses os.Stdout", assertSame(os.Stdout, cli.out))
		t.Run("uses os.Stderr", assertSame(os.Stderr, cli.err))
		t.Run("uses os.Exit", assertSame(os.Exit, cli.exit))
		t.Run("uses os.Getenv", assertSame(os.Getenv, cli.env))
	})
}

func Test_main(t *testing.T) {
	var called bool
	var test app
	test, ui = ui, appFunc(func() { called = true })
	defer func() { ui = test }()
	main()
	t.Run("calls ui.main", assert(called))
}

func Test_getenv(t *testing.T) {
	t.Run(
		"it reads existing env vars",
		assertEqual(os.Getenv("USER"), getenv("USER", "WRONG")),
	)
	t.Run(
		"it returns the fallback for missing env vars",
		assertEqual("fallback", getenv("✗", "fallback")),
	)
}

func Test_cors(t *testing.T) {
	t.Run("preflight request", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Set("Origin", "example.com")
		r.Header.Set("Access-Control-Request-Method", http.MethodGet)
		cors().ServeHTTP(w, r)
		t.Run(
			"allows the origin",
			assertHeader(w, "Access-Control-Allow-Origin", "example.com"),
		)
		t.Run(
			"allows the methods properly",
			assertHeader(w, "Access-Control-Allow-Methods", http.MethodGet),
		)
		t.Run(
			"allows credentials",
			assertHeader(w, "Access-Control-Allow-Credentials", "true"),
		)
	})
	t.Run("preflight request without an origin", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodOptions, "/", nil)
		r.Header.Set("Access-Control-Request-Method", http.MethodGet)
		cors().ServeHTTP(w, r)
		t.Run(
			"it gives a BadRequest response",
			assertStatusCode(w.Result(), http.StatusBadRequest),
		)
	})
	t.Run("normal request", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		cors().ServeHTTP(w, r)
		// TODO: more tests
	})
}

func Test_youtube(t *testing.T) {
	// TODO: test me
}

func Test_oauth(t *testing.T) {
	// TODO: test me
}

func assert(v bool) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		if !v {
			t.Error("expected true")
		}
	}
}

func assertEqual(expected, actual interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		if actual != expected {
			t.Errorf("expected %q, got %q", expected, actual)
		}
	}
}

func assertSame(expected, actual interface{}) func(t *testing.T) {
	return assertEqual(fmt.Sprintf("%p", expected), fmt.Sprintf("%p", actual))
}

func assertHeader(w http.ResponseWriter, name, value string) func(t *testing.T) {
	return assertEqual(value, w.Header().Get(name))
}

func assertStatusCode(r *http.Response, status int) func(t *testing.T) {
	return assertEqual(strconv.Itoa(status), strconv.Itoa(r.StatusCode))
}
