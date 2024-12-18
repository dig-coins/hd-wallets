package hdwallets

import (
	"github.com/dig-coins/hd-wallets/helpers/ethhelper"
)

func init() {
	coins[ETH] = newETH
}

type eth struct {
	name   string
	symbol string
	key    *Key

	// eth token
	contract string
}

func newETH(key *Key) Wallet {
	return &eth{
		name:   "Ethereum",
		symbol: "ETH",
		key:    key,
	}
}

func (c *eth) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *eth) GetName() string {
	return c.name
}

func (c *eth) GetSymbol() string {
	return c.symbol
}

func (c *eth) GetKey() *Key {
	return c.key
}

func (c *eth) GetAddress() (string, error) {
	return ethhelper.PubkeyToAddress(*c.key.PublicECDSA).Hex(), nil
}

func (c *eth) GetPrivateKey() (string, error) {
	return c.key.PrivateHex(), nil
}
