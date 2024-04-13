package hd

import "testing"

// go test -count=1 -v methd
func TestDerive(t *testing.T) {
	mnemonic := "length toddler champion supply hockey orange oil satisfy wisdom hedgehog scene nominee radar cactus immune"
	path := "m/44'/60'/2/x"
	passphrase := "12"
	start := uint(3)
	count := uint(5)
	output, err := Derive(mnemonic, passphrase, path, start, count)
	if err != nil {
		t.Fatalf("derive error: %s", err)
	}
	t.Logf("derive result: %v", output.String())

}
