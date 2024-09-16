package address

import (
	cmd "met/cmd/address"
	"met/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

// CheckSum represents the checksum command
var CheckSum = &cobra.Command{
	Use:   "checksum",
	Short: "ethereum address related",
	Run:   checksum,
}

var (
	address *string
)

func init() {
	cmd.AddressCmd.AddCommand(CheckSum)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	address = CheckSum.Flags().String("address", "", "ethereum address to be checksumed")
}

func checksum(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("checksum")

	utils.ExitWhen(logger, len(*address) == 0, "need address")
	logger.Info().Msgf("address: %s", *address)

	isValid := utils.IsValidAddress(*address)
	utils.ExitWhen(logger, !isValid, "invalid ethereum address")

	checksumed := common.HexToAddress(*address).Hex()
	logger.Info().Msgf("checksumed address: %s", checksumed)

}
