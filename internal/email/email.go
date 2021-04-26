/*
Package email helps us to use the email as the tool to identify each user
*/
package email

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/GGP1/adak/internal/bufferpool"
	"github.com/GGP1/adak/internal/logger"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Emailer contains emails templates and the sender information.
type Emailer struct {
	addr     string
	host     string
	name     string
	username string
	password string

	validation  *template.Template
	changeEmail *template.Template
}

// Items is a struct that keeps the values passed to the templates.
type Items struct {
	ID       string
	Name     string
	Email    string
	Token    string
	NewEmail string
}

// New returns a new emailer.
func New() Emailer {
	staticFS := viper.Get("static.fs").(embed.FS)
	validation, err := template.ParseFS(staticFS, "static/templates/validation.html")
	if err != nil {
		logger.Log.Fatalf("Failed parsing validation template")
	}
	changeEmail, err := template.ParseFS(staticFS, "static/templates/changeEmail.html")
	if err != nil {
		logger.Log.Fatalf("Failed parsing change email template")
	}
	host := viper.GetString("email.host")

	return Emailer{
		addr:        net.JoinHostPort(host, viper.GetString("email.port")),
		host:        host,
		name:        "Adak",
		username:    viper.GetString("email.sender"),
		password:    viper.GetString("email.password"),
		validation:  validation,
		changeEmail: changeEmail,
	}
}

// SendValidation sends a validation email to the user.
func (e *Emailer) SendValidation(ctx context.Context, username, email, token string) error {
	// Email content
	from := mail.Address{Name: e.name, Address: e.username}
	to := mail.Address{Name: username, Address: email}
	items := Items{
		Name:  username,
		Email: email,
		Token: token,
	}

	headers := make(map[string]string, 4)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = "Validation email"
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := bufferpool.Get()
	defer bufferpool.Put(message)

	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	buf := bufferpool.Get()
	if err := e.validation.Execute(buf, items); err != nil {
		return err
	}
	message.Write(buf.Bytes())
	bufferpool.Put(buf)

	// Connect to smtp
	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	if err := smtp.SendMail(e.addr, auth, from.Address, []string{to.Address}, message.Bytes()); err != nil {
		logger.Log.Errorf("Couldn't send the validation email: %v.\nAddr: %s\nEmail: %s", err, e.addr, to.Address)
		return errors.Wrap(err, "couldn't send the email")
	}

	logger.Log.Infof("Successfully sent email to: %s", to.Address)
	return nil
}

// SendChangeConfirmation sends a confirmation email to the user.
func (e *Emailer) SendChangeConfirmation(id, username, email, newEmail, token string) error {
	// Email content
	from := mail.Address{Name: e.name, Address: e.username}
	to := mail.Address{Name: username, Address: email}
	items := Items{
		ID:       id,
		Name:     username,
		Token:    token,
		NewEmail: newEmail,
	}

	headers := make(map[string]string, 4)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = "Email change confirmation"
	headers["Content-Type"] = `text/html; charset="UTF-8"`

	message := bufferpool.Get()
	defer bufferpool.Put(message)

	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	buf := bufferpool.Get()
	if err := e.changeEmail.Execute(buf, items); err != nil {
		return err
	}
	message.Write(buf.Bytes())
	bufferpool.Put(buf)

	// Connect to smtp
	auth := smtp.PlainAuth("", e.username, e.password, e.host)

	if err := smtp.SendMail(e.addr, auth, from.Address, []string{to.Address}, message.Bytes()); err != nil {
		logger.Log.Errorf("Couldn't send the change confirmation email: %v.\nAddr: %s\nEmail: %s", err, e.addr, to.Address)
		return errors.Wrap(err, "couldn't send the email")
	}

	logger.Log.Infof("Successfully sent email to: %s", to.Address)
	return nil
}
