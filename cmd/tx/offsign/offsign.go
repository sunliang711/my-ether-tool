/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package tx

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"met/cmd/tx"
	database "met/database"
	transaction "met/transaction"
	ttypes "met/types"
	utils "met/utils"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// offsignCmd represents the offsign command
var offsignCmd = &cobra.Command{
	Use:   "offsign",
	Short: "build tx to be signed by other tool then send it",
	Long:  `build a transaction to be signed by other tool,then combine tx and the signature signed by other tool to a raw tx, then send it`,
	Run:   offsign,
}

var (
	network *string
	// rpc   *string
	from  *string
	to    *string
	value *string

	data    *string
	abi     *string
	abiArgs *[]string

	nonce   *string
	chainID *int

	gasLimit *uint64
	eip1559  *bool
	gasPrice *string
	tipCap   *string
	feeCap   *string

	// explorer *string
)

func init() {
	tx.TxCmd.AddCommand(offsignCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// offsignCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// offsignCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// rpc = offsignCmd.Flags().String("rpc", "", "rpc url")
	network = offsignCmd.Flags().String("network", "", "network name")
	from = offsignCmd.Flags().String("from", "", "from address")
	to = offsignCmd.Flags().String("to", "", "receiver address")
	value = offsignCmd.Flags().String("value", "0", "value (uint: eth)")

	// data or abi + args
	data = offsignCmd.Flags().String("data", "", "data of transaction, conflict with --abi")
	abi = offsignCmd.Flags().String("abi", "", "abi string, eg: transfer(address,uint256), conflict with --data")
	abiArgs = offsignCmd.Flags().StringArray("args", nil, "arguments of abi( --args 0x... --args 200)")

	gasLimit = offsignCmd.Flags().Uint64("gasLimit", 0, "gas limit")
	nonce = offsignCmd.Flags().String("nonce", "", "nonce")
	chainID = offsignCmd.Flags().Int("chainID", 0, "chain id")
	gasPrice = offsignCmd.Flags().String("gasPrice", "", "gas price(gwei)")
	tipCap = offsignCmd.Flags().String("tipCap", "", "tipCap(gwei)")
	feeCap = offsignCmd.Flags().String("feeCap", "", "feeCap(gwei)")
	eip1559 = offsignCmd.Flags().Bool("eip1559", true, "eip1559 switch")

	// explorer = offsignCmd.Flags().String("explorer", "", "explorer url")
}

func offsign(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("offsign")
	// utils.ExitWithMsgWhen(*rpc == "", "need rpc\n")
	utils.ExitWhen(logger, *from == "", "need from")
	utils.ExitWhen(logger, *to == "", "need to")
	// utils.ExitWithMsgWhen(*value == "", "need value")

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenErr(logger, err, "load network error: %s", err)
	rpc := net.Rpc

	fmt.Printf("environment info:\n")
	fmt.Printf("%-20s:%s\n", "network name", net.Name)
	fmt.Printf("%-20s:%s\n", "network rpc", net.Rpc)

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	client, err := utils.DialRpc(ctx, net.Rpc)
	utils.ExitWhenErr(logger, err, "dial rpc error: %v", err)
	defer client.Close()

	tx, err := transaction.BuildTransaction(ctx, client, *from, *to, value, *data, *abi, *abiArgs, *gasLimit, *nonce, *chainID, "", *gasPrice, *tipCap, *feeCap, *eip1559, false)
	utils.ExitWhenErr(logger, err, "build transaction error: %s", err)

	signer := types.NewCancunSigner(tx.ChainId())
	txHash := signer.Hash(tx)
	fmt.Printf("Hash to be signed: %s\n", txHash)

	txJsonBytes, err := tx.MarshalJSON()
	utils.ExitWhenErr(logger, err, "Marshal transaction to json error: %s", err)
	fmt.Printf("Transaction json: %s\n", string(txJsonBytes))

	// read signature
	rd := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter signature: ")
	signature, err := rd.ReadString('\n')
	utils.ExitWhenErr(logger, err, "Read signature error: %s", err)

	signature = strings.TrimSpace(signature)
	signature = strings.TrimPrefix(signature, "0x")

	signatureBytes, err := hex.DecodeString(signature)
	utils.ExitWhenErr(logger, err, "Invalid signature: %s", err)

	// tx + signature
	tx, err = tx.WithSignature(signer, signatureBytes)
	utils.ExitWhenErr(logger, err, "Combine signature to transaction error: %s", err)

	txBytes, err := tx.MarshalBinary()
	utils.ExitWhenErr(logger, err, "Marshal transaction to binary error: %s", err)

	txHex := "0x" + hex.EncodeToString(txBytes)
	id, err := uuid.NewUUID()
	utils.ExitWhenErr(logger, err, "create uuid error: %s", err)

	jsonRpcData := ttypes.JsonRpcData{
		JsonRpc: "2.0",
		Method:  "eth_sendRawTransaction",
		Params:  []string{txHex},
		Id:      id.String(),
	}
	// send txHex to rpc
	httpClient := utils.NewHttpClient(rpc, 3)
	resp, err := httpClient.PostStruct(nil, &jsonRpcData)
	utils.ExitWhenErr(logger, err, "Send raw transaction error: %s", err)

	var jsonRpcResult ttypes.JsonRpcResult
	err = json.NewDecoder(resp.Body).Decode(&jsonRpcResult)
	utils.ExitWhenErr(logger, err, "decode json rpc result error: %s", err)

	utils.ExitWhen(logger, jsonRpcResult.Id != id.String(), "json rpc id not match")

	explorer := net.Explorer
	if explorer != "" {
		explorer = strings.TrimSuffix(explorer, "/")
		// fmt.Printf("Transaction link: %s/tx/%s\n", explorer, jsonRpcResult.Result)
		logger.Info().Msgf("Transaction link: %s/tx/%x", explorer, jsonRpcResult.Result)
	} else {
		json.NewEncoder(os.Stdout).Encode(&jsonRpcResult)
	}
}
