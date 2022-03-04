package tools

import (
	"fmt"
	"math/big"
	"testing"
)

func TestAdd(t *testing.T) {
	i := big.NewInt(100000000000)
	d := 18
	AddDecimal(i, d, 10)
	fmt.Println(i)
	fmt.Println(i.String())
}

func TestCut(t *testing.T) {

	i := big.NewInt(100000000000001)
	d := 0
	f := CutDecimal(i, d, 0)
	t.Log(f)
}

func TestFloatAdd(t *testing.T) {
	f, _ := new(big.Float).SetString("23123.222")
	i := FloatAddDecimal(f, 18, 10)
	t.Log(i)
}
