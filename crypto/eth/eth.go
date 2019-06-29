package essentials

import (
	"crypto/ecdsa"
	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common/hexutil"
	crypto "github.com/ethereum/go-ethereum/crypto"
)

type EthereumWallet struct {
	Address    string
	PublicKey  string
	PrivateKey string
}

func GenEthWallet() (EthereumWallet, error) {

	var x EthereumWallet
	ecdsaPrivkey, err := crypto.GenerateKey()
	if err != nil {
		return x, errors.Wrap(err, "could not generate an ethereum keypair, quitting!")
	}

	privateKeyBytes := crypto.FromECDSA(ecdsaPrivkey)
	x.PrivateKey = hexutil.Encode(privateKeyBytes)[2:]
	x.Address = crypto.PubkeyToAddress(ecdsaPrivkey.PublicKey).Hex()

	publicKeyECDSA, ok := ecdsaPrivkey.Public().(*ecdsa.PublicKey)
	if !ok {
		return x, errors.Wrap(err, "error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	x.PublicKey = hexutil.Encode(publicKeyBytes)[4:] // an ethereum address is 65 bytes long and hte first byte is 0x04 for DER encoding, so we omit that

	if crypto.PubkeyToAddress(*publicKeyECDSA).Hex() != x.Address {
		return x, errors.Wrap(err, "addresses don't match, quitting!")
	}

	return x, nil
}
