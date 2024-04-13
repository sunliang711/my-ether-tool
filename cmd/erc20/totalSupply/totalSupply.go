package totalSupply

import (
	"context"
	"fmt"
	"met/cmd/erc20"
	utils "met/utils"
	"time"

	"github.com/spf13/cobra"
)

var totalSupply = &cobra.Command{
	Use:   "totalSupply",
	Short: "get totalSupply of erc20 contract",
	Long:  "get totalSupply of erc20 contract",
	Run:   getTotalSupply,
}

var (
	account      *string
	accountIndex *uint
	network      *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(totalSupply)

	account = totalSupply.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = totalSupply.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = totalSupply.Flags().String("network", "", "used network, use current if empty")

	contract = totalSupply.Flags().String("contract", "", "contract address")
}

func getTotalSupply(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*contract == "", "need contract address")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// read totalSupply
	totalSupply, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20TotalSupply, "", "")
	utils.ExitWhenError(err, "get token total supply error: %v", err)

	// read decimals
	decimals, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Decimals, "", "")
	utils.ExitWhenError(err, "get token decimals error: %v", err)

	humanTotalSupply, err := utils.Erc20AmountToHuman(totalSupply, decimals)
	utils.ExitWhenError(err, "convert totalSupply error: %v", err)

	fmt.Printf("token total supply: %s\n", humanTotalSupply)

}
