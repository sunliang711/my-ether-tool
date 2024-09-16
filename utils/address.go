package utils

import (
	"fmt"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
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
`,
		receipt.TxHash,
		receipt.BlockNumber,
		receipt.BlockHash,
		receipt.ContractAddress,
		receipt.GasUsed,
		receipt.EffectiveGasPrice,
		receipt.Status,
		receipt.TransactionIndex,
		receipt.Type)
	logger.Info().Msg(receiptInfo)
}

// GetContractAddress calculates the contract address from a sender address and nonce
func GetContractAddress(sender string, nonce uint64) (string, error) {
    // Convert sender address to bytes
    senderAddress := common.HexToAddress(sender)

    // RLP encode the sender address and nonce
    rlpEncoded, err := rlp.EncodeToBytes([]interface{}{senderAddress, nonce})
    if err != nil {
        return "", fmt.Errorf("failed to RLP encode: %v", err)
    }

    // Compute the Keccak256 hash of the RLP encoded bytes
    hash := crypto.Keccak256(rlpEncoded)

    // The contract address is the last 20 bytes of the hash
    contractAddress := common.BytesToAddress(hash[12:])

    return contractAddress.Hex(), nil
}
