package stablecoin

import (
	"os"
)

var (
	// HomeDir is the hone directory of the openx platform
	HomeDir = os.Getenv("HOME") + "/.openx"
	// StablecoinCode is the code of the test stablecoin
	StablecoinCode = "STABLEUSD"
	// StablecoinPublicKey is the publickey of the stablecoin
	StablecoinPublicKey = ""
	// StablecoinSeed is the seed of the stablecoin
	StablecoinSeed = ""
	// StableCoinSeedFile denotes the file location of the stablecoin's seed file
	StableCoinSeedFile = HomeDir + "/stablecoinseed.hex"
	// StableCoinAddress denotes the address of the stablecoin
	StableCoinAddress = "GDJE64WOXDXLEK7RDURVYEJ5Y5XFHS6OQZCS3SHO4EEMTABEIJXF6SZ5"
	// StablecoinTrustLimit denotes the trust limit of the stablecoin
	StablecoinTrustLimit = float64(1000000000)
	// AnchorUSDCode is the code of Anchor's stablecoin
	AnchorUSDCode = "USD"
	// AnchorUSDAddress denotes the address of Anchor's stablecoin
	AnchorUSDAddress = "GCKFBEIYV2U22IO2BJ4KVJOIP7XPWQGQFKKWXR6DOSJBV7STMAQSMTGG"
	// AnchorUSDTrustLimit is the default trust limit for trusting Anchor
	AnchorUSDTrustLimit = float64(1000000)
)
