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
)

func setConfigFromEnv(config *AppConfig) error {
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

	return nil
}

func parsePrivateKey(keyData string) (*ecdsa.PrivateKey, error) {
	if strings.TrimSpace(keyData) == "" {
		return nil, fmt.Errorf("private key data is empty")
	}

	// Support both raw PEM and base64-encoded PEM
	keyData = strings.ReplaceAll(keyData, "\\n", "\n")

	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	if block.Type != "EC PRIVATE KEY" && block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("expected EC PRIVATE KEY or PRIVATE KEY, got %s", block.Type)
	}

	var privateKey *ecdsa.PrivateKey
	var err error

	if block.Type == "EC PRIVATE KEY" {
		privateKey, err = x509.ParseECPrivateKey(block.Bytes)
	} else {
		key, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", parseErr)
		}

		var ok bool
		privateKey, ok = key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("not an ECDSA private key")
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse EC private key: %w", err)
	}

	if privateKey.Curve.Params().BitSize < 256 {
		return nil, fmt.Errorf("private key must use at least P-256 curve for security")
	}

	return privateKey, nil
}

func parsePublicKey(keyData string) (*ecdsa.PublicKey, error) {
	if strings.TrimSpace(keyData) == "" {
		return nil, fmt.Errorf("public key data is empty")
	}

	keyData = strings.ReplaceAll(keyData, "\\n", "\n")

	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("expected PUBLIC KEY, got %s", block.Type)
	}

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}

	if publicKey.Curve.Params().BitSize < 256 {
		return nil, fmt.Errorf("public key must use at least P-256 curve for security")
	}

	return publicKey, nil
}
