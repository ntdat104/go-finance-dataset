package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"sync"
)

type KeyBuffer struct {
	mu            sync.Mutex
	keys512       chan KeyPair
	keys1024      chan KeyPair
	keys2048      chan KeyPair
	keys4096      chan KeyPair
	minBufferSize int
	maxBufferSize int
}

func NewKeyBuffer(min, max int) *KeyBuffer {
	kb := &KeyBuffer{
		keys512:       make(chan KeyPair, max),
		keys1024:      make(chan KeyPair, max),
		keys2048:      make(chan KeyPair, max),
		keys4096:      make(chan KeyPair, max),
		minBufferSize: min,
		maxBufferSize: max,
	}
	go kb.fillBuffer(512)
	go kb.fillBuffer(1024)
	go kb.fillBuffer(2048)
	go kb.fillBuffer(4096)
	return kb
}

func (kb *KeyBuffer) GetKeyPair(keySize int) (KeyPair, error) {
	var ch chan KeyPair

	switch keySize {
	case 512:
		ch = kb.keys512
	case 1024:
		ch = kb.keys1024
	case 2048:
		ch = kb.keys2048
	case 4096:
		ch = kb.keys4096
	default:
		return KeyPair{}, fmt.Errorf("unsupported key size: %d", keySize)
	}

	select {
	case kp := <-ch:
		if len(ch) < kb.minBufferSize {
			go kb.fillBuffer(keySize)
		}
		return kp, nil
	default:
		log.Printf("Buffer empty for %d-bit. Generating one synchronously...", keySize)
		kp, err := generateSingleKeyPair(keySize)
		if err != nil {
			return KeyPair{}, err
		}
		go kb.fillBuffer(keySize)
		return kp, nil
	}
}

func (kb *KeyBuffer) fillBuffer(keySize int) {
	kb.mu.Lock()
	defer kb.mu.Unlock()

	var ch chan KeyPair
	switch keySize {
	case 512:
		ch = kb.keys512
	case 1024:
		ch = kb.keys1024
	case 2048:
		ch = kb.keys2048
	case 4096:
		ch = kb.keys4096
	default:
		log.Printf("Unsupported key size: %d", keySize)
		return
	}

	for len(ch) < kb.maxBufferSize {
		kp, err := generateSingleKeyPair(keySize)
		if err != nil {
			log.Printf("Error generating %d-bit RSA key: %v", keySize, err)
			continue
		}
		ch <- kp
	}
}

func generateSingleKeyPair(keySize int) (KeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return KeyPair{}, err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return KeyPair{}, err
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})

	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return KeyPair{}, err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})

	return KeyPair{PublicKey: string(pubPEM), PrivateKey: string(privPEM)}, nil
}

// ParsePrivateKey decodes a PEM-encoded private key string into *rsa.PrivateKey
func ParsePrivateKey(pemStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("invalid private key PEM")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("not an RSA private key")
	}
	return priv, nil
}

// ParsePublicKey decodes a PEM-encoded public key string into *rsa.PublicKey
func ParsePublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("invalid public key PEM")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return pub, nil
}

// getHashFunc returns the hash function and its identifier by algorithm name.
func getHashFunc(algorithm string) (crypto.Hash, error) {
	switch algorithm {
	case "SHA256":
		return crypto.SHA256, nil
	case "SHA384":
		return crypto.SHA384, nil
	case "SHA512":
		return crypto.SHA512, nil
	default:
		return 0, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

// SignMessage signs a message with the private key and returns base64 signature.
func SignMessage(privateKeyPEM, message, algorithm string) (string, error) {
	privateKey, err := ParsePrivateKey(privateKeyPEM)
	if err != nil {
		return "", err
	}

	hashFunc, err := getHashFunc(algorithm)
	if err != nil {
		return "", err
	}

	h := hashFunc.New()
	h.Write([]byte(message))
	hashed := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, hashFunc, hashed)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySignature verifies a message with the public key and base64 signature.
func VerifySignature(publicKeyPEM, message, signatureBase64, algorithm string) error {
	publicKey, err := ParsePublicKey(publicKeyPEM)
	if err != nil {
		return err
	}

	hashFunc, err := getHashFunc(algorithm)
	if err != nil {
		return err
	}

	h := hashFunc.New()
	h.Write([]byte(message))
	hashed := h.Sum(nil)

	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return err
	}

	return rsa.VerifyPKCS1v15(publicKey, hashFunc, hashed, signature)
}
