package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-default:"development"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"0.0.0.0"`
	Port        string        `yaml:"port" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	_ = godotenv.Load()

	cfg := Config{
		Env:         "production",
		StoragePath: "",
		HTTPServer:  HTTPServer{},
	}
	Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseDSN: os.Getenv("DATABASE_URL"),
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("Error opening config file: %s", err)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	return &cfg
}
