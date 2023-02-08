package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	//"time"

	//"time"

	//"net/http"

	"donutBackend/config"
	"donutBackend/models/emailSender"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	gomail "gopkg.in/gomail.v2"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	//"google.golang.org/api/option"
)

var Email string=""
var gmailToken *oauth2.Token = nil
var GoogleOauthConfig *oauth2.Config = &oauth2.Config{
    ClientID:     config.Auth.Google.ClientId,     
    ClientSecret: config.Auth.Google.ClientSecret, 
    Scopes:       []string{"https://www.googleapis.com/auth/gmail.send","https://www.googleapis.com/auth/gmail.labels","openid","profile", "email"},
    Endpoint:     google.Endpoint,
}

func RefreshAccessToken() error {

    email,err:=emailsender.GetEmail()
    if err != nil {
        return err
    }
    
    jwttoken,err:=emailsender.GetToken(email)
    if err != nil {
        return err
    }

    //decoed jwt token
    decodedToken, err := jwt.Parse(jwttoken, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(config.Auth.Google.ClientSecret), nil
    })

    //create google oauth token
    googleOauthToken := &oauth2.Token{
        AccessToken: decodedToken.Claims.(jwt.MapClaims)["access_token"].(string),
        RefreshToken: decodedToken.Claims.(jwt.MapClaims)["refresh_token"].(string),
        TokenType: decodedToken.Claims.(jwt.MapClaims)["token_type"].(string),
    }

    //fmt.Println("GoogleOauthToken-",googleOauthToken)
    token, err := GoogleOauthConfig.TokenSource(context.TODO(), googleOauthToken).Token()
    if err != nil {
        return err
    }
    // fmt.Println("error           -",err)
    // fmt.Println("Token           -",token)
    // if(token.AccessToken==googleOauthToken.AccessToken){
    //     fmt.Println("Tokens not refreshed")
    // }      

    gmailToken = token

    return nil

}

func SetTokenAndConfig(token *oauth2.Token, config *oauth2.Config) {
    gmailToken = token
    GoogleOauthConfig = config
}


func SendMail(to, subject,bodyType, body string) error {
	var message gmail.Message

    if(gmailToken==nil){
        return fmt.Errorf("Gmail Token is not set")
    }

    gmailClient:=GoogleOauthConfig.Client(context.Background(), gmailToken)

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
