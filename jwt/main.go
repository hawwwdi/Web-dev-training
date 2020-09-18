package main

import (
	"log"
	"net/http"
)

var key = []byte("123456")

func main() {
	http.Handle("/res", isAuthorized(helloWorld))
	http.HandleFunc("/req", makeReq)
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
