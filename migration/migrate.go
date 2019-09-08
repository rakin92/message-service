package migration

import (
	"github.com/jinzhu/gorm"
	"github.com/rakin92/message-service/models"
)

// Migrate : will create and migrate the tables, and then make the some relationships if necessary
func Migrate(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&models.Email{}, &models.SMS{})
	return db
}
