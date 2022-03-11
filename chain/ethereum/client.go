package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mgintoki/go-web3"
	"github.com/mgintoki/go-web3/abi"
	"github.com/mgintoki/go-web3/jsonrpc"
	"github.com/mgintoki/multichain/api/client"
	"github.com/mgintoki/multichain/api/fee"
	"github.com/mgintoki/multichain/api/provider"
	"github.com/mgintoki/multichain/api/tx"
	"github.com/mgintoki/multichain/api/txbuilder"
	"github.com/mgintoki/multichain/errno"
	"github.com/mgintoki/multichain/tools"
	"math/big"
	"strings"
	"time"
)

const (
	DefaultAddress = "0x3F43E75Aaba2c2fD6E227C10C6E7DC125A93DE3c"
)

type Client struct {
	provider *jsonrpc.Client
	private  *ecdsa.PrivateKey
	nodeUrl  string
	ctb      *ContractTxBuilder
	tb       *TxBuilder
}

func NewClient(provider provider.CommonProvider) (*Client, error) {
	c := &Client{}
	err := c.SetProvider(provider)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Client) SetProvider(provider provider.CommonProvider) (err error) {
	web3Provider, err := jsonrpc.NewClient(provider.ProviderUrl)
	if err != nil {
		return err
	}
	c.provider = web3Provider
	c.nodeUrl = provider.ProviderUrl

	c.ctb, err = NewContractTxBuilder(provider)
	if err != nil {
		return err
	}
	c.tb, err = NewTxBuilder(provider)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetPrivate(hexPrivate string) (err error) {
	private, err := crypto.HexToECDSA(hexPrivate)
	if err != nil {
		return err
	}
	c.private = private
	return nil
}
func (c *Client) GetAccount() string {
	if c.private == nil {
		return ""
	}
	return crypto.PubkeyToAddress(c.private.PublicKey).String()
}

func (c *Client) GetChainID() (string, error) {
	chainID, err := c.provider.Eth().ChainID()
	if err != nil {
		return "", err
	} else {
		return chainID.String(), nil
	}
}

func (c *Client) BalanceOf(address string, optionAsset *client.OptionAsset) (amount *big.Int, err error) {
	return c.provider.Eth().GetBalance(web3.HexToAddress(address), web3.Latest)
}

func (c *Client) Transfer(to string, amount *big.Int, optionAsset *client.OptionAsset, optionFee *fee.OptionFee) (txHash string, err error) {

	if c.private == nil {
		return "", fmt.Errorf("need private key")
	}

	txn, err := c.tb.BuildTx(txbuilder.BuildTxParam{
		//PrivateHex: "",
		From:  c.GetAccount(),
		To:    to,
		Value: amount,
	})
	if err != nil {
		return "", err
	}

	return c.SendTx(txn, optionFee)

}

func (c *Client) QueryTx(txHash string, isWait bool) (txData *tx.TxData, err error) {

	if txHash == "0x0000000000000000000000000000000000000000000000000000000000000000" || txHash == "" {
		return nil, fmt.Errorf("invalid tx hash:[%v]", txHash)
	}

	web3Hash := web3.HexToHash(txHash)

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	txData = &tx.TxData{
		TxHash: txHash,
	}
	isPending := false

	for {
		receipt, err := c.provider.Eth().GetTransactionReceipt(web3Hash)
		if err != nil {
			if err.Error() == "not found" {
				isPending = true
			} else {
				if err.Error() == "field 'removed' not found" && receipt.BlockNumber != 0 && receipt.TransactionIndex == 0 {
					isPending = false
					//todo 以太坊不同版本返回的txIndex分别是uint64和string,临时这样处理错误，后面改以太坊环境测试
				} else {
					return nil, err
				}
			}
		}
		if receipt != nil {
			isPending = false
			txData.ContractAddress = receipt.ContractAddress.String()
			txData.GasUsed = receipt.GasUsed
		} else {
			isPending = true
		}

		if !isWait {
			if isPending {
				return txData, nil
			} else {
				web3Tx, err := c.provider.Eth().GetTransactionByHash(web3Hash)
				if err != nil {
					return nil, err
				}
				convertTxInfo(txData, *web3Tx)
				return txData, nil
			}
		}

		if !isPending {
			web3Tx, err := c.provider.Eth().GetTransactionByHash(web3Hash)
			if err != nil {
				return nil, err
			}
			convertTxInfo(txData, *web3Tx)
			return txData, nil
		}

		<-ticker.C

	}
}

func convertTxInfo(txData *tx.TxData, web3Tx web3.Transaction) {

	txData.From = web3Tx.From.String()
	if web3Tx.To != nil {
		txData.To = web3Tx.To.String()
	}

	txData.Data = web3Tx.Input
	txData.GasPrice = web3Tx.GasPrice
	txData.GasLimit = web3Tx.Gas
	txData.Nonce = web3Tx.Nonce
	txData.BlockNumber = web3Tx.BlockNumber
	txData.Value = web3Tx.Value
	txData.Raw = []byte(tools.FastMarshal(web3Tx))
}

func (c *Client) SendTx(tx tx.Tx, feeOption *fee.OptionFee) (txHash string, err error) {

	if c.private == nil {
		return "", fmt.Errorf("need private key")
	}

	t, ok := tx.(*Txn)
	if !ok {
		return "", errno.InvalidTxType
	}

	t.From = web3.HexToAddress(crypto.PubkeyToAddress(c.private.PublicKey).String())

	var gasLimit, gasPrice uint64

	if feeOption != nil {
		gasLimit = feeOption.GasLimit
		gasPrice = feeOption.GasPrice
	}

	if gasPrice != 0 {
		t.GasPrice = gasPrice
	} else {
		if t.GasPrice == 0 {
			t.GasPrice, err = c.provider.Eth().GasPrice()
			if err != nil {
				return "", err
			}
		}
	}

	if gasLimit != 0 {
		t.GasLimit = gasLimit
	} else {
		if t.GasLimit == 0 {
			f, err := c.EstimateGas(tx)
			if err != nil {
				return "", err
			} else {
				t.GasLimit = f.GasLimit
			}
		}
	}

	if t.Nonce == 0 {
		t.Nonce, err = c.provider.Eth().GetNonce(web3.HexToAddress(c.GetAccount()), web3.Latest)
		if err != nil {
			return "", err
		}
	}

	chainID, err := c.GetChainID()
	if err != nil {
		return "", err
	}

	if err := t.SignTx(hex.EncodeToString(crypto.FromECDSA(c.private)), chainID); err != nil {
		return "", err
	}

	return c.SendSignedTx(t)
}

func (c *Client) SendSignedTx(signedTx tx.Tx) (txHash string, err error) {

	t, ok := signedTx.(*Txn)
	if !ok {
		return "", errno.InvalidTxType
	}

	web3Hash, err := c.provider.Eth().SendRawTransaction(t.SignedTx)
	if err != nil {
		return "", err
	} else {
		return web3Hash.String(), nil
	}

}

func (c *Client) QueryContract(req client.CallContractParam) (res *client.CallContractRes, err error) {
	abiIns, err := abi.NewABI(req.Abi)
	if err != nil {
		return nil, err
	}
	contractIns := NewContract(web3.HexToAddress(req.ContractAddress), abiIns, c.provider)
	from := req.From
	if from == "" {
		from = DefaultAddress
	}
	contractIns.SetFrom(web3.HexToAddress(from))
	calledFunc := strings.Trim(req.CalledFunc, "()")
	rawRes, contractRes, err := contractIns.Call(calledFunc, web3.Latest, req.Params...)
	if err != nil {
		return nil, err
	} else {
		return &client.CallContractRes{
			RawRes:    rawRes,
			DecodeRes: contractRes,
		}, nil
	}
}

func (c *Client) EstimateGas(tx tx.Tx) (feeRes *fee.OptionFee, err error) {

	ethTx, ok := tx.(*Txn)
	if !ok {
		return nil, errno.InvalidTxType
	}

	var gasPrice, gasLimit uint64

	gasPrice, err = c.provider.Eth().GasPrice()
	if err != nil {
		return nil, err
	}
	//如果目标地址为空，说明为部署合约交易,否则为普通交易类型
	gasLimit, err = ethTx.EstimateGas()

	feeRes = &fee.OptionFee{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
	}
	return feeRes, nil
}
