package api

import "html/template"

var tmpl = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta name="format-detection" content="telephone=no">
    <title>Rune Seer</title>
    <script src="https://unpkg.com/htmx.org@1.9.5"></script>
	<link rel="stylesheet" href="/public/styles.css">
    <link rel="icon" type="image/svg+xml" href="/public/triangle_eye.svg">
</head>
<body>
    <h1>Rune Seer</h1>
    <p>Enter text to see the UTF-8 encoding of each rune.</p>
    <form hx-post="/analyze" 
			hx-target="#result" 
			hx-swap="innerHTML" 
			hx-on::before-request="document.getElementById('details').innerHTML = ''">
        <input type="text" name="input" id="inputField" placeholder="Enter text here" required>
        <button type="submit" class="submit-button">Submit</button>
        <button type="reset" class="reset-button" onclick="resetForm();">Reset</button>
    </form>
	<div id="result" class="rune-container"></div>
	<div id="details" class="details-container"></div>

    <script>
        function resetForm() {
            document.getElementById('result').innerHTML = '';
            document.getElementById('details').innerHTML = '';
            var inputField = document.getElementById('inputField');
            inputField.value = ''; 
            inputField.placeholder = 'Enter text here';  
        }
    </script>
</body>
</html>
`))

var resultTmpl = template.Must(template.New("result").Parse(`
{{range .}}
    <form hx-post="/details" 
          hx-target="#details" 
          hx-swap="innerHTML" 
          class="rune-form">
        <input type="hidden" name="char" value="{{.Char}}">
        <button type="submit" class="rune-box">
            <div class="char">{{.Char}}</div>
            <div class="bytes-info">
                {{range .RuneBytes}}
                    <div class="byte-box">
                        {{.Binary}}
                    </div>
                {{end}}
            </div>
        </button>
    </form>
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
	<span class="code-point-text">Unicode Code Point: </span>
    <span class="code-point-number">{{.CodePoint}}</span>
</div>
`))
