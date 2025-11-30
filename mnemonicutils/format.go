package mnemonicutils

import (
	"fmt"
	"strings"

	"github.com/tyler-smith/go-bip39/wordlists"
	"golang.org/x/exp/slices"
)

func MnemonicFormat(mnemonics string) (string, error) {
	mnemonics = strings.Trim(mnemonics, "\r\n\t ")

	ps := strings.Split(mnemonics, " ")

	for i := 0; i < len(ps); {
		ss := strings.Trim(ps[i], "\r\n\t ")
		if len(ss) == 0 {
			ps = append(ps[:i], ps[i+1:]...)
		} else {
			if !slices.Contains(wordlists.English, ss) {
				return "", fmt.Errorf("invalid %s onmnemonic", ss)
			}

			i++
		}
	}

	if len(ps) != 12 {
		return "", fmt.Errorf("invalid mnemonic number %d", len(ps))
	}

	return strings.Join(ps, " "), nil
}
