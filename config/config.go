package config

import "log/slog"

type Config struct {
	LogLevel    string
	Environment map[string]Environment
}

type Environment struct {
	ServerAddr  string
	DatabaseURL string
}

func (c *Config) Env() Environment {
	return Environment{
		ServerAddr:  "localhost:8080",
		DatabaseURL: "",
	}
}

func ConfigureLogger(config *Config) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
