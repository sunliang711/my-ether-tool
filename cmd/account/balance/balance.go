package balance

import (
	"context"
	"fmt"
	"met/cmd/account"
	"met/consts"
	"met/database"
	"met/types"
	utils "met/utils"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

var balanceCmd = &cobra.Command{
	Use:        "balance",
	ArgAliases: []string{"b"},
	Short:      "account balance",
	Long:       "account balance",
	Run:        getBalance,
}

var (
	accountName  *string
	accountIndex *uint
	networkName  *string
)

func init() {
	account.AccountCmd.AddCommand(balanceCmd)

	accountName = balanceCmd.Flags().String("account", "", "account name")
	accountIndex = balanceCmd.Flags().Uint("account-index", 0, "account index")
	networkName = balanceCmd.Flags().String("network", "", "network name")
}

func getBalance(cmd *cobra.Command, args []string) {
	var (
		err    error
		logger = utils.GetLogger("getBalance")
	)

	account, err := database.QueryAccountOrCurrent(*accountName, *accountIndex)
	utils.ExitWhenErr(logger, err, "query account: %v error: %v", *accountName, err)

	network, err := database.QueryNetworkOrCurrent(*networkName)
	utils.ExitWhenErr(logger, err, "query ntwork: %v error: %v", *networkName, err)

	client, err := ethclient.Dial(network.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc: %v error: %v", network.Rpc, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*consts.DefaultTimeout)
	defer cancel()

	accountDetails, err := types.AccountToDetails(account)
	utils.ExitWhenErr(logger, err, "get account details error: %v", err)

	address := common.HexToAddress(accountDetails.Address)
	balance, err := client.BalanceAt(ctx, address, nil)
	utils.ExitWhenErr(logger, err, "query account balance error: %v", err)

	humanBalance, err := utils.FormatUnits(balance.String(), utils.UnitEth)
	utils.ExitWhenErr(logger, err, "format balance: %v error: %v", balance.String(), err)

	nonce, err := client.PendingNonceAt(ctx, address)
	utils.ExitWhenErr(logger, err, "query nonce error: %v", err)

	logger.Info().Msgf("account: %v account index: %v", accountDetails.Name, accountDetails.CurrentIndex)
	logger.Info().Msgf("address: %v balance: %v %v", accountDetails.Address, humanBalance, network.Symbol)
	logger.Info().Msgf("nonce: %v", nonce)
	logger.Info().Msgf("address link: %v", fmt.Sprintf("%v/address/%v", network.Explorer, accountDetails.Address))
}
