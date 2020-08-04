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

// Items represents a struct with the values passed to the template.
type Items struct {
	Name  string
	Token string
}

// SendValidation sends a validation email to the user.
func SendValidation(user model.User, token string) error {
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

	t, err := template.ParseFiles("../pkg/auth/email/email.html")
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, items)
	if err != nil {
		return err
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
		return errors.Wrap(err, "couldn't send the email")
	}

	return nil
}
