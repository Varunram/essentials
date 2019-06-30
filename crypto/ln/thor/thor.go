package main

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

type GetLnChannelInventoryReturn struct {
	Operator struct {
		Name     string `json:"name"`
		Slug     string `json:"slug"`
		Currency string `json:"currency"`
		Packages []struct {
			Value        string `json:"value"`
			EuroPrice    string `json:"eurprice"`
			SatoshiPrice string `json:"satoshiPrice"`
			USDPrice     string `json:"usdPrice"`
			UserPrice    string `json:"userPrice"`
		} `json:"packages"`
	} `json:"operator"`
}

type ThorOrderReturn struct {
	Id             string `json:"id"`
	Email          string `json:"email"`
	Expired        bool   `json:"expired"`
	Value          string `json:"value"`
	Product        string `json:"product"`
	Price          int64  `json:"price"`
	PartialPayment bool   `json:"partialPayment"`
	UserRef        string `json:"userRef"`
	Status         string `json:"status"`
	Payment        struct {
		Address          string `json:"address"`
		LightningInvoice string `json:"lightningInvoice"`
		SatoshiPrice     string `json:"satoshiPrice"`
		AltcoinCode      string `json:"altcoinCode"`
	} `json:"payment"`
	ThorInfo struct {
		Link  string `json:"link"`
		K1    string `json:"k1"`
		LnURL string `json:"lnurl"`
		Other string `json:"other"`
	} `json:"thorInfo"`
}

func GetLnChannelInventory() (GetLnChannelInventoryReturn, error) {
	var x GetLnChannelInventoryReturn
	baseURL := "https://api.bitrefill.com/v1/inventory/lightning-channel"
	data, err := GetRequest(baseURL)
	if err != nil {
		return x, errors.Wrap(err, "could not make get request to bitrefill api, quitting")
	}

	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, errors.Wrap(err, "could not unmarshal json struct, quitting")
	}

	if x.Operator.Name != "Lightning Channel" {
		return x, errors.New("name from get request does not match Lightning Channel, quitting")
	}

	return x, err
}

func GetTurboLnChannelInventory() (GetLnChannelInventoryReturn, error) {
	var x GetLnChannelInventoryReturn
	baseURL := "https://api.bitrefill.com/v1/inventory/lightning-channel"
	data, err := GetRequest(baseURL)
	if err != nil {
		return x, errors.Wrap(err, "could not make get request to bitrefill api, quitting")
	}

	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, errors.Wrap(err, "could not unmarshal json struct, quitting")
	}

	if x.Operator.Name != "Turbo Lightning Channell" {
		return x, errors.New("name from get request does not match Lightning Channel, quitting")
	}

	return x, err
}

func PostOrder(paymentMode string, sendEmail string, email string) (ThorOrderReturn, error) {
	var x ThorOrderReturn
	baseURL := "https://api.bitrefill.com/v1/order"
	data := url.Values{}

	operatorSlug := "lightning-channel"
	valuePackage := "2,000,000 sats capacity"
	email = "varunramg@bithyve.com"
	sendEmail = "true"

	var paymentMethod string

	switch paymentMode {
	case "ln":
		paymentMethod = "lightning"
	case "ln-ltc":
		paymentMethod = "lightning-ltc"
	case "btc":
		paymentMethod = "bitcoin"
	case "eth":
		paymentMethod = "ethereum"
	case "ltc":
		paymentMethod = "litecoin"
	case "dash":
		paymentMethod = "dash"
	case "doge":
		paymentMethod = "dogecoin"
	default:
		return x, errors.New("payment method not supported, please PR, thanks!")
	}

	refundAddress := ""
	webhookUrl := ""
	userRef := ""

	data.Set("operatorSlug", operatorSlug)
	data.Set("valuePackage", valuePackage)
	data.Set("email", email)
	data.Set("sendEmail", sendEmail)
	data.Set("paymentMethod", paymentMethod)
	data.Set("refund_address", refundAddress)
	data.Set("webhook_url", webhookUrl)
	data.Set("userRef", userRef) // this is an internal id used to reference payments within the system

	payload := strings.NewReader(data.Encode())

	// send post request to url
	returnData, err := PostRequest(baseURL, payload)
	if err != nil {
		return x, errors.New("could not post to URL, quitting")
	}

	err = json.Unmarshal(returnData, &x)
	if err != nil {
		return x, errors.New("could not unmarshal returned data, quitting")
	}

	return x, nil
}

func PostPurchase(paymentMode string, sendEmail string, email string) (ThorOrderReturn, error) {
	var x ThorOrderReturn
	baseURL := "https://api.bitrefill.com/v1/purchase"
	data := url.Values{}

	operatorSlug := "lightning-channel"
	valuePackage := "2,000,000 sats capacity"
	email = "varunramg@bithyve.com"
	sendEmail = "true"
	webhookUrl := ""
	userRef := ""

	data.Set("operatorSlug", operatorSlug)
	data.Set("valuePackage", valuePackage)
	data.Set("email", email)
	data.Set("sendEmail", sendEmail)
	data.Set("webhook_url", webhookUrl)
	data.Set("userRef", userRef) // this is an internal id used to reference payments within the system

	payload := strings.NewReader(data.Encode())

	// send post request to url
	returnData, err := PostRequest(baseURL, payload)
	if err != nil {
		return x, errors.New("could not post to URL, quitting")
	}

	err = json.Unmarshal(returnData, &x)
	if err != nil {
		return x, errors.New("could not unmarshal returned data, quitting")
	}

	return x, nil
}

func GetOrderId(orderId string) (ThorOrderReturn, error) {
	var x ThorOrderReturn

	baseURL := "https://api.bitrefill.com/v1/order/" + orderId
	data, err := GetRequest(baseURL)
	if err != nil {
		return x, errors.Wrap(err, "could not make get request to bitrefill api, quitting")
	}

	err = json.Unmarshal(data, &x)
	if err != nil {
		return x, errors.Wrap(err, "could not unmarshal json struct, quitting")
	}

	return x, err
}
