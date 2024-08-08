package transaction

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"met/database"
	utils "met/utils"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrCancel = errors.New("cancel")
)

// 多返回一个types.Transaction是为了当不需要receipt(confirmations=0)时，能知道tx hash
func SendTx(client *ethclient.Client, from string, tx *types.Transaction, ledger bool, ledgerWallet accounts.Wallet, ledgerAccount *accounts.Account, privateKey *ecdsa.PrivateKey, net *database.Network, noconfirm bool, confirmations int8) (*types.Receipt, *types.Transaction, error) {
	var err error
	logger := utils.GetLogger("SendTx")

	signer := types.LatestSignerForChainID(tx.ChainId())
	txHash := signer.Hash(tx)
	logger.Debug().Msgf("tx hash to be signed: %s", txHash)

	// Sign tx
	logger.Debug().Msgf("sign transaction")
	if ledger {
		chainID, err := client.ChainID(context.Background())
		if err != nil {
			logger.Error().Msgf("get chain id error: %v", err)
			return nil, nil, err
		}
		fmt.Printf("confirm on your ledger device..\n")
		tx, err = ledgerWallet.SignTx(*ledgerAccount, tx, chainID)
		if err != nil {
			logger.Error().Msgf("sign tx error: %v", err)
			return nil, nil, err
		}

	} else {
		tx, err = types.SignTx(tx, signer, privateKey)
		if err != nil {
			return nil, nil, err
		}
	}

	logger.Debug().Msgf("tx hash: %v", tx.Hash())

	gasPrice, err := utils.Wei2Gwei(tx.GasPrice().String())
	if err != nil {
		return nil, nil, err
	}
	tipCap, err := utils.Wei2Gwei(tx.GasTipCap().String())
	if err != nil {
		return nil, nil, err
	}
	feeCap, err := utils.Wei2Gwei(tx.GasFeeCap().String())
	if err != nil {
		return nil, nil, err
	}

	value, err := utils.FormatUnits(tx.Value().String(), utils.UnitEth)
	if err != nil {
		return nil, nil, err
	}

	to := "EMPTY (contract creation)"
	if tx.To() != nil {
		to = tx.To().String()
	}

	accountType := "local wallet"
	if ledger {
		accountType = "ledger"
	}

	txInfo := fmt.Sprintf(`
Transaction to be sent
From:                %s (%s)
To:                  %s
Value:               %s (%s %s)
Data:                0x%s
Nonce:               %v
ChainId:             %s
GasLimit:            %v
GasPrice:            %s (%s Gwei)
GasTipCap:           %s (%s Gwei)
GasFeeCap:           %s (%s Gwei)
`,
		from, accountType,
		to,
		tx.Value().String(), value, net.Symbol,
		hex.EncodeToString(tx.Data()),
		tx.Nonce(),
		tx.ChainId().String(),
		tx.Gas(),
		tx.GasPrice().String(), gasPrice,
		tx.GasTipCap().String(), tipCap,
		tx.GasFeeCap().String(), feeCap)

	logger.Info().Msgf(txInfo)

	if !noconfirm {
		input, err := utils.ReadChar("Send ? [y/N] ")
		if err != nil {
			return nil, nil, err
		}

		if input != 'y' {
			return nil, nil, ErrCancel
		}

	}

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	// Send Tx
	err = client.SendTransaction(ctx, tx)
	if err != nil {
		return nil, nil, err
	}

	ctx2, cancel2 := utils.DefaultTimeoutContext()
	defer cancel2()

	receipt, err := waitTx(ctx2, client, tx, confirmations)
	if err != nil {
		logger.Error().Err(err).Msgf("wait tx")
	}

	return receipt, tx, nil
}

func waitTx(ctx context.Context, client *ethclient.Client, tx *types.Transaction, confirmations int8) (*types.Receipt, error) {
	logger := utils.GetLogger("waitTx")

	if confirmations < 0 {
		logger.Debug().Msgf("confirmations < 0,do not get query receipt")
		return nil, nil
	}

	logger.Info().Msgf("query receipt")
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return nil, err
	}
	logger.Debug().Msgf("get receipt: %+v", receipt)

	if confirmations > 0 {
		logger.Info().Msgf("waiting for %v confirmations..", confirmations)

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		ctx2, cancel := context.WithTimeout(context.Background(), time.Minute*5)
		defer cancel()

	OUTTER:
		for {
			// get latest block
			h, err := client.HeaderByNumber(ctx2, nil)
			if err != nil {
				logger.Error().Err(err).Msgf("query latest block header")
			} else {
				diff := new(big.Int).Sub(h.Number, receipt.BlockNumber)
				logger.Debug().Msgf("diff block number: %v (latest: %v mined: %v)", diff.String(), h.Number, receipt.BlockNumber)
				if diff.Cmp(big.NewInt(int64(confirmations))) >= 0 {
					logger.Debug().Msgf("confirmations meet")
					break OUTTER
				}
			}

			select {
			case <-ctx2.Done():
				logger.Warn().Msgf("context done before confirmations completed")
				break OUTTER
			case <-ticker.C:
			}
		}
	}

	return receipt, nil
}
