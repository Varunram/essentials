package main

import (
	"github.com/pkg/errors"
	"log"

	bech32 "github.com/Varunram/essentials/crypto/btc/bech32"
	bip39 "github.com/Varunram/essentials/crypto/btc/bip39"
	hdwallet "github.com/Varunram/essentials/crypto/btc/hdwallet"
	paynym "github.com/Varunram/essentials/crypto/btc/paynym"
	sss "github.com/Varunram/essentials/sss"
	utils "github.com/Varunram/essentials/utils"
)

// parse all the inputs that the user might provide

func ParseInput(cmd []string) error {
	if len(cmd) == 0 {
		log.Println("type help to know about all the commands")
		return errors.New("no command provided")
	}

	command := cmd[0]

	switch command {
	case "help":
		ColorOutput("list of supported commands: sss, new, paynym, combine", CyanColor)
		return nil
	case "sss":
		if len(cmd) < 2 {
			return errors.New("USAGE: sss <new/combine>")
		}

		subcommand := cmd[1]

		switch subcommand {
		case "new":
			cmd = cmd[2:] // slice off sss new
			if len(cmd) != 3 {
				return errors.New("USAGE: sss new <secret> <min_shares> <total_shares>")
			}
			secret := cmd[0]
			minShares := cmd[1]
			maxShares := cmd[2]

			minInt, err := utils.ToInt(minShares)
			if err != nil {
				return errors.Wrap(err, "minshares not int")
			}

			maxInt, err := utils.ToInt(maxShares)
			if err != nil {
				return errors.Wrap(err, "maxshares not int")
			}

			shares, err := sss.Create(minInt, maxInt, secret)
			if err != nil {
				return errors.Wrap(err, "could not create shares")
			}

			for _, share := range shares {
				ColorOutput("CREATED SHARE: "+share, GreenColor)
			}

		case "combine":
			var shares = cmd[2:] // slice off sss combine
			combinedSecret, err := sss.Combine(shares)
			if err != nil {
				return errors.Wrap(err, "could not combine shares")
			}
			ColorOutput("RETRIEVED SECRET: "+combinedSecret, GreenColor)

		default:
			return errors.New("USAGE: sss <new/combine>")
		}
		// send of sss
	case "new":
		if len(cmd) != 2 {
			return errors.New("USAGE: new <p2pkh / p2wpkh>")
		}

		cmd = cmd[1:]
		switch cmd[0] {
		case "p2wpkh":
			address, err := bech32.GetNewp2wpkh()
			if err != nil {
				return errors.Wrap(err, "could not generate p2wpkh address")
			}
			ColorOutput("ADDRESS: "+address, GreenColor)

		case "p2pkh":
			address, err := bech32.GetNewp2wpkh()
			if err != nil {
				return errors.Wrap(err, "could not generate p2wpkh address")
			}

			base58Addr, err := bech32.Bech32ToBase58Addr(address[0:2], address)
			if err != nil {
				return errors.Wrap(err, "could not convert bech32 to base58")
			}

			ColorOutput("ADDRESS: "+base58Addr, GreenColor)
		}
		// end of new

	case "paynym":
		// first we need a hd wallet
		xpriv, _, err := paynym.SetupWallets()
		if err != nil {
			return errors.Wrap(err, "failed to setup a hd wallet")
		}

		bobPubkey := make([]byte, 32)
		outpoint := make([]byte, 36)

		paynymCode, err := paynym.GenPaynym(xpriv.Key[1:], bobPubkey, xpriv.Chaincode, outpoint)
		if err != nil {
			return errors.Wrap(err, "could not generate paynym, quitting")
		}

		ColorOutput("PAYNYM CODE: "+paynymCode, GreenColor)

		// end of paynym

	case "mnemonic":
		// generate a mnemonic here
		if len(cmd) != 3 {
			return errors.New("USAGE: mnemonic <12/15/18/21/24> passphrase")
		}

		wordSize, err := utils.ToInt(cmd[1])
		if err != nil {
			return errors.Wrap(err, "could not convert input into string, returning")
		}
		wordSizeMap := make(map[int]int, 5)

		wordSizeMap[12] = 128
		wordSizeMap[15] = 160
		wordSizeMap[18] = 192
		wordSizeMap[21] = 224
		wordSizeMap[24] = 256

		entropy, _ := bip39.NewEntropy(wordSizeMap[wordSize])
		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			log.Fatal(err)
		}
		ColorOutput("MNEMONIC: "+mnemonic, GreenColor)

		passphrase := cmd[2]
		seed, err := bip39.NewSeed(mnemonic, passphrase)
		if err != nil {
			return errors.Wrap(err, "could not get new seed, quitting")
		}

		masterKey := hdwallet.MasterKey(seed)
		publicKey := masterKey.Pub()

		ColorOutput("xpub: "+publicKey.String(), GreenColor)
		// end of mnemonic

	case "recover":
		// recover seed from mnemonic
		if len(cmd) < 2 {
			return errors.New("USAGE: recover passphrase <mnemonic>")
		}
		var mnemonic string
		for _, strings := range cmd[2:] {
			mnemonic = mnemonic + " " + strings
		}

		passphrase := cmd[1]
		ColorOutput("PASSPHRASE: "+passphrase, CyanColor)
		mnemonic = mnemonic[1:] // get the first " " out
		ColorOutput("MNEMONIC: "+mnemonic, CyanColor)
		seed, err := bip39.NewSeed(mnemonic, passphrase)
		if err != nil {
			return errors.Wrap(err, "could not get new seed, quitting")
		}

		masterKey := hdwallet.MasterKey(seed)
		publicKey := masterKey.Pub()

		ColorOutput("xpub: "+publicKey.String(), GreenColor)

	default:
		return errors.New("command not recognized")
	}

	return nil
}
