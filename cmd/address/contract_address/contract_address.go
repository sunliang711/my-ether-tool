package contractAddress

import (
	cmd "met/cmd/address"
	"met/utils"

	"github.com/spf13/cobra"
)

// ContractAddress represents the checksum command
var ContractAddress = &cobra.Command{
	Use:   "contractAddress",
	Short: "ethereum contract address",
	Run:   contractAddress,
}

var (
	fromAddress *string
	nonce       *uint64
)

func init() {
	cmd.AddressCmd.AddCommand(ContractAddress)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	fromAddress = ContractAddress.Flags().String("fromAddress", "", "ethereum address deploy contract")
	nonce = ContractAddress.Flags().Uint64("nonce", 0, "nonce of the transaction")
}

func contractAddress(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("checksum")

	utils.ExitWhen(logger, len(*fromAddress) == 0, "need address")
	logger.Info().Msgf("address: %s", *fromAddress)

	isValid := utils.IsValidAddress(*fromAddress)
	utils.ExitWhen(logger, !isValid, "invalid ethereum address")

	contractAddress, err := utils.GetContractAddress(*fromAddress, *nonce)
	utils.ExitWhenErr(logger, err, "get contract address failed")

	logger.Info().Msgf("contract address: %s", contractAddress)

}
