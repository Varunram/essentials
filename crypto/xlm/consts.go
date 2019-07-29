package xlm

import (
	"net/http"

	horizon "github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/network"
)

// constants that are imported by other packages

// TestNetClient defines the horizon client to connect to
var TestNetClient = &horizon.Client{
	HorizonURL: "https://horizon-testnet.stellar.org/",
	HTTP:       http.DefaultClient,
}

// Passphrase defines the stellar network passphrase
var Passphrase = network.TestNetworkPassphrase

// RefillAmount defines the default stellar refill amount
var RefillAmount float64

func SetConsts(amount float64) {
	RefillAmount = amount
}
