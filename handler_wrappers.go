package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type recorderResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *recorderResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *recorderResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func Logger(h http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rrw := &recorderResponseWriter{w, 200}
		h.ServeHTTP(rrw, r)
		log.Printf(
			"\t%d\t%s\t%s\t%s\t%s",
			rrw.statusCode,
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

func Recoverer(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintln(w, "{ \"error\" : \"")
				fmt.Fprintf(w, "%s", debug.Stack())
				fmt.Fprintln(w, "\"}")
			}
		}()
		h.ServeHTTP(w, r)
	})
}
