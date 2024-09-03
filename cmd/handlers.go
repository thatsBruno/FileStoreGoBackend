package main

import (
	"encoding/json"
	"go-api/db"
	"go-api/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func FileHandlers(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "file/"):
		// Files
		downloadFile(w, r)
	case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "files/"):
		getFiles(w, r)
	case r.Method == http.MethodPost:
		saveFile(w, r)
	case r.Method == http.MethodDelete:
		deleteFile(w, r)
	default:
	}
}

// Files
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

// takes request as form-data
func saveFile(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("Name")
	log.Println("File Name:", name)

	ownerId, err := strconv.Atoi(r.FormValue("OwnerId"))
	if err != nil {
		log.Println("Error converting OwnerId:", err)
		http.Error(w, "Invalid OwnerId", http.StatusBadRequest)
		return
	}
	log.Println("Owner ID:", ownerId)

	file, _, err := r.FormFile("File")
	if err != nil {
		log.Println("Error retrieving file:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("File size:", len(fileBytes))

	fileRecord := models.File{
		FileName: name,
		Data:     fileBytes,
		OwnerID:  uint(ownerId),
	}

	if err := db.SaveFileToDb(&fileRecord); err != nil {
		log.Println("Error saving file to database:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully!"))
}

func getFiles(w http.ResponseWriter, r *http.Request) {
	ownerId, err := strconv.Atoi(r.PathValue("ownerId"))
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

func downloadFile(w http.ResponseWriter, r *http.Request) {
	fileId, err := strconv.Atoi(r.PathValue("fileId"))
	if err != nil {
		http.Error(w, "Missing file name", http.StatusBadRequest)
		return
	}

	// Retrieve the file record from the database
	file, err := db.GetFileById(uint(fileId))
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Set the appropriate headers for file download
	w.Header().Set("Content-Disposition", "attachment; filename="+file.FileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(file.Data)))

	// Write the file data to the response
	if _, err := w.Write(file.Data); err != nil {
		http.Error(w, "Error writing file to response", http.StatusInternalServerError)
		return
	}
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
