package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rakin92/message-service/config"
	"github.com/rakin92/message-service/database"
	"github.com/rakin92/message-service/email"
	"github.com/rakin92/message-service/migration"
	"github.com/rakin92/message-service/sms"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Print("Message service starting up.")

	config.SetupEnv()

	db, err := database.ConnectDB()

	if err != nil {
		log.Panic(err)
	}

	db = migration.Migrate(db)

	defer db.Close()

	r := config.SetupRouter()

	emailAuthorized := r.Group("/email")

	emailAuthorized.POST("/send", email.SendEmail)

	smsAuthorized := r.Group("/sms", gin.BasicAuth(gin.Accounts{
		os.Getenv("USER"): os.Getenv("PASSWORD"),
	}))

	smsAuthorized.POST("/send", sms.SendSMS)

	log.Info("running on port 8081")
	r.Run(":8081")
}
