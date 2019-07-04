package main

import (
	"github.com/pkg/errors"
	"log"

	bech32 "github.com/Varunram/essentials/crypto/btc/bech32"
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
		log.Println("list of supported commands: ")
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

			minInt, err := utils.StoICheck(minShares)
			if err != nil {
				return errors.Wrap(err, "minshares not int")
			}

			maxInt, err := utils.StoICheck(maxShares)
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
	default:
		return errors.New("command not recognized")
	}

	return nil
}
