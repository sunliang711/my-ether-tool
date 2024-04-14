package allowance

import (
	"context"
	"met/cmd/erc20"
	"met/consts"
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
	network *string

	contract *string
	owner    *string
	spender  *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(allowanceCmd)

	network = allowanceCmd.Flags().String("network", "", "used network, use current if empty")

	contract = allowanceCmd.Flags().String("contract", "", "contract address")
	owner = allowanceCmd.Flags().String("owner", "", "owner")
	spender = allowanceCmd.Flags().String("spender", "", "spender")
}

func getAllowance(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("getAllowance")

	utils.ExitWhen(logger, *contract == "", "need contract address")
	utils.ExitWhen(logger, *owner == "", "missing owner address")
	utils.ExitWhen(logger, *spender == "", "missing spender address")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*consts.DefaultTimeout)
	defer cancel()

	data, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Allowance, *owner, *spender)
	utils.ExitWhenErr(logger, err, "get allowance error: %v", err)

	logger.Info().Msgf("allowance: %s", data)

}
