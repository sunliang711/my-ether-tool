package transfer

import (
	"context"
	"fmt"
	"met/cmd/erc20"
	"met/consts"
	utils "met/utils"
	"time"

	"github.com/spf13/cobra"
)

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "transfer erc20 token",
	Long:  "transfer erc20 token",
	Run:   transferToken,
}

var (
	account      *string
	accountIndex *uint
	network      *string

	contract *string
	to       *string
	amount   *string
	decimals *uint8
)

func init() {
	erc20.Erc20Cmd.AddCommand(transferCmd)

	account = transferCmd.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = transferCmd.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = transferCmd.Flags().String("network", "", "used network, use current if empty")

	contract = transferCmd.Flags().String("contract", "", "contract address")
	to = transferCmd.Flags().String("to", "", "receiver address")
	amount = transferCmd.Flags().String("amount", "", "token amount")
	decimals = transferCmd.Flags().Uint8("decimals", 0, "token decimals(optional)")
}

func transferToken(cmd *cobra.Command, args []string) {
	var (
		err    error
		logger = utils.GetLogger("transferToken")
	)

	utils.ExitWhen(logger, *contract == "", "need contract address")
	utils.ExitWhen(logger, *to == "", "need to address")
	utils.ExitWhen(logger, *amount == "", "need token amount")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*consts.DefaultTimeout)
	defer cancel()

	// tokenName, err := erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Name, "", "")
	// utils.ExitWhenError(err, "get token symbol error: %v", err)

	decimalsStr := fmt.Sprintf("%v", *decimals)
	// decimals
	if *decimals == 0 {
		// read decimals
		decimalsStr, err = erc20.ReadErc20(ctx, *contract, *network, erc20.Erc20Decimals, "", "")
		utils.ExitWhenErr(logger, err, "get token decimals error: %v", err)
	}

	realAmount, err := utils.Erc20AmountFromHuman(*amount, decimalsStr)
	utils.ExitWhenErr(logger, err, "convert amount error: %v", err)

	hash, err := erc20.WriteErc20(ctx, *contract, *network, *account, *accountIndex, erc20.Erc20Transfer, *to, realAmount, "")
	utils.ExitWhenErr(logger, err, "transfer token error: %v", err)

	// fmt.Printf("tx hash: %s\n", hash)
	logger.Info().Msgf("tx hash: %s", hash)

}
