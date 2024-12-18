package hdwallets

import (
	"github.com/dig-coins/hd-wallets/helpers/ethhelper"
)

func init() {
	coins[ETC] = newETC
}

type etc struct {
	name   string
	symbol string
	key    *Key
}

func newETC(key *Key) Wallet {
	return &etc{
		name:   "Ethereum Classic",
		symbol: "ETC",
		key:    key,
	}
}

func (c *etc) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *etc) GetName() string {
	return c.name
}

func (c *etc) GetSymbol() string {
	return c.symbol
}

func (c *etc) GetKey() *Key {
	return c.key
}

func (c *etc) GetAddress() (string, error) {
	return ethhelper.PubkeyToAddress(*c.key.PublicECDSA).Hex(), nil
}

func (c *etc) GetPrivateKey() (string, error) {
	return c.key.PrivateHex(), nil
}
