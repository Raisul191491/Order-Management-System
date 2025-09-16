package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var Conf Config

type Config struct {
	DBUserRead          string        `mapstructure:"DB_USER_READ"`
	DBUserWrite         string        `mapstructure:"DB_USER_WRITE"`
	DBPassword          string        `mapstructure:"DB_PASSWORD"`
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBName              string        `mapstructure:"DB_NAME"`
	DBHostRead          string        `mapstructure:"DB_HOST_READ"`
	DBHostWrite         string        `mapstructure:"DB_HOST_WRITE"`
	DBPortRead          string        `mapstructure:"DB_PORT_READ"`
	DBPortWrite         string        `mapstructure:"DB_PORT_WRITE"`
	DBMaxOpenConnection int           `mapstructure:"DB_MAX_OPEN_CONNECTION"`
	DBMaxIdleConnection int           `mapstructure:"DB_MAX_IDLE_CONNECTION"`
	DBConnMaxLife       time.Duration `mapstructure:"DB_CONN_MAX_LIFE"`
	Port                string        `mapstructure:"PORT"`
	RedisHost           string        `mapstructure:"REDIS_HOST"`
	RedisPort           string        `mapstructure:"REDIS_PORT"`
}

func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatal("Error reading env file", err)
	}

	return &Conf
}
