package bech32

import (
	"log"
	"github.com/btcsuite/btcd/btcec"
)

var Curve *btcec.KoblitzCurve = btcec.S256()

func GetNewBech32Address() (string, error) {
	localPriv, err := btcec.NewPrivateKey(Curve)
	if err != nil {
		return "", nil
	}

	pubkey := localPriv.PubKey().SerializeCompressed()
	log.Println("BECH #@ Pubkey: ", pubkey, len(pubkey))
	conv, err := ConvertBits(pubkey, 8, 5, true) // bit conversion stuff
	if err != nil {
		return "", nil
	}

	bech32Pubkey, err := Encode("bc", conv)
	if err != nil {
		return "", nil
	}

	return bech32Pubkey, nil
}

func StringToBech32(address []byte) (string, error) {
	conv, err := ConvertBits([]byte(address), 8, 5, true) // bit conversion stuff
	if err != nil {
		return "", nil
	}

	bech32Pubkey, err := Encode("bc", conv)
	if err != nil {
		return "", nil
	}

	return bech32Pubkey, nil
}
