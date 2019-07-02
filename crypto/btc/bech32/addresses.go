package bech32

import (
	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
	"github.com/btcsuite/btcd/btcec"
	// "github.com/btcsuite/btcutil/base58"
	// "log"
)

var Curve *btcec.KoblitzCurve = btcec.S256()

func GetNewp2wpkh() (string, error) {
	segwitVersion := 0
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
	return SegwitAddrEncode("bc", segwitVersion, program)
}
