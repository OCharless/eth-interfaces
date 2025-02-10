package royalties

// Package royalties provides functions to interact with ERC721 royalties.

import (
	"math/big"

	"github.com/OCharless/eth-interfaces/inferences/ERC721Complete"
	"github.com/OCharless/eth-interfaces/internal/base"
	"github.com/OCharless/eth-interfaces/internal/customerrors"
	"github.com/OCharless/eth-interfaces/internal/nft"
	"github.com/OCharless/eth-interfaces/internal/utils"
	"github.com/ethereum/go-ethereum/common"
)

// IERC721RoyaltiesInteractions wraps interactions with the IERC721Royalties contract.
type IERC721RoyaltiesInteractions struct {
	*nft.ERC721Interactions
	ierc721Royalties *ERC721Complete.ERC721CompleteSession
	callError        func(string, error) *base.CallError
}

// RoyaltyInfos holds the royalty receiver and amount for a token sale.
type RoyaltyInfos struct {
	Receiver      common.Address
	RoyaltyAmount *big.Int
}

// NewIERC721RoyaltiesInteractions creates a new instance of IERC721RoyaltiesInteractions.
func NewERC721RoyaltiesInteractions(baseIERC721 *nft.ERC721Interactions, signatures []IERC721RoyaltiesSignature) (*IERC721RoyaltiesInteractions, error) {
	ierc721Royalties, err := ERC721Complete.NewERC721Complete(baseIERC721.GetAddress(), baseIERC721.Client)
	session := ERC721Complete.ERC721CompleteSession{
		Contract:     ierc721Royalties,
		CallOpts:     baseIERC721.GetSession().CallOpts,
		TransactOpts: baseIERC721.GetSession().TransactOpts,
	}
	if err != nil {
		return nil, customerrors.WrapinterfacingError("ierc721Royalties", err)
	}

	var converted []utils.Signature
	for _, sig := range signatures {
		converted = append(converted, sig)
	}

	err = baseIERC721.CheckSignatures(baseIERC721.GetAddress(), converted)
	if err != nil {
		return nil, customerrors.WrapinterfacingError("ierc721Royalties", err)
	}

	callError := func(field string, err error) *base.CallError {
		return baseIERC721.WrapCallError(ERC721Complete.ERC721CompleteABI, "nft.RoyaltyInfo()", err)
	}

	return &IERC721RoyaltiesInteractions{baseIERC721, &session, callError}, nil
}

// RoyaltiesInfos retrieves the royalty information for a given token and sale price.
func (e *IERC721RoyaltiesInteractions) RoyaltiesInfos(tokenID *big.Int, salePrice *big.Int) (RoyaltyInfos, error) {
	rInfos, err := e.ierc721Royalties.RoyaltyInfo(tokenID, salePrice)
	if err != nil {
		return RoyaltyInfos{}, e.callError("nft.RoyaltyInfo()", err)
	}
	return rInfos, nil
}
