package api

import (
	"net/http"
	"unicode/utf8"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
	}
}

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	input := r.FormValue("input")
	response := processStringInput(input)
	w.Header().Set("Content-Type", "text/html")
	if err := resultTmpl.Execute(w, response); err != nil {
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
		return
	}
}

func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	char := r.FormValue("char")
	if char == "" {
		http.Error(w, "Character not provided", http.StatusBadRequest)
		return
	}

	runeVal, size := utf8.DecodeRuneInString(char)
	if runeVal == utf8.RuneError && size == 1 {
		http.Error(w, "Invalid UTF-8 character", http.StatusBadRequest)
		return
	}
	runeInfo := processRune(runeVal)

	w.Header().Set("Content-Type", "text/html")
	if err := detailsTmpl.Execute(w, runeInfo); err != nil {
		http.Error(w, "Template rendering failed", http.StatusInternalServerError)
		return
	}
}
