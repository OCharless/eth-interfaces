package contractextension

import (
	"context"

	"github.com/OCharless/eth-interfaces/utils"
	"github.com/ethereum/go-ethereum/common"
)

type ContractImplementation interface {
	GetAddress() common.Address
	VerifyTransaction(ctx context.Context, to common.Address, data []byte, value int64) error
}

func SimulateCall(ctx context.Context,
	abiPath,
	selector string,
	contract ContractImplementation,
	params ...interface{},
) error {
	encodedFunction, err := utils.GetEncodedFunction(abiPath, selector, params...)
	if err != nil {
		return err
	}
	return contract.VerifyTransaction(ctx, contract.GetAddress(), encodedFunction, 0)
}
