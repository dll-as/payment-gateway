package blockchain

import "github.com/rezatg/payment-gateway/pkg/errors"

// BlockchainService defines the interface for blockchain operations
type BlockchainService interface {
	GenerateAddress() (address string, encryptedPrivateKey string, err error)
	GetBalance(address string) (float64, error)
	SendTransaction(fromPrivKey, toAddress string, amount float64) (string, error)
	CheckTxConfirmed(txHash string) (bool, error)
}

// Factory defines the interface for creating blockchain services
type Factory interface {
	GetService(currency string) (BlockchainService, error)
}

// blockchainFactory implements Factory
type blockchainFactory struct {
	ethService  BlockchainService
	tronService BlockchainService
}

// NewFactory creates a new blockchain factory
func NewFactory(encryptionKey []byte) Factory {
	return &blockchainFactory{
		ethService:  NewEthService(),
		tronService: NewTronService(encryptionKey),
	}
}

// GetService returns the appropriate blockchain service based on currency
func (f *blockchainFactory) GetService(currency string) (BlockchainService, error) {
	switch currency {
	case "ETH":
		return f.ethService, nil
	case "TRX":
		return f.tronService, nil
	default:
		return nil, errors.NewBadRequestError("Unsupported currency: " + currency)
	}
}
