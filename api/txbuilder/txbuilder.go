package txbuilder

import (
	"math/big"
	"multichain/api/tx"
)

// BuildTxParam 是构造一个交易的参数
type BuildTxParam struct {
	From     string   //交易的发起地址
	To       string   //是交易的目标地址
	Value    *big.Int //交易的原生币数量
	Nonce    uint64   //可选，不传则内部计算nonce
	GasLimit uint64   //可选，不传则内部计算推荐值并使用
	GasPrice uint64   //可选，不传则内部计算推荐值并使用
	Payload  []byte   //交易负载数据
}

// TxBuilder 是交易的构造器
type TxBuilder interface {
	BuildTx(req BuildTxParam) (tx tx.Tx, err error)
	// DecodeTx 解析一个序列化后的交易
	DecodeTx(encodedTx string) (tx tx.Tx, err error)
}

// BuildDeployTxReq 定义了一种特定的交易类型-部署合约交易
type BuildDeployTxReq struct {
	From     string        //交易的发起方
	Abi      string        //智能合约的应用程序二进制接口
	ByteCode string        //智能合约编译后得到的16进制字节码
	Params   []interface{} //部署合约的参数
	Nonce    uint64        //可选，不传则内部计算nonce
	Value    *big.Int      //可选的交易的Value字段
	GasLimit uint64        //可选，不传则内部计算推荐值并使用
	GasPrice uint64        //可选，不传则内部计算推荐值并使用
}

// BuildInvokeTxReq 定义了一种特定的交易类型-调用合约交易
type BuildInvokeTxReq struct {
	From            string        //交易的发起方
	Abi             string        //智能合约的应用程序二进制接口
	Method          string        //调用的交易方法（需要在Abi中有定义）
	Params          []interface{} //调用合约的参数
	Nonce           uint64        //可选，不传则内部计算nonce
	Value           *big.Int      //可选的交易的Value字段
	ContractAddress string        //合约地址
	GasLimit        uint64        //可选，不传则内部计算推荐值并使用
	GasPrice        uint64        //可选，不传则内部计算推荐值并使用
}

// ContractTxBuilder 是 TxBuilder 之上的一层封装，主要用于构建智能合约相关的交易
// 一般来说, ContractTxBuilder 内部会调用 TxBuilder
type ContractTxBuilder interface {
	// BuildDeployTx 构造部署智能合约交易
	BuildDeployTx(req BuildDeployTxReq) (tx tx.Tx, err error)
	// BuildInvokeTx 构建调用合约的交易
	BuildInvokeTx(req BuildInvokeTxReq) (tx tx.Tx, err error)
}
