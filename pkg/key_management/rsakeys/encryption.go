package rsakeys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Encrypt шифрует сообщение публичным ключом и возвращает байты PEM блока типа MESSAGE
func Encrypt(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	if publicKey == nil {
		return message, nil
	}

	msgLen := len(message)
	hash := sha256.New()
	step := publicKey.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, message[start:finish], nil)
		if err != nil {
			return nil, err
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}

	return encryptedBytes, nil
}

// Decrypt расшифровывает шифр приватным ключом
func Decrypt(privateKey *rsa.PrivateKey, cipher []byte) ([]byte, error) {
	if privateKey == nil {
		return cipher, nil
	}

	msgLen := len(cipher)
	step := privateKey.PublicKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, cipher[start:finish], nil)
		if err != nil {
			return nil, err
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
