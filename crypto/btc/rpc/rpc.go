package main

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"

	utils "github.com/Varunram/essentials/utils"
)

type RPCReq struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

var BitcoindURL = "http://localhost:18443/" // for regtest

func PostReq(username string, password string, payload RPCReq) ([]byte, error) {
	var req *http.Request
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal json, quitting")
	}

	req, err = http.NewRequest("POST", BitcoindURL, bytes.NewBuffer(payloadJson))
	if err != nil {
		return nil, errors.Wrap(err, "did not POST to bitcoind")
	}

	req.SetBasicAuth(username, password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "did not make http request to bitcoind")
	}

	defer res.Body.Close()
	x, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "did not read from ioutil")
	}

	return x, nil
}

func GetBestBlockHash(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getbestblockhash"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetBlock(username string, password string, blockhash string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblock"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{blockhash}

	return PostReq(username, password, payload)
}

func GetBlockCount(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockcount"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetBlockHash(username string, password string, blockNumber uint32) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockhash"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []uint32{blockNumber}

	return PostReq(username, password, payload)
}

func GetBlockHeader(username string, password string, blockhash string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockheader"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{blockhash}

	return PostReq(username, password, payload)
}

func GetBlockStats(username string, password string, hashOrHeight string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockstats"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	height, err := utils.StoICheck(hashOrHeight)
	if err != nil {
		payload.Params = []string{hashOrHeight}
	} else {
		payload.Params = []int{height}
	}

	return PostReq(username, password, payload)
}

func GetChainTips(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getchaintips"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetChainTxStats(username string, password string, nBlocks string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getchaintxstats"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	nBlocksInt, err := utils.StoICheck(nBlocks)
	if err != nil {
		return nil, errors.New("input block height not integer")
	}
	payload.Params = []int{nBlocksInt}

	return PostReq(username, password, payload)
}

func GetDifficulty(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getdifficulty"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetMempoolAncestors(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmempoolancestors"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{txid}

	return PostReq(username, password, payload)
}

func GetMempoolEntry(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmempoolentry"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{txid}

	return PostReq(username, password, payload)
}

func GetMempoolInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmempoolinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetRawMempool(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrawmempool"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetTxOut(username string, password string, txid string, n int) ([]byte, error) {
	var payload RPCReq
	payload.Method = "gettxout"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = [2]interface{}{txid, n}

	return PostReq(username, password, payload)
}

// TODO: fix this route
func GetTxOutProof(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "gettxoutproof"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	temp := make(map[string]interface{})
	temp["txids"] = txid
	payload.Params = temp

	return PostReq(username, password, payload)
}

func GetTxOutSetInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "gettxoutsetinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func PreciousBlock(username string, password string, blockhash string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "preciousblock"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{blockhash}

	return PostReq(username, password, payload)
}

func PruneBlockchain(username string, password string, height string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getchaintxstats"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	heightI, err := utils.StoICheck(height)
	if err != nil {
		return nil, errors.New("input block height not integer")
	}
	payload.Params = []int{heightI}

	return PostReq(username, password, payload)
}

func SaveMempool(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "savemempool"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO: implement this route
func ScanTxOutset(username string, password string) ([]byte, error) {
	return nil, nil
}

func VerifyChain(username string, password string, nBlocks string, checkLevel string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "verifychain"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	nBlocksInt, err := utils.StoICheck(nBlocks)
	if err == nil {
		payload.Params = []int{nBlocksInt}
	}

	checkLevelInt, err := utils.StoICheck(checkLevel)
	if err == nil {
		payload.Params = []int{checkLevelInt}
	}

	return PostReq(username, password, payload)
}

// TODO: implement this route
func VerifyTxOutProof(username string, password string, txproof string) ([]byte, error) {
	return nil, nil
}

func GetMemoryInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmemoryinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetRPCInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrpcinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Help(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "help"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Logging(username string, password string, params ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "logging"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	var arr []string
	for _, param := range params {
		arr = append(arr, param)
	}
	payload.Params = arr

	return PostReq(username, password, payload)
}

func Stop(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "stop"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Uptime(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "uptime"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Generate(username string, password string, nBlocks string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "generate"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	nBlocksInt, err := utils.StoICheck(nBlocks)
	if err == nil {
		payload.Params = []int{nBlocksInt}
	}

	return PostReq(username, password, payload)
}

func GenerateToAddress(username string, password string, nBlocks string, address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "generatetoaddress"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	nBlocksInt, err := utils.StoICheck(nBlocks)
	if err == nil {
		payload.Params = [2]interface{}{nBlocksInt, address}
	}

	return PostReq(username, password, payload)
}

func SubmitBlock(username string, password string, hexdata string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "submitblock"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{hexdata}

	return PostReq(username, password, payload)
}

func SubmitHeader(username string, password string, hexdata string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "submitheader"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{hexdata}

	return PostReq(username, password, payload)
}

func AddNode(username string, password string, node string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "addnode"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{node}

	return PostReq(username, password, payload)
}

func ClearBanned(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "clearbanned"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func DisconnectNode(username string, password string, address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "disconnectnode"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{address}

	return PostReq(username, password, payload)
}

func GetAddedNodeInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getaddednodeinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetConnectionCount(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getconnectioncount"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetNetTotals(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnettotals"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetNetworkInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnetworkinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetNodeAddresses(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnodeaddresses"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetPeerInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getpeerinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func ListBanned(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listbanned"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func Ping(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "ping"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO: SetBan

func SetNetworkActive(username string, password string, state bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "setnetworkactive"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []interface{}{state}

	return PostReq(username, password, payload)
}

func AnalyzePSBT(username string, password string, psbt string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "analyzepsbt"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{psbt}

	return PostReq(username, password, payload)
}

func CombinePSBT(username string, password string, psbts ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "analyzepsbt"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	var arr []string
	for _, psbt := range psbts {
		arr = append(arr, psbt)
	}
	payload.Params = arr

	return PostReq(username, password, payload)
}

func CombineRawTransaction(username string, password string, psbts ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "combinerawtransaction"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	var arr []string
	for _, psbt := range psbts {
		arr = append(arr, psbt)
	}
	payload.Params = arr

	return PostReq(username, password, payload)
}

func ConvertToPsbt(username string, password string, psbt string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "converttopsbt"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{psbt}

	return PostReq(username, password, payload)
}

// TODO: Implement CreatePSBT, CreateRawTransaction routes here

func DecodePsbt(username string, password string, psbt string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "decodepsbt"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{psbt}

	return PostReq(username, password, payload)
}

func DecodeRawTransaction(username string, password string, rawtx string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "decoderawtransaction"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{rawtx}

	return PostReq(username, password, payload)
}

func DecodeScript(username string, password string, hexString string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "decodescript"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{hexString}

	return PostReq(username, password, payload)
}

// TODO: add optional param here
func FinalizePSBT(username string, password string, psbt string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "finalizepsbt"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{psbt}

	return PostReq(username, password, payload)
}

// TODO: add optional param here
func FundRawTransaction(username string, password string, hexString string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "fundrawtransaction"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{hexString}

	return PostReq(username, password, payload)
}

// TODO: add optional param here
func GetRawTransaction(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrawtransaction"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{txid}

	return PostReq(username, password, payload)
}

func JoinPSBTs(username string, password string, psbts ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrawtransaction"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	var arr []string
	for _, psbt := range psbts {
		arr = append(arr, psbt)
	}
	payload.Params = arr

	return PostReq(username, password, payload)
}

// TODO: add optional param here
func SendRawTransaction(username string, password string, hexString string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "sendrawtransaction"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{hexString}

	return PostReq(username, password, payload)
}

// TODO: add signrawtransactionwithkey, testmempoolaccept methods

func UtxoUpdatePSBT(username string, password string, psbt string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "utxoupdatepsbt"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{psbt}

	return PostReq(username, password, payload)
}

// TODO: add checks here
func CreateMultisig(username string, password string, n int, pubkeys ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "createmultisig"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	var arr []string
	for _, pubkey := range pubkeys {
		arr = append(arr, pubkey)
	}
	payload.Params = []interface{}{n, arr}

	return PostReq(username, password, payload)
}

func DeriveAddresses(username string, password string, descriptor string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "deriveaddresses"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{descriptor}

	return PostReq(username, password, payload)
}

// TODO: add optional param
func EstimateSmartFee(username string, password string, confTarget string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "estimatesmartfee"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	confTargetInt, err := utils.StoICheck(confTarget)
	if err != nil {
		return nil, errors.New("input block height not integer")
	}
	payload.Params = []int{confTargetInt}

	return PostReq(username, password, payload)
}

func GetDescriptorInfo(username string, password string, descriptor string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getdescriptorinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{descriptor}

	return PostReq(username, password, payload)
}

func SignMessageWithPrivkey(username string, password string, privkey string, message string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "signmessagewithprivkey"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{privkey, message}

	return PostReq(username, password, payload)
}

func ValidateAddress(username string, password string, address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "validateaddress"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{address}

	return PostReq(username, password, payload)
}

func VerifyMessage(username string, password string, address string, signature string, message string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "verifymessage"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{address, signature, message}

	return PostReq(username, password, payload)
}

func GetZmqNotifications(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getzmqnotifications"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func AbandonTransaction(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "abandontransaction"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{txid}

	return PostReq(username, password, payload)
}

func AbortRescan(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "abortrescan"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func AddMultisigAddress(username string, password string, n string, keys ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "addmultisigaddress"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	nInt, err := utils.StoICheck(n)
	if err != nil {
		return nil, errors.Wrap(err, "input not integer")
	}

	var arr []string
	for _, key := range keys {
		arr = append(arr, key)
	}

	payload.Params = []interface{}{nInt, arr}

	return PostReq(username, password, payload)
}

func BackupWallet(username string, password string, destination string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "backupwallet"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{destination}

	return PostReq(username, password, payload)
}

// TODO: add options here
func BumpFee(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "bumpfee"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{txid}

	return PostReq(username, password, payload)
}

// TODO: add optional params
func CreateWallet(username string, password string, walletName string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "createwallet"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{walletName}

	return PostReq(username, password, payload)
}

func DumpPrivKey(username string, password string, address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "dumpprivkey"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{address}

	return PostReq(username, password, payload)
}

func EncryptWallet(username string, password string, passphrase string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "encryptwallet"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{passphrase}

	return PostReq(username, password, payload)
}

func GetAddressesByLabel(username string, password string, label string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getaddressesbylabel"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{label}

	return PostReq(username, password, payload)
}

func GetAddressesInfo(username string, password string, address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getaddressinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{address}

	return PostReq(username, password, payload)
}

// TODO: optional params
func GetBalance(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getbalance"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO: optional params
func GetNewAddress(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnewaddress"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetRawChangeAddress(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrawchangeaddress"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetReceivedByLabel(username string, password string, label string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getreceivedbylabel"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{label}

	return PostReq(username, password, payload)
}

// TODO: add options here
func GetTransaction(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "gettransaction"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{txid}

	return PostReq(username, password, payload)
}

func GetUnconfirmedBalance(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getunconfirmedbalance"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func GetWalletInfo(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getwalletinfo"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO: add options here
func ImportAddress(username string, password string, address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "importaddress"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{address}

	return PostReq(username, password, payload)
}

// TODO: implement importmulti

// TODO: add options here
func ImportPrunedFunds(username string, password string, rawtx string, txoutproof string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "importprunedfunds"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{rawtx, txoutproof}

	return PostReq(username, password, payload)
}

// TODO: add options here
func ImportPubkey(username string, password string, pubkey string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "importpubkey"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{pubkey}

	return PostReq(username, password, payload)
}

func ImportWallet(username string, password string, name string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "importwallet"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{name}

	return PostReq(username, password, payload)
}

// TODO: add options here
func KeypoolRefill(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "keypoolrefill"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func ListAddressGroupings(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listaddressgroupings"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO option
func ListLabels(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listlabels"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func ListLockUnspent(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listlockunspent"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO option
func ListReceivedByAddress(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listreceivedbyaddress"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO option
func ListReceivedByLabel(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listreceivedbylabel"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO option
func ListSinceBlock(username string, password string, blockhash string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listsinceblock"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{blockhash}

	return PostReq(username, password, payload)
}

// TODO option
func ListTransactions(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listtranscations"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO option
func ListUnspent(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listunspent"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func ListWalletDir(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listwalletdir"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

func ListWallets(username string, password string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listwallets"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	return PostReq(username, password, payload)
}

// TODO: implement lockunspent

func RemovePrunedFunds(username string, password string, txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "removeprunedfunds"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"
	payload.Params = []string{txid}

	return PostReq(username, password, payload)
}

func RescanBlockchain(username string, password string, startHeight string, stopHeight string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "rescanblockchain"
	payload.ID = "curltext"
	payload.JsonRPC = "1.0"

	startHeightI, err := utils.StoICheck(startHeight)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert string to int")
	}

	stopHeightI, err := utils.StoICheck(stopHeight)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert string to int")
	}

	payload.Params = []int{startHeightI, stopHeightI}

	return PostReq(username, password, payload)
}

func main() {
	data, err := ListWallets("user", "password")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(data))
}
