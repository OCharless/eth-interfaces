package burnable

// Package enumerable provides functions to interact with ERC721 enumerable properties.

import (
	"math/big"

	"github.com/OCharless/eth-interfaces/inferences/ERC20Burnable"
	"github.com/OCharless/eth-interfaces/internal/base"
	"github.com/OCharless/eth-interfaces/internal/customerrors"
	"github.com/OCharless/eth-interfaces/internal/erc20"
	"github.com/OCharless/eth-interfaces/internal/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// IERC20Burnable wraps interactions with an IERC20Burnable contract, extending basic NFT interactions.
type IERC20BurnableInteractions struct {
	*erc20.ERC20Interactions
	erc20Burnable *ERC20Burnable.ERC20BurnableSession
	callError     func(string, error) *base.CallError
}

// NewIERC20Burnable creates a new enumerable interaction instance using the provided base NFT interactions.
func NewIERC20Burnable(baseIERC20 *erc20.ERC20Interactions, signatures []ERC20BurnableSignatures) (*IERC20BurnableInteractions, error) {
	erc20Burnable, err := ERC20Burnable.NewERC20Burnable(baseIERC20.GetAddress(), baseIERC20.Client)
	if err != nil {
		return nil, customerrors.WrapinterfacingError("ierc20Burnable", err)
	}
	session := ERC20Burnable.ERC20BurnableSession{
		Contract:     erc20Burnable,
		CallOpts:     baseIERC20.GetSession().CallOpts,
		TransactOpts: baseIERC20.GetSession().TransactOpts,
	}

	var converted []utils.Signature
	for _, sig := range signatures {
		converted = append(converted, sig)
	}

	err = baseIERC20.CheckSignatures(baseIERC20.GetAddress(), converted)
	if err != nil {
		return nil, customerrors.WrapinterfacingError("ierc20Burnable", err)
	}

	callError := func(field string, err error) *base.CallError {
		return baseIERC20.WrapCallError(ERC20Burnable.ERC20BurnableABI, field, err)
	}

	return &IERC20BurnableInteractions{baseIERC20, &session, callError}, nil
}

// Burn destroys the specified token from the owner's balance.
func (e *IERC20BurnableInteractions) Burn(qty *big.Int) (*types.Transaction, error) {
	tx, err := e.erc20Burnable.Burn(qty)
	if err != nil {
		return nil, e.callError("erc20.Burn()", err)
	}
	return tx, nil
}

// BurnFrom is a wrapper for Burn that calls the token's burnFrom function instead.
func (e *IERC20BurnableInteractions) BurnFrom(from common.Address, qty *big.Int) (*types.Transaction, error) {
	tx, err := e.erc20Burnable.BurnFrom(from, qty)
	if err != nil {
		return nil, e.callError("nft.BurnFrom()", err)
	}
	return tx, nil
}
