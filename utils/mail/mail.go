package mail

import (
    "encoding/base64"
    "fmt"

    "net/http"
    
    
	"donutBackend/config"
	gomail "gopkg.in/gomail.v2"
    //"golang.org/x/oauth2"
    //"golang.org/x/oauth2/google"
    "google.golang.org/api/gmail/v1"
    //"google.golang.org/api/option"

)

var Email string=""
var gmailClient *http.Client

func SetClient(client *http.Client) {
    
    gmailClient = client
}


func SendMail(to, subject,bodyType, body string) error {
	var message gmail.Message

    if(gmailClient==nil){
        return fmt.Errorf("Gmail Client is not set")
    }

    srv, err := gmail.New(gmailClient)

    subject="=?utf-8?B?"+base64.StdEncoding.EncodeToString([]byte(subject))+"?="

	message.Raw = base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nContent-Type: %s; charset=UTF-8\r\n\r\n%s", to, subject,bodyType, body)))
	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
        return err
    }

	return nil
}

func SendMailBySMTP(email string,subject string, bodyType string, body string) error {
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
