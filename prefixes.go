package hdwallets

import "github.com/btcsuite/btcd/chaincfg"

var (
	prefixes map[string]string
)

func init() {
	prefixes = make(map[string]string)
	prefixes[chaincfg.MainNetParams.Name] = "bitcoincash"
	prefixes[chaincfg.TestNet3Params.Name] = "bchtest"
	prefixes[chaincfg.RegressionNetParams.Name] = "bchreg"
}
