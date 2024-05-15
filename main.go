package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/karl-gustav/runlogger"
)

var log *runlogger.Logger

func init() {
	if os.Getenv("K_SERVICE") != "" { // Check if running in cloud run
		log = runlogger.StructuredLogger()
	} else {
		log = runlogger.PlainLogger()
	}
}

func main() {
	r := chi.NewRouter()
	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		log.Info(
			r.Method+" "+fullURL(r),
			log.Field("url", fullURL(r)),
			log.Field("method", r.Method),
			log.Field("body", string(body)),
			log.Field("headers", r.Header),
		)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]any{
			"url":     fullURL(r),
			"method":  r.Method,
			"body":    string(body),
			"headers": r.Header,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Debug("Serving http://localhost:" + port)
	log.Critical(http.ListenAndServe(":"+port, r))
}

func fullURL(r *http.Request) string {
	scheme := "https"
	if strings.HasPrefix(r.Host, "localhost") {
		scheme = "http"
	}
	var query string
	if r.URL.RawQuery != "" {
		query = "?" + r.URL.RawQuery
	}
	var fragment string
	if r.URL.Fragment != "" {
		fragment = "#" + r.URL.Fragment
	}
	return fmt.Sprintf("%s://%s%s%s%s", scheme, r.Host, r.URL.Path, query, fragment)
}
