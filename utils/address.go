package utils

import (
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
)

var re = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

func IsValidAddress(address string) bool {
	return re.MatchString(address)
}

func ShowReceipt(logger zerolog.Logger, receipt *types.Receipt) {
	receiptInfo := fmt.Sprintf(`
Transaction Receipt
Tx Hash:             %v
Block Number:        %v
Block Hash:          %v
Contract Address:    %v
Gas Used:            %v
Gas Price:           %v
Status:              %v
Tx Index:            %v
Type:                %v
`, receipt.TxHash, receipt.BlockNumber, receipt.BlockHash, receipt.ContractAddress, receipt.GasUsed, receipt.EffectiveGasPrice, receipt.Status, receipt.TransactionIndex, receipt.Type)
	logger.Info().Msg(receiptInfo)
}
