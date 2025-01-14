package mailer

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/mail.v2"
)

func Send(address string) {

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "rick0j1ang@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", address)

	// Set E-Mail subject
	m.SetHeader("Subject", "Time to read")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "2 days, such a long time you have not read your books, open your favorite book, have some tea, enjoy your time!!")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "rick0j1ang@gmail.com", "occe jexr cmkv tquz")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return
}
