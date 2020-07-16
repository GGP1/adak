/*
Package email helps us to use the email as the tool to identify each user
*/
package email

import (
	"log"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendValidation sends a validation email to the user
func SendValidation(user model.User, token string) {
	// Mail content
	from := mail.NewEmail("Palo API", cfg.EmailSender)
	subject := "Validation email"
	to := mail.NewEmail(user.Firstname+" "+user.Lastname, user.Email)
	plainTextContent := "Palo dev team"
	htmlContent := "<h2>Palo email validation</h2><br><h4>Thanks for joining us, " + user.Firstname + " " + user.Lastname + "!</h4><br><p>Please validate your account by clicking the following link:</p><br><a>http://localhost:4000/email/" + token + "</a>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(cfg.SengridKey)

	// Send mail
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}
}
