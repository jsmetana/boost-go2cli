package boost

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"
)

type FIL int64

type Cmd struct {
	name string
	args []string
}

func Init(name string, args ...string) *Cmd {
	finalArgs := append([]string{name}, args...)
	finalArgs = append(finalArgs, "--json")

	return &Cmd{
		name: name,
		args: finalArgs,
	}
}

func (c *Cmd) run(arg ...string) ([]byte, error) {
	cmd := exec.Cmd{
		Path: c.name,
		Args: append(c.args, arg...),
	}
	fmt.Println(cmd.Args)
	return cmd.CombinedOutput()
}

func (c *Cmd) runStruct(args []string, s interface{}) {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		// v.Field(i).Type().Kind() == reflect.String
		// v.Type().Field(i).Name

		val := v.Field(i).Interface()
		switch val {
		case false:
		case 0:
		case uint64(0):
		case int64(0):
		case "":
		default:
			args = append(args, fmt.Sprintf("%v%v", v.Type().Field(i).Tag, val))
		}
	}

	c.run(args...)
}

func (c *Cmd) Init() {
	c.run("init")
}

type DealArgs struct {
	HttpUrl            string "--http-url="
	HttpHeaders        string "--http-headers="
	Provider           string "--provider="
	Commp              string "--commp="
	PieceSize          uint64 "--piece-size="
	CarSize            uint64 "--car-size="
	PayloadCid         string "--payload-cid="
	StartEpoch         int    "--start-epoch="
	Duration           int    "--duration="
	ProviderCollateral int    "--provider-collateral="
	StoragePrice       int64  "--storage-price="
	Verified           bool   "--verified="
	Wallet             string "--wallet="
}

func (c *Cmd) Deal(args DealArgs) {
	c.runStruct([]string{"deal"}, args)
}

type DealStatusArgs struct {
	Provider string "--provider="
	DealUuid string "--deal-uuid="
	Wallet   string "--wallet="
}

func (c *Cmd) DealStatus(args DealStatusArgs) {
	c.runStruct([]string{"deal-status"}, args)
}

type OfflineDealArgs struct {
	Provider           string "--provider="
	Commp              string "--commp="
	PieceSize          uint64 "--piece-size="
	CarSize            uint64 "--car-size="
	PayloadCid         string "--payload-cid="
	StartEpoch         int    "--start-epoch="
	Duration           int    "--duration="
	ProviderCollateral int    "--provider-collateral="
	StoragePrice       int64  "--storage-price="
	Verified           bool   "--verified="
	Wallet             string "--wallet="
}

func (c *Cmd) OfflineDeal(args OfflineDealArgs) {
	c.runStruct([]string{"offline-deal"}, args)
}

func (c *Cmd) ProviderLibp2pInfo(providerAddr string) {
	c.run("provider", "libp2p-info", providerAddr)
}

func (c *Cmd) ProviderStorageAsk() {
}

func (c *Cmd) ProviderRetrievalAsk() {
}

func (c *Cmd) ProviderRetrievalTransports() {
}

type walletAddrOut struct {
	Address string
}

type WalletType string

const (
	WalletTypeSecp256k1 WalletType = "secp256k1"
	WalletTypeBls                  = "bls"
)

func (c *Cmd) WalletNew(t WalletType) string {
	// TODO: should there be a default for empty t?
	var address walletAddrOut
	out, err := c.run("wallet", "new", string(t))
	if err != nil {
		// TODO: add error checking
	}
	json.Unmarshal(out, &address)
	return address.Address
}

type WalletListOut struct {
	Address string
	Balance FIL
	Default bool
	Id      string
	Nonce   int
}

func (c *Cmd) WalletList(lookupIds bool) []WalletListOut {
	var list []WalletListOut

	out, err := c.run("wallet", "list", "--id="+strconv.FormatBool(lookupIds))

	if err != nil {
		// TODO: add error checking
	}
	json.Unmarshal(out, &list)
	return list
}

type WalletBalanceOut struct {
	Balance FIL
	Warning string
	// "warning": "may display 0 if chain sync in progress" is returned from boost when balance == 0
}

func (c *Cmd) WalletBalance(address string) FIL {
	// empty address is the default address
	var balance WalletBalanceOut
	out, err := c.run("wallet", "balance", address)
	if err != nil {
		// TODO: add error checking
	}
	json.Unmarshal(out, &balance)
	return balance.Balance
}

type walletExportOut struct {
	Key string
}

func (c *Cmd) WalletExport(address string) string {
	var key walletExportOut
	out, err := c.run("wallet", "export", address)

	if err != nil {
		// TODO: add error checking
	}

	json.Unmarshal(out, &key)
	return key.Key
}

type WalletFormat string

const (
	WalletFormatHexLotus  WalletFormat = "hex-lotus"
	WalletFormatJsonLotus              = "json-lotus"
	WalletFormatGfcJson                = "gfc-json"
)

// TODO: stdin input, test
func (c *Cmd) WalletImport(path string, format WalletFormat, asDefault bool) string {
	var address walletAddrOut

	f := string(format)
	d := strconv.FormatBool(asDefault)

	out, err := c.run("wallet", "import", "--format="+f, "--as-default="+d, path)
	if err != nil {
		// TODO: add error checking
	}

	json.Unmarshal(out, &address)
	return address.Address
}

func (c *Cmd) WalletDefault() string {
	var address walletAddrOut

	out, err := c.run("wallet", "default")
	if err != nil {
		// TODO: add error checking
	}

	json.Unmarshal(out, &address)
	return address.Address
}

func (c *Cmd) WalletSetDefault(address string) {
	if address != "" { // TODO: maybe give an error?
		c.run("wallet", "set-default", address)
	}
}

func (c *Cmd) WalletDelete(address string) {
	c.run("wallet", "delete", address)
}

type walletSignatureOut struct {
    Signature string
}

func (c *Cmd) WalletSign(address, hexMessage string) string {
    var signature walletSignatureOut

    out, err := c.run("wallet", "sign", address, hexMessage)
    if err != nil {
        // TODO: add error checking
    }

    json.Unmarshal(out, &signature)
    return signature.Signature
}

