package mailer

import (
	"bytes"
	"html/template"
	"net/smtp"
)

type SMTPMailer struct {
	fromEmail string
	password  string
}

func NewSMTPMailer(fromEmail, password string) *SMTPMailer {
	return &SMTPMailer{
		fromEmail: fromEmail,
		password:  password,
	}
}

func (m *SMTPMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	from := m.fromEmail
	password := m.password

	auth := smtp.PlainAuth(
		"",
		from,
		password,
		"smtp.gmail.com",
	)

	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return err
	}

	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return err
	}

	message := []byte(
		"From: " + from + "\r\n" +
			"To: " + email + "\r\n" +
			"Subject: " + subject.String() + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body.String() + "\r\n",
	)

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		from,
		[]string{email},
		message,
	)
}
