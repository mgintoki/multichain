package okex

import (
	"github.com/mgintoki/multichain/api/provider"
	"github.com/mgintoki/multichain/chain/ethereum"
)

type TxBuilder = ethereum.TxBuilder

func NewTxBuilder(provider provider.CommonProvider) (*TxBuilder, error) {
	return ethereum.NewTxBuilder(provider)
}

type ContractTxBuilder = ethereum.ContractTxBuilder

func NewContractTxBuilder(provider provider.CommonProvider) (*ContractTxBuilder, error) {
	return ethereum.NewContractTxBuilder(provider)
}
