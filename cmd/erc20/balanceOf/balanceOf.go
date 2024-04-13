package balanceOf

import (
	"context"
	"fmt"
	"met/cmd/erc20"
	utils "met/utils"
	"time"

	"github.com/spf13/cobra"
)

var balanceOfCmd = &cobra.Command{
	Use:   "balanceOf",
	Short: "get erc20 balance",
	Long:  "get erc20 balance",
	Run:   getBalance,
}

var (
	account      *string
	accountIndex *uint
	network      *string

	contract *string
	owner    *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(balanceOfCmd)

	account = balanceOfCmd.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = balanceOfCmd.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = balanceOfCmd.Flags().String("network", "", "used network, use current if empty")

	contract = balanceOfCmd.Flags().String("contract", "", "contract address")
	owner = balanceOfCmd.Flags().String("owner", "", "owner address")
}

func getBalance(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*contract == "", "need contract address")
	utils.ExitWithMsgWhen(*owner == "", "need owner address")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// read balance
	balance, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20BalanceOf, *owner, "")
	utils.ExitWhenError(err, "get token balance error: %v", err)

	// read decimals
	decimals, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Decimals, "", "")
	utils.ExitWhenError(err, "get token decimals error: %v", err)

	humanBalance, err := utils.Erc20AmountToHuman(balance, decimals)
	utils.ExitWhenError(err, "convert balance error: %v", err)

	fmt.Printf("token balance: %s\n", humanBalance)

}
