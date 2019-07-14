package escrow

import (
	"log"

	stablecoin "github.com/Varunram/essentials/crypto/stablecoin"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	multisig "github.com/Varunram/essentials/crypto/xlm/multisig"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	utils "github.com/Varunram/essentials/utils"
	"github.com/pkg/errors"
)

// escrow implements an escrow based off Stellar

// InitEscrow creates a new keypair and stores it in a file
func InitEscrow(projIndex int, seedpwd string, recpPubkey string, mySeed string, platformSeed string) (string, error) {
	platformPubkey, err := wallet.ReturnPubkey(platformSeed)
	if err != nil {
		return "", errors.Wrap(err, "could not get pubkey from seed")
	}

	pubkey, err := initMultisigEscrow(recpPubkey, platformPubkey)
	if err != nil {
		return pubkey, errors.Wrap(err, "error while initializing multisig escrow, quitting!")
	}

	log.Println("successfully initialized multisig escrow")
	// define two seeds that are needed for signing transactions from the escrow
	seed1 := platformSeed
	seed2 := mySeed

	log.Println("stored escrow pubkey successfully")
	err = multisig.AuthImmutable2of2(pubkey, seed1, seed2)
	if err != nil {
		return pubkey, errors.Wrap(err, "could not set auth immutable on account, quitting!")
	}

	log.Println("set auth immutable on account successfully")
	multisig.TrustAssetTx(stablecoin.StablecoinCode, stablecoin.StablecoinPublicKey, "10000000000", pubkey, seed1, seed2)
	if err != nil {
		return pubkey, errors.Wrap(err, "could not trust stablecoin, quitting!")
	}

	return pubkey, nil
}

// TransferFundsToEscrow transfers a specific amount of currency to the escrow. Usually called by the platform or recipient
func TransferFundsToEscrow(amount float64, projIndex int, escrowPubkey string, platformSeed string) error {
	// we have the wallet pubkey, transfer funds to the escrow now
	aS, _ := utils.ToString(amount)
	_, txhash, err := assets.SendAsset(stablecoin.StablecoinCode, stablecoin.StablecoinPublicKey, escrowPubkey,
		aS, platformSeed, "escrow init")
	if err != nil {
		return errors.Wrap(err, "could not fund escrow, quitting!")
	}

	log.Println("tx hash for funding project escrow is: ", txhash)
	return nil
}

// InitMultisigEscrow initializes a multisig escrow
func initMultisigEscrow(pubkey1 string, pubkey2 string) (string, error) {
	return multisig.New2of2(pubkey1, pubkey2)
}

// SendFundsFromEscrow sends funds to a destination address from the project escrow
func SendFundsFromEscrow(escrowPubkey string, destination string, signer1 string, signer2 string, amount string, memo string) error {
	log.Println("ESCROW PUBKEY: ", escrowPubkey, "destination: ", destination, "signer1: ", signer1, "amount: ", amount, "memo: ", memo)
	return multisig.Tx2of2(escrowPubkey, destination, signer1, signer2, amount, memo)
}
