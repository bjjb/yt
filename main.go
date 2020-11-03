package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const name = "yt"
const version = "0.0.1"

// init sets up ui to be a cli app with all the standard os goodies
func init() {
	ui = &cli{
		args: os.Args,
		in:   os.Stdin,
		out:  os.Stdout,
		err:  os.Stderr,
		env:  os.Getenv,
		exit: os.Exit,
	}
}

// main simply calls the main function of the ui
func main() { ui.main() }

// An app simply exposes a main function
type app interface{ main() }

// a cli wraps os variables, and implements app
type cli struct {
	args     []string
	in       io.Reader
	out, err io.Writer
	exit     func(int)
	env      func(string) string
}

// main parses flags and takes some action based on its setup
func (c *cli) main() {
	fmt.Fprintln(c.out, "hello!")
}

// ui is the app that will be invoked by main
var ui app

// getenv is a helper function to get a value from the environment with a
// fallback in case it's blank or unset
func getenv(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

// oauth returns is a middleware to provide OAuth2 tokens and to grant
// permissions to requests which contain valid tokens.
func oauth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})
}

// cors returns a CORS middleware
func cors() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		get := func(k string) string { return r.Header.Get(k) }
		set := func(k, v string) { w.Header().Set(k, v) }
		add := func(k, v string) { w.Header().Add(k, v) }
		err := func(msg string, status int) { http.Error(w, msg, status) }
		cut := func(v string) []string { return trim(strings.Split(v, ",")) }
		reqMethod := get("Access-Control-Request-Method")
		add("Vary", "Origin")
		add("Vary", "Access-Control-Request-Method")
		add("Vary", "Access-Control-Request-Headers")
		if method == http.MethodOptions && reqMethod != "" {
			if get("Origin") == "" {
				err("missing Origin", http.StatusBadRequest)
				return
			}
			reqHeaders := cut(r.Header.Get("Access-Control-Request-Headers"))
			set("Access-Control-Allow-Origin", get("Origin"))
			set("Access-Control-Allow-Methods", reqMethod)
			set("Access-Control-Allow-Headers", strings.Join(reqHeaders, ", "))
			set("Access-Control-Allow-Credentials", "true")
		}
	})
}

// youtube returns a http.Handler which can search for things on YouTube
func youtube() http.Handler {
	return new(http.ServeMux)
}

// trim is a utility function which trims whitespace from all the values
func trim(values []string) []string {
	result := []string{}
	for _, v := range values {
		result = append(result, strings.TrimSpace(v))
	}
	return result
}

// a doneWriter is a http.ResponseWriter whose done function is called
// whenever something is written (which means the response is finished).
type doneWriter struct {
	http.ResponseWriter
	done func()
}

// WriteHeader implements ResponseWriter.WriteHeader calling the done function
func (w *doneWriter) WriteHeader(status int) {
	w.done()
	w.ResponseWriter.WriteHeader(status)
}

// WriteHeader implements ResponseWriter.Write calling the done function
func (w *doneWriter) Write(b []byte) (int, error) {
	w.done()
	return w.ResponseWriter.Write(b)
}

// Chains a list of handlers together - the returned http.Handler will invoke
// each in turn until either w.WriteHeader or w.Write has been called (which
// indicates that the chain stops).
func chain(handlers ...http.Handler) http.Handler {
	var done bool
	wrap := func(w http.ResponseWriter) http.ResponseWriter {
		return &doneWriter{w, func() { done = true }}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, handler := range handlers {
			handler.ServeHTTP(wrap(w), r)
			if done {
				return
			}
		}
	})
}
