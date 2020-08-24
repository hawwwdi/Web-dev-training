package main

import (
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
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
var sessionsMap map[string]string
var usersMap map[string]User

type User struct {
	Id, Password string
	IsAdmin      bool
}

func init() {
	sessionsMap = make(map[string]string)
	usersMap = make(map[string]User)
	adminUUID := uuid.Must(uuid.NewV4())
	sessionsMap[adminUUID.String()] = "admin"
	usersMap["admin"] = User{
		Id:       "admin",
		Password: "admin",
		IsAdmin:  true,
	}
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	fmt.Println("running...")
	mux := httprouter.New()
	//mux.Handler("GET", "/", http.FileServer(http.Dir("./templates")))
	mux.GET("/", index)
	mux.Handler("GET", "/favicon.ico", http.FileServer(http.Dir("")))
	mux.POST("/panel", postLogin)
	mux.GET("/panel", cookieLogin)
	mux.GET("/changePass", showChangePass)
	//	mux.POST("/", changePassword)
	mux.GET("/show/:pic", show)
	mux.POST("/show", showPic)
	//mux.Handler("GET", "/files/",http.StripPrefix("/files", http.FileServer(http.Dir("./"))))
	mux.ServeFiles("/files/*filepath", http.Dir("./"))
	err := http.ListenAndServe("localhost:8080", mux)
	handleErr(os.Stdout, err)
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if _, err := r.Cookie("session"); err == nil {
		http.Redirect(w, r, "/panel", http.StatusSeeOther)
		return
	}
	err1 := tpl.ExecuteTemplate(w, "index.html", nil)
	handleErr(w, err1)
}

func cookieLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie, _ := r.Cookie("session")
	UUID := cookie.Value
	username, exists := sessionsMap[UUID]
	if !exists {
		fmt.Fprintln(w, "invalid cookie!")
		return
	}
	user := usersMap[username]
	login(w, &user)
}

func postLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(64)
	handleErr(w, err)
	username, pass, rememberMe := r.PostFormValue("user"), r.PostFormValue("pass"), r.PostFormValue("rememberMe")
	//fmt.Println("user: ", user, " pass: ", pass)
	user, err2 := getUser(username, pass)
	if err2 != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if rememberMe == "true" {
		writeSession(w, username)
	}
	/*if user != USER || pass != PASS {
		handleErr(w, tpl.ExecuteTemplate(w, "index.html", true))
		return
	}*/
	login(w, user)
}

func login(w http.ResponseWriter, user *User) {
	if user.IsAdmin {
		http.SetCookie(w, &http.Cookie{
			Name:  "last-seen",
			Value: time.Now().Format("15:04:05"),
		})
		data := struct {
			User, Pass, LastSeen string
		}{user.Id, user.Password, "Now"}
		err := tpl.ExecuteTemplate(w, "panel.html", data)
		handleErr(w, err)
	} else {
		_, err := fmt.Fprintf(w, "welcome %v", user.Id)
		handleErr(w, err)
	}
}

func writeSession(w http.ResponseWriter, id string) {
	userUUID := uuid.Must(uuid.NewV4())
	sessionsMap[userUUID.String()] = id
	cookie := http.Cookie{
		Name:  "session",
		Value: userUUID.String(),
	}
	http.SetCookie(w, &cookie)
}

func getUser(username, pass string) (*User, error) {
	user, err := usersMap[username]
	if !err {
		return nil, errors.New("invalid username")
	}
	if user.Password != pass {
		return nil, errors.New("invalid password")
	}
	return &user, nil
}

/*func changePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	handleErr(w, r.ParseForm())
	PASS = r.PostFormValue("pass")
	handleErr(w, tpl.ExecuteTemplate(w, "index.html", nil))
}*/

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

func showPic(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, details, err := r.FormFile("file")
	handleErr(w, err)
	defer file.Close()
	bytes, err1 := ioutil.ReadAll(file)
	handleErr(w, err1)
	toSave, err2 := os.Create(filepath.Join("./files/", details.Filename))
	defer toSave.Close()
	handleErr(w, err2)
	_, err3 := toSave.Write(bytes)
	handleErr(w, err3)
	err = tpl.ExecuteTemplate(w, "show.html", string(bytes))
	handleErr(w, err)
}

func handleErr(w io.Writer, err error) {
	if err != nil {
		if _, err1 := io.Copy(w, strings.NewReader(err.Error())); err1 != nil {
			log.Fatalln(err1)
		}
	}
}
