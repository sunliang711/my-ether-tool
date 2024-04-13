package utils

import (
	"regexp"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

var re = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

func IsValidAddress(address string) bool {
	return re.MatchString(address)
}

func ShowReceipt(logger zerolog.Logger, receipt *types.Receipt) {
	logger.Info().Msgf("[Receipt] Block Number: %v", receipt.BlockNumber)
	logger.Info().Msgf("[Receipt] Block Hash: %v", receipt.BlockHash)
	logger.Info().Msgf("[Receipt] Contract Address: %v", receipt.ContractAddress)
	logger.Info().Msgf("[Receipt] Gas Used: %v", receipt.GasUsed)
	logger.Info().Msgf("[Receipt] Status: %v", receipt.Status)
	logger.Info().Msgf("[Receipt] Tx Index: %v", receipt.TransactionIndex)
	logger.Info().Msgf("[Receipt] Type: %v", receipt.Type)
}
