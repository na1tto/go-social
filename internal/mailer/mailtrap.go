package mailer

import (
	"bytes"
	"errors"
	"html/template"

	gomail "github.com/go-mail/mail/v2"
)

type mailTrapClient struct {
	fromEmail       string
	apiKey          string
	sandboxUsername string
	password        string
}

func NewMailTrapClient(apiKey, fromEmail, sandboxUsername, password string) (mailTrapClient, error) {
	if apiKey == "" {
		return mailTrapClient{}, errors.New("api key is required")
	}

	return mailTrapClient{
		apiKey:          apiKey,
		fromEmail:       fromEmail,
		sandboxUsername: sandboxUsername,
		password:        password,
	}, nil
}

func (m mailTrapClient) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := gomail.NewMessage()
	if m.fromEmail == "" {
		return -1, errors.New("mailer: from email is empty. Check your env file")
	}
	message.SetAddressHeader("From", m.fromEmail, FromName)
	message.SetAddressHeader("To", email, username)

	message.SetHeader("Subject", subject.String())
	message.SetBody("text/html", body.String())

	// change sandbox for live later
	dialer := gomail.NewDialer("sandbox.smtp.mailtrap.io", 587, m.sandboxUsername, m.password)

	//TO DO: implement a form of retry that can be reused across both implementations

	if err := dialer.DialAndSend(message); err != nil {
		return -1, err
	}

	return 200, nil
}
