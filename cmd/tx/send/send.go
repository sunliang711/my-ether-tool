package tx

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"my-ether-tool/cmd/tx"
	"my-ether-tool/database"
	"my-ether-tool/transaction"
	ttypes "my-ether-tool/types"
	"my-ether-tool/utils"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
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

	noconfirm = sendCmd.Flags().BoolP("noconfirm", "y", false, "do not need to confirm")
}

func sendTransaction(cmd *cobra.Command, args []string) {

	account, err := database.QueryAccountOrCurrent(*account, *accountIndex)
	utils.ExitWhenError(err, "load account error: %s\n", err)

	details, err := ttypes.AccountToDetails(account)
	utils.ExitWhenError(err, "calculate address error: %s\n", err)

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(details.PrivateKey, "0x"))
	utils.ExitWhenError(err, "parse privateKey error: %s\n", err)

	from := details.Address

	net, err := database.QueryNetworkOrCurrent(*network)
	utils.ExitWhenError(err, "load network error: %s\n", err)
	rpc := net.Rpc

	fmt.Printf("environment info:\n")
	fmt.Printf("%-20s:%s\n", "network name", net.Name)
	fmt.Printf("%-20s:%s\n", "network rpc", net.Rpc)
	fmt.Printf("%-20s:%s\n", "account name", details.Name)
	if details.Type == ttypes.MnemonicType {
		fmt.Printf("%-20s:%s\n", "hd path", details.Path)
	}
	fmt.Printf("%-20s:%s\n", "address", details.Address)

	// build tx
	tx, err := transaction.BuildTransaction(rpc, from, *to, *value, *data, *abi, *abiArgs, *gasLimit, *nonce, *chainID, *gasPrice, *tipCap, *feeCap, *eip1559)
	utils.ExitWhenError(err, "build tx error: %s\n", err)

	signer := types.LatestSignerForChainID(tx.ChainId())
	txHash := signer.Hash(tx)
	fmt.Printf("Hash to be signed: %s\n", txHash)

	// sign tx
	fmt.Printf("Sign transaction..\n")
	tx, err = types.SignTx(tx, signer, privateKey)
	utils.ExitWhenError(err, "sign tx error: %s\n", err)

	fmt.Printf("Transaction to be sent:\n")
	fmt.Printf("%-20s:%s\n", "From", from)
	fmt.Printf("%-20s:%s\n", "To", tx.To())
	fmt.Printf("%-20s:%s %s\n", "Value", *value, net.Symbol)
	if *abi != "" {
		fmt.Printf("%-20s:%s\n", "abi", *abi)
	}
	if len(*abiArgs) != 0 {
		fmt.Printf("%-20s:%v\n", "abi args", *abiArgs)
	}
	fmt.Printf("%-20s:%s\n", "Data", hex.EncodeToString(tx.Data()))
	fmt.Printf("%-20s:%d\n", "Nonce", tx.Nonce())
	fmt.Printf("%-20s:%s\n", "ChainId", tx.ChainId())
	fmt.Printf("%-20s:%d\n", "GasLimit", tx.Gas())

	gasPrice, err := utils.Wei2Gwei(tx.GasPrice().String())
	utils.ExitWhenError(err, "convert wei to gwei error: %s\n", err)
	tipCap, err := utils.Wei2Gwei(tx.GasTipCap().String())
	utils.ExitWhenError(err, "convert wei to gwei error: %s\n", err)
	feeCap, err := utils.Wei2Gwei(tx.GasFeeCap().String())
	utils.ExitWhenError(err, "convert wei to gwei error: %s\n", err)

	fmt.Printf("%-20s:%s gwei\n", "GasPrice", gasPrice)
	fmt.Printf("%-20s:%s gwei\n", "GasTipCap", tipCap)
	fmt.Printf("%-20s:%s gwei\n", "GasFeeCap", feeCap)

	// send transaction
	txBytes, err := tx.MarshalBinary()
	utils.ExitWhenError(err, "Marshal transaction to binary error: %s\n", err)

	txHex := "0x" + hex.EncodeToString(txBytes)
	fmt.Printf("%-20s:%s\n", "Raw transaction", txHex)
	id, err := uuid.NewUUID()
	utils.ExitWhenError(err, "create uuid error: %s\n", err)

	if !*noconfirm {
		input, err := utils.ReadChar("Send? [y/N] ")
		utils.ExitWhenError(err, "read input error: %s\n", err)

		if input != 'y' {
			os.Exit(0)
		}

	}

	jsonRpcData := ttypes.JsonRpcData{
		JsonRpc: "2.0",
		Method:  "eth_sendRawTransaction",
		Params:  []string{txHex},
		Id:      id.String(),
	}
	// send txHex to rpc
	httpClient := utils.NewHttpClient(rpc, 3)
	fmt.Printf("Send tx: %s to rpc..\n", txHex)
	resp, err := httpClient.PostStruct(nil, &jsonRpcData)
	utils.ExitWhenError(err, "Send raw transaction error: %s", err)

	var jsonRpcResult ttypes.JsonRpcResult
	fmt.Printf("Decode result..\n")
	err = json.NewDecoder(resp.Body).Decode(&jsonRpcResult)
	utils.ExitWhenError(err, "decode json rpc result error: %s", err)

	utils.ExitWithMsgWhen(jsonRpcResult.Id != id.String(), "json rpc id not match")

	utils.ExitWithMsgWhen(jsonRpcResult.Result == "", "json rpc no result: %+v\n", jsonRpcResult)

	explorer := net.Explorer
	if explorer != "" {
		explorer = strings.TrimSuffix(explorer, "/")
		fmt.Printf("Transaction link: %s/tx/%s\n", explorer, jsonRpcResult.Result)
	} else {
		json.NewEncoder(os.Stdout).Encode(&jsonRpcResult)
	}

}
