package tools

import (
	"math/big"
)

// AddDecimal i = i*p^d, p默认为10
func AddDecimal(i *big.Int, d int, p int64) {
	if p == 0 {
		p = 10
	}
	var decimals, pow = big.NewInt(int64(d)), big.NewInt(p)
	pow.Exp(pow, decimals, nil)
	i.Mul(i, pow)
}

// CutDecimal format i， 得到 i/p^d ,p默认10
func CutDecimal(i *big.Int, d int, p int64) *big.Float {

	if p == 0 {
		p = 10
	}

	var decimals, pow = big.NewInt(int64(d)), big.NewInt(p)
	pow.Exp(pow, decimals, nil)

	bigF := new(big.Float).SetInt(i)
	bigF.Quo(bigF, new(big.Float).SetInt(pow))

	return bigF
	//str := strconv.FormatFloat(f, 'f', precious, 64)
}

// FloatAddDecimal i = i*p^d, p默认为10
func FloatAddDecimal(f *big.Float, d int, p int64) *big.Int {
	if p == 0 {
		p = 10
	}
	var pow, dec = big.NewInt(p), big.NewInt(int64(d))
	pow.Exp(pow, dec, nil)
	f.Mul(f, new(big.Float).SetInt(pow))
	i := new(big.Int)
	i, _ = f.Int(i)
	return i
}
