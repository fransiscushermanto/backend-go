package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/fransiscushermanto/backend/internal/utils"
)

func main() {
	// Generate ECDSA key pair
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		utils.Log().Fatal().Err(err).Msg("Failed to generate ECDSA key pair")
	}

	// Marshal private key
	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		utils.Log().Fatal().Err(err).Msg("Failed to marshal private key")
	}

	// Create PEM block for private key
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// Marshal public key
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		utils.Log().Fatal().Err(err).Msg("Failed to marshal public key")
	}

	// Create PEM block for public key
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	// Write to files
	err = os.WriteFile("private_key.pem", privateKeyPEM, 0600)
	if err != nil {
		utils.Log().Fatal().Err(err).Msg("Failed to write private key to file")
	}

	err = os.WriteFile("public_key.pem", publicKeyPEM, 0644)
	if err != nil {
		utils.Log().Fatal().Err(err).Msg("Failed to write public key to file")
	}

	fmt.Println("Private Key:")
	fmt.Println(string(privateKeyPEM))
	fmt.Println("\nPublic Key:")
	fmt.Println(string(publicKeyPEM))

	fmt.Println("\nKeys saved to private_key.pem and public_key.pem")
	fmt.Println("\nFor production, set environment variable:")
	fmt.Printf("export JWT_PRIVATE_KEY=\"%s\"\n", string(privateKeyPEM))
}
