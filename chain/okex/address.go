package okex

import "multichain/chain/ethereum"

type AddressEncoder = ethereum.AddressEncoder

func NewAddressEncoder() *AddressEncoder {
	return ethereum.NewAddressEncoder()
}
