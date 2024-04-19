package tx

import (
	"encoding/hex"
	"fmt"
	"met/cmd/tx"
	"met/consts"
	database "met/database"
	transaction "met/transaction"
	ttypes "met/types"
	utils "met/utils"
	"strings"

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
	abi = sendCmd.Flags().String("abi", "", "abi JSON string, conflict with --data")
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
}

func sendTransaction(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("sendTransaction")

	account, err := database.QueryAccountOrCurrent(*account, *accountIndex)
	utils.ExitWhenErr(logger, err, "load account error: %s", err)

	details, err := ttypes.AccountToDetails(account)
	utils.ExitWhenErr(logger, err, "calculate address error: %s", err)

	privateKeyStr, err := details.PrivateKey()
	utils.ExitWhenErr(logger, err, "get account private key error: %s", err)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(privateKeyStr, "0x"))
	utils.ExitWhenErr(logger, err, "parse privateKey error: %s", err)

	from, err := details.Address()
	utils.ExitWhenErr(logger, err, "get account address error: %s", err)

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "load network error: %s", err)

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)
	defer client.Close()

	logger.Info().Msgf("Network Name: %s", net.Name)
	logger.Info().Msgf("Network RPC: %s", net.Rpc)
	logger.Info().Msgf("Account Name: %s", details.Name)
	logger.Info().Msgf("Account Index: %v", details.CurrentIndex)
	logger.Info().Msgf("Address: %s", from)

	utils.ExitWhen(logger, *data != "" && (*abi != "" || len(*abiArgs) > 0), "--data conflicts with --abi and --args")
	var input []byte
	if *data != "" {
		input, err = hex.DecodeString(*data)
		utils.ExitWhenErr(logger, err, "decode data: %v error: %v", *data, err)
	}

	if *abi != "" {
		abiJson := *abi
		// built-in abi
		switch *abi {
		case consts.Erc20:
			logger.Debug().Msgf("use built-in %v abi", *abi)
			abiJson = consts.Erc20Abi
		case consts.Erc721:
			logger.Debug().Msgf("use built-in %v abi", *abi)
			abiJson = consts.Erc721Abi
		case consts.Erc1155:
			logger.Debug().Msgf("use built-in %v abi", *abi)
			abiJson = consts.Erc1155Abi
		default:
			logger.Debug().Msgf("use custom abi")
		}
		abiObj, err := transaction.ParseAbi(abiJson)
		utils.ExitWhenErr(logger, err, "parse abi error: %v", err)

		logger.Debug().Msgf("method: %v", *method)
		logger.Debug().Msgf("args: %v", *abiArgs)

		methodName, paramNames, realArgs, err := transaction.AbiArgs(abiObj, *method, *abiArgs...)
		utils.ExitWhenErr(logger, err, "AbiArgs error: %v", err)

		paramsStr := strings.Join(paramNames, ",")
		functionSignature := fmt.Sprintf("%s(%s)", methodName, paramsStr)
		logger.Info().Msgf("function signature: %v", functionSignature)
		if len(*abiArgs) != 0 {
			logger.Info().Msgf("abi args: %s", *abiArgs)
		}

		input, err = abiObj.Pack(methodName, realArgs...)
		utils.ExitWhenErr(logger, err, "abi pack error: %v", err)

		logger.Trace().Msgf("abi: %s", abiJson)
	}

	mode := ttypes.GasMode(ttypes.GasMode_value[*gasMode])
	// build tx
	tx, err := transaction.BuildTx(ctx, client, from, *to, value, input, mode, *nonce, *chainID, *gasLimit, *gasLimitRatio, *gasRatio, *gasPrice, *tipCap, *feeCap, *all)
	utils.ExitWhenErr(logger, err, "build tx error: %s", err)

	receipt, tx2, err := transaction.SendTx(client, from, tx, privateKey, net, *noconfirm, *confirmations)
	utils.ExitWhenErr(logger, err, "send transaction error: %v", err)

	if receipt != nil {
		utils.ShowReceipt(logger, receipt)
	}

	link := fmt.Sprintf("%v/tx/%v", net.Explorer, tx2.Hash())
	logger.Info().Msgf("tx link: %v", link)

}
