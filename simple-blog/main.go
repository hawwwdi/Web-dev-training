package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
	"net/http"
)

var tpl *template.Template

const (
	USER = "admin"
	PASS = "admin"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	mux := httprouter.New()
	mux.GET("/", index)
	mux.POST("/", login)
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "login.html", nil)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(64)
	if err != nil {
		log.Println(err)
		fmt.Fprintln(w, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, pass := r.PostFormValue("user"), r.PostFormValue("pass")
	fmt.Println("user: ", user, " pass: ", pass)
	if user != USER || pass != PASS {
		tpl.ExecuteTemplate(w, "login.html", true)
		return
	}
	data := struct {
		User, Pass string
	}{user, pass}
	tpl.ExecuteTemplate(w, "panel.html", data)
}
