package main

import (
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

	emailAPI := r.Group("/email")
	emailAPI.POST("/send", email.SendEmail)

	smsAPI := r.Group("/sms")
	smsAPI.POST("/send", sms.SendSMS)

	log.Info("running on port 8081")
	r.Run(":8081")
}
