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
	"110xxxxx",
	"1110xxxx",
	"11110xxx",
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
	Utf8Remainder string
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
        .submit-button {
            padding: 8px 12px;
            font-size: 1em;
            border: none;
            background-color: #5a9bd5;
            color: white;
            border-radius: 5px;
            cursor: pointer;
        }
		.submit-button:hover {
			background-color: #4a8bc4; /* Darker blue on hover */
		}
		.reset-button {
			background-color: #e89f71; /* Red color */
			color: white;
			border: none;
			padding: 8px 12px;
			font-size: 1em;
			border-radius: 5px;
			cursor: pointer;
		}
		.reset-button:hover {
    		background-color: #d88a5f; /* Darker red on hover */
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
		.rune-box:hover {
    		background-color: #b8dfc2; /* Slightly darker green */
		}
        .char {
            font-size: 1.5em;
            font-weight: bold;
            min-height: 30px;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .details-container {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            justify-content: center;
			margin-top: 10px;
        }
		.details-box {
            border: 1px solid black;
            padding: 10px;
            text-align: center;
            background-color:rgb(212, 219, 237);
            border-radius: 5px;
            cursor: pointer;
            min-width: 50px;
            display: flex;
            flex-direction: column;
            align-items: center;
            position: relative;
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
		.mask-box {
            border: 1px solid gray;
            padding: 5px;
            background-color: #f4f4f4;
            border-radius: 3px;
            font-size: 0.9em;
            text-align: center;
            width: 80px;
            height: 80px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
			gap: 5px;
		}
		.mask-box span {
			display: block; /* Each element will be on a new line */
		}
		.mask-box .utf8-mask {
			font-size: 1em;
		}
		.mask-box .binary {
			font-size: 1em;
		}
		.mask-box .utf8-remainder {
			font-size: 1em;
			text-align: right; /* Right-align remainder */
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
        <input type="text" name="input" id="inputField" placeholder="Enter text here" required>
        <button type="submit" class="submit-button">Submit</button>
        <button type="reset" 
            class="reset-button"
            onclick="resetForm();">
            Reset
        </button>
    </form>
    <div id="all-results">
        <div id="result" class="rune-container"></div>
        <div id="details" class="details-container"></div>
    </div>

    <script>
        function resetForm() {
            document.getElementById('result').innerHTML = '';
            document.getElementById('details').innerHTML = '';
            // Manually reset the input field and reapply the placeholder
            var inputField = document.getElementById('inputField');
            inputField.value = '';  // Clear the input value
            inputField.placeholder = 'Enter text here';  // Reset the placeholder
        }
    </script>
</body>
</html>
`))

var resultTmpl = template.Must(template.New("result").Parse(`
{{range .}}
	<div class="rune-box" 
		hx-get="/details?char={{.Char}}" 
		hx-target="#details" 
		hx-swap="innerHTML">
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
<div class="details-box">
    <h3>Rune: {{.Char}}</h3>
    <div class="bytes-info">
        {{range .RuneBytes}}
            <div class="mask-box">
				<span class="utf8-mask">{{.Utf8Mask}}</span>
				<span class="binary">{{.Binary}}</span>
				<span></span>
				<span class="utf8-remainder">{{.Utf8Remainder}}</span>
            </div>
        {{end}}
    </div>
    <div style="margin-top: 10px; font-weight: bold;">
        Unicode Code Point:
    </div>
    <div style="font-size: 1.2em; font-weight: bold;">
        {{.CodePoint}}
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

	http.HandleFunc("/reset", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(""))
		// w.Write([]byte(`<div id="result"></div><div id="details"></div>`))
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
	var mask string
	var remainder string

	rune := RuneInfo{
		Char:      string(r),
		CodePoint: fmt.Sprintf("%d", r),
	}

	utf8Bytes := []byte(string(r))
	utf8ByteBinary := fmt.Sprintf("%08b", utf8Bytes[0])

	n := len(utf8Bytes)
	if n > 1 {
		mask = utf8Mask[n]
		remainder = string(utf8ByteBinary[n+1:])
	} else {
		mask = utf8Mask[0]
		remainder = string(utf8ByteBinary[1:])
	}

	fmt.Println(remainder)

	for j, b := range utf8Bytes {
		utf8ByteBinary = fmt.Sprintf("%08b", b)
		if j > 0 {
			mask = utf8Mask[1]
			remainder = string(utf8ByteBinary[2:])
		}
		rune.RuneBytes = append(rune.RuneBytes, RuneByte{
			RuneByteIndex: j,
			Byte:          b,
			Binary:        utf8ByteBinary,
			Utf8Mask:      mask,
			Utf8Remainder: remainder,
		})
	}
	return rune
}
