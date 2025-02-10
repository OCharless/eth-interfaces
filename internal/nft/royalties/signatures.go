package royalties

import (
	"encoding/hex"

	"github.com/OCharless/eth-interfaces/internal/nft"
	"github.com/ethereum/go-ethereum/crypto"
)

type IERC721RoyaltiesSignature nft.BaseNFTSignature

const (
	RoyaltyInfo IERC721RoyaltiesSignature = "royaltyInfo(uint256,uint256)"
)

func (s IERC721RoyaltiesSignature) GetHex() string {
	hash := crypto.NewKeccakState()
	hash.Write([]byte(s))
	selector := hash.Sum(nil)[:4]
	return hex.EncodeToString(selector)
}
