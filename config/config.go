package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port            string        `mapstructure:"PORT"`
	Connection      string        `mapstructure:"MYSQL_CONNECTION"`
	ConnRetrySleep  time.Duration `default:"3s"`
	ConnMaxLifetime time.Duration
	MaxIdleConns    uint   `mapstructure:"MYSQL_MAX_IDLE"`
	MaxOpenConns    uint   `mapstructure:"MYSQL_MAX_CONN"`
	RabbitUrl       string `mapstructure:"RABBIT_URL"`
	OutboxQueue     string `mapstructure:"OUTBOX_QUEUE"`
	OutboxInterval  string `mapstructure:"OUTBOX_INTERVAL"`
}

type MySQL struct {
	Connection      string
	ConnRetrySleep  time.Duration `default:"3s"`
	ConnMaxLifetime time.Duration
	MaxIdleConns    uint `validate:"min=0"`
	MaxOpenConns    uint `validate:"min=0"`
}

func GetConf() *AppConfig {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		fmt.Printf("%v", err)
	}

	conf := &AppConfig{}
	err = viper.Unmarshal(conf)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}

	return conf
}
