package stablecoin

import (
	"log"
	"sync"

	"github.com/pkg/errors"

	xlm "github.com/Varunram/essentials/xlm"
	assets "github.com/Varunram/essentials/xlm/assets"
	// utils "github.com/Varunram/essentials/utils"
)

// Exchange exchanges xlm for STABLEUSD
func Exchange(recipientPK string, recipientSeed string, convAmount float64) error {

	if Mainnet {
		return errors.New("Exchange in mainent needs to be done through dex")
	}

	if !xlm.AccountExists(recipientPK) {
		return errors.New("account does not exist, quitting")
	}

	var balance, trustLimit float64

	var wg sync.WaitGroup
	err1 := make(chan error, 0)
	err2 := make(chan error, 0)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		balance = xlm.GetNativeBalance(recipientPK)
		if balance < convAmount {
			err1 <- errors.New("balance is less than amount requested")
		}
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		trustLimit = xlm.GetAssetTrustLimit(recipientPK, StablecoinCode)
		if trustLimit < convAmount && trustLimit != 0 {
			err2 <- errors.New("trust limit doesn't warrant investment")
		}
	}(&wg)

	wg.Wait()

	select {
	case err := <-err1:
		log.Println(err)
		return err
	case err := <-err2:
		log.Println(err)
		return err
	default:
		break
	}

	_, err := assets.TrustAsset(StablecoinCode, StablecoinPublicKey, StablecoinTrustLimit, recipientSeed)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "couldn't trust asset")
	}
	log.Println("stableUSD trustline created")

	_, _, err = xlm.SendXLM(StablecoinPublicKey, convAmount, recipientSeed, "Exchange XLM for stablecoin")
	if err != nil {
		log.Println("error while sending XLM", StablecoinPublicKey, convAmount, recipientSeed, err)
		return errors.Wrap(err, "couldn't send xlm")
	}
	log.Println("sent xlm / waiting to receive stableUSD")

	return nil
}

/*
// OfferExchange offers to exchange user's xlm balance for stableusd if the user does not have enough
// stableUSD to complete the payment
func OfferExchange(publicKey string, seed string, invAmount float64) error {

	if Mainnet {
		return errors.New("Exchange offers in mainnet need to be done through dex")
	}

	balance := xlm.GetAssetBalance(publicKey, StablecoinCode)
	if balance < invAmount {
		log.Println("Offering xlm to stableusd exchange to investor")
		// user's stablecoin balance is less than the amount he wishes to invest, get stablecoin
		// equal to the amount he wishes to exchange
		diff := invAmount - balance + 10 // the extra 1 is to cover for fees
		// checking whether the user has enough xlm balance to cover for the exchange is done by Exchange()
		xlmBalance := xlm.GetNativeBalance(publicKey)
		totalUSD := tickers.ExchangeXLMforUSD(xlmBalance) // amount in stablecoin that the user would receive for diff
		if totalUSD < diff {
			return errors.New("User does not have enough funds to complete this transaction")
		}

		// now we need to exchange XLM equal to diff in stablecoin
		exchangeRate := tickers.ExchangeXLMforUSD(1)
		// 1 xlm can fetch exchangeRate USD, how much xlm does diff USD need?
		amountToExchange := diff / exchangeRate
		log.Println(diff, exchangeRate, amountToExchange)
		err := Exchange(publicKey, seed, amountToExchange)
		if err != nil {
			return errors.Wrap(err, "Unable to exchange XLM for USD and automate payment. Please get more STABLEUSD to fulfil the payment")
		}
		time.Sleep(10 * time.Second) // 5 seconds for issuing stalbeusd to the person who's requested for it
	} else {
		log.Println("User has sufficient stablecoin balance, not exchanging xlm for usd")
	}

	return nil
}
*/
