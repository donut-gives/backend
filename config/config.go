package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	App    App
	Server Server
	DB     DB
	Auth   Auth
	Env    string
}

type App struct {
	Name    string
	Version string
	Desc    string
}

type Server struct {
	Host string
	Port string
}

type DB struct {
	Url string
}

type Auth struct {
	JWTSecret string
	Google    Google
}

type Google struct {
	ClientId     string
	ClientSecret string
}

var Configs *Config

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

	err := viper.Unmarshal(&Configs)
	if err != nil {
		log.Fatalf("Error decoding config, %v", err)
	}
}
