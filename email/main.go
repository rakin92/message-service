package email

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakin92/message-service/database"
	"github.com/rakin92/message-service/models"
	log "github.com/sirupsen/logrus"
)

// SendEmail : sends an email
func SendEmail(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Error("Error reading request body")
	}

	var newEmail = models.Email{
		Status: "PENDING",
	}

	jsonErr := json.Unmarshal(body, &newEmail)

	if jsonErr != nil {
		log.Error(jsonErr)
	}

	database.DB.Create(&newEmail)

	defer c.Request.Body.Close()

	err = newEmail.SendSendgridEmail()

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, "Successfully sent email")
	}
}
