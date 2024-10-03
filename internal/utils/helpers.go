package utils

import "strings"

// GetTokensFromTradingPair returns the two tokens from a trading pair
func GetTokensFromTradingPair(tradingPair string) (string, string) {
	v := strings.SplitN(tradingPair, "-", 2)
	return v[0], v[1]
}
