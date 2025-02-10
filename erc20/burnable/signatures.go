package burnable

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

type ERC20BurnableSignatures string

const (
	Burn     ERC20BurnableSignatures = "burn(uint256)"
	BurnFrom ERC20BurnableSignatures = "burnFrom(address,uint256)"
)

func (s ERC20BurnableSignatures) GetHex() string {
	hash := crypto.NewKeccakState()
	hash.Write([]byte(s))
	selector := hash.Sum(nil)[:4]
	return hex.EncodeToString(selector)
}
