package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

var urls = make(map[string]string)

func main() {
	http.HandleFunc("/", echo)
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/go/", redirectHandler)
	fmt.Println("http://localhost:8080")
	loadFromFile()
	http.ListenAndServe(":8080", nil)

}

func echo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="ru">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>URL Shortener</title>
	<style>
		* {
			box-sizing: border-box;
			margin: 0;
			padding: 0;
			font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
		}
		body {
			background: linear-gradient(135deg, #1e293b 0%, #0f172a 100%);
			height: 100vh;
			display: flex;
			justify-content: center;
			align-items: center;
			color: #f8fafc;
		}
		.card {
			background: #1e293b;
			border: 1px solid #334155;
			padding: 40px;
			border-radius: 16px;
			box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.3);
			width: 100%;
			max-width: 480px;
			text-align: center;
		}
		h1 {
			font-size: 24px;
			margin-bottom: 8px;
			font-weight: 700;
		}
		p {
			color: #94a3b8;
			font-size: 14px;
			margin-bottom: 24px;
		}
		form {
			display: flex;
			gap: 10px;
		}
		input[type="text"] {
			flex: 1;
			padding: 12px 16px;
			border-radius: 8px;
			border: 1px solid #475569;
			background: #0f172a;
			color: #f8fafc;
			font-size: 15px;
			outline: none;
			transition: border-color 0.2s;
		}
		input[type="text"]:focus {
			border-color: #3b82f6;
		}
		input.input-error {
			border-color: #ef4444;
		}
		button {
			padding: 12px 24px;
			border-radius: 8px;
			border: none;
			background: #3b82f6;
			color: white;
			font-size: 15px;
			font-weight: 600;
			cursor: pointer;
			transition: background 0.2s;
		}
		button:hover {
			background: #2563eb;
		}
		.error-text {
			color: #ef4444;
			font-size: 13px;
			margin-top: 10px;
			text-align: left;
			display: none;
		}
	</style>
</head>
<body>
	<div class="card">
		<h1>Сокращатель ссылок</h1>
		<p>Вставьте длинный URL, чтобы получить компактную ссылку</p>
		<form id="shorten-form" action="/save" method="GET">
			<input type="text" id="url-input" name="url" placeholder="https://example.com...">
			<button type="submit">Сократить</button>
		</form>
		<div id="error-msg" class="error-text">Пожалуйста, вставьте ссылку перед отправкой</div>
	</div>

	<script>
		const form = document.getElementById('shorten-form');
		const input = document.getElementById('url-input');
		const errorMsg = document.getElementById('error-msg');

		form.addEventListener('submit', function(e) {
			if (input.value.trim() === '') {
				e.preventDefault();
				errorMsg.style.display = 'block';
				input.classList.add('input-error');
			}
		});

		input.addEventListener('input', function() {
			errorMsg.style.display = 'none';
			input.classList.remove('input-error');
		});
	</script>
</body>
</html>`)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	longURL := r.URL.Query().Get("url")
	if longURL == "" {
		http.Error(w, "введите ссылку ", http.StatusBadRequest)
		return
	}

	key := generateKey()
	urls[key] = longURL
	saveToFile()
	fmt.Fprintf(w, `ссылка сохранена по коду: <a href="/go/%s"><pre>http://localhost:8080/go/%s</pre></a>`, key, key)

}

func redirectHandler(w http.ResponseWriter, r *http.Request) {

	key := r.URL.Path[len("/go/"):]
	longURL := urls[key]
	http.Redirect(w, r, longURL, http.StatusFound)

}

func generateKey() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]

	}
	return string(b)
}

func saveToFile() {
	data, _ := json.Marshal(urls)
	os.WriteFile("urls.json", data, 0644)
}

func loadFromFile() {
	data, err := os.ReadFile("urls.json")
	if err != nil {
		return
	}
	json.Unmarshal(data, &urls)
}
