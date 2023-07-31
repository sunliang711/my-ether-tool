package hd

import (
	"fmt"

	"github.com/tyler-smith/go-bip39"
)

const ENTROPY_BIT_SIZE_12 = 32 * 4
const ENTROPY_BIT_SIZE_15 = 32 * 5
const ENTROPY_BIT_SIZE_18 = 32 * 6
const ENTROPY_BIT_SIZE_21 = 32 * 8
const ENTROPY_BIT_SIZE_24 = 32 * 8

func CreateMnemonic(words uint8) (string, error) {

	bitSize := 0

	switch words {
	case 12:
		bitSize = ENTROPY_BIT_SIZE_12
	case 15:
		bitSize = ENTROPY_BIT_SIZE_15
	case 18:
		bitSize = ENTROPY_BIT_SIZE_18
	case 21:
		bitSize = ENTROPY_BIT_SIZE_21
	case 24:
		bitSize = ENTROPY_BIT_SIZE_24
	default:
		return "", fmt.Errorf("no support words: %d", words)
	}

	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return "", fmt.Errorf("new entropy error: %s", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("new mnemonic error: %s", err)
	}

	return mnemonic, nil

}
