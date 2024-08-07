package tx

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"met/cmd/tx"
	"met/consts"
	database "met/database"
	transaction "met/transaction"
	ttypes "met/types"
	utils "met/utils"
	"strings"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "send tx",
	Long:  "send transaction",
	Run:   sendTransaction,
}
var (
	account      *string
	accountIndex *uint
	network      *string

	to    *string
	value *string
	// 忽略value，发送所有ether
	all *bool

	data    *string
	abi     *string
	method  *string
	abiArgs *[]string

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
	tx.TxCmd.AddCommand(sendCmd)

	// sendCmd.Flags()
	account = sendCmd.Flags().String("account", "", "account to be used to send tx,use current if empty")
	accountIndex = sendCmd.Flags().Uint("account-index", 0, "account index to be used to send tx")
	network = sendCmd.Flags().String("network", "", "used network, use current if empty")

	to = sendCmd.Flags().String("to", "", "transaction receiver")
	value = sendCmd.Flags().String("value", "0", "value (uint: eth)")
	all = sendCmd.Flags().Bool("all", false, "send all ether")

	// data or abi + args
	data = sendCmd.Flags().String("data", "", "data of transaction, conflict with --abi")
	abi = sendCmd.Flags().String("abi", "", "abi JSON string, conflict with --data, available builtin abi: erc20 erc721 erc1155")
	method = sendCmd.Flags().String("method", "", "methodName, conflict with --data")
	abiArgs = sendCmd.Flags().StringArray("args", nil, "arguments of abi( --args 0x... --args 200)")

	nonce = sendCmd.Flags().String("nonce", "", "nonce")
	chainID = sendCmd.Flags().String("chainId", "", "chain id")

	gasLimit = sendCmd.Flags().String("gasLimit", "", "gas limit")
	gasLimitRatio = sendCmd.Flags().String("gasLimitRatio", "", "gas limit ratio")

	gasMode = sendCmd.Flags().String("gasMode", "auto", "gas mode(eg: auto,legacy,1559)")
	gasRatio = sendCmd.Flags().String("gasRatio", "", "gasRatio")
	gasPrice = sendCmd.Flags().String("gasPrice", "", "gas price(gwei)")
	tipCap = sendCmd.Flags().String("tipCap", "", "tipCap(gwei)")
	feeCap = sendCmd.Flags().String("feeCap", "", "feeCap(gwei)")

	noconfirm = sendCmd.Flags().BoolP("noconfirm", "y", false, "do not need to confirm")

	confirmations = sendCmd.Flags().Int8("confirmations", 0, "blocks of confirmation (N<0: send tx without receipt. 0: send tx with receipt. N>0: send tx with receipt and N blocks confirmations)")

	blockHeight = sendCmd.Flags().String("height", "", "send tx after block height")
	blockHeightInterval = sendCmd.Flags().Uint("heightInterval", 2, "check block height interval(unit: ms)")
	blockHeightTimeout = sendCmd.Flags().Uint("heightTimeout", 600, "check block height timeout(unit: s)")

	ledger = sendCmd.Flags().Bool("ledger", false, "use ledger to sign tx")
	ledgerDerivePath = sendCmd.Flags().String("ledgerDerivePath", "m/44'/60'/0'/0/0", "ledger derive path")
}

func sendTransaction(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("sendTransaction")

	var (
		err          error
		privateKey   *ecdsa.PrivateKey
		from         string
		accountName  string
		accoutnIndex uint

		ledgerWallet  accounts.Wallet
		ledgerAccount *accounts.Account
	)

	walletType := consts.WalletTypeNormal
	if *ledger {
		walletType = consts.WalletTypeLedger

		ledgerWallet, ledgerAccount, err = utils.ConnectLedger(*ledgerDerivePath)
		utils.ExitWhenErr(logger, err, "connect ledger error: %s", err)
		defer ledgerWallet.Close()

		accountName = "ledger"
		from = ledgerAccount.Address.Hex()

	} else {
		// account
		account, err := database.QueryAccountOrCurrent(*account, *accountIndex)
		utils.ExitWhenErr(logger, err, "load account error: %s", err)

		details, err := ttypes.AccountToDetails(account)
		utils.ExitWhenErr(logger, err, "calculate address error: %s", err)

		privateKeyStr, err := details.PrivateKey()
		utils.ExitWhenErr(logger, err, "get account private key error: %s", err)

		privateKey, err = crypto.HexToECDSA(strings.TrimPrefix(privateKeyStr, "0x"))
		utils.ExitWhenErr(logger, err, "parse privateKey error: %s", err)

		from, err = details.Address()
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
	logger.Info().Msgf("Address: %s", from)
	logger.Info().Msgf("Network Name: %s", net.Name)
	logger.Info().Msgf("Network RPC: %s", net.Rpc)

	utils.ExitWhen(logger, *data != "" && (*abi != "" || len(*abiArgs) > 0), "--data conflicts with --abi and --args")
	var input []byte
	if *data != "" {
		input, err = hex.DecodeString(*data)
		utils.ExitWhenErr(logger, err, "decode data: %v error: %v", *data, err)
	}

	input, err = transaction.ParseAbi(*abi, *method, *abiArgs...)
	utils.ExitWhenErr(logger, err, "parse abi error: %v", err)

	mode := ttypes.GasMode(ttypes.GasMode_value[*gasMode])

	// wait block height
	err = transaction.WaitBlock(client, *blockHeight, *blockHeightInterval, *blockHeightTimeout)
	utils.ExitWhenErr(logger, err, "WaitBlock error: %v", err)

	ctx2, cancel2 := utils.DefaultTimeoutContext()
	defer cancel2()
	// build tx
	tx, err := transaction.BuildTx(ctx2, client, from, *to, value, input, *ledger, mode, *nonce, *chainID, *gasLimit, *gasLimitRatio, *gasRatio, *gasPrice, *tipCap, *feeCap, *all)
	utils.ExitWhenErr(logger, err, "build tx error: %s", err)

	// send tx
	receipt, tx2, err := transaction.SendTx(client, from, tx, walletType, ledgerWallet, ledgerAccount, privateKey, net, *noconfirm, *confirmations)
	utils.ExitWhenErr(logger, err, "send transaction error: %v", err)

	if receipt != nil {
		utils.ShowReceipt(logger, receipt)
	}

	link := fmt.Sprintf("%v/tx/%v", net.Explorer, tx2.Hash())
	logger.Info().Msgf("tx link: %v", link)

}
