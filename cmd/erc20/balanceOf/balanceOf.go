package balanceOf

import (
	"met/cmd/erc20"
	"met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var balanceOfCmd = &cobra.Command{
	Use:   "balanceOf",
	Short: "get erc20 balance",
	Long:  "get erc20 balance",
	Run:   getBalance,
}

var (
	network *string

	contract *string
	owner    *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(balanceOfCmd)

	network = balanceOfCmd.Flags().String("network", "", "used network, use current if empty")

	contract = balanceOfCmd.Flags().String("contract", "", "contract address")
	owner = balanceOfCmd.Flags().String("owner", "", "owner address")
}

func getBalance(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("getBalance")

	utils.ExitWhen(logger, *contract == "", "need contract address")
	utils.ExitWhen(logger, *owner == "", "need owner address")

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "query network error: %v", err)

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)

	// read balance
	balance, err := erc20.ReadErc20(ctx, *contract, client, net, erc20.Erc20BalanceOf, *owner, "")
	utils.ExitWhenErr(logger, err, "get token balance error: %v", err)

	// read decimals
	decimals, err := erc20.ReadErc20(ctx, *contract, client, net, erc20.Erc20Decimals, "", "")
	utils.ExitWhenErr(logger, err, "get token decimals error: %v", err)

	humanBalance, err := utils.Erc20AmountToHuman(balance, decimals)
	utils.ExitWhenErr(logger, err, "convert balance error: %v", err)

	// fmt.Printf("token balance: %s\n", humanBalance)
	logger.Info().Msgf("token balance: %s", humanBalance)

}
