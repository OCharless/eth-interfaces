package erc20

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

type BaseERC20Signature string

const (
	Name             BaseERC20Signature = "name()"
	Symbol           BaseERC20Signature = "symbol()"
	Decimals         BaseERC20Signature = "decimals()"
	BalanceOf        BaseERC20Signature = "balanceOf(address)"
	TotalSupply      BaseERC20Signature = "totalSupply()"
	TokenURI         BaseERC20Signature = "tokenURI(uint256)"
	Approve          BaseERC20Signature = "approve(address,uint256)"
	TransferFrom     BaseERC20Signature = "transferFrom(address,address,uint256)"
	SafeTransferFrom BaseERC20Signature = "safeTransferFrom(address,address,uint256)"
)

func (s BaseERC20Signature) GetHex() string {
	hash := crypto.NewKeccakState()
	hash.Write([]byte(s))
	selector := hash.Sum(nil)[:4]
	return hex.EncodeToString(selector)
}
