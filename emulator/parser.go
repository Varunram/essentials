package main

import (
	"github.com/pkg/errors"
	"log"

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
				return errors.Wrap(err, "could not combine shares, quitting")
			}
			ColorOutput("RETRIEVED SECRET: "+combinedSecret, GreenColor)

		default:
			return errors.New("USAGE: sss <new/combine>")
		}

	default:
		return errors.New("command not recognized, quitting")
	}

	return nil
}
