/*
Package email contains the methods relevant to email confirmation
*/
package email

import (
	"log"
	"os"

	"github.com/GGP1/palo/internal/utils/env"
	"github.com/GGP1/palo/pkg/model"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Confirmation sends a confirmation email to the user
func Confirmation(user *model.User) {
	env.LoadEnv()

	// Mail content
	from := mail.NewEmail("Palo API", os.Getenv("EMAIL_SENDER"))
	subject := "Email confirmation"
	to := mail.NewEmail(user.Firstname+" "+user.Lastname, user.Email)
	plainTextContent := "Palo dev team"
	htmlContent := "<h4>Palo email confirmation</h4><br><p>Please validate your account by clicking the following link:</p><br><a>http://localhost:4000/verify</a>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_KEY"))

	// Send mail
	_, err := client.Send(message)
	if err != nil {
		log.Println(err)
	}
}
