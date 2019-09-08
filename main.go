package main

import (
	"github.com/rakin92/message-service/config"
	"github.com/rakin92/message-service/database"
	"github.com/rakin92/message-service/migration"
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

	email := r.Group("/email")
	email.POST("/send", email.SendEmail)

	sms := r.Group("/sms")
	sms.POST("/send", sms.SendSMS)

	log.Info("running on port 8081")
	r.Run(":8081")
}
