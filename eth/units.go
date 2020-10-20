package eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

var (
	etherFloat  = big.NewFloat(params.Ether)
	etherInt, _ = etherFloat.Int(nil)
)

func EthToWeiFloat(ethv float64) *big.Int {
	weifl := big.NewFloat(ethv)
	weifl.Mul(weifl, etherFloat)
	wei, _ := weifl.Int(nil)
	return wei
}

func EthToWeiInt(ethv int64) *big.Int {
	wei := big.NewInt(ethv)
	wei.Mul(wei, etherInt)
	return wei
}

func WeiToEthInt(weiv *big.Int) *big.Int {
	return new(big.Int).Div(weiv, etherInt)
}

func WeiToEthFloat(weiv *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(weiv), etherFloat)
}
