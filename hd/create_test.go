package hd

import "testing"

// go test -count=1 -v my-ether-tool/hd

func TestCreateMnemonic(t *testing.T) {
	mnemonic, err := CreateMnemonic(12)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("mnemonic: %s\n", mnemonic)
}
