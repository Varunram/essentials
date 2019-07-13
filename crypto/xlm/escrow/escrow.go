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

// escrow defines the escrow of asset for a specific project. We should generate a
// new seed and public key pair for each project that is at stage 3, so this would be
// automated at that stage. Once an investor has finished investing in the project,
// we need to send the recipient DebtAssets and then set all weights to zero in order
// to lock the account and prevent any further transactions from being authorized.
// One can stil send fund to the frozen account but the account can not use them
// this serves our purpose since we only want receipt of debt assets and want to freeze
// issuance so that anybody who hacks us can not print more tokens.

// In financial terms, an escrow is a special purpose vehicle (kind of cool that we have SPV in finance)

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

// InitMultisigEscrow initializes a multisig escrow with one signer as the recipient and the other as the platform
func initMultisigEscrow(pubkey1 string, platformPubkey string) (string, error) {
	// recpPubkey is the public key of the recipient
	// the seed of the escrow is needed to init the first tx that will change options
	// we now have the two public keys that are needed to authorize this transaction. Construct a 2of2 multisig
	return multisig.New2of2(pubkey1, platformPubkey)
}

// SendFundsFromEscrow sends funds to a destination address from the project escrow
func SendFundsFromEscrow(escrowPubkey string, destination string, signer1 string, signer2 string, amount string, memo string) error {
	log.Println("ESCROW PUBKEY: ", escrowPubkey, "destination: ", destination, "signer1: ", signer1, "amount: ", amount, "memo: ", memo)
	return multisig.Tx2of2(escrowPubkey, destination, signer1, signer2, amount, memo)
}
