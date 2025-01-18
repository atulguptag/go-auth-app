package utils

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, templateFile string, data interface{}) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	email := os.Getenv("EMAIL_ADDRESS")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Fatalf("Error parsing email template: %v", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		log.Fatalf("Error executing email template: %v", err)
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", email)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body.String())

	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		log.Fatalf("Invalid SMTP port: %v", err)
	}

	dialer := gomail.NewDialer(smtpHost, port, email, password)

	// Send the email
	if err := dialer.DialAndSend(mailer); err != nil {
		log.Printf("Failed to send email: %v", err)
	}
}
