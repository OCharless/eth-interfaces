package merged

// Package merged provides a unified interface that combines multiple NFT interaction extensions such as enumerable and royalties.

import (
	"math/big"

	"github.com/OCharless/eth-interfaces/models"
	"github.com/OCharless/eth-interfaces/nft"
	"github.com/OCharless/eth-interfaces/nft/enumerable"
	"github.com/OCharless/eth-interfaces/nft/royalties"
	"github.com/OCharless/eth-interfaces/utils"
	"github.com/ethereum/go-ethereum/common"
)

// IERC721SummedInteractions aggregates NFT interactions from various extensions (e.g., royalties and enumerable) into a single interface.
type IERC721SummedInteractions struct {
	*nft.ERC721Interactions
	*royalties.IERC721RoyaltiesInteractions
	*enumerable.ERC721EnumerableInteractions
}

// ExtensionEnum denotes the types of NFT interaction extensions to be included in the summed interactions.
type ExtensionEnum int

const (
	// Enumerable represents the enumerable extension.
	Enumerable ExtensionEnum = iota
	// Royalties represents the royalties extension.
	Royalties
)

// NewERC721SummedInteractions creates a new instance of IERC721SummedInteractions by initializing the specified extensions from the base NFT interactions.
func NewERC721SummedInteractions(baseIERC721 *nft.ERC721Interactions, signatures []utils.Signature, extensions ...ExtensionEnum) (*IERC721SummedInteractions, error) {
	var enum *enumerable.ERC721EnumerableInteractions = nil
	var roy *royalties.IERC721RoyaltiesInteractions = nil
	var err error

	err = baseIERC721.CheckSignatures(baseIERC721.GetAddress(), signatures)
	if err != nil {
		return nil, err
	}

	for _, extension := range extensions {
		switch extension {
		case Enumerable:
			enum, err = enumerable.NewERC721EnumerableInteractions(
				baseIERC721,
				[]enumerable.IERC721EnumerableSignature{},
			)
			if err != nil {
				return nil, err
			}
		case Royalties:
			roy, err = royalties.NewERC721RoyaltiesInteractions(
				baseIERC721,
				[]royalties.IERC721RoyaltiesSignature{},
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return &IERC721SummedInteractions{baseIERC721, roy, enum}, nil
}

// AllInfos retrieves combined information for a given token, including base metadata, total supply, and royalty information.
func (s *IERC721SummedInteractions) AllInfos(tokenIDs ...*big.Int) (*models.TokenMeta, *big.Int, *royalties.RoyaltyInfos, error) {
	var tokenID *big.Int
	if len(tokenIDs) == 0 {
		tokenID = common.Big0
	} else {
		tokenID = tokenIDs[0]
	}

	baseInfos, err := s.TokenMetaInfos(tokenID)
	if err != nil {
		return nil, nil, nil, err
	}

	supply, err := s.TotalSupply()
	if err != nil {
		return baseInfos, nil, nil, err
	}

	royalties, err := s.RoyaltiesInfos(tokenID, big.NewInt(1))
	if err != nil {
		return baseInfos, supply, nil, err
	}

	return baseInfos, supply, &royalties, nil
}
