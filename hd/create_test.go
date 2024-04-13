package hd

import "testing"

// go test -count=1 -v methd

func TestCreateMnemonic(t *testing.T) {
	mnemonic, err := CreateMnemonic(12)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("mnemonic: %s\n", mnemonic)
}
