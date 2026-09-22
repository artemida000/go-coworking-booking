package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Confing struct {
	Port string
	Env string

	DBHost string
	DBPort string
	DBUser string
	DBPassword string
	DBName string
	DBSslmode string 

	JWTSecret string 
}

func Load() (*Confing, error) {

	_ = godotenv.Load()
	
	cfg := &Confing{
		Port: os.Getenv("PORT"),
		Env: os.Getenv("ENV"),
		DBHost: os.Getenv("DB_HOST"),
		DBPort: os.Getenv("DB_PORT"),
		DBUser: os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName: os.Getenv("DB_NAME"),
		DBSslmode: os.Getenv("DB_SSLMODE"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	return cfg, nil
}