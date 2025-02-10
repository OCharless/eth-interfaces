package enumerable_test

// Package enumerable_test contains tests for enumerable interactions.

import (
	"testing"

	"github.com/OCharless/eth-interfaces/base"
	"github.com/OCharless/eth-interfaces/inferences/ERC721Complete"
	"github.com/OCharless/eth-interfaces/nft"
	"github.com/OCharless/eth-interfaces/nft/enumerable"
	"github.com/OCharless/eth-interfaces/utils"
	"github.com/stretchr/testify/assert"
)

// Test_GetAddressOwnedTokens tests retrieval of tokens owned by an address and validates that the number of tokens matches the balance.
func Test_GetAddressOwnedTokens(t *testing.T) {
	backend, auth, contractAddr, privKey, err := utils.SetupBlockchain(t,
		ERC721Complete.ERC721CompleteABI,
		ERC721Complete.ERC721CompleteBin,
		"MyNFT",
		"MNFT",
	)
	assert.Nil(t, err)
	defer backend.Close()

	baseinteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)
	nftA, err := nft.NewERC721Interactions(baseinteractions, *contractAddr, []nft.BaseNFTSignature{nft.BalanceOf})
	assert.Nil(t, err)

	enum, err := enumerable.NewERC721EnumerableInteractions(nftA, []enumerable.IERC721EnumerableSignature{
		enumerable.TokenOfOwnerByIndex,
	})
	assert.Nil(t, err)

	tokens, err := enum.GetAddressOwnedTokens(auth.From)
	assert.Nil(t, err)

	balance, err := nftA.BalanceOf(auth.From)
	assert.Nil(t, err)

	assert.Equal(t, balance.Int64(), int64(len(tokens)), "number of tokens should equal balance")
}

// Test_GetAllTokenIDs tests retrieval of all token IDs and validates that the total supply equals the number of tokens returned.
func Test_GetAllTokenIDs(t *testing.T) {
	backend, _, contractAddr, privKey, err := utils.SetupBlockchain(t,
		ERC721Complete.ERC721CompleteABI,
		ERC721Complete.ERC721CompleteBin,
		"MyNFT",
		"MNFT",
	)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer backend.Close()

	baseinteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)
	nftA, err := nft.NewERC721Interactions(baseinteractions, *contractAddr, []nft.BaseNFTSignature{nft.BalanceOf})
	if err != nil {
		t.Fatal(err.Error())
	}

	enum, err := enumerable.NewERC721EnumerableInteractions(nftA, []enumerable.IERC721EnumerableSignature{
		enumerable.TokenByIndex,
	})
	if err != nil {
		t.Fatal(err.Error())
	}

	tokens, err := enum.GetAllTokenIDs()
	assert.Nil(t, err)

	supply, err := nftA.TotalSupply()
	assert.Nil(t, err)

	assert.Equal(t, supply.Int64(), int64(len(tokens)), "number of tokens should equal total supply")
}
