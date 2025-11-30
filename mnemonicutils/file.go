package mnemonicutils

import (
	"os"

	"github.com/GizmoVault/gotools/crypt/edfile"
	"gopkg.in/yaml.v3"
)

func GetMnemonicFromMnemonicFile(f string) (mnemonic string, err error) {
	d, err := os.ReadFile(f)
	if err != nil {
		return
	}

	mnemonic, err = MnemonicFormat(string(d))

	return
}

func GetMnemonicFromSecMnemonicFile(f, key string) (mnemonic string, err error) {
	d, err := edfile.ReadSecFile(f, key)
	if err != nil {
		return
	}

	mnemonic, err = MnemonicFormat(string(d))

	return
}

func WriteSecMnemonicFile(mnemonic, f, key string) error {
	return edfile.WriteSecFile(f, key, []byte(mnemonic))
}

func GetMnemonicsFromSecMnemonicFile(f, key string) (mnemonics []string, err error) {
	d, err := edfile.ReadSecFile(f, key)
	if err != nil {
		return
	}

	var ss []string

	err = yaml.Unmarshal(d, &ss)
	if err != nil {
		return
	}

	mnemonics = make([]string, 0, len(ss))

	var mnemonic string

	for _, s := range ss {
		mnemonic, err = MnemonicFormat(s)
		if err != nil {
			return
		}

		mnemonics = append(mnemonics, mnemonic)
	}

	return
}
