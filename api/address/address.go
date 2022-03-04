package address

// Address 是一个链上账户地址
type Address interface{}

// Encoder 处理 Address 类型和string类型的格式转换
type Encoder interface {
	// AddressToHex 将账户地址从 Address 格式转换为string格式
	AddressToHex(addr Address) string
	// HexToAddress 将账户地址从string格式转换为 Address 格式
	HexToAddress(addr string) Address
}
