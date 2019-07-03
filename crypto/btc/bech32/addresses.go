package bech32

import (
	"github.com/Varunram/essentials/crypto/btc/base58"
	btcutils "github.com/Varunram/essentials/crypto/btc/utils"
	"github.com/btcsuite/btcd/btcec"
	"log"

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

	log.Println("LENDATA: ", len(data))
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
	log.Println("ADDR len: ", len(address))
	chksum := btcutils.DoubleSha256(address)
	return base58.Encode(append(address, chksum[:4]...)), nil
}

func Base58ToBech32Address(addr string) (string, error){
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
		hrp="tb"
	default:
		return "", errors.New("prefix not recognized")
	}

	return SegwitAddrEncode(hrp, SegwitVersion, ByteArrToInt(byteString))

}
