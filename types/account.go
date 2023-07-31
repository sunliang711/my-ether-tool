package types

import (
	"errors"
	"fmt"
	"my-ether-tool/database"
	"my-ether-tool/hd"
	"strings"
)

const (
	MnemonicType   = "mnemonic"
	PrivateKeyType = "private key"

	DefaultHDPath = "m/44'/60'/0'/0/x"
)

type FullAccount struct {
	database.Account

	PrivateKey string
	Address    string
	Path       string
}

func AccountToFullAccount(account *database.Account) (*FullAccount, error) {
	var privateKey string
	var address string
	var path string

	switch account.Type {
	case MnemonicType:
		//derive
		path = strings.Replace(account.PathFormat, "x", fmt.Sprintf("%d", account.CurrentIndex), 1)
		out, err := hd.Derive(account.Value, account.Passphrase, path, uint(account.CurrentIndex), 1)
		if err != nil {
			return nil, err
		}
		if len(out.Keys) != 1 {
			return nil, errors.New("derive menmonic error: length not 1")
		}
		privateKey = out.Keys[0].PrivateKey
		address = out.Keys[0].EthereumAddress
	case PrivateKeyType:
		privateKey = account.Value
		pubkey, err := hd.PrivateKeyToPublicKey(privateKey)
		if err != nil {
			return nil, err
		}
		address, err = hd.PubkeyToAddress(pubkey)
	default:
		return nil, errors.New("invalid account type")
	}

	fullAccount := FullAccount{
		Account:    *account,
		PrivateKey: privateKey,
		Address:    address,
		Path:       path,
	}

	return &fullAccount, nil
}
