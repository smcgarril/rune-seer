package api

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	oldTmpl := tmpl
	tmpl = template.Must(template.New("test").Parse("<title>Mock Rune Seer</title>"))
	defer func() { tmpl = oldTmpl }()

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	RootHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", res.Status)
	}

	body := w.Body.String()
	if !strings.Contains(body, "<title>Mock Rune Seer</title>") {
		t.Fatalf("Expected response to contain mocked page title, but got: %s", body)
	}
}

func TestAnalyzeHandler(t *testing.T) {
	oldResultTmpl := resultTmpl
	resultTmpl = template.Must(template.New("test").Parse(`<div class="rune-box">Mock Rune</div>`))
	defer func() { resultTmpl = oldResultTmpl }()

	formData := "input=test"
	req := httptest.NewRequest("POST", "/analyze", strings.NewReader(formData))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	AnalyzeHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", res.Status)
	}

	body := w.Body.String()
	if !strings.Contains(body, `<div class="rune-box">Mock Rune</div>`) {
		t.Fatalf("Expected response to contain mocked rune box, but got: %s", body)
	}
}

func TestDetailsHandler(t *testing.T) {
	oldDetailsTmpl := detailsTmpl
	detailsTmpl = template.Must(template.New("test").Parse(`<div class="details-box">Mock Details</div>`))
	defer func() { detailsTmpl = oldDetailsTmpl }()

	req := httptest.NewRequest("GET", "/details?char=A", nil)
	w := httptest.NewRecorder()

	DetailsHandler(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK; got %v", res.Status)
	}

	body := w.Body.String()
	if !strings.Contains(body, `<div class="details-box">Mock Details</div>`) {
		t.Fatalf("Expected response to contain mocked details box, but got: %s", body)
	}
}
