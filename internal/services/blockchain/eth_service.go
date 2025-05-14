package blockchain

type EthereumService struct{}

func (e *EthereumService) GenerateAddress() (string, error) {
	// generate ETH address
	return "0xABC123...", nil
}

func (e *EthereumService) GetBalance(address string) (float64, error) {
	// use Infura / etherscan API
	return 0.18, nil
}

func (e *EthereumService) SendTransaction(fromPrivKey, toAddress string, amount float64) (string, error) {
	// web3 send transaction
	return "eth_tx_hash", nil
}

func (e *EthereumService) CheckTxConfirmed(txHash string) (bool, error) {
	// check confirmations using etherscan/infura
	return true, nil
}
