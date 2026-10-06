package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"icmongolang/pkg/jwe"
)

func main() {
	kp, err := jwe.GenerateKeyPair()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating key pair: %v\n", err)
		os.Exit(1)
	}

	signingPrivB64 := base64.StdEncoding.EncodeToString(kp.SigningPrivate)
	signingPubB64 := base64.StdEncoding.EncodeToString(kp.SigningPublic)
	kexPrivBytes := kp.KEXPrivate.Bytes()
	kexPrivB64 := base64.StdEncoding.EncodeToString(kexPrivBytes)
	kexPubBytes := kp.KEXPublic.Bytes()
	kexPubB64 := base64.StdEncoding.EncodeToString(kexPubBytes)

	fmt.Println("# JWE Keys for config/config.yml")
	fmt.Println("# Add these to your jwt section:")
	fmt.Println()
	fmt.Printf("jwe_signing_private_key: \"%s\"\n", signingPrivB64)
	fmt.Printf("jwe_signing_public_key: \"%s\"\n", signingPubB64)
	fmt.Printf("jwe_kex_private_key: \"%s\"\n", kexPrivB64)
	fmt.Printf("jwe_kex_public_key: \"%s\"\n", kexPubB64)
	fmt.Println()
	fmt.Println("# Also generate a pepper for API keys (min 32 bytes):")
	fmt.Println("# api_key_pepper: \"your-32-byte-pepper-here!!\"")
}
