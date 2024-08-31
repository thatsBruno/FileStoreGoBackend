package main

import (
	"encoding/json"
	"go-api/db"
	"go-api/models"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func UserHandlers(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		getUser(w, r)
	case r.Method == http.MethodPost:
		createUser(w, r)
	case r.Method == http.MethodDelete:
		deleteUser(w, r)
	default:
	}
}

// Files
func FileHandlers(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		getFiles(w, r)
	case r.Method == http.MethodPost:
		saveFile(w, r)
	case r.Method == http.MethodDelete:
		deleteFile(w, r)
	default:
	}
}

func deleteFile(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func saveFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusBadRequest)
		return
	}

	var file models.File

	fileData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	file.FileName = r.FormValue("file_name")
	file.Data = fileData
	file.OwnerID = 2

	if err := db.SaveFileToDb(&file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully"))
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	ownerId, err := strconv.Atoi(r.PathValue("ownerid"))
	if err != nil {
		http.Error(w, "Invalid id passed", http.StatusBadRequest)
		return
	}
	files, err := db.GetFilesFromDb(uint(ownerId))

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

// Users
func createUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusCreated)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid id passed", http.StatusBadRequest)
		return
	}
	user, err := db.GetUserByID(uint(userId))

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid id passed", http.StatusBadRequest)
		return
	}

	if err := db.DeleteUser(uint(userId)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
