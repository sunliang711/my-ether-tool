package decimals

import (
	"context"
	"fmt"
	"my-ether-tool/cmd/erc20"
	"my-ether-tool/utils"
	"time"

	"github.com/spf13/cobra"
)

var decimals = &cobra.Command{
	Use:   "decimals",
	Short: "get decimals of erc20 contract",
	Long:  "get decimals of erc20 contract",
	Run:   getDecimals,
}

var (
	account      *string
	accountIndex *uint
	network      *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(decimals)

	account = decimals.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = decimals.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = decimals.Flags().String("network", "", "used network, use current if empty")

	contract = decimals.Flags().String("contract", "", "contract address")
}

func getDecimals(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*contract == "", "need contract address")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tokenName, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Decimals, "", "")
	utils.ExitWhenError(err, "get token decimals error: %v", err)

	fmt.Printf("token decimals: %s\n", tokenName)

}
