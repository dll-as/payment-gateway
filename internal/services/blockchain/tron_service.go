package blockchain

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/bytedance/sonic"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mr-tron/base58"
	"github.com/rezatg/payment-gateway/pkg/errors"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/sha3"
)

const tronGridAPI = "https://api.trongrid.io/v1/accounts"

// TronService handles Tron blockchain operations
type TronService struct {
	encryptionKey []byte // Key for encrypting private keys
}

// NewTronService creates a new TronService with encryption key
func NewTronService(encryptionKey []byte) *TronService {
	if len(encryptionKey) != 32 { // AES-256 requires 32-byte key
		panic("Encryption key must be 32 bytes")
	}
	return &TronService{encryptionKey: encryptionKey}
}

// GenerateAddress generates a new Tron address and encrypted private key
func (t *TronService) GenerateAddress() (string, string, error) {
	// Generate private key
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", errors.NewInternalServerError("Failed to generate private key: " + err.Error())
	}

	// Get public key
	publicKey := privateKey.PublicKey
	publicKeyBytes := crypto.FromECDSAPub(&publicKey)

	// Hash public key
	hash := sha3.NewLegacyKeccak256()
	hash.Write(publicKeyBytes[1:]) // Remove 0x04 prefix
	hashedPublicKey := hash.Sum(nil)

	// Create address (0x41 prefix + last 20 bytes of hash)
	address := append([]byte{0x41}, hashedPublicKey[len(hashedPublicKey)-20:]...)

	// Generate checksum
	firstHash := sha256.Sum256(address)
	secondHash := sha256.Sum256(firstHash[:])
	checksum := secondHash[:4]

	// Create final address
	fullAddress := append(address, checksum...)
	base58Address := base58.Encode(fullAddress)

	// Validate address
	if !isValidTronAddress(base58Address) {
		return "", "", errors.NewInternalServerError("Generated invalid Tron address")
	}

	// Encrypt private key
	privateKeyBytes := crypto.FromECDSA(privateKey)
	encryptedPrivKey, err := t.encryptPrivateKey(privateKeyBytes)
	if err != nil {
		return "", "", errors.NewInternalServerError("Failed to encrypt private key: " + err.Error())
	}

	return base58Address, encryptedPrivKey, nil
}

// encryptPrivateKey encrypts the private key using AES-256-GCM
func (t *TronService) encryptPrivateKey(privateKey []byte) (string, error) {
	block, err := aes.NewCipher(t.encryptionKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, privateKey, nil)
	return hex.EncodeToString(ciphertext), nil
}

// decryptPrivateKey decrypts the private key (for SendTransaction)
func (t *TronService) decryptPrivateKey(encryptedKey string) ([]byte, error) {
	ciphertext, err := hex.DecodeString(encryptedKey)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(t.encryptionKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.NewInternalServerError("Invalid ciphertext")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.NewInternalServerError("Failed to decrypt private key")
	}

	return plaintext, nil
}

// isValidTronAddress validates a Tron address
func isValidTronAddress(address string) bool {
	decoded, err := base58.Decode(address)
	if err != nil || len(decoded) != 25 || decoded[0] != 0x41 {
		return false
	}

	firstHash := sha256.Sum256(decoded[:21])
	secondHash := sha256.Sum256(firstHash[:])
	checksum := secondHash[:4]

	for i, b := range checksum {
		if b != decoded[21+i] {
			return false
		}
	}

	return true
}

func (t *TronService) GetBalance(address string) (float64, error) {
	statusCode, body, err := fasthttp.Get(nil, fmt.Sprintf("%s/%s", tronGridAPI, address))
	if err != nil {
		return 0, err
	}

	if statusCode != fasthttp.StatusOK {
		return 0, fmt.Errorf("error retrieving account information, HTTP status: %d", statusCode)
	}

	var response TronResponse
	if err := sonic.Unmarshal(body, &response); err != nil || len(response.Data) == 0 {
		return 0, err
	}

	return float64(response.Data[0].Balance) / 1_000_000, nil
}

func (t *TronService) SendTransaction(fromPrivKey, toAddress string, amount float64) (string, error) {
	return "trx_hash_123", nil
}

func (t *TronService) CheckTxConfirmed(txHash string) (bool, error) {
	// call TronGrid to check confirmations
	return true, nil
}
