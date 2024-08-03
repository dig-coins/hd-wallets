package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	hdwallet "github.com/dig-coins/hd-wallet"
	"github.com/howeyc/gopass"
)

const mnemonicFile = "mnemonic.txt"

func main() {
	var seedPassword string
	flag.StringVar(&seedPassword, "p", "", "password of seed")
	flag.Parse()

	if seedPassword == "<null>" {
		v, e := gopass.GetPasswdPrompt("input your password", true, os.Stdin, os.Stdout)
		if e != nil {
			panic(e)
		}

		seedPassword = string(v)
	}

	d, err := os.ReadFile(mnemonicFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}

		mnemonic, _ := hdwallet.NewMnemonic(12, "")
		_ = os.WriteFile(mnemonicFile, []byte(mnemonic), os.ModePerm)

		d = []byte(mnemonic)
	}

	s := string(d)
	s = strings.Trim(s, " \r\n\t")

	seed, err := hdwallet.NewSeed(s, seedPassword, "")
	if err != nil {
		panic(err)
	}

	master, err := hdwallet.NewKey(
		hdwallet.Seed(seed),
	)
	if err != nil {
		panic(err)
	}

	fnPrintBTC := func(addressIndex, coinType uint32, coinName string) {
		wallet49, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+49), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex))

		var address, addressSegWit, wifPrivateKey, addressSegWitNative, addressTaproot string

		address, err = wallet49.GetAddress()
		if err != nil {
			panic(err)
		}

		addressSegWit, err = wallet49.GetKey().AddressP2WPKHInP2SH()
		if err != nil {
			panic(err)
		}

		wifPrivateKey, err = wallet49.GetKey().PrivateWIF(true)
		if err != nil {
			panic(err)
		}

		privateKeyHex := wallet49.GetKey().PrivateHex()

		pubKeyHex := wallet49.GetKey().PublicHex(true)

		fmt.Printf("%s Index %d:\n  privateKeyWIF:%s\n  privateKeyHex:%s\n  publicKey:%s\n  address:%s\n  addressSegWit:%s\n----------\n",
			coinName+"-Nested SegWit[P2SH]", addressIndex, wifPrivateKey, privateKeyHex, pubKeyHex, address, addressSegWit)

		//

		wallet84, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+84), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex))

		addressSegWitNative, err = wallet84.GetKey().AddressP2WPKH()
		if err != nil {
			panic(err)
		}

		wifPrivateKey, err = wallet84.GetKey().PrivateWIF(true)
		if err != nil {
			panic(err)
		}

		privateKeyHex = wallet84.GetKey().PrivateHex()

		pubKeyHex = wallet84.GetKey().PublicHex(true)
		fmt.Printf("%s Index %d:\n  privateKeyWIF:%s\n  privateKeyHex:%s\n  publicKey:%s\n  addressSegWitNative:%s\n----------\n",
			coinName+"-Native SegWit[Bech32]", addressIndex, wifPrivateKey, privateKeyHex, pubKeyHex, addressSegWitNative)

		//

		wallet86, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+86), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex))

		addressTaproot, err = wallet86.GetKey().AddressTapRoot()
		if err != nil {
			panic(err)
		}

		wifPrivateKey, err = wallet86.GetKey().PrivateWIF(true)
		if err != nil {
			panic(err)
		}

		privateKeyHex = wallet86.GetKey().PrivateHex()

		pubKeyHex = wallet86.GetKey().PublicHex(true)

		fmt.Printf("%s Index %d:\n  privateKeyWIF:%s\n  privateKeyHex:%s\n  publicKey:%s\n  addressTaproot:%s\n----------\n",
			coinName+"-Taproot", addressIndex, wifPrivateKey, privateKeyHex, pubKeyHex, addressTaproot)
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
		fnPrintBTC(idx, hdwallet.BTC, "BTC")
	}

	for idx := uint32(0); idx < 4; idx++ {
		fnPrintBTC(idx, hdwallet.BTCTestnet, "BTCTestnet")
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
