package hdwallets

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

func GetFingerprint(compressedPubKey []byte) uint32 {
	sha256Hash := sha256.Sum256(compressedPubKey)

	ripeMd160Hasher := ripemd160.New()
	ripeMd160Hasher.Write(sha256Hash[:])
	fingerprintBytes := ripeMd160Hasher.Sum(nil)

	return binary.BigEndian.Uint32(fingerprintBytes[:4])
}

func GetFingerprintString(compressedPubKey []byte) string {
	return fmt.Sprintf("0x%08x", GetFingerprint(compressedPubKey))
}
