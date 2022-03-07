package binance

import (
	"github.com/mgintoki/multichain/api/provider"
	"github.com/mgintoki/multichain/chain/ethereum"
)

type Client = ethereum.Client

func NewClient(provider provider.CommonProvider) (*Client, error) {
	return ethereum.NewClient(provider)
}
