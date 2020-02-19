package xlm

import (
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// GetXLM makes an API call to the stellar friendbot, which gives 10000 testnet XLM
func GetXLM(PublicKey string) error {
	if Mainnet {
		return errors.New("no friendbot on mainnet, quitting")
	}
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + PublicKey)
	if err != nil || resp.Status != "200 OK" {
		return errors.Wrap(err, "API Request did not succeed")
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}
