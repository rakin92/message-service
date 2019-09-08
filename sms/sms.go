package sms

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakin92/message-service/database"
	"github.com/rakin92/message-service/models"
	log "github.com/sirupsen/logrus"
)

// SendSMS : sends an email
func SendSMS(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Error("Error reading request body")
	}

	var newSMS = models.SMS{
		Status: "PENDING",
	}

	jsonErr := json.Unmarshal(body, &newSMS)

	if jsonErr != nil {
		log.Error(jsonErr)
	}

	database.DB.Create(&newSMS)

	defer c.Request.Body.Close()

	err = newSMS.SendTwilioSMS()

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "Successfully sent email")
	}
}
