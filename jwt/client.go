package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type httpClient struct {
	client http.Client
	sync.Mutex
}

var client = &httpClient{}

func makeReq(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := generateJWT()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}
	req, _ := http.NewRequest("GET", address, nil)
	req.Header.Set("Token", jwtToken)
	client.Lock()
	res, err1 := client.client.Do(req)
	client.Unlock()
	if err1 != nil {
		fmt.Fprint(w, err1.Error())
		return
	}
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Fprint(w, string(body))
}

func generateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = "hawwwdi"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	fmt.Printf("token is: %v", tokenString)
	return tokenString, nil
}