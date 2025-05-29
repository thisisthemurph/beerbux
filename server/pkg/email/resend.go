package email

import (
	"beerbux/internal/api/config"
	"github.com/resend/resend-go/v2"
)

type ResendEmailSender struct {
	*resend.Client
	devSendToEmail string
}

type Sender interface {
	Send(to, subject, html string) (string, error)
}

func New(conf config.ResendConfig) Sender {
	return &ResendEmailSender{
		Client:         resend.NewClient(conf.Key),
		devSendToEmail: conf.DevelopmentSendToEmail,
	}
}

func (r *ResendEmailSender) Send(to, subject, html string) (string, error) {
	sendToEMail := to
	if r.devSendToEmail != "" {
		sendToEMail = r.devSendToEmail
	}

	params := &resend.SendEmailRequest{
		From:    "Beerbux <onboarding@resend.dev>",
		To:      []string{sendToEMail},
		Html:    html,
		Subject: subject,
	}

	sent, err := r.Client.Emails.Send(params)
	if err != nil {
		return "", err
	}
	return sent.Id, nil
}
