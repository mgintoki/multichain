package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/mgintoki/go-web3"
	"github.com/mgintoki/go-web3/abi"
	"github.com/mgintoki/go-web3/jsonrpc"
	"github.com/mgintoki/multichain/api/fee"
	"github.com/mgintoki/multichain/errno"
	"math/big"
	"strconv"
)

type Txn struct {
	Private    *ecdsa.PrivateKey `json:"private"`
	From       web3.Address      `json:"from"`
	Nonce      uint64            `json:"nonce"`
	Addr       *web3.Address     `json:"to"`
	Value      *big.Int          `json:"value"`
	GasPrice   uint64            `json:"gasPrice"`
	GasLimit   uint64            `json:"gas"`
	Data       []byte            `json:"input"`
	Provider   *jsonrpc.Client   `json:"provider"`
	Method     *abi.Method       `json:"method"`
	Args       []interface{}     `json:"args"`
	Bin        []byte            `json:"bin"`
	SignedHash []byte            `json:"signedHash"`
	SignedTx   []byte            `json:"signedTx"`
	Hash       web3.Hash         `json:"hash"`
	Receipt    *web3.Receipt     `json:"receipt"`
}

func (t *Txn) GetHash() string {
	return new(AddressEncoder).AddressToHex(t.Hash)
}

func (t *Txn) GetNonce() uint64 {
	return t.Nonce
}

func (t *Txn) GetFrom() string {
	return new(AddressEncoder).AddressToHex(t.From)
}

func (t *Txn) GetTo() string {
	if t.Addr == nil {
		return ""
	} else {
		return new(AddressEncoder).AddressToHex(t.Addr)
	}
}

func (t *Txn) GetValue() *big.Int {
	return t.Value
}

func (t *Txn) GetPayload() []byte {
	return t.Data
}

func (t *Txn) SignTx(privateHex string, chainID string) error {

	private := t.Private
	if private == nil {
		p, err := crypto.HexToECDSA(privateHex)
		if err != nil {
			return err
		} else {
			private = p
		}
	}

	chainIDInt, err := strconv.Atoi(chainID)
	if err != nil {
		return err
	}

	web3Tx := &web3.Transaction{
		Nonce:    t.Nonce,
		To:       t.Addr,
		Value:    t.Value,
		Gas:      t.GasLimit,
		GasPrice: t.GasPrice,
		Input:    t.Data,
	}

	web3Tx, err = SignTx(web3Tx, private, uint64(chainIDInt))
	if err != nil {
		return err
	}

	t.SignedTx = web3Tx.MarshalRLP()

	return nil
}

func (t *Txn) SignHash(privateHex string, chainID string, hexHash string) (hexSignature string, err error) {

	private := t.Private
	if private == nil {
		private, err = crypto.HexToECDSA(privateHex)
		if err != nil {
			return "", err
		}
	}

	hash, err := hex.DecodeString(hexHash)
	if err != nil {
		return "", err
	}

	sig, err := Sign(private, hash)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}

func (t *Txn) EncodeTx() (txEncoded string, err error) {
	b, err := json.Marshal(t)
	if err != nil {
		return "", err
	} else {
		return hex.EncodeToString(b), nil
	}
}

func (t *Txn) GetFee() *fee.OptionFee {
	return &fee.OptionFee{
		GasPrice: t.GasPrice,
		GasLimit: t.GasLimit,
	}
}

func (t *Txn) SetFee(fee *fee.OptionFee) {
	if fee.GasLimit != 0 {
		t.GasLimit = fee.GasLimit
	}

	if fee.GasPrice != 0 {
		t.GasPrice = fee.GasPrice
	}
}

func (t *Txn) GetTxHash(chainID string) (string, error) {
	ci, err := strconv.Atoi(chainID)
	if err != nil {
		return "", err
	}
	web3Tx := &web3.Transaction{
		Nonce:    t.Nonce,
		To:       t.Addr,
		Value:    t.Value,
		Gas:      t.GasLimit,
		GasPrice: t.GasPrice,
		Input:    t.Data,
	}
	hash := signHash(web3Tx, uint64(ci))
	return hex.EncodeToString(hash), nil
}

func (t *Txn) InjectSignature(hexSignature string, chainID string) (err error) {

	sig, err := hex.DecodeString(hexSignature)
	if err != nil {
		return err
	}

	web3Tx := &web3.Transaction{
		Nonce:    t.Nonce,
		To:       t.Addr,
		Value:    t.Value,
		Gas:      t.GasLimit,
		GasPrice: t.GasPrice,
		Input:    t.Data,
	}

	chainIDBig, ok := new(big.Int).SetString(chainID, 10)
	if !ok {
		return errno.InvalidStringToBigNum
	}

	vv := uint64(sig[64]) + 35 + chainIDBig.Uint64()*2

	web3Tx.R, web3Tx.S, err = trimLeadingZero(sig[:32], sig[32:64])
	if err != nil {
		return err
	}

	web3Tx.V = new(big.Int).SetUint64(vv).Bytes()

	t.SignedTx = web3Tx.MarshalRLP()

	return nil
}

//func (t *Txn) SignTx2(private *ecdsa.PrivateKey, chainID uint64) error {
//
//	web3Tx := &web3.Transaction{
//		Nonce:    t.Nonce,
//		To:       t.Addr,
//		Value:    t.Value,
//		Gas:      t.GasLimit,
//		GasPrice: t.GasPrice,
//		Input:    t.Data,
//	}
//
//	web3Tx, err := SignTx(web3Tx, private, chainID)
//	if err != nil {
//		return err
//	}
//
//	data := web3Tx.MarshalRLP()
//	t.SignedTx = data
//
//	return nil
//}

// Validate validates the arguments of the transaction
func (t *Txn) Validate() error {
	if t.Data != nil {
		// Already validated
		return nil
	}
	if t.isContractDeployment() {
		t.Data = append(t.Data, t.Bin...)
	}
	if t.Method != nil {
		data, err := abi.Encode(t.Args, t.Method.Inputs)
		if err != nil {
			return fmt.Errorf("failed to encode arguments: %v", err)
		}
		if !t.isContractDeployment() {
			t.Data = append(t.Method.ID(), data...)
		} else {
			t.Data = append(t.Data, data...)
		}
	}
	return nil
}

// SetGasPrice sets the gas price of the transaction
func (t *Txn) SetGasPrice(gasPrice uint64) *Txn {
	t.GasPrice = gasPrice
	return t
}

// SetGasLimit sets the gas limit of the transaction
func (t *Txn) SetGasLimit(gasLimit uint64) *Txn {
	t.GasLimit = gasLimit
	return t
}

// Wait waits till the transaction is mined
func (t *Txn) Wait() error {

	if (t.Hash == web3.Hash{}) {
		panic("transaction not executed")
	}
	var err error
	for {
		t.Receipt, err = t.Provider.Eth().GetTransactionReceipt(t.Hash)
		if err != nil {
			if err.Error() != "not found" {
				return err
			}
		}
		if t.Receipt != nil {
			break
		}
	}
	return nil
}

// Receipt returns the receipt of the transaction after wait
func (t *Txn) GetReceipt() *web3.Receipt {
	return t.Receipt
}

func (t *Txn) isContractDeployment() bool {
	return t.Addr == nil
}

// AddArgs is used to set the arguments of the transaction
func (t *Txn) AddArgs(args ...interface{}) *Txn {
	t.Args = args
	return t
}

// SetValue sets the value for the txn
func (t *Txn) SetValue(v *big.Int) *Txn {
	t.Value = new(big.Int).Set(v)
	return t
}

// EstimateGas estimates the gas for the call
func (t *Txn) EstimateGas() (uint64, error) {
	if err := t.Validate(); err != nil {
		return 0, err
	}
	return t.estimateGas()
}

func (t *Txn) estimateGas() (uint64, error) {
	if t.isContractDeployment() {
		// 部署合约时，data要带0x前缀
		return t.Provider.Eth().EstimateGasContractWithFrom(t.From, t.Data)
	}
	msg := &web3.CallMsg{
		From:  t.From,
		To:    t.Addr,
		Data:  t.Data,
		Value: t.Value,
	}
	return t.Provider.Eth().EstimateGas(msg)
}
