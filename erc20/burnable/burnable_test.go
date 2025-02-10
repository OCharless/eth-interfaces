package burnable_test

// Package burnable_test contains tests for burnable interactions.

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/OCharless/eth-interfaces/base"
	"github.com/OCharless/eth-interfaces/erc20"
	"github.com/OCharless/eth-interfaces/erc20/burnable"
	"github.com/OCharless/eth-interfaces/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

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
			Name:          "KO - ERC20 doesn't implement the interface",
			ExpectError:   true,
			ContractAddr:  erc20Contract,
			ExpectedError: "interface setup error function CheckSignatures, error :",
		},
	}

	baseInteractions := base.NewBaseInteractions(backend.Client(), privKey, nil)
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			erc20Interactions, err := erc20.NewIERC20Interactions(baseInteractions, tt.ContractAddr, []erc20.BaseERC20Signature{})
			if err != nil {
				t.Fatalf("failed to create interactions interface, error: %s", err.Error())
			}
			_, err = burnable.NewIERC20Burnable(
				erc20Interactions,
				[]burnable.ERC20BurnableSignatures{
					burnable.Burn,
					burnable.BurnFrom,
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

// Test_Burn tests the burn functionality and ensures that the token burn behaves as expected.
func Test_Burn(t *testing.T) {
	backend, _, contractAddress, privKey, err := utils.SetupBlockchain(t,
		"/build/ERC20Burnable.abi",
		"/build/ERC20Burnable.bin",
	)
	if err != nil {
		t.Fatal(err)
	}

	type transferArgs struct {
		pk  *ecdsa.PrivateKey
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
			Name: "OK - Successfully burn 10 tokens",
			args: transferArgs{
				qty: big.NewInt(10),
			},
			ContractAddr: *contractAddress,
		},
		{
			Name: "KO - Unsufficient balance",
			args: transferArgs{
				pk: func() *ecdsa.PrivateKey {
					key, _ := crypto.GenerateKey()
					return key
				}(),
				qty: big.NewInt(1),
			},
			ContractAddr:  *contractAddress,
			ExpectError:   true,
			ExpectedError: "call error on erc20.Burn(): ERC20InsufficientBalance",
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
			session, err := erc20.NewIERC20Interactions(baseInteractions, tt.ContractAddr, []erc20.BaseERC20Signature{erc20.Name, erc20.BalanceOf})
			if err != nil {
				t.Fatal("setting up should not fail")
			}
			zeroAddressBalance, err := session.GetBalance()
			if err != nil {
				t.Fatal("failed to get zero address balance")
			}

			burn, err := burnable.NewIERC20Burnable(session, []burnable.ERC20BurnableSignatures{burnable.Burn})
			if err != nil {
				t.Fatal("setting up should not fail")
			}
			_, err = burn.Burn(tt.args.qty)
			backend.Commit()
			if tt.ExpectError {
				if err == nil {
					t.Error("expected error but there's none")
					return
				}
				assert.Contains(t, err.Error(), tt.ExpectedError)
			} else {
				assert.Nil(t, err)
				bal, err := session.GetBalance()
				if err != nil {
					t.Fatal("failed to get owner")
				}
				assert.Equal(t, -1, bal.Cmp(zeroAddressBalance))
			}
		})
	}
}
