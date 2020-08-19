package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	mux.Handler("GET", "/", http.FileServer(http.Dir("./templates")))
	mux.Handler("GET", "/favicon.ico", http.FileServer(http.Dir("")))
	mux.POST("/panel", login)
	mux.GET("/changePass", showChangePass)
	mux.POST("/", changePassword)
	mux.GET("/show/:pic", show)
	mux.POST("/show", showPic)
	//mux.Handler("GET", "/files/",http.StripPrefix("/files", http.FileServer(http.Dir("./"))))
	mux.ServeFiles("/files/*filepath", http.Dir("./"))
	err := http.ListenAndServe("localhost:8089", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	handleErr(w, err)
}

func login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(64)
	handleErr(w, err)
	user, pass := r.PostFormValue("user"), r.PostFormValue("pass")
	//fmt.Println("user: ", user, " pass: ", pass)
	if user != USER || pass != PASS {
		handleErr(w, tpl.ExecuteTemplate(w, "index.html", true))
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: "last-seen",
		Value: time.Now().Format("15:04:05"),
	})
	var lastseen string
	lastSeen, _:= r.Cookie("last-seen")
	if lastSeen != nil{
		lastseen = lastSeen.Value
	} else {
		lastseen = "Now"
	}
	data := struct {
		User, Pass, LastSeen string
	}{user, pass, lastseen}
	err = tpl.ExecuteTemplate(w, "panel.html", data)
	handleErr(w, err)
}

func changePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	handleErr(w, r.ParseForm())
	PASS = r.PostFormValue("pass")
	handleErr(w, tpl.ExecuteTemplate(w, "index.html", nil))
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

func showPic(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	file, details, err := r.FormFile("file")
	handleErr(w, err)
	defer file.Close()
	bytes, err1 := ioutil.ReadAll(file)
	handleErr(w, err1)
	toSave, err2 := os.Create(filepath.Join("./files/", details.Filename))
	defer toSave.Close()
	handleErr(w, err2)
	toSave.Write(bytes)
	err = tpl.ExecuteTemplate(w, "show.html", string(bytes))
	handleErr(w, err)
}

func handleErr(w io.Writer, err error) {
	if err != nil {
		if _, err1 := io.Copy(w, strings.NewReader(err.Error())); err1 != nil{
			log.Fatalln(err1)
		}
	}
}
