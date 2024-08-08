package transfer

import (
	"crypto/ecdsa"
	"fmt"
	"met/cmd/erc20"
	"met/consts"
	"met/database"
	"met/transaction"
	ttypes "met/types"
	"met/utils"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
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
	receiver *string // erc20 receiver
	amount   *string // erc20 amount
	decimals *string
	symbol   *string
	value    *string // ether value

	nonce   *string
	chainID *string

	gasLimit      *string
	gasLimitRatio *string

	gasMode  *string
	gasRatio *string
	gasPrice *string
	tipCap   *string
	feeCap   *string

	noconfirm *bool

	confirmations *int8

	blockHeight         *string
	blockHeightInterval *uint
	blockHeightTimeout  *uint

	ledger           *bool
	ledgerDerivePath *string
)

func init() {
	erc20.Erc20Cmd.AddCommand(transferCmd)

	account = transferCmd.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = transferCmd.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = transferCmd.Flags().String("network", "", "used network, use current if empty")

	contract = transferCmd.Flags().String("contract", "", "contract address")
	receiver = transferCmd.Flags().String("to", "", "receiver address")
	amount = transferCmd.Flags().String("amount", "", "token amount, eg: 20.5 (USDT)")
	decimals = transferCmd.Flags().String("decimals", "", "token decimals(optional)")
	symbol = transferCmd.Flags().String("symbol", "", "token symbol(optional)")
	value = transferCmd.Flags().String("value", "0", "ether value(optional)")

	nonce = transferCmd.Flags().String("nonce", "", "nonce")
	chainID = transferCmd.Flags().String("chainId", "", "chain id")

	gasLimit = transferCmd.Flags().String("gasLimit", "", "gas limit")
	gasLimitRatio = transferCmd.Flags().String("gasLimitRatio", "", "gas limit ratio")

	gasMode = transferCmd.Flags().String("gasMode", "auto", "gas mode(eg: auto,legacy,1559)")
	gasRatio = transferCmd.Flags().String("gasRatio", "", "gasRatio")
	gasPrice = transferCmd.Flags().String("gasPrice", "", "gas price(gwei)")
	tipCap = transferCmd.Flags().String("tipCap", "", "tipCap(gwei)")
	feeCap = transferCmd.Flags().String("feeCap", "", "feeCap(gwei)")

	noconfirm = transferCmd.Flags().Bool("noconfirm", false, "noconfirm")

	confirmations = transferCmd.Flags().Int8("confirmations", 0, "blocks of confirmation (N < 0: send tx without receipt. 0: send tx with receipt. N > 0: send tx with receipt and N blocks confirmations)")

	blockHeight = transferCmd.Flags().String("height", "", "send tx after block height")
	blockHeightInterval = transferCmd.Flags().Uint("heightInterval", 2, "check block height interval(unit: ms)")
	blockHeightTimeout = transferCmd.Flags().Uint("heightTimeout", 600, "check block height timeout(unit: s)")

	ledger = transferCmd.Flags().Bool("ledger", false, "use ledger to sign tx, this flag will ignore --account and --account-index")
	ledgerDerivePath = transferCmd.Flags().String("ledgerDerivePath", "m/44'/60'/0'/0/0", "ledger derive path, works only when --ledger is true")
}

func transferToken(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("transferToken")

	var (
		err          error
		privateKey   *ecdsa.PrivateKey
		sender       string
		accountName  string
		accoutnIndex uint

		ledgerWallet  accounts.Wallet
		ledgerAccount *accounts.Account
	)

	if *ledger {
		// 启用ledger时
		ledgerWallet, ledgerAccount, err = utils.ConnectLedger(*ledgerDerivePath)
		utils.ExitWhenErr(logger, err, "connect ledger error: %s", err)
		defer ledgerWallet.Close()

		accountName = "ledger"
		sender = ledgerAccount.Address.Hex()

	} else {
		// 使用普通账户时
		account, err := database.QueryAccountOrCurrent(*account, *accountIndex)
		utils.ExitWhenErr(logger, err, "load account error: %s", err)

		details, err := ttypes.AccountToDetails(account)
		utils.ExitWhenErr(logger, err, "calculate address error: %s", err)

		privateKeyStr, err := details.PrivateKey()
		utils.ExitWhenErr(logger, err, "get account private key error: %s", err)

		privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(privateKeyStr, "0x"))
		utils.ExitWhenErr(logger, err, "parse privateKey error: %s", err)

		sender, err = details.Address()
		utils.ExitWhenErr(logger, err, "get account address error: %s", err)

		accountName = details.Name
		accoutnIndex = details.CurrentIndex
	}

	// network
	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "load network error: %s", err)

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)
	defer client.Close()

	// print
	logger.Info().Msgf("Account Name: %s", accountName)
	logger.Info().Msgf("Account Index: %v", accoutnIndex)
	logger.Info().Msgf("Address: %s", sender)
	logger.Info().Msgf("Network Name: %s", net.Name)
	logger.Info().Msgf("Network RPC: %s", net.Rpc)

	input, err := transaction.ParseErc20Input(client, *contract, sender, *symbol, *decimals, consts.Erc20Transfer, *receiver, *amount)
	utils.ExitWhenErr(logger, err, "%v", err)

	mode := ttypes.GasMode(ttypes.GasMode_value[*gasMode])

	// wait block height
	err = transaction.WaitBlock(client, *blockHeight, *blockHeightInterval, *blockHeightTimeout)
	utils.ExitWhenErr(logger, err, "WaitBlock error: %v", err)

	// build tx
	tx, err := transaction.BuildTx(client, sender, *contract, value, input, *ledger, mode, *nonce, *chainID, *gasLimit, *gasLimitRatio, *gasRatio, *gasPrice, *tipCap, *feeCap, false)
	utils.ExitWhenErr(logger, err, "build tx error: %s", err)

	// send tx
	receipt, tx, err := transaction.SendTx(client, sender, tx, *ledger, ledgerWallet, ledgerAccount, privateKey, net, *noconfirm, *confirmations)
	utils.ExitWhenErr(logger, err, "send transaction error: %v", err)

	if receipt != nil {
		utils.ShowReceipt(logger, receipt)
	}

	link := fmt.Sprintf("%v/tx/%v", net.Explorer, tx.Hash())
	logger.Info().Msgf("tx link: %v", link)

}
