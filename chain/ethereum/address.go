package ethereum

import (
	"github.com/mgintoki/go-web3"
	"multichain/api/address"
)

type AddressEncoder struct {
}

func NewAddressEncoder() *AddressEncoder {
	return &AddressEncoder{}
}

func (a *AddressEncoder) AddressToHex(addr address.Address) string {
	ethAddress, ok := addr.(*web3.Address)
	if !ok {
		adr, ok := addr.(web3.Address)
		if !ok {
			return ""
		} else {
			return adr.String()
		}
	} else {
		return ethAddress.String()
	}
}

func (a *AddressEncoder) HexToAddress(addr string) address.Address {
	return web3.HexToAddress(addr)
}
