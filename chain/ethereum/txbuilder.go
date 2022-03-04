package ethereum

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mgintoki/go-web3"
	"github.com/mgintoki/go-web3/abi"
	"github.com/mgintoki/go-web3/jsonrpc"
	"multichain/api/provider"
	"multichain/api/tx"
	"multichain/api/txbuilder"
	"multichain/errno"
	"strings"
)

type TxBuilder struct {
	provider *jsonrpc.Client
}

func NewTxBuilder(provider provider.CommonProvider) (*TxBuilder, error) {
	p, err := jsonrpc.NewClient(provider.ProviderUrl)
	if err != nil {
		return nil, err
	}
	txBuilder := &TxBuilder{p}
	return txBuilder, nil
}

func (t *TxBuilder) BuildTx(req txbuilder.BuildTxParam) (tx tx.Tx, err error) {

	from := req.From
	if from == "" {
		return nil, errno.TxFromNotSet
	}

	var toAddr *web3.Address
	if req.To != "" {
		tmp := web3.HexToAddress(req.To)
		toAddr = &tmp
	}

	nonce := req.Nonce
	if nonce == 0 {
		nonce, err = t.provider.Eth().GetNonce(web3.HexToAddress(from), web3.Latest)
		if err != nil {
			return nil, err
		}
	}

	txn := &Txn{
		Provider: t.provider,
		From:     web3.HexToAddress(from),
		Addr:     toAddr,
		Data:     req.Payload,
		Value:    req.Value,
		Nonce:    nonce,
	}

	if req.GasPrice != 0 {
		txn.GasPrice = req.GasPrice
	} else {
		txn.GasPrice, err = t.provider.Eth().GasPrice()
		if err != nil {
			return nil, err
		}
	}

	if req.GasLimit != 0 {
		txn.GasLimit = req.GasLimit
	} else {
		txn.GasLimit, err = txn.EstimateGas()
		if err != nil {
			return nil, err
		}
	}

	return txn, err
}

func (t *TxBuilder) DecodeTx(encodedTx string) (tx tx.Tx, err error) {
	txByte, err := hex.DecodeString(encodedTx)
	if err != nil {
		return nil, err
	}
	var txn *Txn
	if err := json.Unmarshal(txByte, &txn); err != nil {
		return nil, err
	} else {
		return txn, nil
	}
}

type ContractTxBuilder struct {
	provider *jsonrpc.Client
}

func NewContractTxBuilder(provider provider.CommonProvider) (*ContractTxBuilder, error) {
	p, err := jsonrpc.NewClient(provider.ProviderUrl)
	if err != nil {
		return nil, err
	}
	txBuilder := &ContractTxBuilder{p}
	return txBuilder, nil
}

//构建部署合约的交易
func (b *ContractTxBuilder) BuildDeployTx(req txbuilder.BuildDeployTxReq) (tx.Tx, error) {

	abiIns, err := abi.NewABI(req.Abi)
	if err != nil {
		return nil, err
	}

	bin, err := hex.DecodeString(strings.TrimPrefix(req.ByteCode, "0x"))
	if err != nil {
		return nil, err
	}

	data, err := abi.Encode(req.Params, abiIns.Constructor.Inputs)
	if err != nil {
		return nil, err
	}

	data = append(bin, data...)

	t := &TxBuilder{b.provider}
	txn, err := t.BuildTx(txbuilder.BuildTxParam{
		From:     req.From,
		Payload:  data,
		Nonce:    req.Nonce,
		Value:    req.Value,
		GasPrice: req.GasPrice,
		GasLimit: req.GasLimit,
	})
	return txn, err
}

//构建调用合约的交易
func (b *ContractTxBuilder) BuildInvokeTx(req txbuilder.BuildInvokeTxReq) (tx.Tx, error) {

	abiIns, err := abi.NewABI(req.Abi)
	if err != nil {
		return nil, err
	}

	m, ok := abiIns.Methods[req.Method]
	if !ok {
		return nil, fmt.Errorf("invalid method")
	}

	data, err := abi.Encode(req.Params, m.Inputs)
	if err != nil {
		return nil, err
	}

	data = append(m.ID(), data...)

	t := &TxBuilder{b.provider}
	return t.BuildTx(txbuilder.BuildTxParam{
		From:     req.From,
		To:       req.ContractAddress,
		Payload:  data,
		Nonce:    req.Nonce,
		Value:    req.Value,
		GasPrice: req.GasPrice,
		GasLimit: req.GasLimit,
	})
}
