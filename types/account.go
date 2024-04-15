package types

import (
	"errors"
	"fmt"
	database "met/database"
	hd "met/hd"
	"strings"
)

const (
	MnemonicType   = "mnemonic"
	PrivateKeyType = "private key"

	DefaultHDPath = "m/44'/60'/0'/0/x"
)

type AccountDetails struct {
	database.Account

	privateKey string
	address    string
	Path       string
}

func (f *AccountDetails) Address() (string, error) {
	if f.Encrypted {
		return "", fmt.Errorf("account: %v locked", f.Name)
	}
	return f.address, nil
}

func (f *AccountDetails) PrivateKey() (string, error) {
	if f.Encrypted {
		return "", fmt.Errorf("account: %v locked", f.Name)
	}
	return f.privateKey, nil

}

func (f AccountDetails) AsString(insecure bool) string {
	var msgArray []string
	msgArray = append(msgArray, fmt.Sprintf("Account Name: %s\n", f.Name))
	msgArray = append(msgArray, fmt.Sprintf("Account Type: %s\n", f.Type))
	switch f.Type {
	case MnemonicType:
		if insecure {
			if f.Encrypted {
				msgArray = append(msgArray, fmt.Sprintf("Mnemonic: locked\n"))
				msgArray = append(msgArray, fmt.Sprintf("Passphrase: %s\n", f.Passphrase))
				msgArray = append(msgArray, fmt.Sprintf("Private key: locked\n"))
			} else {
				msgArray = append(msgArray, fmt.Sprintf("Mnemonic: %s\n", f.Value))
				msgArray = append(msgArray, fmt.Sprintf("Passphrase: %s\n", f.Passphrase))
				msgArray = append(msgArray, fmt.Sprintf("Private key: %s\n", f.privateKey))
			}

		}
		msgArray = append(msgArray, fmt.Sprintf("Path Format: %s\n", f.PathFormat))
		msgArray = append(msgArray, fmt.Sprintf("Path: %s\n", f.Path))

	case PrivateKeyType:
		if insecure {
			if f.Encrypted {
				msgArray = append(msgArray, fmt.Sprintf("Private Key: locked\n"))
			} else {
				msgArray = append(msgArray, fmt.Sprintf("Private Key: %s\n", f.Value))

			}
		}
	default:
		return "invalid account type"
	}

	if f.Encrypted {
		msgArray = append(msgArray, fmt.Sprintf("Address: locked\n"))
		msgArray = append(msgArray, fmt.Sprintf("Is Current: %v\n", f.Current))
		msgArray = append(msgArray, fmt.Sprintf("Current Index: %d\n", f.CurrentIndex))
	} else {
		address, _ := f.Address()
		msgArray = append(msgArray, fmt.Sprintf("Address: %s\n", address))
		msgArray = append(msgArray, fmt.Sprintf("Is Current: %v\n", f.Current))
		msgArray = append(msgArray, fmt.Sprintf("Current Index: %d\n", f.CurrentIndex))

	}

	return strings.Join(msgArray, "")
}

func AccountToDetails(account *database.Account) (*AccountDetails, error) {
	var privateKey string
	var address string
	var path string

	if account.Encrypted {
		return &AccountDetails{Account: *account}, nil
	}

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
		if err != nil {
			return nil, fmt.Errorf("pubkeyToAddress error: %v", err)
		}
	default:
		return nil, errors.New("invalid account type")
	}

	fullAccount := AccountDetails{
		Account:    *account,
		privateKey: privateKey,
		address:    address,
		Path:       path,
	}

	return &fullAccount, nil
}
