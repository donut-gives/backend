package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	App    AppStruct
	Server ServerStruct
	Paytm  PaytmStruct
	DB     DBStruct
	Auth   AuthStruct
	Env    string
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

type PaytmStruct struct {
	MerchantID  string
	MerchantKey string
}

type DBStruct struct {
	Url string
}

type AuthStruct struct {
	JWTSecret string `mapstructure:"jwt_secret"`
	Google    GoogleStruct
}

type GoogleStruct struct {
	ClientId     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

var configs *Config

var App *AppStruct
var Server *ServerStruct
var DB *DBStruct
var Paytm *PaytmStruct
var Auth *AuthStruct
var Env *string

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
	Paytm = &configs.Paytm
	Auth = &configs.Auth
	Env = &configs.Env

}
