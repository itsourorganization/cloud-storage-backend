package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env      string `yaml:"env" env-default:"local"`
	Http     `yaml:"http"`
	Jwt      `yaml:"jwt"`
	Database `yaml:"database"`
}

type Http struct {
	BindAddr    string        `yaml:"bind_addr" env-default:":8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"30s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"30s"`
}

type Jwt struct {
	AccessExpire  time.Duration `yaml:"access_expire" env-default:"20m"`
	AccessSecret  string        `env:"ACCESS_SECRET" env-required:"true"`
	RefreshExpire time.Duration `yaml:"refresh_expire" env-default:"144h"`
	RefreshSecret string        `env:"REFRESH_SECRET" env-required:"true"`
}

type Database struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	DbName   string `yaml:"db_name" env-required:"true"`
	Password string `env:"DATABASE_PASSWORD" env-required:"true"`
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
