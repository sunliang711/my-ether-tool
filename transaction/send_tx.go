package transaction

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"met/database"
	utils "met/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrCancel = errors.New("cancel")
)

func SendTx(client *ethclient.Client, from string, tx *types.Transaction, privateKey *ecdsa.PrivateKey, net *database.Network, noconfirm bool) (*types.Receipt, error) {
	logger := utils.GetLogger("SendTx")

	signer := types.LatestSignerForChainID(tx.ChainId())
	txHash := signer.Hash(tx)
	logger.Debug().Msgf("tx hash to be signed: %s", txHash)

	// Sign tx
	logger.Debug().Msgf("sign transaction")
	tx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		return nil, err
	}

	gasPrice, err := utils.Wei2Gwei(tx.GasPrice().String())
	if err != nil {
		return nil, err
	}
	tipCap, err := utils.Wei2Gwei(tx.GasTipCap().String())
	if err != nil {
		return nil, err
	}
	feeCap, err := utils.Wei2Gwei(tx.GasFeeCap().String())
	if err != nil {
		return nil, err
	}

	value, err := utils.FormatUnits(tx.Value().String(), utils.UnitEth)
	if err != nil {
		return nil, err
	}

	logger.Info().Msgf("Transaction to be sent")
	logger.Info().Msgf("From: %s", from)
	logger.Info().Msgf("To: %s", tx.To().String())
	logger.Info().Msgf("Value: %s (%s %s)", tx.Value().String(), value, net.Symbol)
	logger.Info().Msgf("Data: 0x%s", hex.EncodeToString(tx.Data()))
	logger.Info().Msgf("Nonce: %v", tx.Nonce())
	logger.Info().Msgf("ChainId: %s", tx.ChainId().String())
	logger.Info().Msgf("GasLimit: %v", tx.Gas())
	logger.Info().Msgf("GasPrice: %s (%s Gwei)", tx.GasPrice().String(), gasPrice)
	logger.Info().Msgf("GasTipCap: %s (%s Gwei)", tx.GasTipCap().String(), tipCap)
	logger.Info().Msgf("GasFeeCap: %s (%s Gwei)", tx.GasFeeCap().String(), feeCap)

	if !noconfirm {
		input, err := utils.ReadChar("Send ? [y/N] ")
		if err != nil {
			return nil, err
		}

		if input != 'y' {
			return nil, ErrCancel
		}

	}

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	// Send Tx
	err = client.SendTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	ctx2, cancel2 := utils.DefaultTimeoutContext()
	defer cancel2()

	// Wait confirmation
	logger.Info().Msgf("waiting for confirmation")
	receipt, err := bind.WaitMined(ctx2, client, tx)
	if err != nil {
		logger.Error().Err(err).Msgf("get receipt")
		return nil, err
	}
	return receipt, nil
}
