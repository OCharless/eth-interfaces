package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
)

func LoadBytecode(path string) ([]byte, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))

	content, err := os.ReadFile(dir + path)
	if err != nil {
		return nil, err
	}
	// Remove any trailing newlines or spaces
	return common.Hex2Bytes(strings.TrimSpace(string(content))), nil
}

func LoadAbi(path string) (abi.ABI, error) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))

	content, err := os.ReadFile(dir + path)
	if err != nil {
		return abi.ABI{}, err
	}
	// Remove any trailing newlines or spaces
	// return content, nil

	contractABI, err := abi.JSON(strings.NewReader(string(content)))
	if err != nil {
		return abi.ABI{}, err
	}
	return contractABI, nil
}

func GetFunctionSelector(signature Signature) string {
	return signature.GetHex()
}

func GetEncodedFunction(abiString, signature string, params ...interface{}) ([]byte, error) {
	contractABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return nil, err
	}
	return contractABI.Pack(signature, params...)
}

func DeployContract(auth *bind.TransactOpts,
	client simulated.Client,
	abiPath string,
	byteCodePath string,
	params ...interface{},
) (common.Address, *types.Transaction, *bind.BoundContract, error) {
	contractABI, err := LoadAbi(abiPath)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	byteCode, err := LoadBytecode(byteCodePath)
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	return bind.DeployContract(
		auth,
		contractABI,
		byteCode,
		client,
		params...,
	)
}
