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

func add_salting(input string, transbits int) string {
	temp := ""
	for _, char := range input {
		text := rune(int(char) - transbits)
		temp += string(text)
	}
	return temp
}

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

func homepage(w http.ResponseWriter, r *http.Request) {
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
		t, _ := template.ParseFiles("src/index.html")
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

func login_failed(w http.ResponseWriter, r *http.Request) {
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
		t, _ := template.ParseFiles("src/failed.html")
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
		password := r.Form.Get("password") + add_salting(r.Form.Get("password"), 32)
		hashedPassword := hashPassword(password)
		fmt.Println("username      : ", username)
		fmt.Println("password      :", password)
		fmt.Println("hash password : ", hashedPassword)
		// exam password and username
		var dbPassword string
		err_exam := db.QueryRow("SELECT password FROM info WHERE username = ?", username).Scan(&dbPassword)
		if err_exam != nil {
			http.Redirect(w, r, "/loginfailed", http.StatusSeeOther)
			return
		}
		if hashedPassword == dbPassword {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, "/loginfailed", http.StatusSeeOther)
			return
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
		password := r.Form.Get("password") + add_salting(r.Form.Get("password"), 32)

		/*begin code*/
		confirm_password := r.Form.Get("confirm-password")
		if confirm_password != password {
			fmt.Println("Error password comfirm!!!!\n")
		}
		/*end code*/
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
			temp, err_src := template.ParseFiles("src/success.html")
			if err_src != nil {
				log.Fatalf("Error parsing success template: %s", err_src)
				return
			}
			err_src = temp.Execute(w, nil)
			if err_src != nil {
				log.Fatalf("Error executing success template: %s", err_src)
				return
			}
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
	http.HandleFunc("/home", homepage)
	http.HandleFunc("/loginfailed", login_failed)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		log.Fatalf("Error starting server %s", err)
		log.Println("=====================FAILED=====================")
	}
}
