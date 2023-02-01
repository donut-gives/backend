package mail

import (
	"donutBackend/config"
	gomail "gopkg.in/gomail.v2"
)

func SendMail(email string,subject string, bodyType string, body string) error {
	msg := gomail.NewMessage()
    msg.SetHeader("From", config.Emailer.Email)
    msg.SetHeader("To", email)
    msg.SetHeader("Subject", subject)
    msg.SetBody(bodyType, body)

    n := gomail.NewDialer("smtp.gmail.com", 587, config.Emailer.Email, config.Emailer.AppPassword)

    // Send the email
    if err := n.DialAndSend(msg); err != nil {
        return err
    }
	return nil
}