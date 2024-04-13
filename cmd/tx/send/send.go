package tx

import (
	"context"
	"encoding/hex"
	"fmt"
	"my-ether-tool/cmd/tx"
	"my-ether-tool/database"
	"my-ether-tool/transaction"
	ttypes "my-ether-tool/types"
	"my-ether-tool/utils"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	abiArgs *[]string

	nonce   *string
	chainID *int

	gasRatio *string
	gasLimit *uint64
	eip1559  *bool
	gasPrice *string
	tipCap   *string
	feeCap   *string

	noconfirm *bool
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
	abi = sendCmd.Flags().String("abi", "", "abi string, eg: transfer(address,uint256), conflict with --data")
	abiArgs = sendCmd.Flags().StringArray("args", nil, "arguments of abi( --args 0x... --args 200)")

	gasLimit = sendCmd.Flags().Uint64("gasLimit", 0, "gas limit")
	nonce = sendCmd.Flags().String("nonce", "", "nonce")
	chainID = sendCmd.Flags().Int("chainID", 0, "chain id")
	gasPrice = sendCmd.Flags().String("gasPrice", "", "gas price(gwei)")
	tipCap = sendCmd.Flags().String("tipCap", "", "tipCap(gwei)")
	feeCap = sendCmd.Flags().String("feeCap", "", "feeCap(gwei)")
	eip1559 = sendCmd.Flags().Bool("eip1559", true, "eip1559 switch")
	gasRatio = sendCmd.Flags().String("gasRatio", "", "gasRatio")

	noconfirm = sendCmd.Flags().BoolP("noconfirm", "y", false, "do not need to confirm")
}

func sendTransaction(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("sendTransaction")

	account, err := database.QueryAccountOrCurrent(*account, *accountIndex)
	utils.ExitWhenErr(logger, err, "load account error: %s", err)

	details, err := ttypes.AccountToDetails(account)
	utils.ExitWhenErr(logger, err, "calculate address error: %s", err)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(details.PrivateKey, "0x"))
	utils.ExitWhenErr(logger, err, "parse privateKey error: %s", err)

	from := details.Address

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "load network error: %s", err)
	rpc := net.Rpc

	logger.Info().Msgf("network name: %s", net.Name)
	logger.Info().Msgf("network rpc: %s", net.Rpc)
	logger.Info().Msgf("account name: %s", details.Name)
	logger.Info().Msgf("address: %s", details.Address)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// build tx
	tx, err := transaction.BuildTransaction(ctx, rpc, from, *to, value, *data, *abi, *abiArgs, *gasLimit, *nonce, *chainID, *gasRatio, *gasPrice, *tipCap, *feeCap, *eip1559, *all)
	utils.ExitWhenErr(logger, err, "build tx error: %s", err)

	client, err := ethclient.Dial(rpc)
	if err != nil {
		return
	}
	defer client.Close()

	signer := types.LatestSignerForChainID(tx.ChainId())
	txHash := signer.Hash(tx)
	logger.Debug().Msgf("tx hash to be signed: %s", txHash)

	// sign tx
	logger.Debug().Msgf("sign transaction")
	tx, err = types.SignTx(tx, signer, privateKey)
	utils.ExitWhenErr(logger, err, "sign tx error: %s", err)

	logger.Info().Msgf("transaction to be sent")
	logger.Info().Msgf("From: %s", from)
	logger.Info().Msgf("To: %s", tx.To())
	logger.Info().Msgf("Value: %s %s", *value, net.Symbol)

	if *abi != "" {
		logger.Info().Msgf("abi: %s", *abi)
	}
	if len(*abiArgs) != 0 {
		logger.Info().Msgf("abi args: %s", *abiArgs)
	}
	logger.Info().Msgf("Data: %s", hex.EncodeToString(tx.Data()))
	logger.Info().Msgf("Nonce: %v", tx.Nonce())
	logger.Info().Msgf("ChainId: %s", tx.ChainId())
	logger.Info().Msgf("GasLimit: %v", tx.Gas())

	gasPrice, err := utils.Wei2Gwei(tx.GasPrice().String())
	utils.ExitWhenErr(logger, err, "convert wei to gwei error: %s", err)
	tipCap, err := utils.Wei2Gwei(tx.GasTipCap().String())
	utils.ExitWhenErr(logger, err, "convert wei to gwei error: %s", err)
	feeCap, err := utils.Wei2Gwei(tx.GasFeeCap().String())
	utils.ExitWhenErr(logger, err, "convert wei to gwei error: %s", err)

	logger.Info().Msgf("GasPrice: %s Gwei", gasPrice)
	logger.Info().Msgf("GasTipCap: %s Gwei", tipCap)
	logger.Info().Msgf("GasFeeCap: %s Gwei", feeCap)

	if !*noconfirm {
		input, err := utils.ReadChar("Send? [y/N] ")
		utils.ExitWhenErr(logger, err, "read input error: %s", err)

		if input != 'y' {
			os.Exit(0)
		}
	}

	logger.Info().Msgf("send tx")
	err = client.SendTransaction(ctx, tx)
	if err != nil {
		fmt.Printf("send tx error: %v\n", err)
		return
	}

	logger.Info().Msgf("waiting for confirmation")
	bind.WaitMined(ctx, client, tx)

	logger.Debug().Msgf("get receipt")
	receipt, err := client.TransactionReceipt(ctx, tx.Hash())
	if err != nil {
		logger.Error().Err(err).Msgf("get receipt")
	} else {
		logger.Info().Msgf("receipt: %v", receipt) // TODO: format receipt
	}

	link := fmt.Sprintf("%v/tx/%v", net.Explorer, tx.Hash())
	logger.Info().Msgf("tx link: %v", link)

	// send transaction
	// txBytes, err := tx.MarshalBinary()
	// utils.ExitWhenError(err, "Marshal transaction to binary error: %s\n", err)

	// txHex := "0x" + hex.EncodeToString(txBytes)
	// fmt.Printf("%-20s:%s\n", "Raw transaction", txHex)
	// id, err := uuid.NewUUID()
	// utils.ExitWhenError(err, "create uuid error: %s\n", err)
	// jsonRpcData := ttypes.JsonRpcData{
	// 	JsonRpc: "2.0",
	// 	Method:  "eth_sendRawTransaction",
	// 	Params:  []string{txHex},
	// 	Id:      id.String(),
	// }
	// // send txHex to rpc
	// httpClient := utils.NewHttpClient(rpc, 3)
	// fmt.Printf("Send tx: %s to rpc..\n", txHex)
	// resp, err := httpClient.PostStruct(nil, &jsonRpcData)
	// utils.ExitWhenError(err, "Send raw transaction error: %s", err)

	// var jsonRpcResult ttypes.JsonRpcResult
	// fmt.Printf("Decode result..\n")
	// err = json.NewDecoder(resp.Body).Decode(&jsonRpcResult)
	// utils.ExitWhenError(err, "decode json rpc result error: %s", err)

	// utils.ExitWithMsgWhen(jsonRpcResult.Id != id.String(), "json rpc id not match")

	// utils.ExitWithMsgWhen(jsonRpcResult.Result == "", "json rpc no result: %+v\n", jsonRpcResult)

	// explorer := net.Explorer
	// if explorer != "" {
	// 	explorer = strings.TrimSuffix(explorer, "/")
	// 	fmt.Printf("Transaction link: %s/tx/%s\n", explorer, jsonRpcResult.Result)
	// } else {
	// 	json.NewEncoder(os.Stdout).Encode(&jsonRpcResult)
	// }

}
