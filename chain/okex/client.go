package okex

import (
	"multichain/api/provider"
	"multichain/chain/ethereum"
)

type Client = ethereum.Client

func NewClient(provider provider.CommonProvider) (*Client, error) {
	return ethereum.NewClient(provider)
}
