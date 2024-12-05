package hdwallets

func init() {
	coins[BTCRegTest] = newBTCRegTest
}

type btcRegTest struct {
	name   string
	symbol string
	key    *Key
}

func newBTCRegTest(key *Key) Wallet {
	key.Opt.Params = &BTCRegTestParams
	return &btcRegTest{
		name:   "Bitcoin RegTest",
		symbol: "BTCRegTest",
		key:    key,
	}
}

func (c *btcRegTest) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *btcRegTest) GetName() string {
	return c.name
}

func (c *btcRegTest) GetSymbol() string {
	return c.symbol
}

func (c *btcRegTest) GetKey() *Key {
	return c.key
}

func (c *btcRegTest) GetAddress() (string, error) {
	return c.key.AddressBTC()
}

func (c *btcRegTest) GetPrivateKey() (string, error) {
	return c.key.PrivateWIF(true)
}
