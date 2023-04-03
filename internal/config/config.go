package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	VultrApiKey string
	Port        int
}

func LoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		panic(err.Error())
	}

	databaseUrl := os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		panic(fmt.Errorf("please provide DATABASE_URL"))
	}

	port := 3000

	if os.Getenv("PORT") != "" {
		if val, err := strconv.Atoi(os.Getenv("PORT")); err != nil {
			panic(err)
		} else {
			port = val
		}
	}

	vultrApiKey := os.Getenv("VULTR_API_KEY")

	if vultrApiKey == "" {
		panic(fmt.Errorf("please provide VULTR_API_KEY"))
	}

	config := &Config{
		DatabaseUrl: databaseUrl,
		VultrApiKey: vultrApiKey,
		Port:        port,
	}

	return config
}
