package rsakeys

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

const (
	publicKeyType  = "PUBLIC KEY"
	privateKeyType = "PRIVATE KEY"
)

// ExportPublicKeyAsPemBytes создаёт байты PEM блока из публичного ключа
func ExportPublicKeyAsPemBytes(publicKey *rsa.PublicKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: publicKeyType, Bytes: x509.MarshalPKCS1PublicKey(publicKey)})
}

// ExportPrivateKeyAsPemBytes создаёт байты PEM блока из приватного ключа
func ExportPrivateKeyAsPemBytes(privateKey *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: privateKeyType, Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})
}

// ImportPublicKeyFromFile возвращает публичный ключ из файла
func ImportPublicKeyFromFile(path string) (*rsa.PublicKey, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read public key file: %w", err)
	}

	return ImportPublicKeyFromBytes(data)
}

// ImportPrivateKeyFromFile возвращает приватный ключ из файла
func ImportPrivateKeyFromFile(path string) (*rsa.PrivateKey, error) {
	if path == "" {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	return ImportPrivateKeyFromBytes(data)
}

// ImportPublicKeyFromBytes возвращает публичный ключ из байтов PEM блока
func ImportPublicKeyFromBytes(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != publicKeyType {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed parsing public key from PEM block: %w", err)
	}

	return pub, nil
}

// ImportPrivateKeyFromBytes возвращает приватный ключ из байтов PEM блока
func ImportPrivateKeyFromBytes(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != privateKeyType {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	pub, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed parsing private key from PEM block: %w", err)
	}

	return pub, nil
}
