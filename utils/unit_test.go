package utils

import "testing"

// go test -count=1 -v  my-ether-tool/utils -run 'TestParseUnits'
func TestParseUnits(t *testing.T) {
	r1, err := ParseUnits("1.2", UnitEth)
	if err != nil {
		t.Logf("error: %s\n", err)
	}
	t.Logf("r1: %s\n", r1)

	r2, err := ParseUnits("3.14", UnitGwei)
	if err != nil {
		t.Logf("error: %s\n", err)
	}
	t.Logf("r2: %s\n", r2)

}

// go test -count=1 -v  my-ether-tool/utils -run 'TestStringMul'
func TestStringMul(t *testing.T) {
	r1, err := StringMul("1.2", "1000000")
	if err != nil {
		t.Fatalf("error: %s", err)
	}

	t.Logf("r1: %s\n", r1.String())

	r2, err := StringMul("1.23456", "100")
	if err != nil {
		t.Fatalf("error: %s\n", err)
	}
	t.Logf("r2: %s\n", r2.String())
}

// go test -count=1 -v  my-ether-tool/utils -run 'TestStringDiv'
func TestStringDiv(t *testing.T) {
	r1, err := StringDiv("1200", "1000000")
	if err != nil {
		t.Fatalf("error: %s\n", err)
	}
	t.Logf("r1: %s\n", r1.String())

	r2, err := StringDiv("10000000", "1_000_000")
	if err != nil {
		t.Fatalf("error: %s\n", err)
	}
	t.Logf("r2: %s\n", r2.String())
}
