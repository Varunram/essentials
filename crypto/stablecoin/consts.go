package stablecoin

import (
	"os"
)

var (
	// XLM Stablecoin constants
	HomeDir              = os.Getenv("HOME") + "/.openx"
	StablecoinCode       = "STABLEUSD"
	StablecoinPublicKey  = ""
	StablecoinSeed       = ""
	StableCoinSeedFile   = HomeDir + "/stablecoinseed.hex"
	StableCoinAddress    = "GDJE64WOXDXLEK7RDURVYEJ5Y5XFHS6OQZCS3SHO4EEMTABEIJXF6SZ5"
	StablecoinTrustLimit = "1000000000"

	// AnchorUSD constants
	AnchorUSDCode       = "USD"
	AnchorUSDAddress    = "GCKFBEIYV2U22IO2BJ4KVJOIP7XPWQGQFKKWXR6DOSJBV7STMAQSMTGG"
	AnchorUSDTrustLimit = "1000000"
)
