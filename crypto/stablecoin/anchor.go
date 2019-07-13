package stablecoin

import (
	"github.com/pkg/errors"
	"log"
	// "time"

	tickers "github.com/Varunram/essentials/crypto/exchangetickers"
	xlm "github.com/Varunram/essentials/crypto/xlm"
	assets "github.com/Varunram/essentials/crypto/xlm/assets"
	utils "github.com/Varunram/essentials/utils"
)

// GetAnchorUSD gets anchorUSD from Anchor
func GetAnchorUSD(recpSeed string, amountUSDs string) (string, error) {
	txhash, err := assets.TrustAsset(AnchorUSDCode, AnchorUSDAddress, AnchorUSDTrustLimit, recpSeed)
	// txhash, err := assets.TrustAsset(Code, StableCoinAddress, StablecoinTrustLimit, recpSeed)
	if err != nil {
		return txhash, errors.Wrap(err, "couldn't trust anchorUSD")
	}
	log.Println("tx hash for trusting stableUSD: ", txhash)
	// now send coins across and see if our tracker detects it
	// the given amount is in USD, we need to convert it into XLM since we're sending XLM
	amountUSD, err := utils.ToFloat(amountUSDs)
	if err != nil {
		return txhash, err
	}

	exchangeRate, err := tickers.XLMUSD()
	if err != nil {
		return txhash, errors.Wrap(err, "error in fetching price from oracle")
	}
	amountXLM := exchangeRate * amountUSD

	log.Println("Exchanging: ", amountXLM, " XLM for anchorUSD")
	amountXLMS, _ := utils.ToString(amountXLM)
	_, txhash, err = xlm.SendXLM(AnchorUSDAddress, amountXLMS, recpSeed, "Exchange XLM for anchorUSD")
	if err != nil {
		return txhash, errors.Wrap(err, "couldn't send xlm")
	}
	log.Println("tx hash for sent xlm: ", txhash)
	return txhash, nil
}
