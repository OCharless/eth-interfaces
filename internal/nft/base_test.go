package nft_test

// Package nft_test contains tests for NFT interactions defined in base.go.

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	"github.com/OCharless/eth-interfaces/internal/base"
	"github.com/OCharless/eth-interfaces/internal/erc20"
	"github.com/OCharless/eth-interfaces/internal/nft"
	"github.com/OCharless/eth-interfaces/internal/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// Test_DeploySuccessfully tests if the blockchain setup and contract deployment succeed without errors.
func Test_DeploySuccessfully(t *testing.T) {
	backend, _, _, _, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT", // Arg 1: name
		"MNFT",  // Arg 2: symbol
	)
	assert.Nil(t, err, "failed to create interactions interface, error: %w", err)
	backend.Close()
}

// Test_Instantiation verifies that the NFT interactions interface is correctly instantiated using various contracts, including a valid NFT contract, an empty contract, and an ERC20 contract.
func Test_Instantiation(t *testing.T) {
	backend, auth, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT", // Arg 1: name
		"MNFT",  // Arg 2: symbol
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	emptyContract, err := utils.DeployEmptyContract(auth, backend)
	if err != nil {
		log.Fatalf("failed to deploy empty contract: %s", err)
	}

	erc20Contract, tx, _, err := utils.DeployContract(
		auth,
		backend.Client(),
		"/build/ERC20.abi",
		"/build/ERC20.bin",
	)
	if err != nil {
		t.Fatalf("failed to deploy ERC20 contract: %s", err)
	}
	backend.Commit()
	receipt, err := backend.Client().TransactionReceipt(context.Background(), tx.Hash())
	if err != nil || receipt.Status != 1 {
		t.Fatalf("failed to deploy ERC20 contract: %s", err)
	}

	testCases := []struct {
		Name           string
		ContractAddr   common.Address
		ExpectedResult string
		ExpectError    bool
		ExpectedError  string
	}{
		{
			Name:           "OK - Successfully instantiated",
			ExpectedResult: "MyNFT",
			ContractAddr:   *contractAddress,
		},
		{
			Name:          "KO - Empty contract doesn't implement interface",
			ExpectError:   true,
			ContractAddr:  *emptyContract,
			ExpectedError: "interface setup error function CheckSignatures, error :",
		},
		{
			Name:          "KO - ERC20 doesn't implement the interface",
			ExpectError:   true,
			ContractAddr:  erc20Contract,
			ExpectedError: "interface setup error function CheckSignatures, error :",
		},
	}

	baseInteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := erc20.NewIERC20Interactions(baseInteractions, tt.ContractAddr, []erc20.BaseERC20Signature{erc20.Name, erc20.Symbol, erc20.TokenURI})
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Contains(t, err.Error(), tt.ExpectedError)
			} else {
				assert.NoError(t, err, "failed to create interactions interface, error: %w", err)
			}
		})
	}
}

// Test_Name verifies that the NFT contract correctly returns its name.
func Test_Name(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT", // Arg 1: name
		"MNFT",  // Arg 2: symbol
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	testCases := []struct {
		Name           string
		ContractAddr   common.Address
		ExpectedResult string
		ExpectError    bool
		ExpectedError  string
	}{
		{
			Name:           "OK - Successfully get NFT name",
			ExpectedResult: "MyNFT",
			ContractAddr:   *contractAddress,
		},
	}

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			session, err := nft.NewERC721Interactions(base, tt.ContractAddr, []nft.BaseNFTSignature{nft.Name})
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Equal(t, tt.ExpectedError, err.Error())
			} else {
				assert.Nil(t, err, "failed to create interactions interface, error: %w", err)
				name, err := session.Name()
				assert.Nil(t, err)
				assert.Equal(t, tt.ExpectedResult, name)
			}
		})
	}
}

// Test_Symbol verifies that the NFT contract correctly returns its symbol.
func Test_Symbol(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT", // Arg 1: name
		"MNFT",  // Arg 2: symbol
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	testCases := []struct {
		Name           string
		ContractAddr   common.Address
		ExpectedResult string
		ExpectError    bool
		ExpectedError  string
	}{
		{
			Name:           "OK - Successfully get NFT symbol",
			ExpectedResult: "MNFT",
			ContractAddr:   *contractAddress,
		},
	}

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			session, err := nft.NewERC721Interactions(base, tt.ContractAddr, []nft.BaseNFTSignature{nft.Symbol})
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Equal(t, tt.ExpectedError, err.Error())
			} else {
				assert.Nil(t, err, "failed to create interactions interface, error: %w", err)
				symbol, err := session.Symbol()
				assert.Nil(t, err)
				assert.Equal(t, tt.ExpectedResult, symbol)
			}
		})
	}
}

// Test_TotalSupply verifies that the total supply of NFTs is correctly reported by the contract.
func Test_TotalSupply(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT", // Arg 1: name
		"MNFT",  // Arg 2: symbol
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	testCases := []struct {
		Name           string
		ContractAddr   common.Address
		ExpectedResult uint64
		ExpectError    bool
		ExpectedError  string
	}{
		{
			Name:           "OK - Successfully get NFT total supply",
			ExpectedResult: 30,
			ContractAddr:   *contractAddress,
		},
	}

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			session, err := nft.NewERC721Interactions(base, tt.ContractAddr, []nft.BaseNFTSignature{nft.TotalSupply})
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Equal(t, tt.ExpectedError, err.Error())
			} else {
				assert.Nil(t, err, "failed to create interactions interface, error: %w", err)
				supply, err := session.TotalSupply()
				assert.Nil(t, err)
				assert.Equal(t, tt.ExpectedResult, supply.Uint64())
			}
		})
	}
}

// Test_OwnerOf verifies that the owner of a given token is correctly identified.
func Test_OwnerOf(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT", // Arg 1: name
		"MNFT",  // Arg 2: symbol
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	testCases := []struct {
		Name           string
		ContractAddr   common.Address
		argument       *big.Int
		ExpectedResult common.Address
		ExpectError    bool
		ExpectedError  string
	}{
		{
			Name:           "OK - Successfully get NFT owner",
			ExpectedResult: crypto.PubkeyToAddress(privKey.PublicKey),
			argument:       common.Big0,
			ContractAddr:   *contractAddress,
		},
		{
			Name:          "OK - Not minted, owner is zero address",
			argument:      big.NewInt(31),
			ContractAddr:  *contractAddress,
			ExpectError:   true,
			ExpectedError: "call error on nft.OwnerOf(): OwnerQueryForNonexistentToken",
		},
	}

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			session, err := nft.NewERC721Interactions(base, tt.ContractAddr, []nft.BaseNFTSignature{nft.OwnerOf})
			if err != nil {
				t.Fatal("setting up should not fail")
			}
			owner, err := session.OwnerOf(tt.argument)
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Equal(t, tt.ExpectedError, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.ExpectedResult, owner)
			}
		})
	}
}

// Test_Transfer tests the transfer functionality and ensures that the token transfer behaves as expected.
func Test_Transfer(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT", // Arg 1: name
		"MNFT",  // Arg 2: symbol
	)
	if err != nil {
		t.Fatal(err)
	}

	type transferArgs struct {
		pk      *ecdsa.PrivateKey
		To      common.Address
		TokenID *big.Int
	}

	testCases := []struct {
		Name          string
		ContractAddr  common.Address
		args          transferArgs
		ExpectError   bool
		ExpectedError string
	}{
		{
			Name: "OK - Successfully get transfer NFT",
			args: transferArgs{
				To:      common.HexToAddress("1"),
				TokenID: big.NewInt(10),
			},
			ContractAddr: *contractAddress,
		},
		{
			Name: "OK - Burn NFT",
			args: transferArgs{
				To:      common.HexToAddress("0"),
				TokenID: big.NewInt(1),
			},
			ContractAddr: *contractAddress,
		},
		{
			Name: "KO - Incorrect owner/Unallowed",
			args: transferArgs{
				pk: func() *ecdsa.PrivateKey {
					key, _ := crypto.GenerateKey()
					return key
				}(),
				To:      crypto.PubkeyToAddress(privKey.PublicKey),
				TokenID: common.Big0,
			},
			ContractAddr:  *contractAddress,
			ExpectError:   true,
			ExpectedError: "call error on nft.TransferFrom(): TransferFromIncorrectOwner",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			baseInteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)
			if tt.args.pk != nil {
				pk := tt.args.pk
				_, err := baseInteractions.TransferETH(crypto.PubkeyToAddress(pk.PublicKey), big.NewInt(1e18))
				if err != nil {
					t.Fatal(err)
				}

				backend.Commit()
				baseInteractions = base.NewBaseInteractions(backend.Client(), pk, nil)
			}
			session, err := nft.NewERC721Interactions(baseInteractions, tt.ContractAddr, []nft.BaseNFTSignature{nft.TransferFrom})
			if err != nil {
				t.Fatal("setting up should not fail")
			}
			_, err = session.TransferTo(tt.args.To, tt.args.TokenID)
			backend.Commit()
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Contains(t, err.Error(), tt.ExpectedError)
			} else {
				assert.Nil(t, err)
				owner, err := session.OwnerOf(tt.args.TokenID)
				if err != nil {
					t.Fatal("failed to get owner")
				}
				assert.Equal(t, tt.args.To, owner)
			}
		})
	}
}

// Test_GetBalance verifies that the NFT balance is correctly returned for an address.
func Test_GetBalance(t *testing.T) {
	backend, auth, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT",
		"MNFT",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	nft, err := nft.NewERC721Interactions(base, *contractAddress, []nft.BaseNFTSignature{nft.BalanceOf}, auth)
	assert.Nil(t, err)

	balance, err := nft.GetBalance()
	assert.Nil(t, err)
	assert.Equal(t, balance.Uint64(), uint64(30))
}

// Test_TransferFirstOwnedTo tests transferring the first owned token to a specified address.
func Test_TransferFirstOwnedTo(t *testing.T) {
	backend, auth, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT",
		"MNFT",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	type transferArgs struct {
		pk *ecdsa.PrivateKey
		To common.Address
	}

	tests := []struct {
		name          string
		args          transferArgs
		expectError   bool
		errorContains string
	}{
		{
			name: "OK -Successful transfer",
			args: transferArgs{
				To: common.HexToAddress("10"),
			},
			expectError: false,
		},
		{
			name: "OK - Burn to zero address",
			args: transferArgs{
				To: common.HexToAddress("0"),
			},
		},
		{
			name: "NOK - No NFTs owned",
			args: transferArgs{
				To: common.HexToAddress("1"),
				pk: func() *ecdsa.PrivateKey {
					key, _ := crypto.GenerateKey()
					return key
				}(),
			},
			expectError:   true,
			errorContains: "no nft found from signer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nftInterface *nft.ERC721Interactions
			baseInteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)
			if tt.args.pk != nil {
				pk := tt.args.pk
				_, err := baseInteractions.TransferETH(crypto.PubkeyToAddress(pk.PublicKey), big.NewInt(1e18))
				if err != nil {
					t.Fatal(err)
				}
				backend.Commit()
				baseInteractions = base.NewBaseInteractions(backend.Client(), pk, nil)
				if err != nil {
					t.Fatal(err)
				}
			} else {
				baseInteractions = base.NewBaseInteractions(backend.Client(), privKey, nil)
			}
			nftInterface, err = nft.NewERC721Interactions(
				baseInteractions,
				*contractAddress,
				[]nft.BaseNFTSignature{nft.TransferFrom, nft.OwnerOf},
				auth,
			)
			assert.Nil(t, err)

			_, err = nftInterface.TransferFirstOwnedTo(tt.args.To)
			backend.Commit()
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// Test_BalanceOf verifies the BalanceOf function for different addresses.
func Test_BalanceOf(t *testing.T) {
	backend, auth, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT",
		"MNFT",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	nft, err := nft.NewERC721Interactions(base, *contractAddress, []nft.BaseNFTSignature{nft.BalanceOf})
	assert.Nil(t, err)

	testCases := []struct {
		Name           string
		Owner          common.Address
		ExpectedResult *uint64
		ExpectError    bool
		ExpectedError  string
	}{
		{
			Name:           "NOK - Zero address",
			Owner:          common.Address{},
			ExpectedResult: nil,
			ExpectError:    true,
			ExpectedError:  "call error on nft.BalanceOf(): BalanceQueryForZeroAddress",
		},
		{
			Name:           "OK - non empty balance",
			Owner:          auth.From,
			ExpectedResult: func() *uint64 { val := uint64(30); return &val }(),
		},
		{
			Name:           "OK - empty balance",
			Owner:          common.HexToAddress("1"),
			ExpectedResult: func() *uint64 { val := uint64(0); return &val }(),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			balance, err := nft.BalanceOf(tt.Owner)
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Contains(t, err.Error(), tt.ExpectedError)
			} else {
				assert.Nil(t, err)
				if tt.ExpectedResult != nil {
					assert.Equal(t, *tt.ExpectedResult, balance.Uint64())
				}
			}
		})
	}
}

// Test_Approve tests the approval functionality for token transfers.
func Test_Approve(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT",
		"MNFT",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	type approveArgs struct {
		pk      *ecdsa.PrivateKey
		To      common.Address
		TokenID *big.Int
	}

	tests := []struct {
		name          string
		expectError   bool
		args          approveArgs
		errorContains string
	}{
		{
			name: "OK -Successful approval",
			args: approveArgs{
				To:      common.HexToAddress("1"),
				TokenID: big.NewInt(1),
			},
			expectError: false,
		},
		{
			name: "NOK - Not owner",
			args: approveArgs{
				To: common.HexToAddress("10"),
				pk: func() *ecdsa.PrivateKey {
					key, _ := crypto.GenerateKey()
					return key
				}(),
				TokenID: big.NewInt(10),
			},
			expectError:   true,
			errorContains: "call error on nft.Approve(): ApprovalCallerNotOwnerNorApproved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseInteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)
			if tt.args.pk != nil {
				pk := tt.args.pk
				_, err := baseInteractions.TransferETH(crypto.PubkeyToAddress(pk.PublicKey), big.NewInt(1e18))
				if err != nil {
					t.Fatal(err)
				}
				backend.Commit()
				baseInteractions = base.NewBaseInteractions(backend.Client(), pk, nil)
			}

			nft, err := nft.NewERC721Interactions(baseInteractions, *contractAddress, []nft.BaseNFTSignature{nft.Approve, nft.GetApproved})
			if err != nil {
				t.Fatal(err)
			}

			_, err = nft.Approve(tt.args.To, tt.args.TokenID)
			backend.Commit()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.Nil(t, err)
				approved, err := nft.GetApproved(tt.args.TokenID)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.To.Hex(), approved.Hex())
			}
		})
	}
}

// Test_TokenMetaInfos verifies that the metadata (name, symbol, and URI) for a token is correctly retrieved.
func Test_TokenMetaInfos(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT",
		"MNFT",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	nft, err := erc20.NewIERC20Interactions(base, *contractAddress, []erc20.BaseERC20Signature{erc20.Name, erc20.Symbol, erc20.TokenURI})
	assert.Nil(t, err)

	// Test meta infos for token
	nftInfo, err := nft.TokenMetaInfos()
	assert.Nil(t, err)
	assert.Equal(t, "MyNFT", nftInfo.Name)
	assert.Equal(t, "MNFT", nftInfo.Symbol)
	assert.Empty(t, nftInfo.URI)
}
