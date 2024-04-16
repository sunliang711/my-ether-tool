package decimals

import (
	"met/cmd/erc20"
	"met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var decimals = &cobra.Command{
	Use:   "decimals",
	Short: "get decimals of erc20 contract",
	Long:  "get decimals of erc20 contract",
	Run:   getDecimals,
}

var (
	network *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(decimals)

	network = decimals.Flags().String("network", "", "used network, use current if empty")

	contract = decimals.Flags().String("contract", "", "contract address")
}

func getDecimals(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("getDecimals")

	utils.ExitWhen(logger, *contract == "", "need contract address")

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "query network error: %v", err)

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)

	tokenDecimals, err := erc20.ReadErc20(ctx, *contract, client, net, erc20.Erc20Decimals, "", "")
	utils.ExitWhenErr(logger, err, "get token decimals error: %v", err)

	logger.Info().Msgf("token decimals: %s", tokenDecimals)

}
