package basetokens

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type BaseToken interface {
	Name() (string, error)
	Symbol() (string, error)
	TokenURI(tokenID *big.Int) (string, error)
	GetBalance() (*big.Int, error)
	TransferTo(to common.Address, tokenID *big.Int) (*types.Transaction, error)
	TotalSupply() (*big.Int, error)
	BalanceOf(owner common.Address) (*big.Int, error)
	Approve(to common.Address, tokenID *big.Int) (*types.Transaction, error)
}
