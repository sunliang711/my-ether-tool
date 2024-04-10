package transferFrom

import (
	"context"
	"fmt"
	"my-ether-tool/cmd/erc20"
	"my-ether-tool/utils"
	"time"

	"github.com/spf13/cobra"
)

var transferFromCmd = &cobra.Command{
	Use:   "transferFrom",
	Short: "transferFrom erc20 token",
	Long:  "transferFrom erc20 token",
	Run:   transferToken,
}

var (
	account      *string
	accountIndex *uint
	network      *string

	contract *string
	from     *string
	to       *string
	amount   *string
	decimals *uint8
)

func init() {
	erc20.Erc20Cmd.AddCommand(transferFromCmd)

	account = transferFromCmd.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = transferFromCmd.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = transferFromCmd.Flags().String("network", "", "used network, use current if empty")

	contract = transferFromCmd.Flags().String("contract", "", "contract address")
	from = transferFromCmd.Flags().String("from", "", "from address")
	to = transferFromCmd.Flags().String("to", "", "to address")
	amount = transferFromCmd.Flags().String("amount", "", "token amount")
	decimals = transferFromCmd.Flags().Uint8("decimals", 0, "token decimals(optional)")
}

func transferToken(cmd *cobra.Command, args []string) {
	var err error
	utils.ExitWithMsgWhen(*contract == "", "need contract address")
	utils.ExitWithMsgWhen(*from == "", "need from address")
	utils.ExitWithMsgWhen(*to == "", "need to address")
	utils.ExitWithMsgWhen(*amount == "", "need token amount")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// tokenName, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Name, "", "")
	// utils.ExitWhenError(err, "get token symbol error: %v", err)

	decimalsStr := fmt.Sprintf("%v", *decimals)
	// decimals
	if *decimals == 0 {
		// read decimals
		decimalsStr, err = erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Decimals, "", "")
		utils.ExitWhenError(err, "get token decimals error: %v", err)
	}

	realAmount, err := utils.Erc20AmountFromHuman(*amount, decimalsStr)
	utils.ExitWhenError(err, "convert amount error: %v", err)

	hash, err := erc20.WriteErc20(ctx, *contract, *network, *account, *accountIndex, erc20.Erc20TransferFrom, *from, *to, realAmount)
	utils.ExitWhenError(err, "transfer token error: %v", err)

	fmt.Printf("tx hash: %s\n", hash)

}
