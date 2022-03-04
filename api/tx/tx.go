package tx

import (
	"math/big"
	"multichain/api/fee"
)

// Tx 是链SDK中的一个交易
type Tx interface {
	// GetHash that uniquely identifies the transaction. Hashes are usually the
	// result of an irreversible hashing function applied to some serialized
	// representation of the transaction.
	GetHash() string

	// GetNonce used to order the transaction with respect to all other
	// transactions signed and submitted by the sender of this transaction.
	GetNonce() uint64

	// GetFrom returns the address from which value is being sent.
	GetFrom() string

	// GetTo returns the address to which value is being sent.
	GetTo() string

	// GetValue being sent from one address to another.
	GetValue() *big.Int

	// GetPayload returns arbitrary data that is associated with the transaction.
	// This payload is often used to send notes between external accounts, or
	// call functions on a contract.
	GetPayload() []byte

	// GetFee 计算推荐的交易费用
	GetFee() *fee.OptionFee

	// SetFee 自定义交易费用
	SetFee(fee *fee.OptionFee)

	// SignTx 传入私钥对交易自身做签名
	SignTx(privateHex string, chainID string) (err error)

	// GetTxHash 获取交易需要被签名的hash
	// 获取代签名hash之后，可以在别处使用私钥签名该hash，然后使用 InjectSignature 将对hash的签名注入到交易里
	GetTxHash(chainID string) (hexHash string, err error)

	// SignHash 签名某个hash
	// 注意：若非在当前上下文中签名，则实际签名的地方，签名算法应当与本SDK特定链类型定义的签名方法一致
	SignHash(privateHex string, chainID string, hexHash string) (hexSignature string, err error)

	// InjectSignature 将交易hash的签名注入交易中
	InjectSignature(hexSignature string, chainID string) (err error)

	// EncodeTx 序列化交易
	EncodeTx() (txEncoded string, err error)
}

// TxData 是交易详情
type TxData struct {
	TxHash          string   `json:"txHash"`
	ContractAddress string   `json:"contractAddress"` //  部署合约交易返回的交易地址
	TxType          string   `json:"txType"`
	From            string   `json:"from"`  // 交易的发起方
	Nonce           uint64   `json:"nonce"` // 请参考以太坊nonce定义
	Data            []byte   `json:"input"`
	To              string   `json:"to"`
	Value           *big.Int `json:"value"`
	GasLimit        uint64   `json:"gasLimit"`
	GasUsed         uint64   `json:"gas"`
	GasPrice        uint64   `json:"gasPrice"`
	BlockNumber     uint64   `json:"blockNumber"`
	Date            string   `json:"date"`
	Raw             []byte   `json:"raw"` //  交易详情原文,使用者可按需解析
}
