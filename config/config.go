package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"go-deduplication/internal/entity"
)

type Config struct {
	Postgres      string
	ByteSize      int
	FileDirectory string
	HashFunction  entity.HashFunction
	Address       string
}

func InitConfig(fileName, path string) (*Config, error) {
	viper.SetConfigName(fileName)
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	var configuration Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "reading config from viper")
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal config from viper")
	}
	return &configuration, nil
}
