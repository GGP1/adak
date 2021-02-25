/*
Package email helps us to use the email as the tool to identify each user
*/
package email

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/GGP1/adak/internal/token"

	"github.com/pkg/errors"
)

var (
	// Both os.Getenv and viper.GetString work
	emailSender   = os.Getenv("EMAIL_SENDER")
	emailPassword = os.Getenv("EMAIL_PASSWORD")
	emailHost     = os.Getenv("EMAIL_HOST")
	emailPort     = os.Getenv("EMAIL_PORT")
)

// Items is a struct that keeps the values passed to the templates.
type Items struct {
	ID       string
	Name     string
	Email    string
	Token    string
	NewEmail string
}

// SendValidation sends a validation email to the user.
func SendValidation(ctx context.Context, username, email string, errCh chan error) {
	token := token.GenerateRunes(20)

	// 	Email content
	from := mail.Address{Name: "Adak", Address: emailSender}
	to := mail.Address{Name: username, Address: email}
	subject := "Validation email"
	items := Items{
		Name:  username,
		Email: email,
		Token: token,
	}

	headers := make(map[string]string, 4)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	t, err := template.ParseFiles("../internal/email/templates/validation.html")
	if err != nil {
		errCh <- err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, items); err != nil {
		errCh <- err
	}

	message += buf.String()

	// Connect to smtp
	auth := smtp.PlainAuth("", emailSender, emailPassword, emailHost)

	if err := smtp.SendMail(emailHost+":"+emailPort, auth, from.Address, []string{to.Address}, []byte(message)); err != nil {
		errCh <- errors.Wrap(err, "couldn't send the email")
	}
}

// SendChangeConfirmation sends a confirmation email to the user.
func SendChangeConfirmation(id, username, email, newEmail, token string, errCh chan error) {
	// 	Email content
	from := mail.Address{Name: "Adak", Address: emailSender}
	to := mail.Address{Name: username, Address: email}
	subject := "Email change confirmation"
	items := Items{
		ID:       id,
		Name:     username,
		Token:    token,
		NewEmail: newEmail,
	}

	headers := make(map[string]string, 4)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	t, err := template.ParseFiles("../internal/email/templates/changeEmail.html")
	if err != nil {
		errCh <- err
	}

	buf := new(bytes.Buffer)
	if err := t.Execute(buf, items); err != nil {
		errCh <- err
	}

	message += buf.String()

	// Connect to smtp
	auth := smtp.PlainAuth("", emailSender, emailPassword, emailHost)

	if err := smtp.SendMail(emailHost+":"+emailPort, auth, from.Address, []string{to.Address}, []byte(message)); err != nil {
		errCh <- errors.Wrap(err, "couldn't send the email")
	}
}
