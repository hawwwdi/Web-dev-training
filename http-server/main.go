package main

import (
	"bufio"
	"fmt"
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
	{
		mux := http.NewServeMux()
		mux.HandleFunc("/", serveHTTP)
		mux.HandleFunc("/hello/", sayHello)
		mux.HandleFunc("/redirect", redirect)
		mux.Handle("/redirects", http.RedirectHandler("/hello/", http.StatusMovedPermanently))
		_ = http.ListenAndServe(":8080", mux)
	}
	/*if err != nil {
		log.Fatal(err)
	}*/

}

func redirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("redirect to /hello & method = ", r.Method)
	w.Header().Set("Location", "/hello/")
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//w.Write([]byte("Location: /hello"))
	w.WriteHeader(http.StatusMovedPermanently)
	return
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:  "say",
		Value: "hello",
	}
	http.SetCookie(w, &cookie)
	fmt.Println("this is domain: ", cookie.Domain)
	fmt.Println(r.Method)
	fmt.Fprintln(w, "hello world :|")
	fmt.Fprintln(w, r.Cookies())
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	cookie := http.Cookie{
		Name:  "wolf",
		Value: "hehehe",
	}
	http.SetCookie(w, &cookie)
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
