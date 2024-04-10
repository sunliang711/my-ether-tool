package approve

import (
	"context"
	"fmt"
	"my-ether-tool/cmd/erc20"
	"my-ether-tool/utils"
	"time"

	"github.com/spf13/cobra"
)

var approveCmd = &cobra.Command{
	Use:   "approve",
	Short: "approve erc20 token",
	Long:  "approve erc20 token",
	Run:   approveToken,
}

var (
	account      *string
	accountIndex *uint
	network      *string

	contract *string
	spender  *string
	amount   *string
	decimals *uint8
)

func init() {
	erc20.Erc20Cmd.AddCommand(approveCmd)

	account = approveCmd.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = approveCmd.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = approveCmd.Flags().String("network", "", "used network, use current if empty")

	contract = approveCmd.Flags().String("contract", "", "contract address")
	spender = approveCmd.Flags().String("spender", "", "token spender")
	amount = approveCmd.Flags().String("amount", "", "token amount")
	decimals = approveCmd.Flags().Uint8("decimals", 0, "token decimals")
}

func approveToken(cmd *cobra.Command, args []string) {
	var err error
	utils.ExitWithMsgWhen(*contract == "", "need contract address")
	utils.ExitWithMsgWhen(*spender == "", "need token spender")
	utils.ExitWithMsgWhen(*amount == "", "need token amount")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	decimalsStr := fmt.Sprintf("%v", *decimals)
	// decimals
	if *decimals == 0 {
		// read decimals
		decimalsStr, err = erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Decimals, "", "")
		utils.ExitWhenError(err, "get token decimals error: %v", err)
	}

	realAmount, err := utils.Erc20AmountFromHuman(*amount, decimalsStr)
	utils.ExitWhenError(err, "convert amount error: %v", err)

	hash, err := erc20.WriteErc20(ctx, *contract, *network, *account, *accountIndex, erc20.Erc20Approve, *spender, realAmount, "")
	utils.ExitWhenError(err, "approve token error: %v", err)

	fmt.Printf("tx hash: %s\n", hash)

}