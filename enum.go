package hdwallets

// mnemonic language
const (
	English            = "english"
	ChineseSimplified  = "chinese_simplified"
	ChineseTraditional = "chinese_traditional"
)

// zero is default of uint32
const (
	Zero      uint32 = 0
	ZeroQuote uint32 = 0x80000000
	BTCToken  uint32 = 0x10000000
	ETHToken  uint32 = 0x20000000
)

// wallet coin type

const (
	// PurposeType44
	// BIP-44
	// Legacy
	// P2PKH,P2SH
	// 1... | [m|n]...
	PurposeType44 = ZeroQuote + 44

	// PurposeType49
	// BIP-49
	// Nested SegWit
	// P2SH-P2WPKH
	// 3... | 2...
	PurposeType49 = ZeroQuote + 49

	// PurposeType84
	// BIP-84
	// Native SegWit
	// P2WPKH
	// bc1q... | tb1q...
	PurposeType84 = ZeroQuote + 84

	// PurposeType86
	// BIP-86
	// Native Taproot
	// P2TR
	// bc1p... | tb1p...
	PurposeType86 = ZeroQuote + 86
)

// wallet account

func HardenAccount(account uint32) uint32 {
	return ZeroQuote + account
}

// wallet purpose type from bip44
const (
	// https://github.com/satoshilabs/slips/blob/master/slip-0044.md#registered-coin-types

	BTC        = ZeroQuote + 0
	BTCTestnet = ZeroQuote + 1
	LTC        = ZeroQuote + 2
	DOGE       = ZeroQuote + 3
	DASH       = ZeroQuote + 5
	ETH        = ZeroQuote + 60
	ETC        = ZeroQuote + 61
	BCH        = ZeroQuote + 145
	QTUM       = ZeroQuote + 2301
	TRX        = ZeroQuote + 195
	BNB        = ZeroQuote + 714
	FIL        = ZeroQuote + 461

	//

	BTCRegTest = ZeroQuote + 6000

	// btc token

	USDT = BTCToken + 1

	// eth token

	IOST = ETHToken + 1
	USDC = ETHToken + 2
)

// network
const (
	MAINNET = "mainnet"
	TESTNET = "testnet"
)

var coinTypes = map[uint32]uint32{
	USDT:       BTC,
	IOST:       ETH,
	USDC:       ETH,
	TRX:        TRX,
	BNB:        BNB,
	FIL:        FIL,
	BTCRegTest: BTCTestnet,
}
