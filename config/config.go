package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type config struct {
	ChatWork struct{
		APIKey string `env:"CHATWORK_TOKEN_KEY" envDefault:""`
		RoomID string `env:"CHATWORK_ROOM_ID" envDefault:""`
	}
	ECR struct{
		Api string `env:"ECR_API" envDefault:""`
		Bff string `env:"ECR_BFF" envDefault:""`
		BffNginx string `env:"ECR_BFF_NGINX" envDefault:""`
	}
	S3 struct{
		BukkenName string `env:"S3_BUKKEN_NAME" envDefault:""`
	}
	K8s struct{
		NameSpace string `env:"K8S_NAMESPACE" envDefault:"dev"`
	}
}

//C config struct
var C config

//LoadConfig load config from environment and parse to struct
func LoadConfig() {
	cfg := &C

	if err := env.Parse(cfg); err != nil {
		log.Println("Application config parsing failed: " + err.Error())
	}
}