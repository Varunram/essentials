package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"log"
	"math/big"
	"bytes"

	"github.com/Varunram/essentials/crypto/btc/base58"
	"github.com/Varunram/essentials/crypto/btc/hdwallet"
	"github.com/btcsuite/btcd/btcec"
)

type PayNym struct {
	Version    byte
	BitMessage byte
	Sign       byte
	XValue     [32]byte
	Chaincode  [32]byte
	Reserved   [13]byte
}

func (p *PayNym) Bytes() []byte {
	var x bytes.Buffer
	_ = x.WriteByte(p.Version)
	_ = x.WriteByte(p.BitMessage)
	_ = x.WriteByte(p.Sign)
	_, _ = x.Write(p.XValue[:])
	_, _ = x.Write(p.Chaincode[:])
	_, _ = x.Write(p.Reserved[:])
	return x.Bytes()
}

// while serializing to base58, version byte should be 0x47 and the payload must
// be the the binary serialization

// m / purpose' / coin_type' / identity'
// m / 47 / 0 / 0 for the first wallet, m/47/0/1 for the second wallet and so on

var Purpose = uint32(0x8000002F)
var CoinType = uint32(0x80000000)
var Identity = uint32(0x80000000)

var AddressVersion = 0x01  // address
var MultisigVersion = 0x02 // bloom-multisig

var Curve *btcec.KoblitzCurve = btcec.S256()

func setupWallets() (*hdwallet.HDWallet, *hdwallet.HDWallet, error) {

	aliceSeed, err := hdwallet.GenSeed(256)
	if err != nil {
		log.Fatal(err)
	}

	aliceMasterprv := hdwallet.MasterKey(aliceSeed)

	// Generate new child key based on private or public key
	alicePriv1, err := aliceMasterprv.Child(Purpose)
	if err != nil {
		log.Fatal(err)
	}
	alicePriv2, err := alicePriv1.Child(CoinType)
	if err != nil {
		log.Fatal(err)
	}
	alicePriv3, err := alicePriv2.Child(Identity)
	if err != nil {
		log.Fatal(err)
	}

	return alicePriv3, alicePriv3.Pub(), nil
}

func main() {
	// first let setup alice and bob's wallets
	// log.Println("MASTER PUBKEY: ", masterpub)
	alicePriv, alicePub, err := setupWallets()
	if err != nil {
		log.Fatal(err)
	}

	bobPriv, bobPub, err := setupWallets()
	if err != nil {
		log.Fatal(err)
	}

	bobNotifAddress := bobPub.Address()

	aliceAddress := alicePub.Address()
	alicePrivKey := alicePriv.Key[1:]

	bobPubkey := bobPub.Key

	aliceI := new(big.Int)
	bobI := new(big.Int)

	log.Println(alicePrivKey)
	aliceI.SetBytes(alicePrivKey)
	bobI.SetBytes(alicePrivKey)

	secretPointX, secretPointY := Curve.ScalarMult(aliceI, bobI, make([]byte, 33))

	var outpoint [36]byte

	hmacHash := hmac.New(sha512.New, []byte(secretPointX.String()))
	hmacHash.Write(outpoint[:])

	log.Println(bobPriv, bobNotifAddress, aliceAddress, bobPubkey, secretPointY)

	log.Println("size of blinding factor: ", hmacHash.Size())
	var paymentCode PayNym

	paymentCode.Version = byte(AddressVersion)
	paymentCode.BitMessage = 0x00
	paymentCode.Sign = 0x02 // get the actual sign

	for i, val := range hmacHash.Sum(nil)[32:63] {
		paymentCode.Chaincode[i] = byte(0x00000000) ^ val // xor b on element of random
	}

	for i, val := range hmacHash.Sum(nil)[0:31] {
		paymentCode.XValue[i] = byte(0x00000000) ^ val // xor b on element of random
	}

	var x [13]byte
	paymentCode.Reserved = x

	log.Println("LEN OF BLAH: ", paymentCode)
	paymentCodeBytes := paymentCode.Bytes()
	versionByteArray := []byte{0x47}

	log.Println("LEN OF paymentCodeBytes: ", len(paymentCodeBytes))
	log.Println(base58.Encode(append(versionByteArray, paymentCodeBytes...)))
}
