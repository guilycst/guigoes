package main

import (
	"log"

	"github.com/guilycst/guigoes/pkg"
	"github.com/wneessen/go-mail"
)

func main() {
	pkg.LoadEnvFile(".prod.env")
	m := mail.NewMsg()
	if err := m.From("no-reply@guigoes.com"); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := m.To("guilycst@gmail.com"); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}

	m.Subject("New post!")
	m.SetBodyString(mail.TypeTextPlain, "Testando o envio de email enviado pelo Go!")

	c, err := mail.NewClient(pkg.SMTP_ENDPOINT,
		mail.WithPort(pkg.SMTP_PORT),
		mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithUsername(pkg.SMTP_USR_NAME),
		mail.WithPassword(pkg.SMTP_USR_PW),
	)
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}

	if err := c.DialAndSend(m); err != nil {
		log.Fatalf("failed to send mail: %s", err)
	}
}
