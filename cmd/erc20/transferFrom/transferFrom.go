package transferFrom

import (
	"fmt"
	"met/cmd/erc20"
	"met/database"
	"met/types"
	utils "met/utils"

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

	noconfirm *bool
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

	noconfirm = transferFromCmd.Flags().Bool("noconfirm", false, "noconfirm")
}

func transferToken(cmd *cobra.Command, args []string) {
	var (
		err    error
		logger = utils.GetLogger("transferToken")
	)

	utils.ExitWhen(logger, *contract == "", "need contract address")
	utils.ExitWhen(logger, *from == "", "need from address")
	utils.ExitWhen(logger, *to == "", "need to address")
	utils.ExitWhen(logger, *amount == "", "need token amount")

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "query network error: %v", err)

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)

	acc, err := database.QueryAccountOrCurrent(*account, *accountIndex)
	utils.ExitWhenErr(logger, err, "query account error: %v", err)

	accountDetails, err := types.AccountToDetails(acc)
	utils.ExitWhenErr(logger, err, "get account details error: %v", err)

	decimalsStr := fmt.Sprintf("%v", *decimals)
	// decimals
	if *decimals == 0 {
		// read decimals
		decimalsStr, err = erc20.ReadErc20(ctx, *contract, client, net, erc20.Erc20Decimals, "", "")
		utils.ExitWhenErr(logger, err, "get token decimals error: %v", err)
	}

	realAmount, err := utils.Erc20AmountFromHuman(*amount, decimalsStr)
	utils.ExitWhenErr(logger, err, "convert amount error: %v", err)

	hash, err := erc20.WriteErc20(ctx, *contract, *noconfirm, client, net, accountDetails, erc20.Erc20TransferFrom, *from, *to, realAmount)
	utils.ExitWhenErr(logger, err, "transfer token error: %v", err)

	// fmt.Printf("tx hash: %s\n", hash)
	logger.Info().Msgf("tx hash: %s", hash)

}
