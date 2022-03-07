package client

import (
	"github.com/mgintoki/multichain/api/fee"
	"github.com/mgintoki/multichain/api/provider"
	"github.com/mgintoki/multichain/api/tx"
	"math/big"
)

// CallContractParam 是查询合约的参数
// From 代表查询的发起地址，如果不传，会随机从链上选择一个账户作为查询发起方
// ContractAddress 是独一无二的合约地址
// CalledFunc 对应的方法名称必须存在与 Abi 中，才能被正确解析
// Params 对应的参数必须和 Abi 中 CalledFunc 对应的方法参数的数量类型相对应
// 注意 如果 Abi 中定义了的某个参数类型为 address 类型，代表参数类型为 address.Address,
// 此时使用string格式的地址是非法的，需要使用 address.Encoder 将string格式的账户地址转换为 address.Address 类型
type CallContractParam struct {
	From            string        //调用发起方
	ContractAddress string        //合约地址
	Abi             string        //string格式的abi
	CalledFunc      string        //方法名
	Params          []interface{} //合约方法参数
}

// CallContractRes 定义了查询合约的返回值
//
// DecodeRes 是根据 Abi定义的合约返回值，解析得到的map。
// 如果调用的方法只有一个匿名返回值，DecodeRes 中会存在一个key为"0"的值代表返回值,
// 如果Abi中定义了返回值的名称，则会以定义的返回值名称作为key
//
// RawRes 定义了查询合约的原始返回，使用者亦可自行对其进行解析
type CallContractRes struct {
	RawRes    string                 //调用链合约的原始返回
	DecodeRes map[string]interface{} //经过ABI编码后的返回, key值对应abi中返回值名称的定义
}

// OptionAsset 是可选的资产参数
// 为一条链多种币的余额和转账接口做预置
type OptionAsset struct {
	InternalAssetType string
}

// Client 定义了 multichain多链SDK提供的功能接口
type Client interface {
	// SetProvider 为 Client 设置一个对于链调用的服务提供者
	SetProvider(provider provider.CommonProvider) (err error)

	// SetPrivate 设置 Client 的私钥
	// 当设置过私钥之后，才可以调用 SendTx 接口发送交易以及 Transfer 转账接口
	SetPrivate(hexPrivate string) (err error)

	// GetAccount 获取当前 Client 的使用账户地址
	// 必须先调用过 SetPrivate 才能获取正确的 GetAccount 返回值
	GetAccount() string

	// GetChainID 获取当前链的ID
	GetChainID() (string, error)

	// BalanceOf 当前会查返回 address 对应账户地址的链原生币余额。
	// optionAsset 是一个占位参数，为一条链多种原生币做预置
	BalanceOf(address string, optionAsset *OptionAsset) (amount *big.Int, err error)

	// Transfer 是链的原生币转账接口，会将amount数量的原生币支付给to对应的账户地址
	// 需要已经设置过私钥
	// 返回值为交易的hash，可调用 QueryTx 接口拿到交易的详情
	Transfer(to string, amount *big.Int, optionAsset *OptionAsset, optionFee *fee.OptionFee) (txHash string, err error)

	// QueryTx 是查询交易详情
	// txHash 是交易的hash
	// isWait 代表是否异步等待查询结果
	// isWait 为false，代表直接返回查询的结果（若交易没有完成或是失败，返回结果中只会有交易hash）
	// isWait 为true，代表使用同步方式，在方法内部做轮询，直到交易完成或失败才会返回
	QueryTx(txHash string, isWait bool) (txData *tx.TxData, err error)

	// SendTx 发送一个交易
	// 方法内部会对交易做签名并发送到链上
	// 需要先调用过 SetPrivate 接口设置过私钥
	SendTx(tx tx.Tx, feeOption *fee.OptionFee) (txHash string, err error)

	// SendSignedTx 与 SendTx 的区别是， SendSignedTx 发送的是签名后的交易
	// 用户构建完成交易后，可以在别处签名，然后使用该接口发送签名后的交易
	// 出于对于私钥安全性的考虑，请谨慎选择 SendTx 和 SendSignedTx
	SendSignedTx(signedTx tx.Tx) (txHash string, err error)

	// QueryContract 查询合约
	// 详情请参考 CallContractParam 和 CallContractRes 的定义
	QueryContract(req CallContractParam) (res *CallContractRes, err error)

	// EstimateGas 计算一个交易所需要的交易费用
	// 可以根据计算得到的交易费用自行指定实际需要的交易费用
	EstimateGas(tx tx.Tx) (feeRes *fee.OptionFee, err error)
}
