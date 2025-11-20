package config

import (
	"errors"
	"flag"
	"os"
)

type Config struct {
	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	JSONSchemaPath string
	Port           string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.DBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DBPort, "db-port", "", "Database port")
	flag.StringVar(&cfg.DBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DBPassword, "db-psw", "", "Database password")
	flag.StringVar(&cfg.DBName, "db-name", "", "Database name")
	flag.StringVar(&cfg.Port, "port", "", "Application port")
	flag.StringVar(&cfg.JSONSchemaPath, "json-schema", "", "Path to JSON schema file")

	flag.Parse()

	if cfg.DBHost == "" {
		cfg.DBHost = getEnv("DB_HOST", "localhost")
	}
	if cfg.DBPort == "" {
		cfg.DBPort = getEnv("DB_PORT", "3306")
	}
	if cfg.DBUser == "" {
		cfg.DBUser = getEnv("DB_USER", "")
	}
	if cfg.DBPassword == "" {
		cfg.DBPassword = getEnv("DB_PSW", "")
	}
	if cfg.DBName == "" {
		cfg.DBName = getEnv("DB_NAME", "")
	}
	if cfg.Port == "" {
		cfg.Port = getEnv("PORT", "8080")
	}
	if cfg.JSONSchemaPath == "" {
		cfg.JSONSchemaPath = getEnv("JSON_SCHEMA", "")
	}

	if cfg.DBName == "" {
		return nil, errors.New("DB_NAME é obrigatória (use --db-name ou variável de ambiente DB_NAME)")
	}
	if cfg.DBUser == "" {
		return nil, errors.New("DB_USER é obrigatária (use --db-user ou variável de ambiente DB_USER)")
	}
	if cfg.JSONSchemaPath == "" {
		return nil, errors.New("JSON_SCHEMA é obrigatária (use --json-schema ou variável de ambiente JSON_SCHEMA)")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
