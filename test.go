package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() { // Open the SQLite database
	db, err := sql.Open("sqlite3", "./activities.db")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer db.Close()
	const createTable = `CREATE TABLE IF NOT EXISTS files (id INTEGER PRIMARY KEY, filename TEXT, contents BLOB);`

	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileData := []byte("Hello, World!")
	_, err = db.Exec(`INSERT INTO files (filename, contents) VALUES (?, ?)`, "example.txt", fileData)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func banana() {
	http.HandleFunc("/api/files", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusBadRequest)
			return
		} // Get the file data from the request body
		var fileData []byte
		r.Body.Read(fileData)

		db, err := sql.Open("sqlite3", "./activities.db")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		_, err = db.Exec(`INSERT INTO files (filename, contents) VALUES (?, ?)`, "example.txt", fileData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	})
}
