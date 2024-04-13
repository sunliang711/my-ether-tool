package symbol

import (
	"context"
	"fmt"
	"met/cmd/erc20"
	utils "met/utils"
	"time"

	"github.com/spf13/cobra"
)

var symbolCmd = &cobra.Command{
	Use:   "symbol",
	Short: "get symbol of erc20 contract",
	Long:  "get symbol of erc20 contract",
	Run:   getSymbol,
}

var (
	account      *string
	accountIndex *uint
	network      *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(symbolCmd)

	account = symbolCmd.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = symbolCmd.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = symbolCmd.Flags().String("network", "", "used network, use current if empty")

	contract = symbolCmd.Flags().String("contract", "", "contract address")
}

func getSymbol(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*contract == "", "need contract address")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	data, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Name, "", "")
	utils.ExitWhenError(err, "get token symbol error: %v", err)

	fmt.Printf("token symbol: %s\n", data)

}
