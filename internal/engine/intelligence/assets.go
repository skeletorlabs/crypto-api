package intelligence

var SupportedAssets = map[string]struct{}{
	"bitcoin": {},
	// "ethereum": {},
	// "solana": {},
}

func IsSupportedAsset(symbol string) bool {
	_, ok := SupportedAssets[symbol]
	return ok
}
