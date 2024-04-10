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

type AccountDetails struct {
	database.Account

	PrivateKey string
	Address    string
	Path       string
}

func (f AccountDetails) AsString(insecure bool) string {
	msg := fmt.Sprintf("Account Name: %s\n", f.Name)
	msg += fmt.Sprintf("Account Type: %s\n", f.Type)
	switch f.Type {
	case MnemonicType:
		if insecure {
			msg += fmt.Sprintf("Mnemonic: %s\n", f.Value)
			msg += fmt.Sprintf("Passphrase: %s\n", f.Passphrase)
			msg += fmt.Sprintf("Private key: %s\n", f.PrivateKey)
		}
		msg += fmt.Sprintf("Path Format: %s\n", f.PathFormat)
		msg += fmt.Sprintf("Path: %s\n", f.Path)
	case PrivateKeyType:
		if insecure {
			msg += fmt.Sprintf("Private Key: %s\n", f.Value)
		}
	default:
		return "invalid account type"
	}
	msg += fmt.Sprintf("Address: %s\n", f.Address)
	msg += fmt.Sprintf("Is Current: %v\n", f.Current)
	msg += fmt.Sprintf("Current Index: %d\n", f.CurrentIndex)

	return msg
}

func AccountToDetails(account *database.Account) (*AccountDetails, error) {
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

	fullAccount := AccountDetails{
		Account:    *account,
		PrivateKey: privateKey,
		Address:    address,
		Path:       path,
	}

	return &fullAccount, nil
}
