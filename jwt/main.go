package main

import (
	"log"
	"net/http"
)

var key = []byte("123456")

func main() {
	http.Handle("/hello", isAuthorized(helloWorld))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
