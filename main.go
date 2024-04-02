package main

import (
	"crypto/sha1"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	PORT = ":8081"
)

func hashPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))
	return hashedPassword
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("src/login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		fmt.Println("username: ", username)
		fmt.Println("password: ", password)
		hashedPassword := hashPassword(password)
		if hashedPassword == "0f224bdbd25c5c3532037fed583e6240f600174f" {
			fmt.Fprintf(w, "Hello %s, welcome to our website!", username)
		} else {
			fmt.Fprintf(w, "Invalid username or password")
		}
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("src/signup.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

func main() {
	log.Println("Server is running on port", PORT)
	log.Println("=====================SUCCESS=====================")
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("src/img"))))
	http.HandleFunc("/", login)
	http.HandleFunc("/signup", signup)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatalf("Error starting server %s", err)
		log.Println("=====================FAILED=====================")
	}
}
