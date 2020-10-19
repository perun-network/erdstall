package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func EthToWeiFloat(ethv float64) *big.Int {
	weifl := big.NewFloat(ethv)
	weifl.Mul(weifl, big.NewFloat(params.Ether))
	wei, _ := weifl.Int(nil)
	return wei
}

func EthToWeiInt(ethv int64) *big.Int {
	wei := big.NewInt(ethv)
	ether, _ := big.NewFloat(params.Ether).Int(nil)
	wei.Mul(wei, ether)
	return wei
}
