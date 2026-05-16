package handlers

import (
	"net/http"
	"strconv"
)

func getUserFromCookie(r *http.Request) (int, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return 0, err
	}
	userID, err := strconv.Atoi(cookie.Value)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	if _, err := getUserFromCookie(r); err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "web/index.html")
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	if _, err := getUserFromCookie(r); err == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.ServeFile(w, r, "web/login.html")
}
