package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
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
	handleErr(w, err)
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(64)
	handleErr(w, err)
	user, pass := r.PostFormValue("user"), r.PostFormValue("pass")
	fmt.Println("user: ", user, " pass: ", pass)
	if user != USER || pass != PASS {
		tpl.ExecuteTemplate(w, "login.html", true)
		return
	}
	data := struct {
		User, Pass string
	}{user, pass}
	err = tpl.ExecuteTemplate(w, "panel.html", data)
	handleErr(w, err)
}

func handleErr(w io.Writer, err error){
	if err != nil{
		io.Copy(w, strings.NewReader(err.Error()))
	}
}