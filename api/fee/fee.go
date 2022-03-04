package fee

// OptionFee 是可选的交易费用参数
// "可选" 代表在有该参数的接口，不指定该参数时，方法会在内部计算交易的推荐费用并使用
// 如果使用者需要自行设置交易费，请清楚你在做什么
type OptionFee struct {
	// GasPrice 是当前链的Gas价格
	GasPrice uint64
	// GasPrice 是完成当前交易需要消耗多少Gas
	GasLimit uint64
}
