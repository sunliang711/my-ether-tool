package balance

import (
	"fmt"
	"met/cmd/account"
	"met/database"
	"met/types"
	utils "met/utils"

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

	logger.Info().Msgf("query account: %v account index: %v", *accountName, *accountIndex)
	account, err := database.QueryAccountOrCurrent(*accountName, *accountIndex)
	utils.ExitWhenErr(logger, err, "query account: %v error: %v", *accountName, err)

	logger.Info().Msgf("query network: %v", *networkName)
	network, err := database.QueryNetworkOrCurrent(*networkName)
	utils.ExitWhenErr(logger, err, "query ntwork: %v error: %v", *networkName, err)

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	logger.Debug().Msgf("dial rpc: %v", network.Rpc)
	client, err := ethclient.DialContext(ctx, network.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc: %v error: %v", network.Rpc, err)

	accountDetails, err := types.AccountToDetails(account)
	utils.ExitWhenErr(logger, err, "get account details error: %v", err)

	addressStr, err := accountDetails.Address()
	utils.ExitWhenErr(logger, err, "get account address error: %v", err)
	address := common.HexToAddress(addressStr)

	logger.Info().Msgf("query balance for address: %v", addressStr)
	balance, err := client.BalanceAt(ctx, address, nil)
	utils.ExitWhenErr(logger, err, "query account balance error: %v", err)

	humanBalance, err := utils.FormatUnits(balance.String(), utils.UnitEth)
	utils.ExitWhenErr(logger, err, "format balance: %v error: %v", balance.String(), err)

	logger.Info().Msgf("query nonce for address: %v", accountDetails.Address)
	nonce, err := client.PendingNonceAt(ctx, address)
	utils.ExitWhenErr(logger, err, "query nonce error: %v", err)

	logger.Info().Msgf("account: %v account index: %v", accountDetails.Name, accountDetails.CurrentIndex)
	logger.Info().Msgf("address: %v balance: %v %v", accountDetails.Address, humanBalance, network.Symbol)
	logger.Info().Msgf("nonce: %v", nonce)
	logger.Info().Msgf("address link: %v", fmt.Sprintf("%v/address/%v", network.Explorer, accountDetails.Address))
}
