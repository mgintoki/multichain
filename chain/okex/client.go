package okex

import (
	"github.com/mgintoki/multichain/api/client"
	"github.com/mgintoki/multichain/api/provider"
	"github.com/mgintoki/multichain/chain/ethereum"
)

type Client = ethereum.Client

func NewClient(provider provider.CommonProvider) (client.Client, error) {
	return ethereum.NewClient(provider)
}
