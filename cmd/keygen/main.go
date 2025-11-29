package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
)

const (
	secretPrefix    = "whsec_"
	defaultKeyBytes = 32 // 256 bits
)

func main() {
	var (
		keyBytes int
		count    int
	)

	flag.IntVar(&keyBytes, "bytes", defaultKeyBytes, "key length in bytes (24-64)")
	flag.IntVar(&count, "n", 1, "number of keys to generate")
	flag.Parse()

	if keyBytes < 24 || keyBytes > 64 {
		fmt.Fprintf(os.Stderr, "error: key length must be between 24 and 64 bytes\n")
		os.Exit(1)
	}

	for i := 0; i < count; i++ {
		key, err := generateKey(keyBytes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to generate key: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(key)
	}
}

func generateKey(bytes int) (string, error) {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return secretPrefix + base64.StdEncoding.EncodeToString(buf), nil
}
