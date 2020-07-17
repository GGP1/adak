/*
Package email helps us to use the email as the tool to identify each user
*/
package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/mail"
	"net/smtp"

	"github.com/GGP1/palo/internal/cfg"
	"github.com/GGP1/palo/pkg/model"
)

// Items represents a struct with the values passed to the template
type Items struct {
	Name  string
	Token string
}

// SendValidation sends a validation email to the user
func SendValidation(user model.User, token string) error {
	// =================
	// 	Email content
	// =================
	username := user.Firstname + " " + user.Lastname

	from := mail.Address{Name: "Palo", Address: cfg.EmailSender}
	to := mail.Address{Name: username, Address: user.Email}
	subject := "Validation email"
	items := Items{
		Name:  user.Firstname + " " + user.Lastname,
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

	t, err := template.ParseFiles("../pkg/auth/email/template.html")
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
	serverName := "smtp.gmail.com:465"
	host := "stmp.gmail.com"

	auth := smtp.PlainAuth("", cfg.EmailSender, cfg.EmailPassword, host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", serverName, tlsConfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	err = client.Auth(auth)
	if err != nil {
		return err
	}

	err = client.Mail(from.Address)
	if err != nil {
		return err
	}

	err = client.Rcpt(to.Address)
	if err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = client.Quit()
	if err != nil {
		return err
	}

	return nil
}
