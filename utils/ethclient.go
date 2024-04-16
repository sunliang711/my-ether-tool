package utils

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
)

func DialRpc(ctx context.Context, rpc string) (*ethclient.Client, error) {
	loggers := GetLogger("DialRpc")

	loggers.Info().Msgf("Dial rpc: %v", rpc)
	client, err := ethclient.DialContext(ctx, rpc)
	if err != nil {
		return nil, err
	}

	return client, nil
}
