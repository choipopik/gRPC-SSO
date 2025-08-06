package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env           string `yaml:"env" env-default:"local"`
	Storage_path  string `yaml:"storage_path" env-required:"true"` //Порт gRPC + timeout
	MigrationPath string
	TokenTTL      time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC          GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config path not exist: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("empty config path " + err.Error())
	}

	return &cfg
}

// Запуск приложения через флаг sso --config=./config/local.yaml
// или через CONFIG_PATH=./config/local.yaml sso
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
