package main

import (
	"errors"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func alreadySignedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false
	}
	_, ok := sessionsMap[cookie.Value]
	return ok
}

func checkUser(username, pass string) (*User, error) {
	user, err := usersMap[username]
	if !err {
		return nil, errors.New("invalid username")
	}
	if bcrypt.CompareHashAndPassword(user.Password, []byte(pass)) != nil {
		return nil, errors.New("invalid password")
	}
	return &user, nil
}

func getUser(r *http.Request) *User {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}
	userId := sessionsMap[cookie.Value]
	user := usersMap[userId]
	return &user
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

func logOut(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
