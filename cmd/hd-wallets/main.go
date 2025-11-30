package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/GizmoVault/gotools/pathx"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	hdwallet "github.com/dig-coins/hd-wallets"
	"github.com/dig-coins/hd-wallets/mnemonicutils"
)

const (
	_mnemonicFile    = "mnemonic.txt"
	_secMnemonicFile = "sec-mnemonic.dat"
	_address1File    = "address1.txt"
)

func main() {
	var workDir, seedPassword, fileKey string

	var op int

	var deepLevel int

	flag.StringVar(&workDir, "w", ".", "working directory")
	flag.StringVar(&seedPassword, "p", "", "password of seed")
	flag.StringVar(&fileKey, "k", "", "key of sec file")
	flag.IntVar(&deepLevel, "d", 1, "key deep level")
	flag.IntVar(&op, "op", 0, `0: dump
1: generate
2: print mnemonic
3: check mnemonic and sec-mnemonic
4. generate by mnemonic file`)
	flag.Parse()

	if deepLevel <= 0 {
		deepLevel = 1
	}

	mnemonicFile := filepath.Join(workDir, _mnemonicFile)
	secMnemonicFile := filepath.Join(workDir, _secMnemonicFile)
	address1File := filepath.Join(workDir, _address1File)

	fnGetSeedPassword := func() string {
		if seedPassword == "<null>" {
			seedPassword = mnemonicutils.TypeString("input your mnemonic password")
		}

		return seedPassword
	}

	fnGetFileKey := func() string {
		if fileKey == "<null>" {
			fileKey = mnemonicutils.TypeString("input your mnemonic file key")
		}

		return fileKey
	}

	if op == 1 {
		fileNotExistsOrPanic(mnemonicFile)
		fileNotExistsOrPanic(secMnemonicFile)

		mnemonic, _ := hdwallet.NewMnemonic(12, "")

		mnemonicC, err := mnemonicutils.MnemonicFormat(mnemonic)
		if err != nil {
			log.Fatalln("invalid mnemonic", mnemonic, err)
		}

		if mnemonic != mnemonicC {
			log.Fatalln("invalid mnemonic format logic", mnemonic, mnemonicC)
		}

		err = mnemonicutils.WriteSecMnemonicFile(mnemonic, secMnemonicFile, fnGetFileKey())
		if err != nil {
			log.Fatalln("write mnemonic failed:", err)
		}

		return
	}

	if op == 2 {
		mnemonic, err := mnemonicutils.GetMnemonicFromSecMnemonicFile(secMnemonicFile, fnGetFileKey())
		if err != nil {
			log.Fatalln("read sec mnemonic file failed:", err)
		}

		log.Printf("mnemonic is %s\n", mnemonic)

		return
	}

	if op == 3 {
		mnemonic, err := mnemonicutils.GetMnemonicFromSecMnemonicFile(secMnemonicFile, fnGetFileKey())
		if err != nil {
			log.Fatalln("read sec mnemonic file failed:", err)
		}

		mnemonic2, err := mnemonicutils.GetMnemonicFromMnemonicFile(mnemonicFile)
		if err != nil {
			log.Fatalln("read mnemonic file failed:", err)
		}

		if mnemonic != mnemonic2 {
			log.Fatalln("mnemonic not match")
		}

		log.Printf("mnemonic matched: %s\n", mnemonic)

		return
	}

	if op == 4 {
		fileNotExistsOrPanic(secMnemonicFile)

		mnemonic, err := mnemonicutils.GetMnemonicFromMnemonicFile(mnemonicFile)
		if err != nil {
			log.Fatalln("read mnemonic file failed:", err)
		}

		err = mnemonicutils.WriteSecMnemonicFile(mnemonic, secMnemonicFile, fnGetFileKey())
		if err != nil {
			log.Fatalln("write mnemonic failed:", err)
		}

		return
	}

	//
	//
	//

	addressExists := true

	d, err := os.ReadFile(address1File)
	if err != nil {
		addressExists = false
	}

	address1 := strings.Trim(string(d), "\r\n\t ")

	var address1Matched bool

	mnemonic, err := mnemonicutils.GetMnemonicFromSecMnemonicFile(secMnemonicFile, fnGetFileKey())
	if err != nil {
		log.Fatalln("read sec mnemonic file failed:", err)
	}

	seed, err := hdwallet.NewSeed(mnemonic, fnGetSeedPassword(), "")
	if err != nil {
		panic(err)
	}

	master, err := hdwallet.NewKey(
		hdwallet.Seed(seed),
	)
	if err != nil {
		panic(err)
	}

	fnGetX := func(extendedKey *hdkeychain.ExtendedKey) (s string) {
		if !extendedKey.IsPrivate() {
			s += fmt.Sprintf("Derivation key[xPub]: %s\n", extendedKey.String())

			return
		}

		s += fmt.Sprintf("Derivation key[xPri]: %s\n", extendedKey.String())

		pubExtendedKey, e := extendedKey.Neuter()
		if e != nil {
			panic(e)
		}

		if pubExtendedKey == nil {
			s += "To Neuter Failed/n"

			return
		}

		s += fmt.Sprintf("Derivation key[xPub]: %s\n", pubExtendedKey.String())

		return
	}

	fnPrintBTC := func(addressIndex, coinType uint32, coinName string) {
		var address, addressSegWit, wifPrivateKey, addressSegWitNative, addressTaproot string

		wallet44, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+44), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex))
		wallet, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+44), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex),
			hdwallet.MaxLevel(hdwallet.PathLevelAccount))

		address, err = wallet44.GetAddress()
		if err != nil {
			panic(err)
		}

		if address == address1 {
			address1Matched = true
		}

		addressSegWit, err = wallet44.GetKey().AddressP2WPKHInP2SH()
		if err != nil {
			panic(err)
		}

		if addressSegWit == address1 {
			address1Matched = true
		}

		wifPrivateKey, err = wallet44.GetKey().PrivateWIF(true)
		if err != nil {
			panic(err)
		}

		privateKeyHex := wallet44.GetKey().PrivateHex()

		pubKeyHex := wallet44.GetKey().PublicHex(true)

		fmt.Printf("%s Index %d:\n  privateKeyWIF:%s\n  privateKeyHex:%s\n  publicKey:%s\n  address:%s\n  addressSegWit:%s\n xKey:%s\n----------\n",
			coinName+"-44", addressIndex, wifPrivateKey, privateKeyHex, pubKeyHex, address, addressSegWit, fnGetX(wallet.GetKey().Extended))
		//

		wallet49, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+49), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex))
		wallet, _ = master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+49), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex),
			hdwallet.MaxLevel(hdwallet.PathLevelAccount))

		address, err = wallet49.GetAddress()
		if err != nil {
			panic(err)
		}

		if address == address1 {
			address1Matched = true
		}

		addressSegWit, err = wallet49.GetKey().AddressP2WPKHInP2SH()
		if err != nil {
			panic(err)
		}

		if addressSegWit == address1 {
			address1Matched = true
		}

		wifPrivateKey, err = wallet49.GetKey().PrivateWIF(true)
		if err != nil {
			panic(err)
		}

		privateKeyHex = wallet49.GetKey().PrivateHex()

		pubKeyHex = wallet49.GetKey().PublicHex(true)

		fmt.Printf("%s Index %d:\n  privateKeyWIF:%s\n  privateKeyHex:%s\n  publicKey:%s\n  address:%s\n  addressSegWit:%s\n xKey:%s\n----------\n",
			coinName+"-49-Nested SegWit[P2SH]", addressIndex, wifPrivateKey, privateKeyHex, pubKeyHex, address, addressSegWit, fnGetX(wallet.GetKey().Extended))

		//

		wallet84, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+84), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex))
		wallet, _ = master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+84), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex),
			hdwallet.MaxLevel(hdwallet.PathLevelAccount))

		addressSegWitNative, err = wallet84.GetKey().AddressP2WPKH()
		if err != nil {
			panic(err)
		}

		if addressSegWitNative == address1 {
			address1Matched = true
		}

		wifPrivateKey, err = wallet84.GetKey().PrivateWIF(true)
		if err != nil {
			panic(err)
		}

		privateKeyHex = wallet84.GetKey().PrivateHex()

		pubKeyHex = wallet84.GetKey().PublicHex(true)
		fmt.Printf("%s Index %d:\n  privateKeyWIF:%s\n  privateKeyHex:%s\n  publicKey:%s\n  addressSegWitNative:%s\nxKey:%s\n----------\n",
			coinName+"-84-Native SegWit[Bech32]", addressIndex, wifPrivateKey, privateKeyHex, pubKeyHex, addressSegWitNative,
			fnGetX(wallet.GetKey().Extended))

		//

		wallet86, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+86), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex))
		wallet, _ = master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+86), hdwallet.CoinType(coinType), hdwallet.AddressIndex(addressIndex),
			hdwallet.MaxLevel(hdwallet.PathLevelAccount))

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

		fmt.Printf("%s Index %d:\n  privateKeyWIF:%s\n  privateKeyHex:%s\n  publicKey:%s\n  addressTaproot:%s\nxKey:%s\n----------\n",
			coinName+"-86-Taproot", addressIndex, wifPrivateKey, privateKeyHex, pubKeyHex, addressTaproot, fnGetX(wallet.GetKey().Extended))
	}

	for idx := uint32(0); idx < uint32(deepLevel); idx++ {
		fnPrintBTC(idx, hdwallet.BTC, "BTC")
	}

	for idx := uint32(0); idx < uint32(deepLevel); idx++ {
		fnPrintBTC(idx, hdwallet.BTCTestnet, "BTCTestnet")
	}

	for idx := uint32(0); idx < uint32(deepLevel); idx++ {
		fnPrintBTC(idx, hdwallet.BTCRegTest, "BTCRegTest")
	}

	if !addressExists {
		log.Println("address file not exists")
	} else if address1Matched {
		log.Println("address matched.")
	} else {
		log.Fatalln("address not matched!!!!")
	}
}

func fileNotExistsOrPanic(f string) {
	ok, err := pathx.IsPathExists(f)
	if err != nil {
		panic("access file " + f + " failed")
	}

	if ok {
		panic(f + " exists, abort!!")
	}
}
