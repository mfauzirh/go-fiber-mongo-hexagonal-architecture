package config

import (
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App         *App
		ProfilingDB *ProfilingDB
		HTTP        *HTTP
	}

	App struct {
		Name string
		Env  string
	}

	ProfilingDB struct {
		URI string
	}

	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}
)

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: os.Getenv("APP_NAME"),
		Env:  os.Getenv("APP_ENV"),
	}

	db := &ProfilingDB{
		URI: os.Getenv("MONGODB_URI"),
	}

	http := &HTTP{
		Env:            os.Getenv("APP_ENV"),
		URL:            os.Getenv("HTTP_URL"),
		Port:           os.Getenv("HTTP_PORT"),
		AllowedOrigins: os.Getenv("HTTP_ALLOWED_ORIGINS"),
	}

	return &Container{
		app,
		db,
		http,
	}, nil
}
