package config

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	instance *Config
	once     sync.Once
)

func Initialize() {
	once.Do(func() {
		log.Println("Init config")
		instance = loadConfig()
	})
}

func GetConfig() *Config {
	if instance == nil {
		panic("Config is not initialized. Call config.Initialize at first.")
	}
	return instance
}

type Settings struct {
	Level        LogLevel      `yaml:"level" env-required:"true"`
	IdleTimeout  time.Duration `yaml:"idleTimeout" env-default:"30s"`
	WriteTimeout time.Duration `yaml:"writeTimeout" env-default:"60s"`
	Concurrency  int           `yaml:"concurrency" env-default:"200"`
	MaxBodySize  int64         `yaml:"maxBodySize" env-default:"107374182400"`
}

type Cli struct {
	Level     LogLevel `yaml:"level" env-required:"true"`
	ChunkSize int64    `yaml:"chunkSize" env-default:"256"`
	ServerURL string   `yaml:"serverURL" env-required:"true"`
}

type Config struct {
	Env      Environment `yaml:"env" env-default:"local"`
	Settings Settings    `yaml:"settings" env-required:"true"`
	Cli      Cli         `yaml:"cli" env-required:"true"`
}

func loadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var conf Config

	if err := cleanenv.ReadConfig(configPath, &conf); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &conf
}
