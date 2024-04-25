package transaction

import (
	"context"
	"fmt"
	"math/big"
	"met/utils"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func WaitBlock(client *ethclient.Client, height string, heightInterval, heightTimeout uint) error {
	logger := utils.GetLogger("WaitBlock")
	if height == "" {
		logger.Info().Msgf("height is empty, do not wait")
		return nil
	}

	logger.Debug().Msgf("check block height interval: %v second", heightInterval)
	logger.Debug().Msgf("check block height timeout: %v second", heightTimeout)

	blockHeight, ok := new(big.Int).SetString(height, 10)
	if !ok {
		msg := fmt.Sprintf("invalid block height: %v", height)
		return fmt.Errorf(msg)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(heightTimeout))
	defer cancel()

	ticker := time.NewTicker(time.Second * time.Duration(heightInterval))
	defer ticker.Stop()

OUTTER:
	for {
		header, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			logger.Error().Err(err).Msgf("get latest block")
		} else {
			logger.Debug().Msgf("waiting for target: %v, current: %v", blockHeight.String(), header.Number.String())
			if header.Number.Cmp(blockHeight) >= 0 {
				logger.Info().Msgf("block meet")
				break
			}
		}

		select {
		case <-ticker.C:
		case <-ctx.Done():
			break OUTTER
		}
	}
	return nil
}
