package hdwallet_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	hdwallet "github.com/dig-coins/hd-wallets"
	"github.com/stretchr/testify/assert"
)

func TestHDWallet(t *testing.T) {
	var ksOption hdwallet.Option
	// 128: 12 phrases
	// 256: 24 phrases
	/*
		mnemonic, _ := hdwallet.NewMnemonic(12, "")
		ksOption = hdwallet.Mnemonic(mnemonic)
		fmt.Println("助记词：", mnemonic)
	*/

	seed, _ := hdwallet.NewSeed(os.Getenv("mnemonic"), "", "")
	ksOption = hdwallet.Seed(seed)

	master, err := hdwallet.NewKey(
		ksOption,
	)
	if err != nil {
		panic(err)
	}

	wallet, _ := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+44 /*44*/), hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(0))
	address, _ := wallet.GetAddress()                               // 1AwEPfoojHnKrhgt1vfuZAhrvPrmz7Rh44
	addressP2WPKH, _ := wallet.GetKey().AddressP2WPKH()             // bc1qdnavt2xqvmc58ktff7rhvtn9s62zylp5lh5sry
	addressP2WPKHInP2SH, _ := wallet.GetKey().AddressP2WPKHInP2SH() // 39vtu9kWfGigXTKMMyc8tds7q36JBCTEHg
	// addressP2WPKHInP2SH的特别说明:这个隔离见证的地址，是属于当前wif私钥的（默认bip44）。
	// 假设你是用生成的助记词导入到imtoken中，对应的隔离见证地址不是这个。
	// 若想和imtoken一致，请在 master.GetWallet 时传入 hdwallet.ZeroQuote+49 （即bip49）得到的隔离见证地址和对应私钥即可
	btcwif, err := wallet.GetKey().PrivateWIF(true)
	if err != nil {
		panic(err)
	}

	fmt.Println("BTC私钥：", btcwif)
	fmt.Println("BTC: ", address, addressP2WPKH, addressP2WPKHInP2SH)

	// BCH: 1CSBT18sjcCwLCpmnnyN5iqLc46Qx7CC91
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.BCH))
	address, _ = wallet.GetAddress()
	addressBCH, _ := wallet.GetKey().AddressBCH()
	fmt.Println("BCH: ", address, addressBCH)

	// LTC: LLCaMFT8AKjDTvz1Ju8JoyYXxuug4PZZmS
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.LTC))
	address, _ = wallet.GetAddress()
	fmt.Println("LTC: ", address)

	// DOGE: DHLA3rJcCjG2tQwvnmoJzD5Ej7dBTQqhHK
	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.DOGE))
	address, _ = wallet.GetAddress()
	fmt.Println("DOGE:", address)

	wallet, _ = master.GetWallet(hdwallet.CoinType(hdwallet.BNB))
	wallet.GetKey().PublicHex(false)
	fmt.Println("BNB私钥：", wallet.GetKey().PrivateHex())
	address, _ = wallet.GetKey().AddressBNB(hdwallet.MAINNET)
	fmt.Println("BNB: ", address)
}

func dumpWallet(t *testing.T, wallet hdwallet.Wallet) {
	address, err := wallet.GetAddress()
	assert.Nil(t, err)

	privateKey, err := wallet.GetPrivateKey()
	if !errors.Is(err, hdwallet.ErrPrivateKeyNotExists) {
		assert.Nil(t, err)
	}

	t.Logf("address: %s\n", address)
	t.Logf("privateKey: %s\n", privateKey)
	t.Logf("publicKey: %s\n", wallet.GetKey().PublicHex(true))
}

func TestHDWallet2(t *testing.T) {
	seed, _ := hdwallet.NewSeed(os.Getenv("mnemonic"), "", "")

	master, err := hdwallet.NewKey(
		hdwallet.Seed(seed),
	)
	assert.Nil(t, err)

	{
		wallet, e := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+49 /*44*/),
			hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(1))
		assert.Nil(t, e)
		dumpWallet(t, wallet)
	}

	wallet, err := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+49 /*44*/),
		hdwallet.CoinType(hdwallet.BTC), hdwallet.MaxLevel(hdwallet.PathLevelAccount))
	assert.Nil(t, err)

	wallet1, err := wallet.GetKey().GetWallet(hdwallet.AddressIndex(1), hdwallet.CoinType(hdwallet.BTC),
		hdwallet.MinLevel(hdwallet.PathLevelChange), hdwallet.MaxLevel(hdwallet.PathLevelAuto))
	assert.Nil(t, err)

	dumpWallet(t, wallet1)

	extPub, err := wallet.GetKey().Extended.Neuter()
	assert.Nil(t, err)

	extPubS := extPub.String()
	t.Log(extPubS)

	master2, err := hdwallet.NewKeyFromString(extPubS)
	assert.Nil(t, err)

	wallet2, err := master2.GetWallet(hdwallet.AddressIndex(1), hdwallet.CoinType(hdwallet.BTC), hdwallet.MinLevel(hdwallet.PathLevelChange),
		hdwallet.MaxLevel(hdwallet.PathLevelAuto))
	assert.Nil(t, err)
	dumpWallet(t, wallet2)
}
