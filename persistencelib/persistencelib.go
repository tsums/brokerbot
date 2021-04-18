package persistencelib

import (
	"strings"
)

var aliasMap = map[string][]string{
	"CRYPTO": {"$BTC", "$ETH", "$LTC", "$LINK", "$BCH", "$ZEC"},
	"MEMES":  {"THCX", "PLUG", "FCEL", "BLDP", "NVDA"},
	"FAANG":  {"FB", "AMZN", "AAPL", "NFLX", "GOOG"},
	"DEFI":   {"$UNI", "$YFI", "$COMP", "$MKR", "$AAVE", "$CRV", "$SUSHI"},
}

func ExpandAlias(s string) (ret []string) {
	trim := strings.TrimPrefix(s, "?")

	if val, ok := aliasMap[trim]; ok {
		return val
	}
	return nil
}
