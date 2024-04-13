package name

import (
	"context"
	"fmt"
	"met/cmd/erc20"
	utils "met/utils"
	"time"

	"github.com/spf13/cobra"
)

var nameCmd = &cobra.Command{
	Use:   "name",
	Short: "get name of erc20 contract",
	Long:  "get name of erc20 contract",
	Run:   getName,
}

var (
	account      *string
	accountIndex *uint
	network      *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(nameCmd)

	account = nameCmd.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = nameCmd.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = nameCmd.Flags().String("network", "", "used network, use current if empty")

	contract = nameCmd.Flags().String("contract", "", "contract address")
}

func getName(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*contract == "", "need contract address")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tokenName, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Name, "", "")
	utils.ExitWhenError(err, "get token name error: %v", err)

	fmt.Printf("token name: %s\n", tokenName)

}
