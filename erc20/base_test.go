package erc20_test

// Package nft_test contains tests for NFT interactions defined in base.go.

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"testing"

	"github.com/OCharless/eth-interfaces/base"
	"github.com/OCharless/eth-interfaces/erc20"
	"github.com/OCharless/eth-interfaces/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

// Test_DeploySuccessfully tests if the blockchain setup and contract deployment succeed without errors.
func Test_DeploySuccessfully(t *testing.T) {
	backend, _, _, _, err := utils.SetupBlockchain(t,
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
	)
	assert.Nil(t, err, "failed to create interactions interface, error: %w", err)
	backend.Close()
}

// Test_Instantiation verifies that the NFT interactions interface is correctly instantiated using various contracts, including a valid NFT contract, an empty contract, and an ERC20 contract.
func Test_Instantiation(t *testing.T) {
	backend, auth, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
	)
	if err != nil {
		t.Fatal(err)
	}

	defer backend.Close()

	emptyContract, err := utils.DeployEmptyContract(auth, backend)
	if err != nil {
		log.Fatalf("failed to deploy empty contract: %s", err)
	}

	erc721CompleteContract, tx, _, err := utils.DeployContract(
		auth,
		backend.Client(),
		"/build/ERC721Complete.abi",
		"/build/ERC721Complete.bin",
		"MyNFT", // Arg 1: name
		"MNFT",  // Arg 2: symbol
	)
	if err != nil {
		t.Fatalf("failed to deploy ERC721A contract: %s", err)
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
			Name:          "KO - ERC721A doesn't implement the interface",
			ExpectError:   true,
			ContractAddr:  erc721CompleteContract,
			ExpectedError: "interface setup error function CheckSignatures, error :",
		},
	}

	baseInteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			_, err := erc20.NewIERC20Interactions(
				baseInteractions,
				tt.ContractAddr,
				[]erc20.BaseERC20Signature{
					erc20.Name,
					erc20.Symbol,
					erc20.Decimals,
				},
			)
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
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
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
			ExpectedResult: "NvidiaToken",
			ContractAddr:   *contractAddress,
		},
	}

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			session, err := erc20.NewIERC20Interactions(base, tt.ContractAddr, []erc20.BaseERC20Signature{erc20.Name})
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
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
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
			ExpectedResult: "NVD",
			ContractAddr:   *contractAddress,
		},
	}

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			session, err := erc20.NewIERC20Interactions(base, tt.ContractAddr, []erc20.BaseERC20Signature{erc20.Symbol})
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
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	testCases := []struct {
		Name           string
		ContractAddr   common.Address
		ExpectedResult *big.Int
		ExpectError    bool
		ExpectedError  string
	}{
		{
			Name:           "OK - Successfully get NFT total supply",
			ExpectedResult: big.NewInt(100_000_000_000),
			ContractAddr:   *contractAddress,
		},
	}

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			session, err := erc20.NewIERC20Interactions(base, tt.ContractAddr, []erc20.BaseERC20Signature{erc20.TotalSupply})
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
				assert.Equal(t, 0, supply.Cmp(big.NewInt(0).Mul(tt.ExpectedResult, big.NewInt(1e18))))
			}
		})
	}
}

// Test_Transfer tests the transfer functionality and ensures that the token transfer behaves as expected.
func Test_Transfer(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
	)
	if err != nil {
		t.Fatal(err)
	}

	type transferArgs struct {
		pk  *ecdsa.PrivateKey
		To  common.Address
		qty *big.Int
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
				To:  common.HexToAddress("1"),
				qty: big.NewInt(10),
			},
			ContractAddr: *contractAddress,
		},
		{
			Name: "KO - Burn NFT",
			args: transferArgs{
				To:  common.HexToAddress("0"),
				qty: big.NewInt(1),
			},
			ContractAddr:  *contractAddress,
			ExpectError:   true,
			ExpectedError: "call error on erc20.Transfer(): ERC20InvalidReceiver",
		},
		{
			Name: "KO - Unsufficient balance",
			args: transferArgs{
				pk: func() *ecdsa.PrivateKey {
					key, _ := crypto.GenerateKey()
					return key
				}(),
				To:  crypto.PubkeyToAddress(privKey.PublicKey),
				qty: big.NewInt(1),
			},
			ContractAddr:  *contractAddress,
			ExpectError:   true,
			ExpectedError: "call error on erc20.Transfer(): ERC20InsufficientBalance",
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
			session, err := erc20.NewIERC20Interactions(baseInteractions, tt.ContractAddr, []erc20.BaseERC20Signature{erc20.TransferFrom})
			if err != nil {
				t.Fatal("setting up should not fail")
			}
			_, err = session.TransferTo(tt.args.To, tt.args.qty)
			backend.Commit()
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Contains(t, err.Error(), tt.ExpectedError)
			} else {
				assert.Nil(t, err)
				bal, err := session.BalanceOf(tt.args.To)
				if err != nil {
					t.Fatal("failed to get owner")
				}
				assert.Zero(t, bal.Cmp(tt.args.qty))
			}
		})
	}
}

// Test_GetBalance verifies that the NFT balance is correctly returned for an address.
func Test_GetBalance(t *testing.T) {
	backend, auth, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	token, err := erc20.NewIERC20Interactions(base, *contractAddress, []erc20.BaseERC20Signature{erc20.BalanceOf}, auth)
	assert.Nil(t, err)

	balance, err := token.GetBalance()
	assert.Nil(t, err)
	assert.Equal(t, 0, balance.Cmp(big.NewInt(0).Mul(big.NewInt(100_000_000_000), big.NewInt(1e18))))
}

// Test_BalanceOf verifies the BalanceOf function for different addresses.
func Test_BalanceOf(t *testing.T) {
	backend, auth, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	token, err := erc20.NewIERC20Interactions(base, *contractAddress, []erc20.BaseERC20Signature{erc20.BalanceOf})
	assert.Nil(t, err)

	testCases := []struct {
		Name           string
		Owner          common.Address
		ExpectedResult int64
		ExpectError    bool
		ExpectedError  string
	}{
		{
			Name:           "OK - non empty balance",
			Owner:          auth.From,
			ExpectedResult: 100_000_000_000,
		},
		{
			Name:           "OK - empty balance",
			Owner:          common.HexToAddress("1"),
			ExpectedResult: 0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			balance, err := token.BalanceOf(tt.Owner)
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Contains(t, err.Error(), tt.ExpectedError)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, 0, balance.Cmp(big.NewInt(0).Mul(big.NewInt(tt.ExpectedResult), big.NewInt(1e18))))
			}
		})
	}
}

// Test_Approve tests the approval functionality for token transfers.
func Test_Approve(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	type approveArgs struct {
		To  common.Address
		qty *big.Int
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
				To:  common.HexToAddress("1"),
				qty: big.NewInt(1),
			},
		},
		{
			name: "NOK - ZeroAddress spender",
			args: approveArgs{
				To:  common.HexToAddress("0"),
				qty: big.NewInt(10),
			},
			expectError:   true,
			errorContains: "call error on erc20.Approve(): ERC20InvalidSpender",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseInteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)

			token, err := erc20.NewIERC20Interactions(baseInteractions, *contractAddress, []erc20.BaseERC20Signature{erc20.Approve})
			if err != nil {
				t.Fatal(err)
			}

			_, err = token.Approve(tt.args.To, tt.args.qty)
			backend.Commit()

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorContains)
			} else {
				assert.Nil(t, err)
				approved, err := token.Allowance(baseInteractions.Address, tt.args.To)
				assert.Nil(t, err)
				assert.Equal(t, tt.args.qty.Int64(), approved.Int64())
			}
		})
	}
}

// Test_TokenMetaInfos verifies that the metadata (name, symbol, and URI) for a token is correctly retrieved.
func Test_TokenMetaInfos(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
	)
	if err != nil {
		t.Fatal(err)
	}
	defer backend.Close()

	base := base.NewBaseInteractions(backend.Client(), privKey, nil)
	token, err := erc20.NewIERC20Interactions(base, *contractAddress, []erc20.BaseERC20Signature{erc20.Name, erc20.Symbol})
	assert.Nil(t, err)

	// Test meta infos for token
	tokenInfo, err := token.TokenMetaInfos()
	assert.Nil(t, err)
	assert.Equal(t, "NvidiaToken", tokenInfo.Name)
	assert.Equal(t, "NVD", tokenInfo.Symbol)
}
