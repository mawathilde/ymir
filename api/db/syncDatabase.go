package db

import "ymir/api/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
