package bech32

import (
	"github.com/btcsuite/btcd/btcec"
	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
	// "github.com/btcsuite/btcutil/base58"
	// "log"
)

var Curve *btcec.KoblitzCurve = btcec.S256()

func GetNewBech32Address() (string, error) {
	localPriv, err := btcec.NewPrivateKey(Curve)
	if err != nil {
		return "", nil
	}

	pubkey := localPriv.PubKey().SerializeCompressed()
	hash160 := btcutils.Hash160(pubkey) // p2wpkh
	var program []int
	for _, vals := range hash160 {
		program = append(program, int(vals))
	}
	return SegwitAddrEncode("bc", 0, program)
}
