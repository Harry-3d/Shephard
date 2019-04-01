package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var key = []byte("super-secret-key")
var store = sessions.NewCookieStore(key)

type TmplData struct {
	Attempts int
}

var data TmplData

func main() {

	data.Attempts = 3

	// Dynamic pages
	http.HandleFunc("/", login)
	http.HandleFunc("/logout", logout)

	// Static resources
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/css"))))

	http.ListenAndServe(":8080", nil)
}

func internal(w http.ResponseWriter, r *http.Request) {
	// get the cookie for this session
	session, _ := store.Get(r, "cookie-name")

	// check if user is authenticated
	auth, ok := session.Values["authenticated"].(bool)
	if auth == false || ok == false {

		// dec Attempts remaining here, stay on login page.
		data.Attempts--
		fmt.Printf("Attempts:%d\n", data.Attempts)

		http.Redirect(w, r, "/", 303)

	} else {
		// user is authenticated
		// internal page
		t, err := template.ParseFiles("internal.html")
		if err != nil {
			log.Printf("internal.html page err:%s", err)
		} else {
			err := t.Execute(w, nil)
			if err != nil {
				log.Printf("internal tmpl exe err:%s\n", err)
			}
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if data.Attempts < 0 {
		http.Error(w, "Access Denied", http.StatusForbidden)
		return
	}

	session, _ := store.Get(r, "cookie-name")

	//read the submitted password
	if r.Method == "POST" {
		val := r.FormValue("pass")
		log.Printf("pass:%s\n", val)

		// if correct password
		if val == "password" {
			session.Values["authenticated"] = true
			session.Save(r, w)

		}
		internal(w, r)
		return
	}
	if r.Method == "GET" {
		t, err := template.ParseFiles("login.html")
		if err != nil {
			log.Printf("login page err:%s", err)
		} else {
			err := t.Execute(w, &data)
			if err != nil {
				log.Printf("login tmpl exe err:%s\n", err)
			}
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", 303)
	data.Attempts = 3
}
