package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type LoggingConfig struct {
	Level    string
	Encoding string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "postgres")
	viper.SetDefault("database.name", "subscriptions")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.encoding", "json")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Конфигурационный файл не найден, используются значения по умолчанию")
		} else {
			return nil, fmt.Errorf("ошибка чтения конфигурационного файла: %w", err)
		}
	}

	viper.AutomaticEnv()

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("server.port"),
			Mode: viper.GetString("server.mode"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("database.host"),
			Port:     viper.GetString("database.port"),
			User:     viper.GetString("database.user"),
			Password: viper.GetString("database.password"),
			Name:     viper.GetString("database.name"),
			SSLMode:  viper.GetString("database.sslmode"),
		},
		Logging: LoggingConfig{
			Level:    viper.GetString("logging.level"),
			Encoding: viper.GetString("logging.encoding"),
		},
	}

	return config, nil
}
