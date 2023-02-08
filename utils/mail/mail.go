package mail

import (
	"encoding/base64"
	"fmt"
    "context"
	//"net/http"

	"donutBackend/config"

	"golang.org/x/oauth2"
	gomail "gopkg.in/gomail.v2"

	//"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	//"google.golang.org/api/option"
)

var Email string=""
var gmailToken *oauth2.Token = nil
var googleOauthConfig *oauth2.Config = nil

func RefreshAccessToken() error {
    if(gmailToken==nil){
        return fmt.Errorf("Gmail Token is not set")
    }

    config := &oauth2.Config{
		ClientID:     config.Auth.Google.ClientId,     
		ClientSecret: config.Auth.Google.ClientSecret, 
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send","https://www.googleapis.com/auth/gmail.labels","openid","profile", "email"},
		Endpoint:     google.Endpoint,
	}

    token, err := config.TokenSource(oauth2.NoContext, gmailToken).Token()
    if err != nil {
        return err
    }

    gmailToken = token

    return nil

}

func SetTokenAndConfig(token *oauth2.Token, config *oauth2.Config) {
    gmailToken = token
    googleOauthConfig = config
}


func SendMail(to, subject,bodyType, body string) error {
	var message gmail.Message

    if(gmailToken==nil){
        return fmt.Errorf("Gmail Token is not set")
    }

    gmailClient:=googleOauthConfig.Client(context.Background(), gmailToken)

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
