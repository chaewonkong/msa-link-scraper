package config

import "os"

// App represents the configuration for the application
type App struct {
	APIHost string `json:"api_host"`
	Queue   queue  `json:"queue"`
}

type queue struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Port     string `json:"port"`
	Name     string `json:"name"`
}

// NewAppConfig creates a new DBConfig
func NewAppConfig() *App {
	apiHost := os.Getenv("API_HOST")
	queue := queue{
		Host:     os.Getenv("QUEUE_HOST"),
		User:     os.Getenv("QUEUE_USER"),
		Password: os.Getenv("QUEUE_PASSWORD"),
		Port:     os.Getenv("QUEUE_PORT"),
		Name:     os.Getenv("QUEUE_NAME"),
	}

	return &App{apiHost, queue}
}
