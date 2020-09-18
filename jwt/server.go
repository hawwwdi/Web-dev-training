package main

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("request received")
	_, _ = fmt.Fprint(w, "hello world")
}

func isAuthorized(endpoint http.HandlerFunc) http.Handler {
	newEndpoint := func(w http.ResponseWriter, r *http.Request) {
		if tokenHeader, ok := r.Header["Token"]; ok {
			token, err := jwt.Parse(tokenHeader[0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("there was an error")
				}
				return key, nil
			})
			if err != nil {
				_, _ = fmt.Fprint(w, err.Error())
				return
			}
			if token.Valid {
				endpoint(w, r)
			}
		} else {
			_, _ = fmt.Fprint(w, "Not Authorized")
		}
	}
	return http.HandlerFunc(newEndpoint)
}