package db

import (
	"log"
	"sync"

	"go-api/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	mu sync.Mutex
)

func InitDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the schema
	DB.AutoMigrate(&models.User{})
}

// TODO
// func UploadFileToDb(){}
// func DownloadFileFromDb(){}

func CreateUser(user *models.User) error {
	mu.Lock()
	defer mu.Unlock()
	result := DB.Create(&user)
	return result.Error
}

func GetUserByID(id uint) (*models.User, error) {
	mu.Lock()
	defer mu.Unlock()
	var user models.User
	result := DB.First(&user, id)
	return &user, result.Error
}

func DeleteUser(id uint) error {
	mu.Lock()
	defer mu.Unlock()
	result := DB.Delete(&models.User{}, id)
	return result.Error
}
