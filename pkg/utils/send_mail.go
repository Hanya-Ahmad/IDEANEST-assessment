package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

// SendMail sends an organization invitation email to a user
func SendMail(email string, name string, id string) error{
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	emailSender := os.Getenv("EMAIL")
	emailPassword := os.Getenv("PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := 587 
	message := gomail.NewMessage()
	message.SetHeader("From", emailSender)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Invitation to IDEANEST's Organization")
	body := fmt.Sprintf("You are invited to IDEANEST's '%s' organization with id: '%s'",name,id)
	message.SetBody("text/plain", body)
	dialer := gomail.NewDialer(smtpHost, smtpPort, emailSender, emailPassword)
	if err := dialer.DialAndSend(message); err != nil {
		return err
	} 
	return nil
}