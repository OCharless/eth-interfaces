package utils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
)

func SetupBlockchain(
	t *testing.T,
	contractABI string,
	byteCode string,
	params ...interface{},
) (
	*simulated.Backend,
	*bind.TransactOpts,
	*common.Address,
	*ecdsa.PrivateKey,
	error,
) {
	privKey, _ := crypto.GenerateKey()
	auth, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(1337))
	if err != nil {
		return nil, nil, nil, nil, err
	}

	testUserAddress := crypto.PubkeyToAddress(privKey.PublicKey)
	alloc := types.GenesisAlloc{
		testUserAddress: {Balance: MAX_UINT256},
	}
	backend := simulated.NewBackend(alloc, simulated.WithBlockGasLimit(9_000_000))

	contractAddr, tx, _, err := DeployContract(
		auth,
		backend.Client(),
		contractABI,
		byteCode,
		params...,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	backend.Commit()

	receipt, err := backend.Client().TransactionReceipt(context.Background(), tx.Hash())
	if err != nil || receipt.Status != 1 {
		t.Fatal("contract deployment failed")
	}

	return backend, auth, &contractAddr, privKey, nil
}

func DeployEmptyContract(auth *bind.TransactOpts, backend *simulated.Backend) (*common.Address, error) {
	contractAddr, tx, _, err := DeployContract(
		auth,
		backend.Client(),
		"/build/EmptyContract.abi",
		"/build/EmptyContract.bin",
	)
	if err != nil {
		return nil, err
	}
	backend.Commit()

	receipt, err := backend.Client().TransactionReceipt(context.Background(), tx.Hash())
	if err != nil || receipt.Status != 1 {
		return nil, fmt.Errorf("empty contract deployment failed: %w", err)
	}
	return &contractAddr, nil
}
