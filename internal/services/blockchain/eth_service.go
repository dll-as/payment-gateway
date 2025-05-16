package blockchain

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rezatg/payment-gateway/config"
	"github.com/rezatg/payment-gateway/pkg/errors"
	"github.com/rezatg/payment-gateway/pkg/logger"
)

// EthService implements BlockchainService for Ethereum
type EthService struct {
	client *ethclient.Client
}

// NewEthService creates a new Ethereum service
func NewEthService() *EthService {
	infuraURL := config.GetEnv("INFURA_URL", "https://mainnet.infura.io/v3/your-infura-key")
	client, err := ethclient.Dial(infuraURL)
	if err != nil {
		logger.Fatal("Failed to connect to Ethereum client", err)
	}
	return &EthService{client: client}
}

// GenerateAddress generates a new Ethereum address
func (s *EthService) GenerateAddress() (string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		logger.Error("Failed to generate Ethereum private key", err)
		return "", errors.NewInternalServerError("Failed to generate address", err.Error())
	}

	publicKey := privateKey.PublicKey
	address := crypto.PubkeyToAddress(publicKey)
	logger.Info("Generated Ethereum address", "address", address.Hex())
	return address.Hex(), nil
}

// GetBalance retrieves the balance of an Ethereum address
func (s *EthService) GetBalance(address string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr := common.HexToAddress(address)
	balance, err := s.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		logger.Error("Failed to get Ethereum balance", err, "address", address)
		return 0, errors.NewInternalServerError("Failed to get balance", err.Error())
	}

	// Convert balance from Wei to Ether
	etherBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	balanceFloat, _ := etherBalance.Float64()
	logger.Info("Retrieved Ethereum balance", "address", address, "balance", balanceFloat)
	return balanceFloat, nil
}

// SendTransaction sends an Ethereum transaction
func (s *EthService) SendTransaction(fromPrivKey, toAddress string, amount float64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	privateKey, err := crypto.HexToECDSA(fromPrivKey)
	if err != nil {
		logger.Error("Invalid private key", err)
		return "", errors.NewBadRequestError("Invalid private key")
	}

	publicKey := privateKey.PublicKey
	fromAddress := crypto.PubkeyToAddress(publicKey)
	nonce, err := s.client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		logger.Error("Failed to get nonce", err, "address", fromAddress.Hex())
		return "", errors.NewInternalServerError("Failed to send transaction", err.Error())
	}

	value := new(big.Int).Mul(big.NewInt(int64(amount*1e18)), big.NewInt(1)) // Convert to Wei
	gasLimit := uint64(21000)
	gasPrice, err := s.client.SuggestGasPrice(ctx)
	if err != nil {
		logger.Error("Failed to get gas price", err)
		return "", errors.NewInternalServerError("Failed to send transaction", err.Error())
	}

	toAddr := common.HexToAddress(toAddress)
	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, nil)

	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		logger.Error("Failed to sign transaction", err)
		return "", errors.NewInternalServerError("Failed to send transaction", err.Error())
	}

	err = s.client.SendTransaction(ctx, signedTx)
	if err != nil {
		logger.Error("Failed to send transaction", err)
		return "", errors.NewInternalServerError("Failed to send transaction", err.Error())
	}

	logger.Info("Sent Ethereum transaction", "tx_hash", signedTx.Hash().Hex(), "from", fromAddress.Hex(), "to", toAddress)
	return signedTx.Hash().Hex(), nil
}

// CheckTxConfirmed checks if an Ethereum transaction is confirmed
func (s *EthService) CheckTxConfirmed(txHash string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hash := common.HexToHash(txHash)
	receipt, err := s.client.TransactionReceipt(ctx, hash)
	if err != nil {
		logger.Error("Failed to get transaction receipt", err, "tx_hash", txHash)
		return false, errors.NewInternalServerError("Failed to check transaction", err.Error())
	}

	confirmed := receipt != nil && receipt.Status == 1
	logger.Info("Checked Ethereum transaction", "tx_hash", txHash, "confirmed", confirmed)
	return confirmed, nil
}
