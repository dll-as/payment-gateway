package blockchain

type BlockchainService interface {
	GenerateAddress() (string, error)
	GetBalance(address string) (float64, error)
	SendTransaction(fromPrivKey, toAddress string, amount float64) (string, error)
	CheckTxConfirmed(txHash string) (bool, error)
}
