package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
)

// Transaction represents the data we want to sign
type Transaction struct {
	AccountID string `json:"account_id"`
	Amount    uint   `json:"amount"`
}

// SignedTransaction includes the transaction data and its signature
type SignedTransaction struct {
	Transaction
	Signature string `json:"signature"` // Base64 encoded signature
}

// ECDSASignature represents the R and S components of an ECDSA signature
type ECDSASignature struct {
	R, S *big.Int
}

// loadPrivateKey loads a PEM encoded EC private key from file
func loadPrivateKey(file string) (*ecdsa.PrivateKey, error) {
	pemBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading private key file: %v", err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing EC private key: %v", err)
	}

	return privateKey, nil
}

// loadPublicKey loads a PEM encoded EC public key from file
func loadPublicKey(file string) (*ecdsa.PublicKey, error) {
	pemBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading public key file: %v", err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	// Parse the public key
	genericPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error parsing EC public key: %v", err)
	}

	// Assert that the key is an ECDSA public key
	publicKey, ok := genericPublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an EC public key")
	}

	return publicKey, nil
}

// signTransaction signs a transaction with an EC private key
func signTransaction(tx Transaction, privateKey *ecdsa.PrivateKey) (*SignedTransaction, error) {
	// Marshal the transaction to JSON
	txJSON, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("error marshaling transaction: %v", err)
	}

	// Calculate SHA-256 hash of the JSON
	hash := sha256.Sum256(txJSON)

	// Sign the hash
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %v", err)
	}

	// Create and marshal ECDSA signature
	signature := ECDSASignature{
		R: r,
		S: s,
	}
	signatureBytes, err := asn1.Marshal(signature)
	if err != nil {
		return nil, fmt.Errorf("error marshaling signature: %v", err)
	}

	// Create the signed transaction
	signedTx := &SignedTransaction{
		Transaction: tx,
		Signature:   base64.StdEncoding.EncodeToString(signatureBytes),
	}

	return signedTx, nil
}

// verifySignature verifies the signature of a signed transaction
func verifySignature(signedTx *SignedTransaction, publicKey *ecdsa.PublicKey) (bool, error) {
	// Create a transaction without the signature to match the original data that was signed
	tx := signedTx.Transaction

	// Marshal the transaction to JSON
	txJSON, err := json.Marshal(tx)
	if err != nil {
		return false, fmt.Errorf("error marshaling transaction: %v", err)
	}

	// Calculate SHA-256 hash of the JSON
	hash := sha256.Sum256(txJSON)

	// Decode the signature from Base64
	signatureBytes, err := base64.StdEncoding.DecodeString(signedTx.Signature)
	if err != nil {
		return false, fmt.Errorf("error decoding signature: %v", err)
	}

	// Unmarshal the signature
	var signature ECDSASignature
	if _, err := asn1.Unmarshal(signatureBytes, &signature); err != nil {
		return false, fmt.Errorf("error unmarshaling signature: %v", err)
	}

	// Verify the signature
	return ecdsa.Verify(publicKey, hash[:], signature.R, signature.S), nil
}

func main() {
	cmdMain()
	// Example usage
	// Create a new transaction
	tx := Transaction{
		AccountID: "9eb95936-3e59-433d-9f71-a054a89f7cd3",
		Amount:    1000,
	}

	// For demonstration purposes, we'll generate a new key pair if the file doesn't exist
	// In production, you'd use your existing key
	privateKeyFile := "./atm_signing_key.dev"
	publicKeyFile := "./atm_signing_key.pub.dev"

	var privateKey *ecdsa.PrivateKey
	var err error

	// Try to load the private key
	privateKey, err = loadPrivateKey(privateKeyFile)
	if err != nil {
		// Generate a new private key
		fmt.Println("Private key file not found or invalid, generating a new key pair...")
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Fatalf("Error generating private key: %v", err)
		}

		// Save the private key
		privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			log.Fatalf("Error marshaling private key: %v", err)
		}
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privateKeyBytes,
		})
		if err := ioutil.WriteFile(privateKeyFile, privateKeyPEM, 0600); err != nil {
			log.Fatalf("Error writing private key to file: %v", err)
		}

		// Save the public key
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			log.Fatalf("Error marshaling public key: %v", err)
		}
		publicKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		})
		if err := ioutil.WriteFile(publicKeyFile, publicKeyPEM, 0644); err != nil {
			log.Fatalf("Error writing public key to file: %v", err)
		}
		fmt.Println("Generated and saved new key pair.")
	}

	// Sign the transaction
	signedTx, err := signTransaction(tx, privateKey)
	if err != nil {
		log.Fatalf("Error signing transaction: %v", err)
	}

	// Print the signed transaction
	signedJSON, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling signed transaction: %v", err)
	}
	fmt.Println("Signed Transaction:")
	fmt.Println(string(signedJSON))

	// Verify the signature
	publicKey, err := loadPublicKey(publicKeyFile)
	if err != nil {
		log.Fatalf("Error loading public key: %v", err)
	}

	valid, err := verifySignature(signedTx, publicKey)
	if err != nil {
		log.Fatalf("Error verifying signature: %v", err)
	}

	if valid {
		fmt.Println("Signature verification successful!")
	} else {
		fmt.Println("Signature verification failed!")
	}
}

// To use this program with command-line arguments
func cmdMain() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: program <account_id> <amount>")
		os.Exit(1)
	}

	accountID := os.Args[1]
	var amount uint
	fmt.Sscanf(os.Args[2], "%d", &amount)

	// Create a transaction
	tx := Transaction{
		AccountID: accountID,
		Amount:    amount,
	}

	// Load the private key
	privateKey, err := loadPrivateKey("./atm_signing_key.dev")
	if err != nil {
		log.Fatalf("Error loading private key: %v", err)
	}

	// Sign the transaction
	signedTx, err := signTransaction(tx, privateKey)
	if err != nil {
		log.Fatalf("Error signing transaction: %v", err)
	}

	// Print the signed transaction
	signedJSON, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling signed transaction: %v", err)
	}
	fmt.Println(string(signedJSON))
}
