/*
	Copyright 2013-present wemeetagain https://github.com/wemeetagain/go-hdwallet
	Copyright 2019-present Varunram Ganesh
*/
package hdwallet

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"github.com/pkg/errors"
	"math/big"

	utils "github.com/Varunram/essentials/utils"
	"github.com/btcsuite/btcutil/base58"
)

// declare version bytes and HMAC key
var (
	MnPubkeyVByte    = []byte{4, 136, 178, 30}  // mainnet magic bytes
	MnPrivkeyVByte   = []byte{4, 136, 173, 228} // mainnet magic bytes
	TestPubkeyVByte  = []byte{4, 53, 135, 207}  // testnet magic bytes
	TestPrivkeyVByte = []byte{4, 53, 131, 148}  // testnet magic bytes
	HmacKey          = []byte("Bitcoin seed")   // HMAC Key
)

// HDWallet defines the components of a hierarchical deterministic wallet
type HDWallet struct {
	VersionBytes []byte //4 bytes
	Depth        uint16 //1 byte
	Fingerprint  []byte //4 bytes
	ChildNumber  []byte //4 bytes
	Chaincode    []byte //32 bytes
	Key          []byte //33 bytes
}

// Child returns the ith child of wallet w. Values of i >= 2^31
// signify private key derivation. Attempting private key derivation
// with a public key will throw an error.
func (w *HDWallet) Child(i uint32) (*HDWallet, error) {
	var fingerprint, childNumber, newkey []byte
	switch {
	case bytes.Compare(w.VersionBytes, MnPrivkeyVByte) == 0, bytes.Compare(w.VersionBytes, TestPrivkeyVByte) == 0:
		mac := hmac.New(sha512.New, w.Chaincode)
		if i >= uint32(0x80000000) {
			mac.Write(append(w.Key, utils.Uint32tB(i)...))
		} else {
			pub := PrivToPub(w.Key)
			mac.Write(append(pub, utils.Uint32tB(i)...))
		}
		childNumber = mac.Sum(nil)
		iL := new(big.Int).SetBytes(childNumber[:32])
		if iL.Cmp(curve.N) >= 0 || iL.Sign() == 0 {
			return &HDWallet{}, errors.New("Invalid Child")
		}
		newkey = AddPrivKeys(childNumber[:32], w.Key)
		fingerprint = Hash160(PrivToPub(w.Key))[:4]

	case bytes.Compare(w.VersionBytes, MnPubkeyVByte) == 0, bytes.Compare(w.VersionBytes, TestPubkeyVByte) == 0:
		mac := hmac.New(sha512.New, w.Chaincode)
		if i >= uint32(0x80000000) {
			return &HDWallet{}, errors.New("Can't do Private derivation on Public key!")
		}
		mac.Write(append(w.Key, utils.Uint32tB(i)...))
		childNumber = mac.Sum(nil)
		iL := new(big.Int).SetBytes(childNumber[:32])
		if iL.Cmp(curve.N) >= 0 || iL.Sign() == 0 {
			return &HDWallet{}, errors.New("Invalid Child")
		}
		newkey = AddPubKeys(PrivToPub(childNumber[:32]), w.Key)
		fingerprint = Hash160(w.Key)[:4]
	}
	return &HDWallet{w.VersionBytes, w.Depth + 1, fingerprint, utils.Uint32tB(i), childNumber[32:], newkey}, nil
}

// Serialize returns the serialized form of the wallet.
func (w *HDWallet) Serialize() []byte {
	depth := utils.Uint16tB(uint16(w.Depth % 256))
	//bindata = VersionBytes||depth||fingerprint||i||chaincode||key
	bindata := make([]byte, 78)
	copy(bindata, w.VersionBytes)
	copy(bindata[4:], depth)
	copy(bindata[5:], w.Fingerprint)
	copy(bindata[9:], w.ChildNumber)
	copy(bindata[13:], w.Chaincode)
	copy(bindata[45:], w.Key)
	chksum := DoubleSha256(bindata)[:4]
	return append(bindata, chksum...)
}

// String returns the base58-encoded string form of the wallet.
func (w *HDWallet) String() string {
	return base58.Encode(w.Serialize())
}

// Pub returns a new wallet which is the public key version of w.
// If w is a public key, Pub returns a copy of w
func (w *HDWallet) Pub() *HDWallet {
	if bytes.Compare(w.VersionBytes, MnPubkeyVByte) == 0 {
		return &HDWallet{w.VersionBytes, w.Depth, w.Fingerprint, w.ChildNumber, w.Chaincode, w.Key}
	} else {
		return &HDWallet{MnPubkeyVByte, w.Depth, w.Fingerprint, w.ChildNumber, w.Chaincode, PrivToPub(w.Key)}
	}
}

// Address returns bitcoin address represented by wallet w.
func (w *HDWallet) Address() string {
	x, y := Expand(w.Key)
	paddedKey := append([]byte{4}, append(x.Bytes(), y.Bytes()...)...) // 04
	var prefix []byte
	if bytes.Compare(w.VersionBytes, TestPubkeyVByte) == 0 || bytes.Compare(w.VersionBytes, TestPrivkeyVByte) == 0 {
		prefix = []byte{111} // 6F for testnet
	} else {
		prefix = []byte{0} // 00 for mainnet
	}
	address := append(prefix, Hash160(paddedKey)...)
	chksum := DoubleSha256(address)
	return base58.Encode(append(address, chksum[:4]...))
}

// GenSeed returns a random seed with a length measured in bytes.
// The length must be at least 128.
func GenSeed(length int) ([]byte, error) {
	b := make([]byte, length)
	if length < 128 {
		return b, errors.New("length must be at least 128 bits")
	}
	_, err := rand.Read(b)
	return b, err
}

// MasterKey returns a new wallet given a random seed.
func MasterKey(seed []byte) *HDWallet {
	mac := hmac.New(sha512.New, HmacKey)
	mac.Write(seed)
	I := mac.Sum(nil)
	secret := I[:len(I)/2]
	chain_code := I[len(I)/2:]
	depth := 0
	i := make([]byte, 4)
	fingerprint := make([]byte, 4)
	zero := make([]byte, 1)
	return &HDWallet{MnPrivkeyVByte, uint16(depth), fingerprint, i, chain_code, append(zero, secret...)}
}
