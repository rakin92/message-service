package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	// gorm postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
	mocket "github.com/selvatico/go-mocket"
	log "github.com/sirupsen/logrus"
)

// DB *grom.DB
var DB *gorm.DB

// ConnectDB : connecting DB
func ConnectDB() (*gorm.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString("DB_HOST"),
		viper.GetInt("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_SSL_MODE"))

	conn, err := gorm.Open("postgres", connectionString)

	if err != nil {
		log.Panic(err)
	}

	if viper.GetString("env") != "PROD" {
		conn.LogMode(true)
	}

	DB = conn

	return DB, nil
}

// GetMockDB : returns a mock db for test purpose
func GetMockDB() (*gorm.DB, error) {

	mocket.Catcher.Register() // Safe register. Allowed multiple calls to save
	mocket.Catcher.Logging = true

	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		viper.GetString("DB_HOST"),
		viper.GetInt("DB_PORT"),
		viper.GetString("DB_USER"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_SSL_MODE"))

	conn, err := gorm.Open(mocket.DriverName, connectionString)

	if err != nil {
		log.Panic(
			fmt.Sprintf("Error connecting to database: %v", err),
			"db.init",
		)
	}

	DB = conn
	return DB, nil
}
