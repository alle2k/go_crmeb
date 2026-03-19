package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type config struct {
	Mysql  mysqlConfig
	Server serverConfig
	Redis  redisConfig
	Token  jwtConfig
}

var AppConfig config

type serverConfig struct {
	Port int
}

type mysqlConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type redisConfig struct {
	Host      string
	Port      int
	Password  string
	Timeout   int
	Database  int
	MaxActive int `mapstructure:"max-active"`
	MaxIdle   int `mapstructure:"max-idle"`
	MinIdle   int `mapstructure:"min-idle"`
}

type jwtConfig struct {
	Header     string
	Secret     string
	ExpireTime int
}

func Init() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "test"
	}

	viper.SetConfigName("config." + env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatal(err)
	}
}
