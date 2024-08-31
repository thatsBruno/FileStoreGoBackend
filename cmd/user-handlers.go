package main

import (
	"encoding/json"
	"go-api/db"
	"go-api/models"
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
	fileid, err := strconv.Atoi(r.PathValue("fileid"))
	if err != nil {
		http.Error(w, "Invalid id passed", http.StatusBadRequest)
		return
	}

	if err := db.DeleteFile(uint(fileid)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func saveFile(w http.ResponseWriter, r *http.Request) {
	var file models.File
	if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.SaveFileToDb(&file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(file)
	w.WriteHeader(http.StatusCreated)
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
		log.Printf("Error decoding JSON: %v", err)
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
