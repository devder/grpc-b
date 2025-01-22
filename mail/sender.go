package mail

import (
	"fmt"

	"gopkg.in/mail.v2"
)

const (
	smtpAuthAddress = "smtp.gmail.com"
	smtpPort        = 465
)

type EmailSender interface {
	SendEmail(
		subject,
		content string,
		to,
		cc,
		bcc,
		attachFiles []string,
	) error
}

type GmailSender struct {
	name,
	fromEmailAddress,
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
}

func (s *GmailSender) SendEmail(
	subject,
	content string,
	to,
	cc,
	bcc,
	attachFiles []string,
) error {
	// create a new message
	m := mail.NewMessage()

	// Set headers
	m.SetHeader("From", fmt.Sprintf("%s <%s>", s.name, s.fromEmailAddress))
	m.SetHeader("Subject", subject)
	for _, t := range to {
		m.SetHeader("To", t)
	}
	for _, c := range cc {
		m.SetHeader("Cc", c)
	}
	for _, b := range bcc {
		m.SetHeader("Bcc", b)
	}
	for _, f := range attachFiles {
		m.Attach(f)
	}

	// Set body
	m.SetBody("text/html", content)

	d := mail.NewDialer(smtpAuthAddress, smtpPort, s.fromEmailAddress, s.fromEmailPassword)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
