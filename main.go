package main

import (
	"log/slog"
	"net/http"

	"github.com/smcgarril/rune-seer/internal/api"
)

func main() {
	slog.Info("Starting server on :8080")

	http.HandleFunc("/", api.RootHandler)
	http.HandleFunc("/analyze", api.AnalyzeHandler)
	http.HandleFunc("/details", api.DetailsHandler)

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":8080", nil)
}
