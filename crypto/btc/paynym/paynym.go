package paynym

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"github.com/pkg/errors"
	"log"
	"math/big"

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

func (p *PayNym) Bytes() ([]byte, error) {
	var x bytes.Buffer
	if !(p.Version == 0x00 || p.Version == 0x01) {
		return nil, errors.New("version should be zero or one")
	}
	_ = x.WriteByte(p.Version)
	_ = x.WriteByte(p.BitMessage)
	if !(p.Sign == 0x02 || p.Sign == 0x03) {
		return nil, errors.New("sign should be 2/3")
	}
	_ = x.WriteByte(p.Sign)
	if !(len(p.XValue) == 32) {
		return nil, errors.New("length of x should be 32 bytes")
	}
	_, _ = x.Write(p.XValue[:])
	if !(len(p.Chaincode) == 32) {
		return nil, errors.New("length of Chaincode should be 32 bytes")
	}
	_, _ = x.Write(p.Chaincode[:])
	if !(len(p.Reserved) == 13) {
		return nil, errors.New("reserved byte  length must be 13")
	}
	_, _ = x.Write(p.Reserved[:])
	return x.Bytes(), nil
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

func SetupWallets() (*hdwallet.HDWallet, *hdwallet.HDWallet, error) {

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

func GenPaynym(alicePrivKey, bobPubkey, chaincode []byte, outpoint []byte) (string, error) {
	if len(alicePrivKey) != 32 {
		return "", errors.New("length of private key not 32 bytes")
	}
	if len(bobPubkey) != 32 {
		return "", errors.New("length of provided public key not 32 bytes")
	}
	if len(chaincode) != 32 {
		return "", errors.New("length of chaincode not 32 bytes")
	}
	if len(outpoint) != 36 {
		return "", errors.New("outpoint length not 36 (32 txhash + 4 vout)")
	}

	// convert the byte arrays to big ints
	aliceI := new(big.Int)
	bobI := new(big.Int)

	aliceI.SetBytes(alicePrivKey)
	bobI.SetBytes(bobPubkey)

	// get the secret point coords
	secretPointX, secretPointY := Curve.ScalarMult(aliceI, bobI, make([]byte, 33))

	// generate the hmac secret s
	hmacHash := hmac.New(sha512.New, []byte(secretPointX.String()))
	hmacHash.Write(outpoint[:])

	// construct paynym
	var paymentCode PayNym

	paymentCode.Version = byte(AddressVersion)
	paymentCode.BitMessage = 0x00

	if secretPointY.Bit(0) == 1 {
		paymentCode.Sign = 0x03
	} else {
		paymentCode.Sign = 0x02
	}

	for i, val := range hmacHash.Sum(nil)[32:63] {
		paymentCode.Chaincode[i] = chaincode[i] ^ val // xor b on element of random
	}

	for i, val := range hmacHash.Sum(nil)[0:31] {
		paymentCode.XValue[i] = chaincode[i] ^ val // xor b on element of random
	}

	var x [13]byte
	paymentCode.Reserved = x

	paymentCodeBytes, err := paymentCode.Bytes()
	if err != nil {
		return "", err
	}

	versionByte := byte(0x47)
	return base58.CheckEncode(paymentCodeBytes, versionByte), nil
}
