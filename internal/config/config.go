package config

import (
	"crypto/ecdsa"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Env             string   `yaml:"env" env:"APP_ENV"`
	Port            int      `yaml:"port" env:"APP_PORT"`
	DatabaseURL     string   `yaml:"database_url" env:"DATABASE_URL"`
	ShutdownTimeout int      `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT"`
	LogLevel        string   `yaml:"log_level" env:"LOG_LEVEL"`
	AllowedOrigins  []string `yaml:"allowed_origins" env:"ALLOWED_ORIGINS"`
	SecretKey       string   `yaml:"secret_key" env:"SECRET_KEY"`
	PrefixApiKey    string   `yaml:"prefix_api_key" env:"PREFIX_API_KEY"`
	LockTimeout     int      `yaml:"lock_timeout" env:"LOCK_TIMEOUT"`
	SSLCertPath     string   `yaml:"ssl_cert_path" env:"SSL_CERT_PATH"`
	SSLKeyPath      string   `yaml:"ssl_key_path" env:"SSL_KEY_PATH"`
}

type CryptoKeys struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func LoadConfig() (*AppConfig, error) {
	config := &AppConfig{
		Env:             "development",
		Port:            8080,
		ShutdownTimeout: 5,
		LogLevel:        "info",
		LockTimeout:     30,
		PrefixApiKey:    "aik_",
		SSLCertPath:     "ssl/cert.pem",
		SSLKeyPath:      "ssl/key.pem",
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

	if err := setConfigFromEnv(config); err != nil {
		return nil, fmt.Errorf("environment configuration failed: %w", err)
	}

	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set via config file or environment variable. This is a mandatory setting")
	}

	return config, nil
}

func LoadCryptoKeys() (*CryptoKeys, error) {
	privateKeyData := os.Getenv("PRIVATE_KEY")
	publicKeyData := os.Getenv("PUBLIC_KEY")

	if privateKeyData == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable is required")
	}

	if publicKeyData == "" {
		return nil, fmt.Errorf("PUBLIC_KEY environment variable is required")
	}

	privateKey, err := parsePrivateKey(privateKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	publicKey, err := parsePublicKey(publicKeyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	if !privateKey.PublicKey.Equal(publicKey) {
		return nil, fmt.Errorf("private key and public key do not form a valid pair")
	}

	return &CryptoKeys{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// Override String method to prevent accidental exposure
func (ck *CryptoKeys) String() string {
	return "CryptoKeys{[REDACTED]}"
}

// Override GoString method (used by fmt.Printf with %#v)
func (ck *CryptoKeys) GoString() string {
	return ck.String()
}

// Custom MarshalJSON to prevent JSON exposure
func (ck *CryptoKeys) MarshalJSON() ([]byte, error) {
	return []byte(`"[REDACTED]"`), nil
}

// Safe method to check if keys are loaded
func (ck *CryptoKeys) IsValid() bool {
	if ck == nil || ck.PrivateKey == nil || ck.PublicKey == nil {
		return false
	}

	// Extra check: Ensure keys still form a valid pair
	return ck.PrivateKey.PublicKey.Equal(ck.PublicKey)
}

// Safe method to get public key info (for debugging)
func (ck *CryptoKeys) PublicKeyInfo() string {
	if ck == nil || ck.PublicKey == nil {
		return "no public key"
	}
	return fmt.Sprintf("ECDSA-%d", ck.PublicKey.Curve.Params().BitSize)
}

// Override String method to prevent accidental exposure of AppConfig secrets
func (ac *AppConfig) String() string {
	return fmt.Sprintf("AppConfig{Env:%s, Port:%d, LogLevel:%s, [SECRETS REDACTED]}",
		ac.Env, ac.Port, ac.LogLevel)
}
