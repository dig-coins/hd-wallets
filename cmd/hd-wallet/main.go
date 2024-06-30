package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	hdwallet "github.com/dig-coins/hd-wallet"
)

const mnemonicFile = "mnemonic.txt"

func main() {
	d, err := os.ReadFile(mnemonicFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		mnemonic, _ := hdwallet.NewMnemonic(12, "")
		_ = os.WriteFile(mnemonicFile, []byte(mnemonic), os.ModePerm)

		d = []byte(mnemonic)
	}

	seed, err := hdwallet.NewSeed(string(d), "", "")
	if err != nil {
		panic(err)
	}

	master, err := hdwallet.NewKey(
		hdwallet.Seed(seed),
	)
	if err != nil {
		panic(err)
	}

	fnPrint := func(purpose, addressIndex, coinType uint32, coinName string) {
		wallet, _ := master.GetWallet(hdwallet.Purpose(purpose), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex))

		var address, addressSegWit, wifPrivateKey, addressSegWitNative string

		address, err = wallet.GetAddress()
		if err != nil {
			panic(err)
		}

		addressSegWit, err = wallet.GetKey().AddressP2WPKHInP2SH()
		if err != nil {
			panic(err)
		}

		addressSegWitNative, err = wallet.GetKey().AddressP2WPKH()
		if err != nil {
			panic(err)
		}

		wifPrivateKey, err = wallet.GetKey().PrivateWIF(true)
		if err != nil {
			panic(err)
		}

		privateKeyHex := wallet.GetKey().PrivateHex()

		pubKeyHex := wallet.GetKey().PublicHex(true)

		fmt.Printf("%s Index %d:\n  privateKeyWIF:%s\n  privateKeyHex:%s\n  publicKey:%s\n  address:%s\n  addressSegWit:%s\n  addressSegWitNative:%s\n----------\n",
			coinName, addressIndex, wifPrivateKey, privateKeyHex, pubKeyHex, address, addressSegWit, addressSegWitNative)
	}

	for idx := uint32(0); idx < 4; idx++ {
		fnPrint(hdwallet.ZeroQuote+49, idx, hdwallet.BTC, "BTC")
	}

	for idx := uint32(0); idx < 4; idx++ {
		fnPrint(hdwallet.ZeroQuote+49, idx, hdwallet.BTCTestnet, "BTCTestnet")
	}

	for idx := uint32(0); idx < 4; idx++ {
		fnPrint(hdwallet.ZeroQuote+44, idx, hdwallet.LTC, "LTC")
	}

	for idx := uint32(0); idx < 4; idx++ {
		fnPrint(hdwallet.ZeroQuote+44, idx, hdwallet.DOGE, "DOGE")
	}

	for idx := uint32(0); idx < 4; idx++ {
		fnPrint(hdwallet.ZeroQuote+44, idx, hdwallet.TRX, "TRX")
	}

	for idx := uint32(0); idx < 4; idx++ {
		fnPrint(hdwallet.ZeroQuote+44, idx, hdwallet.ETH, "ETH")
	}

	time.Sleep(time.Hour)
}
