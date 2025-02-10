package base

// Package baseinteractions provides core utilities for interacting with the blockchain, including transaction setup, contract calls, and error handling.

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/OCharless/eth-interfaces/customerrors"
	"github.com/OCharless/eth-interfaces/inferences/IERC165"
	Disperse "github.com/OCharless/eth-interfaces/inferences/disperse"
	"github.com/OCharless/eth-interfaces/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
)

// BaseInteractions holds the context, client, sender address, private key, disperse contract, and explorer URL.
type BaseInteractions struct {
	Ctx      context.Context
	Client   simulated.Client
	Address  common.Address
	pk       *ecdsa.PrivateKey
	disperse *Disperse.Disperse
	explorer *string
}

// IBaseInteractions defines the interface for verifying transactions.
type IBaseInteractions interface {
	VerifyTransaction(ctx context.Context, to common.Address, data []byte, value int64) error
}

// NewBaseInteractions creates a new instance of BaseInteractions for blockchain interaction.
func NewBaseInteractions(client simulated.Client, pk *ecdsa.PrivateKey, explorer *string) *BaseInteractions {
	ctx := context.TODO()
	_, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := pk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return &BaseInteractions{ctx, client, fromAddress, pk, nil, explorer}
}

// SetDisperse initializes the disperse contract for multi-address fund transfers.
func (b *BaseInteractions) SetDisperse(address string) error {
	var err error
	b.disperse, err = Disperse.NewDisperse(common.HexToAddress(address), b.Client)
	return err
}

// BaseTxSetup sets up transaction options (nonce, gas price, chain ID, etc.) for sending a transaction.
func (b *BaseInteractions) BaseTxSetup() (*bind.TransactOpts, error) {
	gasPrice, err := b.Client.SuggestGasPrice(b.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %v", err)
	}
	nonce, err := b.Client.PendingNonceAt(b.Ctx, b.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get user nonce: %v", err)
	}

	chainID, err := b.Client.ChainID(b.Ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	opts, err := bind.NewKeyedTransactorWithChainID(b.pk, chainID)
	if err != nil {
		return nil, err
	}

	opts.From = b.Address
	opts.GasPrice = gasPrice
	opts.Nonce = big.NewInt(int64(nonce))

	return opts, nil
}

// BaseCallSetup returns the call options for read-only contract operations.
func (b *BaseInteractions) BaseCallSetup() *bind.CallOpts {
	return &bind.CallOpts{
		From:    b.Address,
		Pending: false,
	}
}

// CatchTx waits for a transaction to be mined and returns its hash or an error message.
func (b *BaseInteractions) CatchTx(tx *types.Transaction, err error) (string, error) {
	if err != nil {
		return FailedTx(err)
	}
	receipt, err := bind.WaitMined(context.Background(), b.Client, tx)
	if receipt == nil {
		return FailedTx(err)
	}
	if b.explorer != nil {
		return SuccessTx(fmt.Sprintf(*b.explorer, "/tx/", receipt.TxHash.Hex()))
	}
	return SuccessTx(receipt.TxHash.Hex())

}

// VerifyTransaction simulates a contract call to verify transaction validity.
func (b *BaseInteractions) VerifyTransaction(ctx context.Context, to common.Address, data []byte, value int64) error {
	callMsg := ethereum.CallMsg{
		From:  b.Address,
		To:    &to,
		Data:  data,
		Value: big.NewInt(value),
	}

	_, err := b.Client.CallContract(ctx, callMsg, nil)
	return err
}

// Disperse uses the disperse contract to send funds to multiple addresses.
func (b *BaseInteractions) Disperse(addresses []common.Address, totalValue uint) (string, error) {
	if b.disperse == nil {
		return FailedTx(fmt.Errorf("disperse contract not initialized"))
	}
	opts, err := b.BaseTxSetup()
	if err != nil {
		return FailedTx(err)
	}
	amounts := []*big.Int{}
	for range addresses {
		amounts = append(amounts, big.NewInt(int64(totalValue)/int64(len(addresses))))
	}
	opts.Value = big.NewInt(int64(totalValue))
	fmt.Println("Dispersing...")
	return b.CatchTx(b.disperse.DisperseEther(opts, addresses, amounts))
}

// SendAllFunds transfers the entire balance to a designated address after fee estimation.
func (b *BaseInteractions) SendAllFunds(to common.Address) (*types.Transaction, error) {
	bn, err := b.Client.BlockNumber(b.Ctx)
	if err != nil {
		return nil, err
	}
	balance, err := b.Client.BalanceAt(b.Ctx, b.Address, big.NewInt(int64(bn)))
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		From:  b.Address,
		To:    &to,
		Value: balance,
		Data:  nil,
	}

	gasLimit, err := b.Client.EstimateGas(b.Ctx, msg)
	if err != nil {
		return nil, err
	}

	gasPrice, err := b.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	value := big.NewInt(0).Add(balance, big.NewInt(int64(gasLimit*gasPrice.Uint64())))
	if value.Int64() < balance.Int64() {
		return b.TransferETH(to, value)
	}
	return nil, fmt.Errorf(
		"fees exceed balances\nfees : %f ETH\nbalance : %f ETH",
		utils.ParseEther(big.NewInt(0).Sub(balance, value)),
		utils.ParseEther(balance),
	)
}

// TransferETH transfers Ether to the specified address, ensuring sufficient balance and proper fee estimation.
func (b *BaseInteractions) TransferETH(to common.Address, value *big.Int) (*types.Transaction, error) {

	balance, err := b.Client.BalanceAt(b.Ctx, to, nil)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{From: b.Address, To: &to, Value: balance, Data: nil}

	gasLimit, err := b.Client.EstimateGas(b.Ctx, msg)
	if err != nil {
		return nil, err
	}

	gasPrice, err := b.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	txCost := big.NewInt(0).Add(value, big.NewInt(int64(gasPrice.Uint64()*gasLimit)))
	if txCost.Uint64() < balance.Uint64() {
		return nil, fmt.Errorf(
			"unsufficient balance for the transfer\n value + fees : %f ETH\nbalance : %f ETH",
			utils.ParseEther(txCost),
			utils.ParseEther(balance),
		)
	}

	nonce, err := b.Client.PendingNonceAt(context.Background(), b.Address)
	if err != nil {
		return nil, err
	}

	// Create the transaction
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
	})
	// Get the chain ID
	chainID, err := b.Client.ChainID(b.Ctx)
	if err != nil {
		return nil, err
	}

	// Sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), b.pk)
	if err != nil {
		return nil, fmt.Errorf("failed to sign the tx: %w", err)
	}

	// Broadcast the transaction
	err = b.Client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return nil, fmt.Errorf("failed to send the tx: %w", err)
	}

	return signedTx, nil
}

// SupportsInterface checks if a contract supports a specific interface.
func (b *BaseInteractions) SupportsInterface(address common.Address, signature [4]byte) (bool, error) {
	ierc165, err := IERC165.NewIERC165(address, b.Client)
	if err != nil {
		return false, err
	}

	callopts := b.BaseCallSetup()
	return ierc165.SupportsInterface(callopts, signature)
}

// CheckSignatures checks if a contract supports specific function signatures.
func (b *BaseInteractions) CheckSignatures(contractAddress common.Address, signatures []utils.Signature) error {
	byteCode, err := b.Client.CodeAt(b.Ctx, contractAddress, nil)
	if err != nil {
		return fmt.Errorf("failed to get contract bytecode: %w", err)
	}
	notSupported := ""
	byteCodeHex := common.Bytes2Hex(byteCode)
	for _, signature := range signatures {
		selector := utils.GetFunctionSelector(signature)
		if !strings.Contains(byteCodeHex, selector) {
			if notSupported != "" {
				notSupported = fmt.Sprintf("%s, %s: %s", notSupported, signature, selector)
			} else {
				notSupported = fmt.Sprintf("%s: %s", signature, selector)
			}
		}
	}
	if len(notSupported) > 0 {
		return customerrors.WrapinterfacingError("CheckSignatures", fmt.Errorf("not supported functions: %s", notSupported))
	}
	return nil
}

// ManageCustomContractError handles custom contract errors.
func (b *BaseInteractions) ManageCustomContractError(abiString string, err error) error {
	if len(abiString) == 0 {
		return err
	}
	errBytes, success := ethclient.RevertErrorData(err)
	if success {
		customErr, err := b.MatchErrors(abiString, errBytes)
		if err != nil {
			return customerrors.WrapinterfacingError("MatchErrors", err)
		}
		return errors.New(customErr)
	}
	return nil
}

// MatchErrors matches error bytes to custom error names.
func (b *BaseInteractions) MatchErrors(abiString string, errBytes []byte) (string, error) {
	abi, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		return "", err
	}
	if len(errBytes) < 4 {
		panic("invalid error data")
	}
	methodID := errBytes[:4]

	// Find matching custom error
	var errorName string
	for _, abiError := range abi.Errors {
		abierr := abiError.ID.Bytes()[:4]
		if bytes.Equal(abierr[:4], methodID) {
			errorName = abiError.Name
			break
		}
	}
	return errorName, nil
}
