package mnemonicutils

import (
	"log"
	"os"

	"github.com/howeyc/gopass"
)

func TypeString(prompt string) string {
start:
	v, err := gopass.GetPasswdPrompt(prompt+": ", true, os.Stdin, os.Stdout)

	if err != nil {
		log.Fatalln("no valid user password", err)
	}

	v2, err := gopass.GetPasswdPrompt(prompt+" again: ", true, os.Stdin, os.Stdout)
	if err != nil {
		log.Fatalln("no valid user password", err)
	}

	s := string(v)
	if s != string(v2) {
		log.Println("not match, try again ...")

		goto start
	}

	log.Println(prompt + ": " + s)

	return s
}
