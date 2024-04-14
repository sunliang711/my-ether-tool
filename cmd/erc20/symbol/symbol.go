package symbol

import (
	"met/cmd/erc20"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var symbolCmd = &cobra.Command{
	Use:   "symbol",
	Short: "get symbol of erc20 contract",
	Long:  "get symbol of erc20 contract",
	Run:   getSymbol,
}

var (
	network *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(symbolCmd)

	network = symbolCmd.Flags().String("network", "", "used network, use current if empty")

	contract = symbolCmd.Flags().String("contract", "", "contract address")
}

func getSymbol(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("getSymbol")

	utils.ExitWhen(logger, *contract == "", "need contract address")

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	data, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Name, "", "")
	utils.ExitWhenErr(logger, err, "get token symbol error: %v", err)

	// fmt.Printf("token symbol: %s\n", data)
	logger.Info().Msgf("token symbol: %s", data)

}
