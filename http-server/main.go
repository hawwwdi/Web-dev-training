package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

type index struct {
	tpl *template.Template
}

func main() {
	template := &index{
		tpl: template.Must(template.ParseFiles("templates/index.html")),
	}
	err := http.ListenAndServe(":8080", template)
	if err != nil {
		log.Fatal(err)
	}
}

func (in *index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	fmt.Println(data)
	in.tpl.Execute(w, data)
}
