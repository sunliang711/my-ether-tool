package symbol

import (
	"met/cmd/erc20"
	"met/database"
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

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "query network error: %v", err)

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)

	data, err := erc20.ReadErc20(ctx, *contract, client, net, erc20.Erc20Name, "", "")
	utils.ExitWhenErr(logger, err, "get token symbol error: %v", err)

	// fmt.Printf("token symbol: %s\n", data)
	logger.Info().Msgf("token symbol: %s", data)

}
