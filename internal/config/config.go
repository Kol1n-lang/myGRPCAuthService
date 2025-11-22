package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"authService/internal/utils"
	"github.com/joho/godotenv"
)

type Config struct {
	JWT             JWTConfig
	RabbitMQ        RabbitMQConfig
	DB              DBConfig
	Email           EmailConfig
	MetricsPort     string
	BrokerConstants struct {
		EmailConfirm string
	}
}

type JWTConfig struct {
	JWTSecret           string
	AccessExpireMinutes int
	RefreshExpireDays   int
	Algorithm           string
}

type RabbitMQConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

type EmailConfig struct {
	Sender      string
	AppPassword string
}

func Init() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Using static env variables")
	}
	config := &Config{}

	config.DB = DBConfig{
		Host:     getEnv("DB_HOST", ""),
		Port:     getEnv("DB_PORT", ""),
		User:     getEnv("DB_USER", ""),
		Password: getEnv("DB_PASSWORD", ""),
		DB:       getEnv("DB_NAME", ""),
	}

	config.RabbitMQ = RabbitMQConfig{
		Host:     getEnv("RABBIT_HOST", ""),
		Port:     getEnv("RABBIT_PORT", ""),
		User:     getEnv("RABBIT_LOGIN", ""),
		Password: getEnv("RABBIT_PASSWORD", ""),
	}

	config.JWT = JWTConfig{
		AccessExpireMinutes: utils.Atoi(getEnv("ACCESS_TOKEN_EXPIRE_MINUTES", "")),
		RefreshExpireDays:   utils.Atoi(getEnv("REFRESH_TOKEN_EXPIRE_DAYS", "")),
		JWTSecret:           getEnv("SECRET_KEY", ""),
		Algorithm:           getEnv("ALGORITHM", ""),
	}

	config.Email = EmailConfig{
		Sender:      getEnv("SENDER", ""),
		AppPassword: getEnv("APP_PASSWORD", ""),
	}

	config.MetricsPort = getEnv("METRICS_PORT", "")
	config.BrokerConstants.EmailConfirm = "email-confirm"

	return config
}

func (dbSettings DBConfig) DBUrl() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbSettings.User,
		dbSettings.Password,
		dbSettings.Host,
		dbSettings.Port,
		dbSettings.DB,
	)
}

func (rabbitMQ RabbitMQConfig) RabbitMQUrl() string {
	port, _ := strconv.Atoi(rabbitMQ.Port)
	return fmt.Sprintf("amqp://%s:%s@%s:%d", rabbitMQ.User, rabbitMQ.Password, rabbitMQ.Host, port)
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
