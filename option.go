package hdwallet

import (
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/chaincfg"
)

const (
	PathLevelAuto    = -1
	PathLevelPurpose = iota - 1
	PathLevelCoinType
	PathLevelAccount
	PathLevelChange
	PathLevelAddressIndex
)

// default options
var (
	DefaultParams = &BTCParams

	// master key options

	DefaultPassword = ""
	DefaultLanguage = ""

	// child key options

	DefaultPurpose      = ZeroQuote + 44
	DefaultCoinType     = Zero
	DefaultAccount      = ZeroQuote
	DefaultChange       = Zero
	DefaultAddressIndex = Zero
)

// Option of key
type Option func(*Options)

// Options of key
type Options struct {
	Params *chaincfg.Params

	// master key options
	Mnemonic string
	Password string
	Language string
	Seed     []byte

	// child key options
	Purpose      uint32
	CoinType     uint32
	Account      uint32
	Change       uint32
	AddressIndex uint32

	// level for derive from public key
	MinPathLevel int // [
	MaxPathLevel int // ]
}

func newOptions(opts ...Option) *Options {
	opt := &Options{
		Params:       DefaultParams,
		Password:     DefaultPassword,
		Language:     DefaultLanguage,
		Purpose:      DefaultPurpose,
		CoinType:     DefaultCoinType,
		Account:      DefaultAccount,
		Change:       DefaultChange,
		AddressIndex: DefaultAddressIndex,
		MinPathLevel: PathLevelPurpose,
		MaxPathLevel: PathLevelAddressIndex,
	}

	for _, o := range opts {
		o(opt)
	}

	if opt.MinPathLevel <= 0 {
		opt.MinPathLevel = PathLevelPurpose
	}

	if opt.MaxPathLevel <= 0 {
		opt.MaxPathLevel = PathLevelAddressIndex
	}

	return opt
}

// GetPath return path in bip44 style
func (o *Options) GetPath() []uint32 {
	return []uint32{
		o.Purpose,
		o.CoinType,
		o.Account,
		o.Change,
		o.AddressIndex,
	}
}

// Params set to options
func Params(p *chaincfg.Params) Option {
	return func(o *Options) {
		o.Params = p
	}
}

// Mnemonic set to options
func Mnemonic(m string) Option {
	return func(o *Options) {
		o.Mnemonic = m
	}
}

// Password set to options
func Password(p string) Option {
	return func(o *Options) {
		o.Password = p
	}
}

// Language set to options
func Language(l string) Option {
	return func(o *Options) {
		o.Language = l
	}
}

// Seed set to options
func Seed(s []byte) Option {
	return func(o *Options) {
		o.Seed = s
	}
}

// Purpose set to options
func Purpose(p uint32) Option {
	return func(o *Options) {
		o.Purpose = p
	}
}

// CoinType set to options
func CoinType(c uint32) Option {
	return func(o *Options) {
		o.CoinType = c
	}
}

// Account set to options
func Account(a uint32) Option {
	return func(o *Options) {
		o.Account = a
	}
}

// Change set to options
func Change(c uint32) Option {
	return func(o *Options) {
		o.Change = c
	}
}

// AddressIndex set to options
func AddressIndex(a uint32) Option {
	return func(o *Options) {
		o.AddressIndex = a
	}
}

func MinLevel(c int) Option {
	return func(o *Options) {
		o.MinPathLevel = c
	}
}

func MaxLevel(c int) Option {
	return func(o *Options) {
		o.MaxPathLevel = c
	}
}

// Path set to options
// example: m/44'/0'/0'/0/0
// example: m/Purpose'/CoinType'/Account'/Change/AddressIndex
func Path(path string) Option {
	return func(o *Options) {
		path = strings.TrimPrefix(path, "m/")
		paths := strings.Split(path, "/")
		if len(paths) != 5 {
			return
		}
		o.Purpose = PathNumber(paths[0])
		o.CoinType = PathNumber(paths[1])
		o.Account = PathNumber(paths[2])
		o.Change = PathNumber(paths[3])
		o.AddressIndex = PathNumber(paths[4])
	}
}

// PathNumber 44' => 0x80000000 + 44
func PathNumber(str string) uint32 {
	num64, _ := strconv.ParseInt(strings.TrimSuffix(str, "'"), 10, 64)
	num := uint32(num64)
	if strings.HasSuffix(str, "'") {
		num += ZeroQuote
	}
	return num
}
