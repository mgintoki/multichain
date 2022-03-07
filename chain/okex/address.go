package okex

import "github.com/mgintoki/multichain/chain/ethereum"

type AddressEncoder = ethereum.AddressEncoder

func NewAddressEncoder() *AddressEncoder {
	return ethereum.NewAddressEncoder()
}
