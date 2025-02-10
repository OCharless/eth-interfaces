package enumerable

import (
	"encoding/hex"

	"github.com/OCharless/eth-interfaces/nft"
	"github.com/ethereum/go-ethereum/crypto"
)

type IERC721EnumerableSignature nft.BaseNFTSignature

const (
	TokenOfOwnerByIndex IERC721EnumerableSignature = "tokenOfOwnerByIndex(address,uint256)"
	TokenByIndex        IERC721EnumerableSignature = "tokenByIndex(uint256)"
)

func (s IERC721EnumerableSignature) GetHex() string {
	hash := crypto.NewKeccakState()
	hash.Write([]byte(s))
	selector := hash.Sum(nil)[:4]
	return hex.EncodeToString(selector)
}
