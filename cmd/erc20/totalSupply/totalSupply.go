package totalSupply

import (
	"met/cmd/erc20"
	"met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var totalSupply = &cobra.Command{
	Use:   "totalSupply",
	Short: "get totalSupply of erc20 contract",
	Long:  "get totalSupply of erc20 contract",
	Run:   getTotalSupply,
}

var (
	network *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(totalSupply)

	network = totalSupply.Flags().String("network", "", "used network, use current if empty")

	contract = totalSupply.Flags().String("contract", "", "contract address")
}

func getTotalSupply(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("getTotalSupply")

	utils.ExitWhen(logger, *contract == "", "need contract address")

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "query network error: %v", err)

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)

	// read totalSupply
	totalSupply, err := erc20.ReadErc20(ctx, *contract, client, net, erc20.Erc20TotalSupply, "", "")
	utils.ExitWhenErr(logger, err, "get token total supply error: %v", err)

	// read decimals
	decimals, err := erc20.ReadErc20(ctx, *contract, client, net, erc20.Erc20Decimals, "", "")
	utils.ExitWhenErr(logger, err, "get token decimals error: %v", err)

	humanTotalSupply, err := utils.Erc20AmountToHuman(totalSupply, decimals)
	utils.ExitWhenErr(logger, err, "convert totalSupply error: %v", err)

	logger.Info().Msgf("token total supply: %s", humanTotalSupply)

}
