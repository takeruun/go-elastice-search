package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		Host     string
		Username string
		Password string
		DBName   string
	}
}

func NewConfig() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("env file 読み込み出来ませんでした。")
	}

	c := new(Config)

	c.DB.Host = os.Getenv("DB_HOST")
	c.DB.Username = os.Getenv("POSTGRES_USER")
	c.DB.Password = os.Getenv("POSTGRES_PASSWORD")
	c.DB.DBName = os.Getenv("POSTGRES_DB")

	return c
}
