package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	// gorm postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// SMS : Base with injected fields `ID`, `CreatedAt`, `UpdatedAt`
type SMS struct {
	gorm.Model
	Status     string `grom:"type:varchar(32);not null"`
	ToNumber   string `grom:"type:varchar(256);not null"`
	FromNumber string `grom:"type:varchar(256);not null"`
	Message    string `grom:"type:varchar(256);not null"`
	TwilioID   string `grom:"type:varchar(256)"`
}

func (s SMS) prepareMessage() []byte {
	// Pack up the data for our message
	msgData := url.Values{}
	msgData.Set("To", s.ToNumber)
	msgData.Set("From", s.FromNumber)
	msgData.Set("Body", s.Message)
	msgDataReader := *strings.NewReader(msgData.Encode())
}

// SendTwilioSMS : sends a email
func (s SMS) SendTwilioSMS() error {
	accountSid := viper.GetString("SENDGRID_API_KEY")
	authToken := viper.GetString("SENDGRID_API_KEY")

	urlStr := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSid)

	// Create HTTP request client
	message := s.prepareMessage()
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &message)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, err := client.Do(req)

	if err != nil {
		s.Status = "FAILED"
		database.DB.Save(&s)
		return err
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)

		if err == nil {
			s.Status = "SENT"
			s.TwilioID = data["sid"])
			database.DB.Save(&s)
		}
	}
	return nil
}
