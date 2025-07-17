package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Env             string           `yaml:"env" env:"APP_ENV"`
	Port            int              `yaml:"port" env:"APP_PORT"`
	DatabaseURL     string           `yaml:"database_url" env:"DATABASE_URL"`
	ShutdownTimeout int              `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`
	LogLevel        string           `yaml:"log_level" env:"LOG_LEVEL"`
	AllowedOrigins  []string         `yaml:"allowed_origins" env:"ALLOWED_ORIGINS"`
	SecretKey       string           `yaml:"secret_key" env:"SECRET_KEY"`
	PrivateKey      ecdsa.PrivateKey `yaml:"private_key" env:"PRIVATE_KEY"`
	PublicKey       ecdsa.PublicKey  `yaml:"public_key" env:"PUBLIC_KEY"`
}

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{
		Env:             "development",
		Port:            8080,
		ShutdownTimeout: 5,
		LogLevel:        "info",
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	if data, err := os.ReadFile(configPath); err == nil {

		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %s %w", configPath, err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to read config file: %s %w", configPath, err)
	}

	setConfigFromEnv(config)

	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set via config file or environment variable. This is a mandatory setting")
	}

	return config, nil
}

func setConfigFromEnv(config *AppConfig) {
	v := reflect.ValueOf(config).Elem()
	t := reflect.TypeOf(config).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		envTag := fieldType.Tag.Get("env")
		if envTag == "" {
			continue
		}

		envValue := os.Getenv(envTag)
		if envValue == "" {
			continue
		}

		if envTag == "PRIVATE_KEY" && envValue != "" {
			privateKey, err := parsePrivateKey(envValue)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to parse private key: %v\n", err)
				continue
			}
			// Set the private key field directly
			if field.CanSet() {
				field.Set(reflect.ValueOf(*privateKey))
			}
			continue
		}

		if envTag == "PUBLIC_KEY" && envValue != "" {
			publicKey, err := parsePublicKey(envValue)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to parse public key: %v\n", err)
				continue
			}
			// Set the public key field directly
			if field.CanSet() {
				field.Set(reflect.ValueOf(*publicKey))
			}
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(envValue)
		case reflect.Int:
			if value, err := strconv.Atoi(envValue); err == nil {
				field.SetInt(int64(value))
			} else {
				fmt.Fprintf(os.Stderr, "Warning: %s '%s' is not a valid integer, using default %d\n",
					envTag, envValue, field.Int())
			}
		case reflect.Slice:
			items := strings.Split(envValue, ",")
			slice := reflect.MakeSlice(field.Type(), len(items), len(items))

			for j, item := range items {
				slice.Index(j).SetString(strings.TrimSpace(item))
			}

			field.Set(slice)
		}
	}
}

// Add these helper functions
func parsePrivateKey(keyData string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	return privateKey, nil
}

func parsePublicKey(keyData string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	return publicKey, nil
}
