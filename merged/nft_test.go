package merged_test

import (
	"fmt"
	"testing"

	"github.com/OCharless/eth-interfaces/base"
	"github.com/OCharless/eth-interfaces/merged"
	"github.com/OCharless/eth-interfaces/nft"
	"github.com/OCharless/eth-interfaces/nft/enumerable"
	"github.com/OCharless/eth-interfaces/nft/royalties"
	"github.com/OCharless/eth-interfaces/utils"
	"github.com/stretchr/testify/assert"
)

// TestAllInfosSuccess tests the successful execution of AllInfos with valid dummyBase and both extensions enabled.
func Test_Instantiation(t *testing.T) {

	type args struct {
		contractPath string
		extensions   []merged.ExtensionEnum
		signatures   []utils.Signature
	}

	testCases := []struct {
		Name          string
		Args          args
		ExpectError   bool
		ExpectedError string
	}{
		{
			Name: "OK - Instantiate NFT with enumerable extension",
			Args: args{
				contractPath: "/build/ERC721A.%s",
				extensions:   []merged.ExtensionEnum{merged.Enumerable},
				signatures: []utils.Signature{
					enumerable.TokenByIndex,
					enumerable.TokenOfOwnerByIndex,
				},
			},
		},
		{
			Name: "OK - Instantiate NFT with royalties extension",
			Args: args{
				contractPath: "/build/ERC721ARoyalties.%s",
				extensions:   []merged.ExtensionEnum{merged.Royalties},
				signatures:   []utils.Signature{royalties.RoyaltyInfo},
			},
		},
		{
			Name: "OK - Instantiate NFT with royalties extension and enumerable extension",
			Args: args{
				contractPath: "/build/ERC721ARoyalties.%s",
				extensions:   []merged.ExtensionEnum{merged.Royalties, merged.Enumerable},
				signatures: []utils.Signature{
					royalties.RoyaltyInfo,
					enumerable.TokenByIndex,
					enumerable.TokenOfOwnerByIndex,
				},
			},
		},
		{
			Name: "NOK - Instantiate NFT with enumerable and royalties extensions",
			Args: args{
				contractPath: "/build/ERC721A.%s",
				extensions:   []merged.ExtensionEnum{merged.Enumerable, merged.Royalties},
				signatures:   []utils.Signature{enumerable.TokenOfOwnerByIndex, royalties.RoyaltyInfo},
			},
			ExpectError:   true,
			ExpectedError: "not supported functions: royaltyInfo(uint256,uint256)",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {

			backend, _, contractAddr, privKey, err := utils.SetupBlockchain(t,
				fmt.Sprintf(tt.Args.contractPath, "abi"),
				fmt.Sprintf(tt.Args.contractPath, "bin"),
				"MyNFT",
				"MNFT",
			)
			assert.Nil(t, err)
			defer backend.Close()

			baseInteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)

			nftA, err := nft.NewERC721Interactions(baseInteractions, *contractAddr, []nft.BaseNFTSignature{})
			assert.Nil(t, err)

			// Create a new summed interactions with both Enumerable and Royalties extensions
			_, err = merged.NewERC721SummedInteractions(nftA, tt.Args.signatures, tt.Args.extensions...)
			if tt.ExpectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.ExpectedError)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}
