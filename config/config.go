package config

import (
	"log"
	"os"
	"runtime"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Redis RedisConfig
	Port  string
}

type RedisConfig struct {
	Password string
	Addr     string
	Username string
	Db       int
}

type ApiContext struct {
	Config Config
	Client *redis.Client
}

func Get() Config {

	db := 0
	if runtime.GOOS == "windows" {
		db = 1
	}

	return Config{
		Redis: RedisConfig{
			Password: getEnv("REDIS_PW"),
			Addr:     getEnv("REDIS_ADDR"),
			Username: getEnv("REDIS_USR"),
			Db:       db - 1,
		},
		Port: getEnv("PORT"),
	}
}

func getEnv(key string) string {
	value, set := os.LookupEnv(key)
	if !set {
		log.Fatalf("Config variable %s was missing\n", key)
	}
	return value
}
