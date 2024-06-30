package bnbhelper

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160" // nolint: staticcheck
)

// AccAddress a wrapper around bytes meant to represent an account address.
// When marshaled to a string or JSON, it uses Bech32.
type AccAddress []byte

func (bz AccAddress) Bytes() []byte {
	return bz
}

type HexBytes []byte
type Address = HexBytes

// PubKeyToAddress returns a Bitcoin style addresses: RIPEMD160(SHA256(pubkey))
func PubKeyToAddress(d []byte) Address {
	hasherSHA256 := sha256.New()
	hasherSHA256.Write(d[:]) // does not error
	sha := hasherSHA256.Sum(nil)

	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha) // does not error
	return hasherRIPEMD160.Sum(nil)
}
