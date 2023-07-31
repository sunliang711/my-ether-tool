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

	"my-ether-tool/cmd/tx"
	"my-ether-tool/database"
	"my-ether-tool/transaction"
	ttypes "my-ether-tool/types"
	"my-ether-tool/utils"

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
	// utils.ExitWithMsgWhen(*rpc == "", "need rpc\n")
	utils.ExitWithMsgWhen(*from == "", "need from\n")
	utils.ExitWithMsgWhen(*to == "", "need to\n")
	// utils.ExitWithMsgWhen(*value == "", "need value")

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenError(err, "load network error: %s\n", err)
	rpc := net.Rpc

	fmt.Printf("environment info:\n")
	fmt.Printf("%-20s:%s\n", "network name", net.Name)
	fmt.Printf("%-20s:%s\n", "network rpc", net.Rpc)

	tx, err := transaction.BuildTransaction(rpc, *from, *to, *value, *data, *abi, *abiArgs, *gasLimit, *nonce, *chainID, *gasPrice, *tipCap, *feeCap, *eip1559)
	utils.ExitWhenError(err, "build transaction error: %s\n", err)

	signer := types.NewCancunSigner(tx.ChainId())
	txHash := signer.Hash(tx)
	fmt.Printf("Hash to be signed: %s\n", txHash)

	txJsonBytes, err := tx.MarshalJSON()
	utils.ExitWhenError(err, "Marshal transaction to json error: %s", err)
	fmt.Printf("Transaction json: %s\n", string(txJsonBytes))

	// read signature
	rd := bufio.NewReader(os.Stdin)
	fmt.Printf("Enter signature: ")
	signature, err := rd.ReadString('\n')
	utils.ExitWhenError(err, "Read signature error: %s", err)

	signature = strings.TrimSpace(signature)
	signature = strings.TrimPrefix(signature, "0x")

	signatureBytes, err := hex.DecodeString(signature)
	utils.ExitWhenError(err, "Invalid signature: %s", err)

	// tx + signature
	tx, err = tx.WithSignature(signer, signatureBytes)
	utils.ExitWhenError(err, "Combine signature to transaction error: %s", err)

	txBytes, err := tx.MarshalBinary()
	utils.ExitWhenError(err, "Marshal transaction to binary error: %s\n", err)

	txHex := "0x" + hex.EncodeToString(txBytes)
	id, err := uuid.NewUUID()
	utils.ExitWhenError(err, "create uuid error: %s\n", err)

	jsonRpcData := ttypes.JsonRpcData{
		JsonRpc: "2.0",
		Method:  "eth_sendRawTransaction",
		Params:  []string{txHex},
		Id:      id.String(),
	}
	// send txHex to rpc
	httpClient := utils.NewHttpClient(rpc, 3)
	resp, err := httpClient.PostStruct(nil, &jsonRpcData)
	utils.ExitWhenError(err, "Send raw transaction error: %s", err)

	var jsonRpcResult ttypes.JsonRpcResult
	err = json.NewDecoder(resp.Body).Decode(&jsonRpcResult)
	utils.ExitWhenError(err, "decode json rpc result error: %s", err)

	utils.ExitWithMsgWhen(jsonRpcResult.Id != id.String(), "json rpc id not match")

	explorer := net.Explorer
	if explorer != "" {
		explorer = strings.TrimSuffix(explorer, "/")
		fmt.Printf("Transaction link: %s/tx/%s\n", explorer, jsonRpcResult.Result)
	} else {
		json.NewEncoder(os.Stdout).Encode(&jsonRpcResult)
	}
}
