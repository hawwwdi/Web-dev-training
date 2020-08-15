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

var USER = "admin"
var PASS = "admin"

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	fmt.Println("run 1")
	mux := httprouter.New()
	mux.GET("/", index)
	mux.POST("/panel", login)
	mux.GET("/changePass", showChangePass)
	mux.POST("/", changePassword)
	mux.GET("/show/:pic", show)
	mux.Handler("GET", "/files/",http.StripPrefix("/files", http.FileServer(http.Dir("./"))))
	err := http.ListenAndServe("localhost:8089", mux)
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
	//fmt.Println("user: ", user, " pass: ", pass)
	if user != USER || pass != PASS {
		handleErr(w, tpl.ExecuteTemplate(w, "login.html", true))
		return
	}
	data := struct {
		User, Pass string
	}{user, pass}
	err = tpl.ExecuteTemplate(w, "panel.html", data)
	handleErr(w, err)
}

func changePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	handleErr(w, r.ParseForm())
	PASS = r.PostFormValue("pass")
	handleErr(w, tpl.ExecuteTemplate(w, "login.html", nil))
}

func showChangePass(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	handleErr(w, tpl.ExecuteTemplate(w, "change.html", nil))
}

func show(w http.ResponseWriter, r *http.Request, sp httprouter.Params) {
	picName := sp.ByName("pic")
	//pic, err := os.Open(picName)
	/*if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "image %v not found !", picName)
		return
	}
	defer pic.Close()*/
	//details , exception := pic.Stat()
	//handleErr(w, exception)
	http.ServeFile(w, r, picName)
}

func handleErr(w io.Writer, err error) {
	if err != nil {
		io.Copy(w, strings.NewReader(err.Error()))
	}
}
