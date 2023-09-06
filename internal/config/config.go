package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Http `yaml:"http"`
}

type Http struct {
	BindAddr string `yaml:"bind_addr"`
}

func MustLoad() Config {
	p, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		log.Fatal("env CONFIG_PATH not provided")
	}
	var c Config
	err := cleanenv.ReadConfig(p, &c)
	if err != nil {
		log.Fatalf("config read error: %s", err)
	}
	return c
}
