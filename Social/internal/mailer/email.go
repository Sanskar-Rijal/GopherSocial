package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

type GmailMailer struct {
	fromEmail string
	fromName  string
	password  string //gmail app password
	host      string //smtp.gmail.com
	port      string //587
}

func NewGmailMailer(fromEmail string, password string) *GmailMailer {
	return &GmailMailer{
		fromEmail: fromEmail,
		fromName:  FromName,
		password:  password,
		host:      "smtp.gmail.com",
		port:      "587",
	}
}

// Send function with retry logic
func (m *GmailMailer) Send(templateFile string, username string, email string, data any, isDevelopment bool) (int, error) {
	//parse templatefile

	tmpl, err := template.ParseFS(FS, "template/"+templateFile)
	if err != nil {
		return -1, err
	}

	//build subject from template
	subject := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return -1, err
	}
	//build body from template
	body := new(bytes.Buffer)
	if err := tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return -1, err
	}

	message := buildMessage(
		m.fromEmail,
		m.fromName,
		email,
		subject.String(),
		body.String(),
	)

	//in development don't send email, just log it
	if isDevelopment {
		fmt.Printf("SANDBOX EMAIL\nTo: %s\nSubject: %s\nBody: %s\n",
			email, subject.String(), body.String())
		return 200, nil
	}

	auth := smtp.PlainAuth(
		"",
		m.fromEmail,
		m.password,
		m.host,
	)

	//if email fails we retry atleast 3 times
	var retryError error
	for i := 0; i < maxRetries; i++ {

		retryError = smtp.SendMail(
			m.host+":"+m.port,
			auth,
			m.fromEmail,
			[]string{email},
			[]byte(message),
		)

		if retryError != nil {
			// exponential backoff — wait longer each retry
			// retry 1 — wait 1 second
			// retry 2 — wait 2 seconds
			// retry 3 — wait 3 seconds
			waitTime := time.Second * time.Duration(i+1)
			time.Sleep(waitTime)
			continue
		}

		//if no error return
		return 200, nil
	}

	return -1, fmt.Errorf("Failed to send email after %d attempts, errors: %v", maxRetries, retryError)

}

// builds the raw email message
// smtp needs specific format
func buildMessage(fromEmail, fromName, toEmail, subject, body string) string {

	return fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		fromName,
		fromEmail,
		toEmail,
		subject,
		body,
	)
}
