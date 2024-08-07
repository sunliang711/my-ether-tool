package utils

import (
	"errors"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/usbwallet"
)

func ConnectLedger(ledgerDerivePath string) (accounts.Wallet, *accounts.Account, error) {
	logger := GetLogger("ConnectLedger")

	var wallet accounts.Wallet

	// 初始化ledger hub
	ledgerHub, err := usbwallet.NewLedgerHub()
	if err != nil {
		return nil, nil, err
	}
	// 等待ledger连接
	logger.Info().Msgf("finding ledger")
	wallets := ledgerHub.Wallets()
	for _, w := range wallets {
		if w.URL().Scheme == "ledger" { // ledger
			logger.Info().Msgf("ledger wallet found")
			wallet = w
			break
		}
	}
	if wallet == nil {
		return nil, nil, errors.New("ledger not found")
	}

	// 打开wallet
	if err := wallet.Open(""); err != nil {
		logger.Error().Msgf("open wallet error: %v", err)
		return nil, nil, err
	}

	if ledgerDerivePath == "" {
		ledgerDerivePath = accounts.DefaultBaseDerivationPath.String()
	}
	logger.Debug().Msgf("ledger derive path: %s", ledgerDerivePath)
	derivedPath, err := accounts.ParseDerivationPath(ledgerDerivePath)
	if err != nil {
		logger.Error().Msgf("parse derivation path error: %v", err)
		return nil, nil, err
	}

	logger.Debug().Msgf("derive account with derive path: %s", derivedPath.String())
	account, err := wallet.Derive(derivedPath, true)
	if err != nil {
		logger.Error().Msgf("derive account error: %v", err)
		return nil, nil, err
	}
	logger.Debug().Msgf("account address: %s", account.Address.Hex())

	return wallet, &account, nil
}
