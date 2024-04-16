package allowance

import (
	"met/cmd/erc20"
	"met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var allowanceCmd = &cobra.Command{
	Use:   "allowance",
	Short: "get allowance of erc20 contract",
	Long:  "get allowance of erc20 contract",
	Run:   getAllowance,
}

var (
	network *string

	contract *string
	owner    *string
	spender  *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(allowanceCmd)

	network = allowanceCmd.Flags().String("network", "", "used network, use current if empty")

	contract = allowanceCmd.Flags().String("contract", "", "contract address")
	owner = allowanceCmd.Flags().String("owner", "", "owner")
	spender = allowanceCmd.Flags().String("spender", "", "spender")
}

func getAllowance(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("getAllowance")

	utils.ExitWhen(logger, *contract == "", "need contract address")
	utils.ExitWhen(logger, *owner == "", "missing owner address")
	utils.ExitWhen(logger, *spender == "", "missing spender address")

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "query network error: %v", err)

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)

	data, err := erc20.ReadErc20(ctx, *contract, client, net, erc20.Erc20Allowance, *owner, *spender)
	utils.ExitWhenErr(logger, err, "get allowance error: %v", err)

	logger.Info().Msgf("allowance: %s", data)

}
