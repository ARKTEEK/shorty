package config

import "os"

type Config struct {
	Addr        string
	DatabaseUrl string
}

func Load() *Config {
	return &Config{
		Addr:        getEnv("ADDR", ":8080"),
		DatabaseUrl: getEnv("DATABASE_URL", "root:Pass123..@tcp(localhost:3306)/shorty?parseTime=true"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
