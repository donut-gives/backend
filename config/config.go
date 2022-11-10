package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	App     AppStruct
	Server  ServerStruct
	DB      DBStruct
	Auth    AuthStruct
	Payment PaymentStruct
	Env     string
	Emailer EmailerStruct
	Cloud  	CloudStorageStruct
}

type CloudStorageStruct struct {
	KeyFile string
	UserBucket string
}

type EmailerStruct struct {
	Email string
	AppPassword  string
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
	Url string
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

var configs *Config

var App *AppStruct
var Server *ServerStruct
var DB *DBStruct
var Auth *AuthStruct
var Payment *PaymentStruct
var Env *string
var Emailer *EmailerStruct
var Cloud *CloudStorageStruct

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading env file, %s", err)
	}

	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.BindEnv("Server.Port", "PORT"); err != nil {
		log.Fatalf("Error binding PORT env var, %s", err)
	}

	viper.SetConfigType("yml")
	if os.Getenv("ENV") == "prod" {
		viper.SetConfigName("config-prod")
	} else {
		viper.SetConfigName("config")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config, %s", err)
	}

	log.Infof("%s", viper.AllKeys())

	err := viper.Unmarshal(&configs)
	if err != nil {
		log.Fatalf("Error decoding config, %v", err)
	}

	App = &configs.App
	Server = &configs.Server
	DB = &configs.DB
	Payment = &configs.Payment
	Auth = &configs.Auth
	Env = &configs.Env
	Emailer = &configs.Emailer
	Cloud = &configs.Cloud
}
