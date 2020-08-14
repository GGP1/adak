/*
Package email helps us to use the email as the tool to identify each user
*/
package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"
	"github.com/pkg/errors"
)

// Items represents a struct with the values passed to the templates.
type Items struct {
	ID       uint
	Name     string
	Token    string
	NewEmail string
}

// SendValidation sends a validation email to the user.
func SendValidation(user model.User, token string, errCh chan error) {
	// =================
	// 	Email content
	// =================
	from := mail.Address{Name: "Palo", Address: cfg.EmailSender}
	to := mail.Address{Name: user.Name, Address: user.Email}
	subject := "Validation email"
	items := Items{
		Name:  user.Name,
		Token: token,
	}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	t, err := template.ParseFiles("../pkg/auth/email/validation.html")
	if err != nil {
		errCh <- err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, items)
	if err != nil {
		errCh <- err
	}

	message += buf.String()

	// =================
	// Connect to smtp
	// =================
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", cfg.EmailSender, cfg.EmailPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from.Address, []string{to.Address}, []byte(message))
	if err != nil {
		errCh <- errors.Wrap(err, "couldn't send the email")
	}
}

// SendChangeConfirmation sends a validation email to the user.
func SendChangeConfirmation(user model.User, token, newEmail string, errCh chan error) {
	// =================
	// 	Email content
	// =================
	from := mail.Address{Name: "Palo", Address: cfg.EmailSender}
	to := mail.Address{Name: user.Name, Address: user.Email}
	subject := "Email change confirmation"
	items := Items{
		ID:       user.ID,
		Name:     user.Name,
		Token:    token,
		NewEmail: newEmail,
	}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	t, err := template.ParseFiles("../pkg/auth/email/changeEmail.html")
	if err != nil {
		errCh <- err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, items)
	if err != nil {
		errCh <- err
	}

	message += buf.String()

	// =================
	// Connect to smtp
	// =================
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", cfg.EmailSender, cfg.EmailPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from.Address, []string{to.Address}, []byte(message))
	if err != nil {
		fmt.Println(err)
		errCh <- errors.Wrap(err, "couldn't send the email")
	}
}
