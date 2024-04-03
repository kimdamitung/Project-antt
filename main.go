package main

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"html/template"
	"log"
	"net/http"
	"strconv"
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
	db, errdb := connectDatabase("database/database.db")
	if errdb != nil {
		log.Fatalf("Error connecting to database: %s", errdb)
		return
	}
	defer db.Close()
	// connect http
	fmt.Println("method: ", r.Method)
	if r.Method == "GET" {
		t, _ := template.ParseFiles("src/signup.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		mssvStr := r.Form.Get("mssv")
		mssv, err := strconv.Atoi(mssvStr)
		if err != nil {
			log.Fatalf("Error converting mssv to int: %s", err)
			return
		}
		name := r.Form.Get("name")
		username := r.Form.Get("username")
		password := r.Form.Get("password")
		hashedPassword := hashPassword(password)
		result, err := db.Exec("INSERT INTO info (mssv, name, username, password) VALUES (?, ?, ?, ?)", mssv, name, username, hashedPassword)
		if err != nil {
			log.Fatalf("Error inserting into database: %s", err)
			return
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatalf("Error getting rows affected: %s", err)
			return
		}
		if rowsAffected > 0 {
			fmt.Println("Insert successful")
		} else {
			fmt.Println("Insert failed")
		}
	}
	// closer database
	erclose := db.Close()
	if erclose != nil {
		log.Fatalf("Error closing database: %s", erclose)
		return
	}
}

func main() {
	// connect database
	db, er := connectDatabase("database/database.db")
	checkErr(er)
	defer db.Close()

	// display database
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
