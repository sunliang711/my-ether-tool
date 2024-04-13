package contract

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestPack(t *testing.T) {
	address := common.HexToAddress("bc3085e76Ba7e83d73BAb362C5cdc79EF2AF3")
	t.Logf("address: %v", address)

	decoded, err := hex.DecodeString("0x1234")
	if err != nil {
		t.Logf("decode error: %v", err)
	} else {
		t.Logf("decoded: %v", decoded)
	}
}
