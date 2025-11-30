package hdwallets_test

import (
	"errors"
	"os"
	"testing"

	hdwallet "github.com/dig-coins/hd-wallets"
	"github.com/stretchr/testify/assert"
)

const (
	testMnemonic = "space list vivid country pledge arrest aspect busy sight column dinosaur parrot"
)

func TestHDWallet(t *testing.T) {
	seed, err := hdwallet.NewSeed(testMnemonic, "", "")
	assert.NoError(t, err)

	master, err := hdwallet.NewKey(
		hdwallet.Seed(seed),
	)
	assert.NoError(t, err)

	wallet, err := master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+44), hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(0))
	assert.NoError(t, err)

	address, err := wallet.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, "1Hm3vkMceAHG6sEJSsY4MMSy9kfEeaGH2W", address)

	addressP2WPKH, err := wallet.GetKey().AddressP2WPKH()
	assert.NoError(t, err)
	assert.Equal(t, "bc1qkltvypavrwu76525kxhmlyxnlx7tyle84q6za6", addressP2WPKH)

	addressP2WPKHInP2SH, err := wallet.GetKey().AddressP2WPKHInP2SH()
	assert.NoError(t, err)
	assert.Equal(t, "3GaGdWzMuyW11PN7CYANcf5JmQLnxUXoTK", addressP2WPKHInP2SH)

	wif, err := wallet.GetKey().PrivateWIF(true)
	assert.NoError(t, err)
	assert.Equal(t, "L1hF8A2UdHMac23TPzPghcB3VwSDTAsExdhYm8FuzRaJfB9mbdeN", wif)

	// addressP2WPKHInP2SH的特别说明:这个隔离见证的地址，是属于当前wif私钥的（默认bip44）。
	// 假设你是用生成的助记词导入到imtoken中，对应的隔离见证地址不是这个。
	// 若想和imtoken一致，请在 master.GetWallet 时传入 hdwallet.ZeroQuote+49 （即bip49）得到的隔离见证地址和对应私钥即可
	wallet, err = master.GetWallet(hdwallet.Purpose(hdwallet.ZeroQuote+49), hdwallet.CoinType(hdwallet.BTC), hdwallet.AddressIndex(0))
	assert.NoError(t, err)

	address, err = wallet.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, "13o5buLARR7skNo8gL9JMwStueco7zzPgb", address)

	addressP2WPKH, err = wallet.GetKey().AddressP2WPKH()
	assert.NoError(t, err)
	assert.Equal(t, "bc1qr6nn5fa0akf7agnksud3tmkwmhr2syj4vj9cez", addressP2WPKH)

	addressP2WPKHInP2SH, err = wallet.GetKey().AddressP2WPKHInP2SH()
	assert.NoError(t, err)
	assert.Equal(t, "3FmiCxg2NhCvXwUBNNNtA4DJLKPxa8oPnz", addressP2WPKHInP2SH)

	wif, err = wallet.GetKey().PrivateWIF(true)
	assert.NoError(t, err)
	assert.Equal(t, "L3s4oBZc7HmAzimnohrhx1HY2cqogJaE2DCyn8YtYdETW6usqH2C", wif)

	wallet, err = master.GetWallet(hdwallet.CoinType(hdwallet.BCH))
	assert.NoError(t, err)

	address, err = wallet.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, "1CMBLeHH3KYBYP8p4tStGcHh3AdkWqwhrX", address)

	address, err = wallet.GetKey().AddressBCH()
	assert.NoError(t, err)
	assert.Equal(t, "bitcoincash:1CMBLeHH3KYBYP8p4tStGcHh3AdkWqwhrX", address)

	wallet, err = master.GetWallet(hdwallet.CoinType(hdwallet.LTC))
	assert.NoError(t, err)

	address, err = wallet.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, "LgjbPWQBw76u8Lnfk3U4Nuu47i8rxVqjer", address)

	wallet, err = master.GetWallet(hdwallet.CoinType(hdwallet.DOGE))
	assert.NoError(t, err)

	address, err = wallet.GetAddress()
	assert.NoError(t, err)
	assert.Equal(t, "DAiYHWjG96tVMvbYCdBWwQToBByW5cib3C", address)

	wallet, err = master.GetWallet(hdwallet.CoinType(hdwallet.BNB))
	assert.NoError(t, err)

	publicKey := wallet.GetKey().PublicHex(false)
	assert.Equal(t, "0430af20343bf2ec557b8f27fb656e98cf7d970a4172fc0b3961eda7f25cd4f48daa064c50e746ee919265d971ac0ff50680534aeb2b6a5ae096eccf77079a876b", publicKey)

	privateKey := wallet.GetKey().PrivateHex()
	assert.Equal(t, "829a3e28d0425113b99c922b620a7a6b8a48f61d6d355d448668d1f4b7ce9eea", privateKey)

	address, err = wallet.GetKey().AddressBNB(hdwallet.MAINNET)
	assert.NoError(t, err)
	assert.Equal(t, "bnb1fgtq995aa099cau967e7h906srj94u4e6efl6t", address)
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
	mnemonic := os.Getenv("MNEMONIC")
	if mnemonic == "" {
		t.SkipNow()
	}

	seed, _ := hdwallet.NewSeed(mnemonic, "", "")

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
