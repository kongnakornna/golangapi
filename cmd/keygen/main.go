// cmd/keygen/main.go
//
// Keygen tool — generates JWE (Ed25519 + X25519) and RS256 key pairs
// for the icmongolang JWT stack, printing ready-to-paste base64 values.
//
// Usage:
//
//	go run ./cmd/keygen                    # JWE only (recommended)
//	go run ./cmd/keygen -mode=rs256        # RS256 only
//	go run ./cmd/keygen -mode=all          # JWE + RS256 + pepper
//	go run ./cmd/keygen -mode=jwe -write   # write to config/config.default.yml
package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"icmongolang/pkg/jwe"
)

func main() {
	mode := flag.String("mode", "jwe", "keys to generate: jwe | rs256 | all")
	write := flag.Bool("write", false, "write generated keys into config/config.default.yml")
	flag.Parse()

	var out strings.Builder

	fmt.Fprintln(&out, "# ========================================")
	fmt.Fprintln(&out, "#  icmongolang — JWT Key Generator")
	fmt.Fprintln(&out, "# ========================================")

	// ---------------- JWE ----------------
	if *mode == "jwe" || *mode == "all" {
		kp, err := jwe.GenerateKeyPair()
		if err != nil {
			log.Fatalf("generate JWE key pair: %v", err)
		}

		signPrivB64 := base64.StdEncoding.EncodeToString(kp.SigningPrivate)
		signPubB64 := base64.StdEncoding.EncodeToString(kp.SigningPublic)
		kexPrivB64 := base64.StdEncoding.EncodeToString(kp.KEXPrivate.Bytes())
		kexPubB64 := base64.StdEncoding.EncodeToString(kp.KEXPublic.Bytes())

		fmt.Fprintln(&out, "\n# ---------- JWE (Ed25519 + X25519) ----------")
		fmt.Fprintf(&out, "JWT_JWE_SIGNING_PRIVATE_KEY=%s\n", signPrivB64)
		fmt.Fprintf(&out, "JWT_JWE_SIGNING_PUBLIC_KEY=%s\n", signPubB64)
		fmt.Fprintf(&out, "JWT_JWE_KEX_PRIVATE_KEY=%s\n", kexPrivB64)
		fmt.Fprintf(&out, "JWT_JWE_KEX_PUBLIC_KEY=%s\n", kexPubB64)
		fmt.Fprintf(&out, "JWT_JWE_ACCESS_TOKEN_TTL=15\n")
		fmt.Fprintf(&out, "JWT_JWE_REFRESH_TOKEN_TTL=1440\n")

		// sanity: signing_pub == last 32 bytes of signing_priv
		ok := len(kp.SigningPrivate) == ed25519.PrivateKeySize &&
			string(kp.SigningPrivate[32:]) == string(kp.SigningPublic)
		fmt.Fprintf(&out, "# SANITY: signing_pub == priv[32:] ? %v\n", ok)
	}

	// ---------------- RS256 ----------------
	if *mode == "rs256" || *mode == "all" {
		accessPrivB64, accessPubB64 := genRSA()
		refreshPrivB64, refreshPubB64 := genRSA()

		fmt.Fprintln(&out, "\n# ---------- RS256 access token ----------")
		fmt.Fprintf(&out, "JWT_ACCESSTOKENPRIVATEKEY=%s\n", accessPrivB64)
		fmt.Fprintf(&out, "JWT_ACCESSTOKENPUBLICKEY=%s\n", accessPubB64)
		fmt.Fprintln(&out, "\n# ---------- RS256 refresh token ----------")
		fmt.Fprintf(&out, "JWT_REFRESHTOKENPRIVATEKEY=%s\n", refreshPrivB64)
		fmt.Fprintf(&out, "JWT_REFRESHTOKENPUBLICKEY=%s\n", refreshPubB64)
	}

	// ---------------- API key pepper ----------------
	if *mode == "all" {
		pepper := make([]byte, 32)
		if _, err := rand.Read(pepper); err != nil {
			log.Fatalf("rand: %v", err)
		}
		fmt.Fprintln(&out, "\n# ---------- API key pepper ----------")
		fmt.Fprintf(&out, "JWT_API_KEY_PEPPER=%s\n", base64.StdEncoding.EncodeToString(pepper))
	}

	// ---------------- Output ----------------
	if *write {
		const target = "config/config.default.yml"
		existing, err := os.ReadFile(target)
		if err != nil {
			log.Fatalf("read %s: %v", target, err)
		}
		fmt.Fprintf(os.Stderr, "!! Auto-write to %s is not implemented yet.\n", target)
		fmt.Fprintf(os.Stderr, "!! Please copy values manually from stdout.\n")
		_ = existing
	}

	fmt.Println(out.String())
	fmt.Println("\n>>> Paste the values into .env or config/config.default.yml")
	fmt.Println(">>> Verify with:  go run ./cmd/keygen -mode=jwe  (must not change)")
}

func genRSA() (privB64, pubB64 string) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("rsa.GenerateKey: %v", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	pubDER, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		log.Fatalf("MarshalPKIXPublicKey: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER})
	return base64.StdEncoding.EncodeToString(privPEM),
		base64.StdEncoding.EncodeToString(pubPEM)
}
