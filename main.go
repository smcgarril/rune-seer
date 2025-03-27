package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"unicode/utf8"

	"github.com/smcgarril/rune-seer/internal/api"
)

func main() {
	slog.Info("Starting server on :8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		api.Tmpl.Execute(w, nil)
	})

	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		input := r.FormValue("input")
		response := api.ProcessInput(input)
		w.Header().Set("Content-Type", "text/html")
		api.ResultTmpl.Execute(w, response)
	})

	http.HandleFunc("/details", func(w http.ResponseWriter, r *http.Request) {
		char := r.URL.Query().Get("char")
		if char == "" {
			http.Error(w, "Character not provided", http.StatusBadRequest)
			return
		}

		// Extract the first rune from the string
		runeVal, _ := utf8.DecodeRuneInString(char)
		runeInfo := api.ProcessRune(runeVal)

		fmt.Println("Rune Info: ", runeInfo)

		w.Header().Set("Content-Type", "text/html")
		api.DetailsTmpl.Execute(w, runeInfo)
	})

	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	http.ListenAndServe(":8080", nil)
}
