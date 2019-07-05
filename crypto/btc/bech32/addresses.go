package bech32

import (
	"github.com/Varunram/essentials/crypto/btc/base58"
	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
	"github.com/btcsuite/btcd/btcec"
	"github.com/skip2/go-qrcode"

	"github.com/pkg/errors"
)

var Curve *btcec.KoblitzCurve = btcec.S256()
var SegwitVersion = 0

func ByteArrToInt(input []byte) []int {
	var output []int
	for _, vals := range input {
		output = append(output, int(vals))
	}
	return output
}

func GetNewp2wpkh() (string, error) {
	localPriv, err := btcec.NewPrivateKey(Curve)
	if err != nil {
		return "", nil
	}

	pubkey := localPriv.PubKey().SerializeCompressed()
	hash160 := btcutils.Hash160(pubkey)
	var program []int
	for _, vals := range hash160 {
		program = append(program, int(vals))
	}
	return SegwitAddrEncode("bc", SegwitVersion, program)
}

func Bech32ToBase58Addr(hrp, addr string) (string, error) {
	_, data, err := SegwitAddrDecode(hrp, addr)
	if err != nil {
		return "", err
	}

	var prefix []byte
	switch hrp {
	case "bc":
		// mainnet, so prefix with 0
		prefix = []byte{0} // 00 for mainnet
	case "tb":
		prefix = []byte{111} // 6F for testnet
	}

	var arr []byte
	for _, vals := range data {
		arr = append(arr, byte(vals))
	}
	address := append(prefix, arr...)
	chksum := btcutils.DoubleSha256(address)
	return base58.Encode(append(address, chksum[:4]...)), nil
}

func Base58ToBech32Address(addr string) (string, error) {
	// decode from base58 to bytetring
	var hrp string
	byteString := base58.Decode(addr)

	prefix := int(byteString[0])
	byteString = byteString[1 : len(byteString)-4] // remove the prefix checksum

	if len(byteString) != 20 {
		return "", errors.New("length of address doesn't match, quitting")
	}

	switch prefix {
	case 0:
		// mainnet
		hrp = "bc"
	case 111:
		hrp = "tb"
	default:
		return "", errors.New("prefix not recognized")
	}

	return SegwitAddrEncode(hrp, SegwitVersion, ByteArrToInt(byteString))
}

func PrivKeyToWIF(network string, compressed bool, privkey []byte) (string, error) {
	if len(privkey) != 32 {
		return "", errors.New("length of private key not 32")
	}

	prefixByte := make([]byte, 1)
	if network == "mainnet" {
		prefixByte = []byte{0x80}
	} else if network == "testnet" {
		prefixByte = []byte{0xef}
	} else {
		return "", errors.New("unknown network")
	}

	var exKey []byte
	if compressed {
		exKey = append(prefixByte, append(privkey, byte(0x01))...)
	} else {
		exKey = append(prefixByte, privkey...)
	}

	doubleSha := btcutils.DoubleSha256(exKey)

	checksum := doubleSha[0:4] // first 4 bytes are the checksum
	exKey = append(exKey, checksum...)
	return base58.Encode(exKey), nil
}

func WIFToPrivateKey(wif string, compressed bool) ([]byte, error) {
	decodedString := base58.Decode(wif)
	if CheckCheckSum(wif) != nil {
		return decodedString, errors.New("checksum doesn't match, quitting")
	}

	decodedString = decodedString[1 : len(decodedString)-4] // drop network byte and checksum
	if compressed {
		decodedString = decodedString[0 : len(decodedString)-1]
	}

	if len(decodedString) != 32 {
		return decodedString, errors.New("private key length not 32")
	}

	return decodedString, nil
}

func CheckCheckSum(wif string) error {
	decodedString := base58.Decode(wif)
	provCheckSum := decodedString[len(decodedString)-4 : len(decodedString)]
	decodedString = decodedString[0 : len(decodedString)-4]

	doubleSha := btcutils.DoubleSha256(decodedString)
	for i, val := range doubleSha[0:4] {
		if val != provCheckSum[i] {
			return errors.New("checksums don't match, quitting")
		}
	}

	return nil
}

func ExportToQrCode(secret string, path string) error {
	err := qrcode.WriteFile(secret, qrcode.High, 1024, path)
	if err != nil {
		return errors.Wrap(err, "failed to generte qr code, quitting")
	}
	return nil
}
