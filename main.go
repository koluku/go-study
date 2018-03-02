package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

// temp1は一つのテンプレートを表します
type templateHandler struct {
	once     sync.Once
	filename string
	temp1    *template.Template
}

// ServeHTTPはHTTPリクエストを処理します
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp1 = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.temp1.Execute(w, nil)
}

func main() {
	// ルート
	http.Handle("/", &templateHandler{filename: "chat.html"})
	// Webサーバーを開始します
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
