package news

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	MongoConfig MongoConfig
	Server      ServerConfig
}

// MongoConfig - config
type MongoConfig struct {
	Collection string `envconfig:"MONGO_COLLECTION"`
	Database   string `envconfig:"MONGO_DATABASE"`
	URI        string `envconfig:"MONGO_URI"`
}

type ServerConfig struct {
	Port string `envconfig:"PORT"`
}

func newConfig() (Config, error) {
	var conf Config

	if err := envconfig.Process("", &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
