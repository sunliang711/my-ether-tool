package name

import (
	"met/cmd/erc20"
	"met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var nameCmd = &cobra.Command{
	Use:   "name",
	Short: "get name of erc20 contract",
	Long:  "get name of erc20 contract",
	Run:   getName,
}

var (
	network *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(nameCmd)

	network = nameCmd.Flags().String("network", "", "used network, use current if empty")

	contract = nameCmd.Flags().String("contract", "", "contract address")
}

func getName(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("getName")

	utils.ExitWhen(logger, *contract == "", "need contract address")

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "query network error: %v", err)

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)

	tokenName, err := erc20.ReadErc20(ctx, *contract, client, net, erc20.Erc20Name, "", "")
	utils.ExitWhenErr(logger, err, "get token name error: %v", err)

	// fmt.Printf("token name: %s\n", tokenName)
	logger.Info().Msgf("token name: %s", tokenName)

}
