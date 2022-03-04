package ethereum

import (
	"encoding/hex"
	"fmt"
	"github.com/mgintoki/go-web3"
	"github.com/mgintoki/go-web3/abi"
	"github.com/mgintoki/go-web3/jsonrpc"
)

const (
	rawContractRes = "raw_contract_res"
)

// Contract is an Ethereum contract
type Contract struct {
	Addr     web3.Address
	From     *web3.Address
	Abi      *abi.ABI
	Provider *jsonrpc.Client
}

// NewContract creates a new contract instance
func NewContract(addr web3.Address, abi *abi.ABI, provider *jsonrpc.Client) *Contract {
	return &Contract{
		Addr:     addr,
		Abi:      abi,
		Provider: provider,
	}
}

// Addr returns the address of the contract
func (c *Contract) GetAddr() web3.Address {
	return c.Addr
}

// SetFrom sets the origin of the calls
func (c *Contract) SetFrom(addr web3.Address) {
	c.From = &addr
}

// EstimateGas estimates the gas for a contract call
func (c *Contract) EstimateGas(method string, args ...interface{}) (uint64, error) {
	return c.Txn(method, args).EstimateGas()
}

// Call calls a method in the contract
func (c *Contract) Call(method string, block web3.BlockNumber, args ...interface{}) (string, map[string]interface{}, error) {
	m, ok := c.Abi.Methods[method]
	if !ok {
		return "", nil, fmt.Errorf("method %s not found", method)
	}

	// Encode input
	data, err := abi.Encode(args, m.Inputs)
	if err != nil {
		return "", nil, err
	}
	data = append(m.ID(), data...)

	// CallReturnRaw function
	msg := &web3.CallMsg{
		To:   &c.Addr,
		Data: data,
	}
	if c.From != nil {
		msg.From = *c.From
	}

	rawStr, err := c.Provider.Eth().Call(msg, block)
	if err != nil {
		return "", nil, err
	}

	// Decode output
	raw, err := hex.DecodeString(rawStr[2:])
	if err != nil {
		return "", nil, err
	}
	if len(raw) == 0 {
		//return nil, fmt.Errorf("empty response")
		return "", nil, nil
	}
	respInterface, err := abi.Decode(m.Outputs, raw)
	if err != nil {
		return "", nil, err
	}

	resp := respInterface.(map[string]interface{})
	return rawStr, resp, nil
}

// Txn creates a new transaction object
func (c *Contract) Txn(method string, args ...interface{}) *Txn {
	m, ok := c.Abi.Methods[method]
	if !ok {
		// TODO, return error
		panic(fmt.Errorf("method %s not found", method))
	}

	return &Txn{
		From:     *c.From,
		Addr:     &c.Addr,
		Provider: c.Provider,
		Method:   m,
		Args:     args,
	}
}

// Event is a solidity event
type Event struct {
	event *abi.Event
}

// Encode encodes an event
func (e *Event) Encode() web3.Hash {
	return e.event.ID()
}

// ParseLog parses a log
func (e *Event) ParseLog(log *web3.Log) (map[string]interface{}, error) {
	return abi.ParseLog(e.event.Inputs, log)
}

// Event returns a specific event
func (c *Contract) Event(name string) (*Event, bool) {
	event, ok := c.Abi.Events[name]
	if !ok {
		return nil, false
	}
	return &Event{event}, true
}
