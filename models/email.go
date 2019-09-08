package models

import (
	"github.com/jinzhu/gorm"
	"github.com/rakin92/message-service/database"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/spf13/viper"

	// gorm postgres dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"
)

// Email : Base with injected fields `ID`, `CreatedAt`, `UpdatedAt`
type Email struct {
	gorm.Model
	Status      string             `grom:"type:varchar(32);not null"`
	TemplateID  string             `grom:"type:varchar(256)"`
	ToEmai      string             `grom:"type:varchar(256);not null"`
	FromEmail   string             `grom:"type:varchar(256);not null"`
	ToName      string             `grom:"type:varchar(256);not null"`
	FromName    string             `grom:"type:varchar(256);not null"`
	Subject     string             `grom:"type:varchar(256)"`
	Message     string             `grom:"type:varchar"`
	DynamicData *map[string]string `grom:"type:json"`
}

func (e Email) plainTextEmailBody() []byte {
	from := mail.NewEmail(e.FromName, e.FromEmail)
	to := mail.NewEmail(e.ToName, e.ToEmai)
	message := mail.NewContent("text/plain", e.Message)
	email := mail.NewV3MailInit(from, e.Subject, to, message)

	return mail.GetRequestBody(email)
}

func (e Email) dynamicTemplateEmailBody() []byte {
	newMail := mail.NewV3Mail()

	emailFrom := mail.NewEmail(e.FromName, e.FromEmail)
	newMail.SetFrom(emailFrom)

	newMail.SetTemplateID(e.TemplateID)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(e.ToName, e.ToEmai),
	}
	p.AddTos(tos...)

	if e.DynamicData != nil {
		for key, value := range *e.DynamicData {
			log.Infof("%s %s", key, value)
			p.SetDynamicTemplateData(key, value)
		}
	}

	newMail.AddPersonalizations(p)
	return mail.GetRequestBody(newMail)
}

// SendSendgridEmail : sends a email
func (e Email) SendSendgridEmail() error {
	request := sendgrid.GetRequest(viper.GetString("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"

	if e.TemplateID != "" {
		request.Body = e.dynamicTemplateEmailBody()
	} else {
		request.Body = e.plainTextEmailBody()
	}

	_, err := sendgrid.API(request)

	if err != nil {
		e.Status = "FAILED"
		database.DB.Save(&e)
		return err
	}
	e.Status = "SENT"
	database.DB.Save(&e)

	return nil
}
