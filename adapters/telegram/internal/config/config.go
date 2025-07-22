package config

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type HTTPConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int64         `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type ServiceConfig struct {
	URL string `mapstructure:"url"`
}

type ServicesConfig struct {
	BotCore ServiceConfig `mapstructure:"bot_core"`
}

type Config struct {
	Services ServicesConfig `mapstructure:"services"`
	HTTP     HTTPConfig     `mapstructure:"http"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
