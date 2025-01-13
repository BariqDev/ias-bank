package mail

import (
	"fmt"

	"github.com/wneessen/go-mail"
)

const (
	smtAuthAddress   = "smtp.gmail.com"
	smtServerAddress = "smtp.gmail.com:587"
)

type MailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachedFiles []string,
	) error
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmail, fromEmailPassword string) MailSender {

	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmail,
		fromEmailPassword: fromEmailPassword,
	}
}

func (sender *GmailSender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachedFiles []string,
) error {

	message := mail.NewMsg()

	from:= fmt.Sprintf("%s <%s>",sender.name,sender.fromEmailAddress)
	if err := message.From(from); err != nil {
		return fmt.Errorf("failed to set From address: %s", err)
	}
	if err := message.To(to...); err != nil {
		return fmt.Errorf("failed to set To address: %s", err)
	}
	if err := message.Cc(cc...); err != nil {
		return fmt.Errorf("failed to set Cc address: %s", err)
	}
	if err := message.Bcc(bcc...); err != nil {
		return fmt.Errorf("failed to set bcc address: %s", err)
	}
	for _, f := range attachedFiles {
		message.AttachFile(f)
	}

	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, content)

	client, err := mail.NewClient(
		smtAuthAddress,
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(sender.fromEmailAddress),
		mail.WithPassword(sender.fromEmailPassword))
	if err != nil {
		return fmt.Errorf("failed to create mail client: %s\n", err)
	}

	return client.DialAndSend(message)
}
