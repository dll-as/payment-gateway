package blockchain

import "errors"

func GetBlockchainService(coin string) (BlockchainService, error) {
	switch coin {
	case "TRX", "USDT-TRC20":
		return &TronService{}, nil
	case "ETH", "USDT-ERC20":
		return &EthereumService{}, nil
	case "BTC":
		// return &BitcoinService{}, nil
		return nil, errors.New("BTC not implemented yet")
	default:
		return nil, errors.New("unsupported coin")
	}
}
