package bitcoinrpc

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
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
var ID = "curltext"
var JsonRPC = "1.0"
var RPCUser = "user"
var RPCPass = "password"

func SetBitcoindURL(url, rpcuser, rpcpass string) {
	BitcoindURL = url
	RPCUser = rpcuser
	RPCPass = rpcpass
}

func PostReq(payload RPCReq) ([]byte, error) {
	var req *http.Request

	payload.ID = ID
	payload.JsonRPC = JsonRPC
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "could not marshal json, quitting")
	}

	req, err = http.NewRequest("POST", BitcoindURL, bytes.NewBuffer(payloadJson))
	if err != nil {
		return nil, errors.Wrap(err, "did not POST to bitcoind")
	}

	req.SetBasicAuth(RPCUser, RPCPass)

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

func GetBestBlockHash() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getbestblockhash"

	return PostReq(payload)
}

func GetBlock(blockhash string) ([]byte, error) {
	if len(blockhash) != 64 {
		return nil, errors.New("length of blockhash not 32 bytes")
	}
	var payload RPCReq
	payload.Method = "getblock"
	payload.Params = []string{blockhash}

	return PostReq(payload)
}

func GetBlockCount() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockcount"

	return PostReq(payload)
}

func GetBlockHash(blockNumber uint32) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockhash"
	payload.Params = []uint32{blockNumber}

	return PostReq(payload)
}

func GetBlockHeader(blockhash string) ([]byte, error) {
	if len(blockhash) != 64 {
		return nil, errors.New("length of blockhash not 32 bytes")
	}
	var payload RPCReq
	payload.Method = "getblockheader"
	payload.Params = []string{blockhash}

	return PostReq(payload)
}

func GetBlockStats(hashOrHeight string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getblockstats"
	height, err := utils.ToInt(hashOrHeight)
	if err != nil {
		// hash and not height
		if len(hashOrHeight) != 64 {
			return nil, errors.New("length of blockhash not 32 bytes")
		}
		payload.Params = []string{hashOrHeight}
	} else {
		// blockheight
		payload.Params = []int{height}
	}

	return PostReq(payload)
}

func GetChainTips() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getchaintips"

	return PostReq(payload)
}

func GetChainTxStats(nBlocks string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getchaintxstats"

	nBlocksInt, err := utils.ToInt(nBlocks)
	if err != nil {
		return nil, errors.New("input block height not integer")
	}
	payload.Params = []int{nBlocksInt}

	return PostReq(payload)
}

func GetDifficulty() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getdifficulty"

	return PostReq(payload)
}

func GetMempoolAncestors(txid string) ([]byte, error) {
	if len(txid) != 64 {
		return nil, errors.New("length of txid not 32 bytes")
	}
	var payload RPCReq
	payload.Method = "getmempoolancestors"
	payload.Params = []string{txid}

	return PostReq(payload)
}

func GetMempoolEntry(txid string) ([]byte, error) {
	if len(txid) != 64 {
		return nil, errors.New("length of txid not 32 bytes")
	}
	var payload RPCReq
	payload.Method = "getmempoolentry"
	payload.Params = []string{txid}

	return PostReq(payload)
}

func GetMempoolInfo() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmempoolinfo"

	return PostReq(payload)
}

func GetRawMempool() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrawmempool"

	return PostReq(payload)
}

func GetTxOut(txid string, n int) ([]byte, error) {
	if len(txid) != 64 {
		return nil, errors.New("length of txid not 32 bytes")
	}
	var payload RPCReq
	payload.Method = "gettxout"
	payload.Params = [2]interface{}{txid, n}

	return PostReq(payload)
}

// TODO: fix this route
func GetTxOutProof(txid string) ([]byte, error) {
	if len(txid) != 64 {
		return nil, errors.New("length of txid not 32 bytes")
	}
	var payload RPCReq
	payload.Method = "gettxoutproof"

	temp := make(map[string]interface{})
	temp["txids"] = txid
	payload.Params = temp

	return PostReq(payload)
}

func GetTxOutSetInfo() ([]byte, error) {
	var payload RPCReq
	payload.Method = "gettxoutsetinfo"

	return PostReq(payload)
}

func PreciousBlock(blockhash string) ([]byte, error) {
	if len(blockhash) != 64 {
		return nil, errors.New("length of blockhash not 32 bytes")
	}
	var payload RPCReq
	payload.Method = "preciousblock"
	payload.Params = []string{blockhash}

	return PostReq(payload)
}

func PruneBlockchain(height string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getchaintxstats"

	heightI, err := utils.ToInt(height)
	if err != nil {
		return nil, errors.New("input block height not integer")
	}
	payload.Params = []int{heightI}

	return PostReq(payload)
}

func SaveMempool() ([]byte, error) {
	var payload RPCReq
	payload.Method = "savemempool"

	return PostReq(payload)
}

// TODO: implement this route
func ScanTxOutset() ([]byte, error) {
	return nil, nil
}

func VerifyChain(nBlocks string, checkLevel string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "verifychain"

	var temp []interface{}
	if nBlocks != "" {
		nBlocksInt, err := utils.ToInt(nBlocks)
		if err != nil {
			return nil, errors.New("nblocks not int")
		}
		temp = append(temp, nBlocksInt)
	}
	if checkLevel != "" {
		checkLevelInt, err := utils.ToInt(checkLevel)
		if err != nil {
			return nil, errors.New("checklevel not int")
		}
		temp = append(temp, checkLevelInt)
	}

	return PostReq(payload)
}

// TODO: implement this route
func VerifyTxOutProof(txproof string) ([]byte, error) {
	return nil, nil
}

func GetMemoryInfo() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getmemoryinfo"

	return PostReq(payload)
}

func GetRPCInfo() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrpcinfo"

	return PostReq(payload)
}

func Help() ([]byte, error) {
	var payload RPCReq
	payload.Method = "help"

	return PostReq(payload)
}

func Logging(params ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "logging"

	var arr []string
	for _, param := range params {
		arr = append(arr, param)
	}
	payload.Params = arr

	return PostReq(payload)
}

func Stop() ([]byte, error) {
	var payload RPCReq
	payload.Method = "stop"

	return PostReq(payload)
}

func Uptime() ([]byte, error) {
	var payload RPCReq
	payload.Method = "uptime"

	return PostReq(payload)
}

func Generate(nBlocks string, maxtries string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "generate"

	var temp []interface{}
	nBlocksInt, err := utils.ToInt(nBlocks)
	if err != nil {
		return nil, errors.New("nblocks not int")
	}

	temp = append(temp, nBlocksInt)
	if maxtries != "" {
		maxTriesInt, err := utils.ToInt(maxtries)
		if err != nil {
			return nil, errors.New("maxtries not int")
		}
		temp = append(temp, maxTriesInt)
	}
	payload.Params = temp
	return PostReq(payload)
}

func GenerateToAddress(nBlocks string, address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "generatetoaddress"

	nBlocksInt, err := utils.ToInt(nBlocks)
	if err != nil {
		return nil, errors.New("nblocks not int")
	}

	payload.Params = [2]interface{}{nBlocksInt, address}

	return PostReq(payload)
}

func SubmitBlock(hexdata string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "submitblock"
	payload.Params = []string{hexdata}

	return PostReq(payload)
}

func SubmitHeader(hexdata string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "submitheader"
	payload.Params = []string{hexdata}

	return PostReq(payload)
}

func AddNode(node string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "addnode"
	payload.Params = []string{node}

	return PostReq(payload)
}

func ClearBanned() ([]byte, error) {
	var payload RPCReq
	payload.Method = "clearbanned"

	return PostReq(payload)
}

func DisconnectNode(address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "disconnectnode"
	payload.Params = []string{address}

	return PostReq(payload)
}

func GetAddedNodeInfo() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getaddednodeinfo"

	return PostReq(payload)
}

func GetConnectionCount() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getconnectioncount"

	return PostReq(payload)
}

func GetNetTotals() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnettotals"

	return PostReq(payload)
}

func GetNetworkInfo() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnetworkinfo"

	return PostReq(payload)
}

func GetNodeAddresses() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnodeaddresses"

	return PostReq(payload)
}

func GetPeerInfo() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getpeerinfo"

	return PostReq(payload)
}

func ListBanned() ([]byte, error) {
	var payload RPCReq
	payload.Method = "listbanned"

	return PostReq(payload)
}

func Ping() ([]byte, error) {
	var payload RPCReq
	payload.Method = "ping"

	return PostReq(payload)
}

func SetBan(subnet string, command string, bantime string, absolute bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "setnetworkactive"
	var temp []interface{}

	temp = append(temp, subnet, command)
	if bantime != "" {
		bantimeInt, err := utils.ToInt(bantime)
		if err != nil {
			return nil, errors.New("ban time not int")
		}
		temp = append(temp, bantimeInt)
	}
	if absolute {
		temp = append(temp, absolute)
	}
	payload.Params = temp
	return PostReq(payload)
}

func SetNetworkActive(state bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "setnetworkactive"
	payload.Params = []interface{}{state}

	return PostReq(payload)
}

func AnalyzePSBT(psbt string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "analyzepsbt"
	payload.Params = []string{psbt}

	return PostReq(payload)
}

func CombinePSBT(psbts ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "analyzepsbt"

	var arr []string
	for _, psbt := range psbts {
		arr = append(arr, psbt)
	}
	payload.Params = arr

	return PostReq(payload)
}

func CombineRawTransaction(psbts ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "combinerawtransaction"

	var arr []string
	for _, psbt := range psbts {
		arr = append(arr, psbt)
	}
	payload.Params = arr

	return PostReq(payload)
}

func ConvertToPsbt(psbt string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "converttopsbt"
	payload.Params = []string{psbt}

	return PostReq(payload)
}

// TODO: Implement CreatePSBT, CreateRawTransaction routes here

func DecodePsbt(psbt string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "decodepsbt"
	payload.Params = []string{psbt}

	return PostReq(payload)
}

func DecodeRawTransaction(rawtx string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "decoderawtransaction"
	payload.Params = []string{rawtx}

	return PostReq(payload)
}

func DecodeScript(hexString string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "decodescript"
	payload.Params = []string{hexString}

	return PostReq(payload)
}

func FinalizePSBT(psbt string, extract bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "finalizepsbt"

	var temp []interface{}
	temp = append(temp, psbt)

	if extract {
		temp = append(temp, extract)
	}
	payload.Params = temp

	return PostReq(payload)
}

func FundRawTransaction(hexString string, includeWatching bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "fundrawtransaction"

	var temp []interface{}
	temp = append(temp, hexString)

	if includeWatching {
		temp = append(temp, includeWatching)
	}
	payload.Params = temp

	return PostReq(payload)
}

func GetRawTransaction(txid string, verbose bool, blockhash string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrawtransaction"

	var temp []interface{}
	temp = append(temp, txid)

	if verbose {
		temp = append(temp, true)
	}
	payload.Params = temp

	if blockhash != "" {
		if len(blockhash) != 32 {
			return nil, errors.New("length of blockhash not 32")
		}
		temp = append(temp, blockhash)
	}
	return PostReq(payload)
}

func JoinPSBTs(psbts ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrawtransaction"
	var arr []string
	for _, psbt := range psbts {
		arr = append(arr, psbt)
	}
	payload.Params = arr

	return PostReq(payload)
}

func SendRawTransaction(hexString string, allowHighFees bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "sendrawtransaction"

	var temp []interface{}
	temp = append(temp, hexString)

	if allowHighFees {
		temp = append(temp, allowHighFees)
	}
	payload.Params = temp
	return PostReq(payload)
}

// TODO: add signrawtransactionwithkey, testmempoolaccept methods

func UtxoUpdatePSBT(psbt string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "utxoupdatepsbt"
	payload.Params = []string{psbt}

	return PostReq(payload)
}

func CreateMultisig(n int, pubkeys ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "createmultisig"
	var arr []string
	for _, pubkey := range pubkeys {
		if len(pubkey) != 64 {
			return nil, errors.New("length of pubkey not 32 bytes")
		}
		arr = append(arr, pubkey)
	}
	payload.Params = []interface{}{n, arr}

	return PostReq(payload)
}

func DeriveAddresses(descriptor string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "deriveaddresses"
	payload.Params = []string{descriptor}

	return PostReq(payload)
}

func EstimateSmartFee(confTarget string, estimateMode string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "estimatesmartfee"
	confTargetInt, err := utils.ToInt(confTarget)
	if err != nil {
		return nil, errors.New("input block height not integer")
	}

	var temp []interface{}
	temp = append(temp, confTargetInt)

	if estimateMode != "" {
		temp = append(temp, estimateMode)
	}
	payload.Params = temp

	return PostReq(payload)
}

func GetDescriptorInfo(descriptor string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getdescriptorinfo"
	payload.Params = []string{descriptor}

	return PostReq(payload)
}

func SignMessageWithPrivkey(privkey string, message string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "signmessagewithprivkey"
	payload.Params = []string{privkey, message}

	return PostReq(payload)
}

func ValidateAddress(address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "validateaddress"
	payload.Params = []string{address}

	return PostReq(payload)
}

func VerifyMessage(address string, signature string, message string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "verifymessage"
	payload.Params = []string{address, signature, message}

	return PostReq(payload)
}

func AbandonTransaction(txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "abandontransaction"
	payload.Params = []string{txid}

	return PostReq(payload)
}

func AbortRescan() ([]byte, error) {
	var payload RPCReq
	payload.Method = "abortrescab"
	return PostReq(payload)
}

func AddMultisigAddress(n string, keys ...string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "addmultisigaddress"
	nInt, err := utils.ToInt(n)
	if err != nil {
		return nil, errors.Wrap(err, "input not integer")
	}

	var arr []string
	for _, key := range keys {
		arr = append(arr, key)
	}

	payload.Params = []interface{}{nInt, arr}

	return PostReq(payload)
}

func BackupWallet(destination string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "backupwallet"
	payload.Params = []string{destination}

	return PostReq(payload)
}

// TODO: add json options here
func BumpFee(txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "bumpfee"
	payload.Params = []string{txid}

	return PostReq(payload)
}

func CreateWallet(walletName string, disablePrivkey bool, blank bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "createwallet"

	var temp []interface{}
	temp = append(temp, walletName)

	if disablePrivkey {
		temp = append(temp, disablePrivkey)
	}

	if blank {
		temp = append(temp, blank)
	}

	payload.Params = temp

	return PostReq(payload)
}

func DumpPrivKey(address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "dumpprivkey"
	payload.Params = []string{address}

	return PostReq(payload)
}

func EncryptWallet(passphrase string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "encryptwallet"
	payload.Params = []string{passphrase}

	return PostReq(payload)
}

func GetAddressesByLabel(label string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getaddressesbylabel"
	payload.Params = []string{label}

	return PostReq(payload)
}

func GetAddressesInfo(address string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getaddressinfo"
	payload.Params = []string{address}

	return PostReq(payload)
}

func GetBalance(dummy string, minConf string, watchOnly bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getbalance"

	var temp []interface{}
	temp = append(temp, dummy)

	if minConf != "" {
		minConfInt, err := utils.ToInt(minConf)
		if err != nil {
			return nil, errors.Wrap(err, "minconf not int")
		}
		temp = append(temp, minConfInt)
	}

	if watchOnly {
		temp = append(temp, watchOnly)
	}

	payload.Params = temp

	return PostReq(payload)
}

func GetNewAddress(label string, addressType string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getnewaddress"

	var temp []interface{}

	if label != "" {
		temp = append(temp, label)
	}

	if addressType != "" {
		temp = append(temp, addressType)
	}

	payload.Params = temp

	return PostReq(payload)
}

func GetRawChangeAddress() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getrawchangeaddress"
	return PostReq(payload)
}

func GetReceivedByLabel(label string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "getreceivedbylabel"
	payload.Params = []string{label}

	return PostReq(payload)
}

func GetTransaction(txid string, watchOnly bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "gettransaction"

	var temp []interface{}

	temp = append(temp, txid)
	if watchOnly {
		temp = append(temp, watchOnly)
	}

	payload.Params = temp

	return PostReq(payload)
}

func GetUnconfirmedBalance() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getunconfirmedbalance"
	return PostReq(payload)
}

func GetWalletInfo() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getwalletinfo"
	return PostReq(payload)
}

func ImportAddress(address string, label string, rescan bool, p2sh bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "importaddress"

	var temp []interface{}
	temp = append(temp, address)
	if label != "" {
		temp = append(temp, label)
	}
	if rescan {
		temp = append(temp, rescan)
	}
	if p2sh {
		temp = append(temp, p2sh)
	}

	payload.Params = temp

	return PostReq(payload)
}

// TODO: implement importmulti

func ImportPrunedFunds(rawtx string, txoutproof string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "importprunedfunds"
	payload.Params = []string{rawtx, txoutproof}

	return PostReq(payload)
}

func ImportPubkey(pubkey string, label string, rescan bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "importpubkey"

	var temp []interface{}
	temp = append(temp, pubkey)
	if label != "" {
		temp = append(temp, label)
	}
	if rescan {
		temp = append(temp, rescan)
	}

	payload.Params = temp
	return PostReq(payload)
}

func ImportWallet(name string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "importwallet"
	payload.Params = []string{name}

	return PostReq(payload)
}

func KeypoolRefill(newSize string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "keypoolrefill"

	if newSize != "" {
		newSizeInt, err := utils.ToInt(newSize)
		if err != nil {
			return nil, errors.Wrap(err, "new size not int")
		}
		payload.Params = []int{newSizeInt}
	}
	return PostReq(payload)
}

func ListAddressGroupings() ([]byte, error) {
	var payload RPCReq
	payload.Method = "listaddressgroupings"
	return PostReq(payload)
}

func ListLabels(purpose string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listlabels"

	if purpose != "" {
		payload.Params = []string{purpose}
	}

	return PostReq(payload)
}

func ListLockUnspent() ([]byte, error) {
	var payload RPCReq
	payload.Method = "listlockunspent"
	return PostReq(payload)
}

func ListReceivedByAddress(minConf string, includeEmpty bool, watchOnly bool, addressFilter string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listreceivedbyaddress"

	var temp []interface{}
	if minConf != "" {
		minConfInt, err := utils.ToInt(minConf)
		if err != nil {
			return nil, errors.Wrap(err, "minconf not int")
		}
		temp = append(temp, minConfInt)
	}
	if includeEmpty {
		temp = append(temp, includeEmpty)
	}
	if watchOnly {
		temp = append(temp, watchOnly)
	}
	if addressFilter != "" {
		temp = append(temp, addressFilter)
	}

	payload.Params = temp

	return PostReq(payload)
}

func ListReceivedByLabel(minConf string, includeEmpty bool, watchOnly bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listreceivedbylabel"

	var temp []interface{}
	if minConf != "" {
		minConfInt, err := utils.ToInt(minConf)
		if err != nil {
			return nil, errors.Wrap(err, "minconf not int")
		}
		temp = append(temp, minConfInt)
	}
	if includeEmpty {
		temp = append(temp, includeEmpty)
	}
	if watchOnly {
		temp = append(temp, watchOnly)
	}

	payload.Params = temp

	return PostReq(payload)
}

func ListSinceBlock(blockhash string, targetConf string, watchOnly bool, includeRemoved bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listsinceblock"

	var temp []interface{}
	if blockhash != "" {
		if len(blockhash) != 64 {
			return nil, errors.New("length of blockhash not 32 bytes")
		}
		temp = append(temp, blockhash)
	}
	if targetConf != "" {
		targetConfInt, err := utils.ToInt(targetConf)
		if err != nil {
			return nil, errors.Wrap(err, "minconf not int")
		}
		temp = append(temp, targetConfInt)
	}
	if watchOnly {
		temp = append(temp, watchOnly)
	}
	if includeRemoved {
		temp = append(temp, includeRemoved)
	}

	payload.Params = temp

	return PostReq(payload)
}

func ListTransactions(label string, count string, skip string, watchOnly bool) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listtranscations"

	var temp []interface{}
	if label != "" {
		temp = append(temp, label)
	}
	if count != "" {
		countInt, err := utils.ToInt(count)
		if err != nil {
			return nil, errors.Wrap(err, "minconf not int")
		}
		temp = append(temp, countInt)
	}
	if skip != "" {
		skipInt, err := utils.ToInt(skip)
		if err != nil {
			return nil, errors.Wrap(err, "minconf not int")
		}
		temp = append(temp, skipInt)
	}
	if watchOnly {
		temp = append(temp, watchOnly)
	}

	payload.Params = temp

	return PostReq(payload)
}

func ListUnspent(minConf string, maxConf string, addresses []string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "listunspent"

	var temp []interface{}

	if minConf != "" {
		minConfInt, err := utils.ToInt(minConf)
		if err != nil {
			return nil, errors.Wrap(err, "minconf not int")
		}
		temp = append(temp, minConfInt)
	}
	if maxConf != "" {
		maxConfInt, err := utils.ToInt(maxConf)
		if err != nil {
			return nil, errors.Wrap(err, "minconf not int")
		}
		temp = append(temp, maxConfInt)
	}
	if len(addresses) != 0 {
		temp = append(temp, addresses)
	}
	payload.Params = temp

	return PostReq(payload)
}

func ListWalletDir() ([]byte, error) {
	var payload RPCReq
	payload.Method = "listwalletdir"

	return PostReq(payload)
}

func ListWallets() ([]byte, error) {
	var payload RPCReq
	payload.Method = "listwallets"

	return PostReq(payload)
}

// TODO: implement lockunspent

func RemovePrunedFunds(txid string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "removeprunedfunds"
	payload.Params = []string{txid}

	return PostReq(payload)
}

func RescanBlockchain(startHeight string, stopHeight string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "rescanblockchain"

	startHeightI, err := utils.ToInt(startHeight)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert string to int")
	}

	stopHeightI, err := utils.ToInt(stopHeight)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert string to int")
	}

	payload.Params = []int{startHeightI, stopHeightI}

	return PostReq(payload)
}

// TODO: impelemnt sendmany endpoint here

func SendToAddress(address string, amount string,
	comment string, commentTo string, subtractFee bool, replaceAble bool, confTarget int,
	estimateMode string) ([]byte, error) {

	var payload RPCReq
	payload.Method = "sendtoaddress"
	amountI, err := utils.ToInt(amount)
	if err != nil {
		return nil, errors.New("could not convert string to int, quitting")
	}

	var temp []interface{}
	temp = append(temp, address, amountI)

	if comment != "" {
		temp = append(temp, comment)
	}
	if commentTo != "" {
		temp = append(temp, commentTo)
	}
	if subtractFee {
		temp = append(temp, true)
	}
	if replaceAble {
		temp = append(temp, true)
	}
	if estimateMode != "" {
		temp = append(temp, estimateMode)
	}
	if confTarget != 0 {
		temp = append(temp, confTarget)
	}

	payload.Params = temp
	return PostReq(payload)
}

func SetHdSeed(newkeypool bool, seed string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "sethdseed"

	var temp []interface{}
	if newkeypool {
		temp = append(temp, true)
	}
	if seed != "" {
		temp = append(temp, seed)
	}

	payload.Params = temp
	return PostReq(payload)
}

func SetLabel(address string, label string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "setlabel"
	payload.Params = []interface{}{address, label}

	return PostReq(payload)
}

func SetTxFee(amount string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "settxfee"

	amountI, err := utils.ToInt(amount)
	if err != nil {
		return nil, errors.New("could not convert string to int, quitting")
	}

	payload.Params = []interface{}{amountI}

	return PostReq(payload)
}

func SignMessage(address string, message string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "signmessage"

	payload.Params = []interface{}{address, message}

	return PostReq(payload)
}

func SignRawTransactionWithWallet(hexString string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "signrawtransactionwithwallet"

	payload.Params = []interface{}{hexString}

	return PostReq(payload)
}

func UnloadWallet(walletName string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "unloadwallet"

	payload.Params = []interface{}{walletName}

	return PostReq(payload)
}

// TODO: add walletcreatefundedpsbt method

func WalletLock() ([]byte, error) {
	var payload RPCReq
	payload.Method = "walletlock"

	return PostReq(payload)
}

func WalletPassphrase(passphrase string, timeout string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "walletpassphrase"

	timeoutI, err := utils.ToInt(timeout)
	if err != nil {
		return nil, errors.New("timeout not integer")
	}

	payload.Params = []interface{}{passphrase, timeoutI}
	return PostReq(payload)
}

func WalletPassphraseChange(old string, new string) ([]byte, error) {
	var payload RPCReq
	payload.Method = "walletpassphrasechange"

	payload.Params = []interface{}{old, new}
	return PostReq(payload)
}

func WalletProcessPSBT(psbt string, sign bool,
	sighashType string, bip32derivs bool) ([]byte, error) {

	var payload RPCReq
	payload.Method = "walletprocesspsbt"

	var temp []interface{}
	temp = append(temp, psbt)

	if sign {
		temp = append(temp, true)
	}
	if sighashType != "" {
		temp = append(temp, sighashType)
	}
	if bip32derivs {
		temp = append(temp, true)
	}

	payload.Params = temp
	return PostReq(payload)
}

func GetZmqNotifications() ([]byte, error) {
	var payload RPCReq
	payload.Method = "getzmqnotifications"

	return PostReq(payload)
}
