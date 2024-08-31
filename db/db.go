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
	DB, err = gorm.Open(sqlite.Open("files.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.File{})
}

func SaveFileToDb(file *models.File) error {
	mu.Lock()
	defer mu.Unlock()
	result := DB.Create(&file)
	return result.Error
}

func GetFilesFromDb(id uint) (*[]models.File, error) {
	mu.Lock()
	defer mu.Unlock()
	var files []models.File
	result := DB.Where(&models.File{OwnerID: id}).Find(&files)

	return &files, result.Error
}

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

func DeleteFile(id uint) error {
	mu.Lock()
	defer mu.Unlock()
	result := DB.Unscoped().Delete(&models.File{}, id)
	return result.Error
}
