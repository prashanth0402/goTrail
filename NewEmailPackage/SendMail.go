package newemailpackage

import (
	"crypto/tls"
	"log"

	"gopkg.in/mail.v2"
)

func SendMail() {

	email := mail.NewMessage()
	email.SetHeader("From", "prashanththala1998@gmail.com")
	email.SetHeader("To", "prashanth.s@fcsonline.co.in")
	email.SetAddressHeader("Cc", "pavithra.v@fcsonline.co.in", "CC Pavithra")
	email.SetHeader("Subject", "Test Email")
	email.SetBody("text/plain", "Hello, this is a test email!")

	// Set up your SMTP server configuration
	d := mail.NewDialer("smtpout.secureserver.net", 587, "prashanth.s@fcsonline.co.in", "Best@1234")

	// Send the email
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(email); err != nil {
		log.Fatal(err)
	}
}
