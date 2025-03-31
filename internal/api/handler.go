package api

import (
	"net/http"
	"unicode/utf8"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	input := r.FormValue("input")
	response := processStringInput(input)
	w.Header().Set("Content-Type", "text/html")
	resultTmpl.Execute(w, response)
}

func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	char := r.URL.Query().Get("char")
	if char == "" {
		http.Error(w, "Character not provided", http.StatusBadRequest)
		return
	}

	runeVal, _ := utf8.DecodeRuneInString(char)
	runeInfo := processRune(runeVal)

	w.Header().Set("Content-Type", "text/html")
	detailsTmpl.Execute(w, runeInfo)
}
