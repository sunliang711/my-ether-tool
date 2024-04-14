package decimals

import (
	"context"
	"met/cmd/erc20"
	"met/consts"
	utils "met/utils"
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
	network *string

	contract *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(decimals)

	network = decimals.Flags().String("network", "", "used network, use current if empty")

	contract = decimals.Flags().String("contract", "", "contract address")
}

func getDecimals(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("getDecimals")

	utils.ExitWhen(logger, *contract == "", "need contract address")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*consts.DefaultTimeout)
	defer cancel()

	tokenDecimals, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Decimals, "", "")
	utils.ExitWhenErr(logger, err, "get token decimals error: %v", err)

	// fmt.Printf("token decimals: %s\n", tokenName)
	logger.Info().Msgf("token decimals: %s", tokenDecimals)

}
