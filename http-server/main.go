package main

import (
	"bufio"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("templates/index.html"))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveHTTP)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	writer := bufio.NewWriter(w)
	if err != nil {
		writer.WriteString(err.Error())
		log.Println(err)
	}
	headers := w.Header()
	headers.Set("name", "test")
	//w.WriteHeader(http.StatusNotAcceptable)
	data := struct {
		Method  string
		URL     *url.URL
		Forms   map[string][]string
		Headers http.Header
	}{
		r.Method,
		r.URL,
		r.Form,
		r.Header,
	}
	tpl.Execute(w, data)
}
