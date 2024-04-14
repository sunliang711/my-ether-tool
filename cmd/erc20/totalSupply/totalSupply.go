package totalSupply

import (
	"context"
	"fmt"
	"met/cmd/erc20"
	"met/consts"
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
	network *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(totalSupply)

	network = totalSupply.Flags().String("network", "", "used network, use current if empty")

	contract = totalSupply.Flags().String("contract", "", "contract address")
}

func getTotalSupply(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*contract == "", "need contract address")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*consts.DefaultTimeout)
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
