package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env              string           `yaml:"env" env-default:"local"`
	StorageConfig    StorageConfig    `yaml:"storage"`
	DadataConfig     DadataConfig     `yaml:"dadata"`
	JWTConfig        JWTConfig        `yaml:"jwt"`
	HTTPServerConfig HTTPServerConfig `yaml:"http_server"`
	EmailConfig      EmailConfig      `yaml:"email"`
	BitrixConfig     BitrixConfig     `yaml:"bitrix"`
}

type StorageConfig struct {
	Path string `yaml:"path"`
}

type DadataConfig struct {
	ApiKey    string `yaml:"api_key"`
	SecretKey string `yaml:"secret_key"`
}

type JWTConfig struct {
	Access  JWTInfo `yaml:"access"`
	Refresh JWTInfo `yaml:"refresh"`
}

type JWTInfo struct {
	SecretKey  string        `yaml:"secret_key"`
	Expiration time.Duration `yaml:"expiration"`
}

type HTTPServerConfig struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"4s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"4s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type EmailConfig struct {
	SMTP EmailInfo `yaml:"smtp"`
	IMAP EmailInfo `yaml:"imap"`
	POP3 EmailInfo `yaml:"pop3"`
}

type EmailInfo struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	PortSSL  int    `yaml:"ssl_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type BitrixConfig struct {
	OutgoingWebhookAuth string `yaml:"outgoing_webhook_auth"`
	IncomingWebhook     string `yaml:"incoming_webhook"`
}

func MustLoad() *Config {
	// Определяем флаг для пути к конфигурационному файлу
	configPath := flag.String("config", "", "Path to the configuration file")
	flag.Parse()

	// Проверяем, задан ли флаг config
	if *configPath == "" {
		log.Fatal("CONFIG_PATH is not set. Use -config flag to specify the path to the configuration file.")
	}

	// Проверяем, существует ли файл
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", *configPath)
	}

	var cfg Config

	// Чтение конфигурации из файла
	if err := cleanenv.ReadConfig(*configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err)
	}

	return &cfg
}
