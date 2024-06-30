package hdwallet

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"github.com/dig-coins/hd-wallet/helpers/ethhelper"
	"math/big"

	"github.com/shengdoushi/base58"
)

func init() {
	coins[TRX] = newTRX
}

type trx struct {
	name   string
	symbol string
	key    *Key

	// trc20 token
	//contract string
}

func newTRX(key *Key) Wallet {
	return &trx{
		name:   "Tron",
		symbol: "TRX",
		key:    key,
	}
}

func (c *trx) GetType() uint32 {
	return c.key.Opt.CoinType
}

func (c *trx) GetName() string {
	return c.name
}

func (c *trx) GetSymbol() string {
	return c.symbol
}

func (c *trx) GetKey() *Key {
	return c.key
}

func (c *trx) GetAddress() (string, error) {
	a := PubkeyToAddress(*c.key.PublicECDSA)

	return a.String(), nil
}

func (c *trx) GetPrivateKey() (string, error) {
	return c.key.PrivateHex(), nil
}

//
//
//

const ( // TronBytePrefix is the hex prefix to address
	TronBytePrefix = byte(0x41)
)

// Address represents the 21 byte address of an Tron account.
type Address []byte

// Bytes get bytes from address
func (a Address) Bytes() []byte {
	return a[:]
}

func Encode(input []byte) string {
	return base58.Encode(input, base58.BitcoinAlphabet)
}

func EncodeCheck(input []byte) string {
	h256h0 := sha256.New()
	h256h0.Write(input)
	h0 := h256h0.Sum(nil)

	h256h1 := sha256.New()
	h256h1.Write(h0)
	h1 := h256h1.Sum(nil)

	inputCheck := input
	inputCheck = append(inputCheck, h1[:4]...)

	return Encode(inputCheck)
}

// String implements fmt.Stringer.
func (a Address) String() string {
	if a[0] == 0 {
		return new(big.Int).SetBytes(a.Bytes()).String()
	}
	return EncodeCheck(a.Bytes())
}

// PubkeyToAddress returns address from ecdsa public key
func PubkeyToAddress(p ecdsa.PublicKey) Address {
	address := ethhelper.PubkeyToAddress(p)

	addressTron := make([]byte, 0)
	addressTron = append(addressTron, TronBytePrefix)
	addressTron = append(addressTron, address.Bytes()...)
	return addressTron
}
