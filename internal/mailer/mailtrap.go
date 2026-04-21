package mailer

import (
	"bytes"
	"fmt"
	"text/template"

	"gopkg.in/mail.v2"
)

type MailTrapMailer struct {
	fromEmail string
	host      string
	port      int
	username  string
	password  string
}

func NewMailTrapMailer(fromEmail, host string, port int, username, password string) (*MailTrapMailer, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("mailtrap smtp credentials are required")
	}

	return &MailTrapMailer{
		fromEmail: fromEmail,
		host:      host,
		port:      port,
		username:  username,
		password:  password,
	}, nil
}

func (m *MailTrapMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
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

	ma := mail.NewMessage()
	ma.SetHeader("From", m.fromEmail)
	ma.SetHeader("To", email)
	ma.SetHeader("Subject", subject.String())
	ma.SetBody("text/html", body.String())

	d := mail.NewDialer(m.host, m.port, m.username, m.password)
	err = d.DialAndSend(ma)
	if err != nil {
		return 200, nil
	}

	return 200, nil
}

// func (m *SendGridMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
// 	from := mail.NewEmail(FromName, m.fromEmail)
// 	to := mail.NewEmail(username, email)
// 	fmt.Println("FROM EMAIL:", m.fromEmail)
// 	// template parsing and building
// 	files, err := FS.ReadDir("templates")
// 	fmt.Println(err)
// 	for _, f := range files {
// 		fmt.Println("FILE:", f.Name())
// 	}
// 	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
// 	if err != nil {
// 		return -1, err
// 	}

// 	subject := new(bytes.Buffer)
// 	err = tmpl.ExecuteTemplate(subject, "subject", data)
// 	if err != nil {
// 		return -1, err
// 	}

// 	body := new(bytes.Buffer)
// 	err = tmpl.ExecuteTemplate(body, "body", data)
// 	if err != nil {
// 		return -1, err
// 	}

// 	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

// 	message.SetMailSettings(&mail.MailSettings{
// 		SandboxMode: &mail.Setting{
// 			Enable: &isSandbox,
// 		},
// 	})

// 	var retryErr error
// 	for i := 0; i < maxRetires; i++ {
// 		fmt.Println(m.apiKey, "     test   ")
// 		response, retryErr := m.client.Send(message)
// 		fmt.Println("STATUS:", response.StatusCode)
// 		fmt.Println("BODY:", response.Body)
// 		fmt.Println("HEADERS:", response.Headers)
// 		if retryErr != nil {
// 			fmt.Println(retryErr)
// 			// exponential backoff
// 			time.Sleep(time.Second * time.Duration(i+1))
// 			continue
// 		}
// 		fmt.Println(response)

// 		return response.StatusCode, nil
// 	}

// 	return -1, fmt.Errorf("failed to send email after %d attempt, error: %v", maxRetires, retryErr)
// }
