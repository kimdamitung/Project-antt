package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
)

const (
	PORT = ":8081"
)

func connectDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func hashPassword(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	hashedPassword := fmt.Sprintf("%x", h.Sum(nil))
	return hashedPassword
}

func login(w http.ResponseWriter, r *http.Request) {
	// connect database
	db, err := connectDatabase("database/database.db")
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
		return
	}
	defer db.Close()
	// connect http
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
	// closer database
	err = db.Close()
	if err != nil {
		log.Fatalf("Error closing database: %s", err)
		return
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	// connect database
	db, err := connectDatabase("database/database.db")
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
		return
	}
	// connect http
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("src/signup.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
	// closer database
	err = db.Close()
	if err != nil {
		log.Fatalf("Error closing database: %s", err)
		return
	}
}

func main() {
	// connect database
	db, er := connectDatabase("database/database.db")
	checkErr(er)
	defer db.Close()

	_, er = db.Exec("INSERT INTO info (mssv, name, username, password) VALUES (?, ?, ?, ?)", 21085062, "Duy Tung", "duytung", "0f224bdbd25c5c3532037fed583e6240f600174f")
	checkErr(er)

	fmt.Println("Thêm dữ liệu thành công!")

	rows, er := db.Query("SELECT * FROM info")
	checkErr(er)
	defer rows.Close()
	for rows.Next() {
		var mssv int
		var name string
		var username string
		var password string
		er = rows.Scan(&mssv, &name, &username, &password)
		checkErr(er)
		fmt.Printf("MSSV: %d, Name: %s, Username: %s, Password: %s\n", mssv, name, username, password)
	}
	// connect  http web server
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
