package ethereum

import (
	"crypto/ecdsa"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec"
	"github.com/mgintoki/go-web3"
	"github.com/umbracle/fastrlp"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strings"
)

func SignTx(tx *web3.Transaction, private *ecdsa.PrivateKey, chainID uint64) (*web3.Transaction, error) {
	hash := signHash(tx, chainID)

	sig, err := Sign(private, hash)
	if err != nil {
		return nil, err
	}

	vv := uint64(sig[64]) + 35 + chainID*2

	tx.R, tx.S, err = trimLeadingZero(sig[:32], sig[32:64])
	if err != nil {
		return nil, err
	}

	tx.V = new(big.Int).SetUint64(vv).Bytes()
	return tx, nil
}

func signHash(tx *web3.Transaction, chainID uint64) []byte {
	a := fastrlp.DefaultArenaPool.Get()

	v := a.NewArray()
	v.Set(a.NewUint(tx.Nonce))
	v.Set(a.NewUint(tx.GasPrice))
	v.Set(a.NewUint(tx.Gas))
	if tx.To == nil {
		v.Set(a.NewNull())
	} else {
		v.Set(a.NewCopyBytes((*tx.To)[:]))
	}
	v.Set(a.NewBigInt(tx.Value))
	v.Set(a.NewCopyBytes(tx.Input))

	// EIP155
	if chainID != 0 {
		v.Set(a.NewUint(chainID))
		v.Set(a.NewUint(0))
		v.Set(a.NewUint(0))
	}

	hash := keccak256(v.MarshalTo(nil))
	fastrlp.DefaultArenaPool.Put(a)
	return hash
}

func Sign(private *ecdsa.PrivateKey, hash []byte) ([]byte, error) {
	var S256 = btcec.S256()
	sig, err := btcec.SignCompact(S256, (*btcec.PrivateKey)(private), hash, false)
	if err != nil {
		return nil, err
	}
	term := byte(0)
	if sig[0] == 28 {
		term = 1
	}
	return append(sig, term)[1:], nil
}

// trimLeadingZero 去掉交易签名R、S开头的0x00
// 见 https://github.com/MOACChain/moac-core/issues/24
// 用来避免 rlp: non-canonical integer (leading zero bytes) for *big.Int, decoding into (types.Transaction)(types.txdata).R
func trimLeadingZero(r []byte, s []byte) ([]byte, []byte, error) {
	rs := strings.TrimPrefix(hex.EncodeToString(r), "0x")
	ss := strings.TrimPrefix(hex.EncodeToString(s), "0x")
	for strings.HasPrefix(rs, "00") {
		rs = strings.TrimPrefix(rs, "00")
	}
	for strings.HasPrefix(ss, "00") {
		ss = strings.TrimPrefix(ss, "00")
	}

	rh, err := hex.DecodeString(rs)
	if err != nil {
		return nil, nil, err
	}
	sh, err := hex.DecodeString(ss)
	if err != nil {
		return nil, nil, err
	}
	//rh32 := make([]byte, 32)
	//rh32 = append(rh)
	return rh, sh, nil
}

func keccak256(buf []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	h.Write(buf)
	b := h.Sum(nil)
	return b
}
