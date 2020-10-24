package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/go-redis/redis"
	"github.com/julienschmidt/httprouter"
)

//var ctx = context.Background()
var rdb *redis.Client
var tpl *template.Template

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:9090",
	})
	tpl = template.Must(template.ParseGlob("templates/*.html"))
}

func main() {
	mux := httprouter.New()
	mux.GET("/", index)
	mux.POST("/set", set)
	mux.POST("/get", get)
	fmt.Println("socket is opened on port 4040")
	err := http.ListenAndServe("localhost:4040", mux)
	check(err, os.Stdout)
}

func index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	check(err, w)
}

func set(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(128)
	check(err, w)
	key, value := r.PostFormValue("key"), r.PostFormValue("value")
	err = rdb.Set(key, value, 0).Err()
	check(err, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseMultipartForm(128)
	check(err, w)
	key := r.PostFormValue("key")
	value, err1 := rdb.Get(key).Result()
	//check(err1, w)
	if err1 != nil {
		fmt.Fprintf(w, "key not found")
		return
	}
	fmt.Fprintf(w, "key %v = %v \n", key, value)
}

func check(err error, w io.Writer) {
	if err != nil {
		fmt.Fprintln(w, err)
		panic(err)
	}
}
