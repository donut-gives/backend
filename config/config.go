package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	App     AppStruct
	Server  ServerStruct
	DB      DBStruct
	Auth    AuthStruct
	Payment PaymentStruct
	Captcha CaptchaStruct
	Env     string
	Emailer EmailerStruct
	Cloud   CloudStorageStruct
	OpenAI  OpenAIStruct `mapstructure:"open_ai"`
}

type OpenAIStruct struct {
	APIKey string `mapstructure:"api_key"`
}

type CloudStorageStruct struct {
	KeyFile    string
	UserBucket string
}

type EmailerStruct struct {
	Email       string
	AppPassword string
}

type AppStruct struct {
	Name    string
	Version string
	Desc    string
}

type ServerStruct struct {
	Host string
	Port string
}

type DBStruct struct {
	Name string
	Url  string
}

type AuthStruct struct {
	JWTSecret string `mapstructure:"jwt_secret"`
	Google    Google
}

type Google struct {
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

type PaymentStruct struct {
	Paytm Paytm
}

type Paytm struct {
	MerchantId  string
	MerchantKey string
}

type CaptchaStruct struct {
	Secret string
}

var configs *Config

var App *AppStruct
var Server *ServerStruct
var DB *DBStruct
var Auth *AuthStruct
var Payment *PaymentStruct
var Captcha *CaptchaStruct
var Env *string
var Emailer *EmailerStruct
var Cloud *CloudStorageStruct
var OpenAI *OpenAIStruct

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Errorf("Error loading env file, %s", err)
	}

	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.BindEnv("Server.Port", "PORT"); err != nil {
		log.Errorf("Error binding PORT env var, %s", err)
	}

	viper.SetConfigType("yml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Error reading config, %s", err)
	}

	log.Infof("%s", viper.AllKeys())

	err := viper.Unmarshal(&configs)
	if err != nil {
		log.Errorf("Error decoding config, %v", err)
	}

	App = &configs.App
	Server = &configs.Server
	DB = &configs.DB
	Payment = &configs.Payment
	Auth = &configs.Auth
	Captcha = &configs.Captcha
	Env = &configs.Env
	Emailer = &configs.Emailer
	Cloud = &configs.Cloud
	OpenAI = &configs.OpenAI
}
