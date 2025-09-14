package mail

import (
	"bytes"
	"path/filepath"
	"text/template"
	"time"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dailer *gomail.Dialer
	sender string
}

func NewMailer(host string, port int, username string, password string, from string, ssl bool, tls bool) *Mailer {
	dailer := gomail.NewDialer(host, port, username, password)
	dailer.SSL = ssl

	return &Mailer{dailer: dailer, sender: from}
}

func (m *Mailer) SendEmail(to []string, templateFile string, data interface{}) error {
	// Use absolute path to avoid path resolution issues
	absPath, err := filepath.Abs(templateFile)
	if err != nil {
		return err
	}

	tmpl, err := template.ParseFiles(absPath)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.sender)
	msg.SetHeader("To", to...)

	if dataMap, ok := data.(map[string]interface{}); ok {
		if subject, exists := dataMap["Subject"]; exists {
			msg.SetHeader("Subject", subject.(string))
		}
	}

	msg.SetBody("text/html", buf.String())

	for i := 0; i < 3; i++ {
		err = m.dailer.DialAndSend(msg)
		if nil == err {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return err
}

// SendWelcomeEmail is a convenience method for sending welcome emails
func (m *Mailer) SendWelcomeEmail(to []string, username, verifyLink string) error {
	data := map[string]interface{}{
		"Subject":    "Welcome to Money Manager - Verify Your Email",
		"Username":   username,
		"VerifyLink": verifyLink,
	}
	return m.SendEmail(to, "internal/mail/templates/userwelcome.tmpl", data)
}

// SendVerificationEmail is a convenience method for sending verification emails
func (m *Mailer) SendVerificationEmail(to []string, username string) error {
	data := map[string]interface{}{
		"Subject":  "Email Verified - Money Manager",
		"Username": username,
	}
	return m.SendEmail(to, "internal/mail/templates/userverfied.tmpl", data)
}

// TestConnection tests the SMTP connection
func (m *Mailer) TestConnection() error {
	s, err := m.dailer.Dial()
	if err != nil {
		return err
	}
	defer s.Close()
	return nil
}
