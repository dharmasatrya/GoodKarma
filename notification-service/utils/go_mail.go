package utils

import (
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmailNotification(to, subject, content string) {
	from := os.Getenv("3RD_PARTY_EMAIL")
	password := os.Getenv("3RD_PARTY_EMAIL_PASSWORD")

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", content)

	// Set up the SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatalf("Could not send email: %v", err)
	}

	log.Println("Email sent successfully!")
}
