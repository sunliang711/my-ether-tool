package utils

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	mnemonic := "hello led"
	passphrase := "1234"

	encrypted := Encrypt(passphrase, mnemonic)
	t.Logf("encrypted: '%v'", encrypted)

	decrypted := Decrypt(passphrase+"cc", encrypted)
	t.Logf("decrypted: '%v'", decrypted)
}
