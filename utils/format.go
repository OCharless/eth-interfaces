package utils

import (
	"fmt"
	"math/big"
	"strings"
)

func FloatTo18z(amount float64) *big.Int {
	s := fmt.Sprintf("%f", amount)
	a := strings.Split(s, ".")
	for i := 0; i < 18; i++ {
		if len(a) > 1 && len(a[1]) > i {
			a[0] += string(a[1][i])
		} else {
			a[0] += "0"
		}
	}
	b, ok := new(big.Int).SetString(a[0], 10)
	if !ok {
		panic("Could not set big.Int string for value " + s)
	}
	return b
}

func ParseEther(a *big.Int) float64 {
	fa, _ := a.Float64()
	return fa / (1e18)
}
