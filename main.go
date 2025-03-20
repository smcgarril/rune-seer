package main

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"unicode/utf8"
)

var utf8Mask = []string{
	"0xxxxxxx",
	"10xxxxxx",
	"1110xxxx",
}

type RuneInfo struct {
	Char      string
	RuneIndex int
	RuneBytes []RuneByte
	CodePoint string
}

type RuneByte struct {
	ByteIndex     int
	RuneByteIndex int
	Byte          byte
	Binary        string
	Utf8Mask      string
}

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Rune Seer</title>
    <script src="https://unpkg.com/htmx.org@1.9.5"></script>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }
        h1 {
            text-align: center;
        }
		p {
			text-align: center;
		}
        form {
            display: flex;
            justify-content: center;
            gap: 10px;
            margin-bottom: 20px;
        }
        input[type="text"] {
            padding: 8px;
            font-size: 1em;
            border: 1px solid #ccc;
            border-radius: 5px;
        }
        button {
            padding: 8px 12px;
            font-size: 1em;
            border: none;
            background-color: #007bff;
            color: white;
            border-radius: 5px;
            cursor: pointer;
        }
        .rune-container {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            justify-content: center;
        }
        .rune-box {
            border: 1px solid black;
            padding: 10px;
            text-align: center;
            background-color: #d4edda;
            border-radius: 5px;
            cursor: pointer;
            min-width: 50px;
            display: flex;
            flex-direction: column;
            align-items: center;
            position: relative;
        }
        .char {
            font-size: 1.5em;
            font-weight: bold;
            min-height: 30px;
            display: flex;
            align-items: center;
            justify-content: center;
        }
		.details-box {
			visibility: hidden; /* Keeps it invisible but allows space to be dynamically adjusted */
			opacity: 0; /* Fully transparent */
			transition: opacity 0.3s ease-in-out; /* Smooth fade-in */
			margin-top: 20px;
			padding: 10px;
			background-color: #d4edda; /* Match the green box */
			border: 1px solid black; /* Consistent with .rune-box */
			border-radius: 5px;
			text-align: center;
			width: fit-content; /* Only as wide as needed */
			max-width: 90%; /* Prevents overflow */
		}
        .bytes-info {
            margin-top: 5px;
            display: flex;
            gap: 5px;
            justify-content: center;
            flex-wrap: wrap;
        }
        .byte-box {
            border: 1px solid gray;
            padding: 5px;
            background-color: #f4f4f4;
            border-radius: 3px;
            font-size: 0.9em;
            text-align: center;
            width: 60px;
            height: 40px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
        }
        .byte-index {
            font-size: 0.8em;
            color: #555;
            margin-top: 2px;
        }
        .mask-container {
            margin-top: 10px;
            padding: 10px;
            background-color: #cce5ff;
            border-radius: 5px;
            display: none;
            position: absolute;
            top: 100%;
            left: 0;
            width: 100%;
        }
    </style>
</head>
<body>
    <h1>Rune Seer</h1>
	<p>Enter text to see the UTF-8 encoding of each rune.</p>
	<form hx-post="/analyze" hx-target="#result" hx-swap="innerHTML">
		<input type="text" name="input" placeholder="Enter text here" required>
		<button type="submit">Submit</button>
	</form>
    <div id="result" class="rune-container"></div>
	<div id="details-box" class="details-box"></div>
</body>
</html>
`))

var resultTmpl = template.Must(template.New("result").Parse(`
{{range .}}
	<div class="rune-box" 
		hx-get="/details?char={{.Char}}" 
		hx-target="#details-box" 
		hx-swap="innerHTML show:top">
		<div class="char">{{.Char}}</div>
		<div class="bytes-info">
			{{range .RuneBytes}}
				<div class="byte-box">
					{{.Binary}}
				</div>
			{{end}}
		</div>
	</div>
{{end}}
`))

var detailsTmpl = template.Must(template.New("details").Parse(`
<div class="details-box" style="visibility: visible; opacity: 1;">
    <h2>Character: {{.Char}}</h2>
    <div class="bytes-info">
        {{range .RuneBytes}}
            <div class="byte-box">
                Mask: {{.Utf8Mask}}<br>
                {{.Binary}}
            </div>
        {{end}}
    </div>
    <div style="margin-top: 10px; font-weight: bold;">
        Unicode Binary: TBD
    </div>
    <div style="font-size: 1.2em; font-weight: bold;">
        Code Point: {{.CodePoint}}
    </div>
</div>
`))

func main() {
	slog.Info("Starting server on :8080")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/analyze", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		input := r.FormValue("input")
		response := processInput(input)
		w.Header().Set("Content-Type", "text/html")
		resultTmpl.Execute(w, response)
	})

	http.HandleFunc("/details", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("HERE!")
		char := r.URL.Query().Get("char")
		if char == "" {
			http.Error(w, "Character not provided", http.StatusBadRequest)
			return
		}

		// Extract the first rune from the string
		runeVal, _ := utf8.DecodeRuneInString(char)
		runeInfo := processRune(runeVal)
		fmt.Println("Rune Info: ", runeInfo)
		w.Header().Set("Content-Type", "text/html")
		detailsTmpl.Execute(w, runeInfo)
	})

	http.ListenAndServe(":8080", nil)
}

func processInput(input string) []RuneInfo {
	var runes []RuneInfo
	n := 0
	for i, r := range input {
		rune := RuneInfo{
			Char:      string(r),
			RuneIndex: i,
			CodePoint: fmt.Sprintf("%d", r)}
		utf8Bytes := []byte(string(r))
		for j, b := range utf8Bytes {
			rune.RuneBytes = append(rune.RuneBytes, RuneByte{
				ByteIndex:     n,
				RuneByteIndex: j,
				Byte:          b,
				Binary:        fmt.Sprintf("%08b", b),
				Utf8Mask:      utf8Mask[j],
			})
			n++
		}
		runes = append(runes, rune)
	}
	return runes
}

func processRune(r rune) RuneInfo {
	rune := RuneInfo{
		Char:      string(r),
		CodePoint: fmt.Sprintf("%d", r),
	}

	utf8Bytes := []byte(string(r))
	for j, b := range utf8Bytes {
		rune.RuneBytes = append(rune.RuneBytes, RuneByte{
			RuneByteIndex: j,
			Byte:          b,
			Binary:        fmt.Sprintf("%08b", b),
			Utf8Mask:      utf8Mask[j],
		})
	}
	return rune
}
