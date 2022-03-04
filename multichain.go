package multichain

import (
	"fmt"
	"multichain/api/address"
	"multichain/api/client"
	"multichain/api/provider"
	"multichain/api/txbuilder"
	"multichain/chain/binance"
	"multichain/chain/ethereum"
	"multichain/chain/okex"
	"multichain/errno"
)

const (
	TypeEthereum = 1 // ethereum
	TypeBinance  = 2 // binance
	TypeOKEx     = 3 // okex
)

// NewClient 新建一个多链客户端
func NewClient(chainType uint, provider provider.CommonProvider) (client.Client, error) {
	var cli client.Client
	var err error
	switch chainType {
	case TypeEthereum:
		cli, err = ethereum.NewClient(provider)
	case TypeBinance:
		cli, err = binance.NewClient(provider)
	case TypeOKEx:
		cli, err = okex.NewClient(provider)
	}

	if err != nil {
		return nil, err
	}

	if cli == nil {
		return nil, errno.NotSupportChainType
	} else {
		return cli, nil
	}
}

// NewAddressManager 新建一个账户地址格式转换器
func NewAddressManager(chainType uint) (address.Encoder, error) {
	switch chainType {
	case TypeEthereum:
		return ethereum.NewAddressEncoder(), nil
	case TypeBinance:
		return binance.NewAddressEncoder(), nil
	case TypeOKEx:
		return okex.NewAddressEncoder(), nil
	default:
		return nil, fmt.Errorf("not support chain type : [%v]", chainType)
	}
}

// NewTxBuilder 新建一个交易构造器
func NewTxBuilder(chainType uint, provider provider.CommonProvider) (txbuilder.TxBuilder, error) {
	switch chainType {
	case TypeEthereum:
		builder, err := ethereum.NewTxBuilder(provider)
		if err != nil {
			return nil, err
		} else {
			return builder, nil
		}
	case TypeBinance:
		return binance.NewTxBuilder(provider)
	case TypeOKEx:
		return okex.NewTxBuilder(provider)
	default:
		return nil, errno.NotSupportChainType
	}
}

// NewContractTxBuilder 新建一个合约相关类型的交易构造器
func NewContractTxBuilder(chainType uint, provider provider.CommonProvider) (txbuilder.ContractTxBuilder, error) {
	switch chainType {
	case TypeEthereum:
		builder, err := ethereum.NewContractTxBuilder(provider)
		if err != nil {
			return nil, err
		} else {
			return builder, nil
		}
	case TypeBinance:
		return binance.NewContractTxBuilder(provider)
	case TypeOKEx:
		return okex.NewContractTxBuilder(provider)
	default:
		return nil, errno.NotSupportChainType
	}
}
