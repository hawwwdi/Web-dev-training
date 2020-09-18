package main

import (
	"log"
	"net/http"
)

const address = "localhost:8080"

var key = []byte("123456")

func main() {
	http.Handle("/hello", isAuthorized(helloWorld))
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
}
