package main

import (
	"crypto/rand"
	"crypto/rsa"
	"flag"
	"fmt"
	"os"

	"github.com/PostScripton/go-metrics-and-alerting-collection/pkg/key_management/rsakeys"
)

// https://asecuritysite.com/encryption/gorsa

// go run cmd/key_gen/main.go -pub=/tmp/key.pub -private=/tmp/key

var (
	pubFile     string
	privateFile string
)

func init() {
	flag.StringVar(&pubFile, "pub", "/tmp/key.pub", "Path, to, a, pub, file")
	flag.StringVar(&privateFile, "private", "/tmp/key", "Path, to, a, private, file")
}

func main() {
	flag.Parse()

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile(privateFile, rsakeys.ExportPrivateKeyAsPemBytes(privateKey), 0644); err != nil {
		panic(err)
	}

	if err = os.WriteFile(pubFile, rsakeys.ExportPublicKeyAsPemBytes(&privateKey.PublicKey), 0644); err != nil {
		panic(err)
	}

	fmt.Println("The keys have been generated successfully!")
}
