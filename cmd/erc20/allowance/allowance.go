package allowance

import (
	"context"
	"fmt"
	"met/cmd/erc20"
	utils "met/utils"
	"time"

	"github.com/spf13/cobra"
)

var allowanceCmd = &cobra.Command{
	Use:   "allowance",
	Short: "get allowance of erc20 contract",
	Long:  "get allowance of erc20 contract",
	Run:   getAllowance,
}

var (
	account      *string
	accountIndex *uint
	network      *string

	contract *string
	owner    *string
	spender  *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(allowanceCmd)

	account = allowanceCmd.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = allowanceCmd.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = allowanceCmd.Flags().String("network", "", "used network, use current if empty")

	contract = allowanceCmd.Flags().String("contract", "", "contract address")
	owner = allowanceCmd.Flags().String("owner", "", "owner")
	spender = allowanceCmd.Flags().String("spender", "", "spender")
}

func getAllowance(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*contract == "", "need contract address")
	utils.ExitWithMsgWhen(*owner == "", "missing owner address")
	utils.ExitWithMsgWhen(*spender == "", "missing spender address")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Allowance, *owner, *spender)
	utils.ExitWhenError(err, "get allowance error: %v", err)

	fmt.Printf("allowance : %s\n", data)

}
