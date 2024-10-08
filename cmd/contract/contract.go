/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package contract

import (
	"context"
	"encoding/hex"
	"fmt"
	cmd "met/cmd"
	database "met/database"
	transaction "met/transaction"
	"met/types"
	utils "met/utils"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

// TxCmd represents the tx command
var ContractCmd = &cobra.Command{
	Use:   "contract",
	Short: "contract related",
	Long:  `contract write or read`,
	Run:   nil,
}

func init() {
	cmd.RootCmd.AddCommand(ContractCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	ContractCmd.PersistentFlags().String("network", "", "network name(empty to use current network)")
	ContractCmd.PersistentFlags().String("contract", "", "contract address")
	ContractCmd.PersistentFlags().String("abi", "", "abi json string(use --abi \"$(cat <FILE>)\" to specify file) or built-in abi(eg: erc20 erc721 erc1155)")
	ContractCmd.PersistentFlags().String("method", "", "method name")
	ContractCmd.PersistentFlags().StringArray("args", nil, "arguments of abi (--args xx1 --args xx2 ...)")
}

// func parseAbi(abiJson string) (*abi.ABI, error) {
// 	abiObj, err := abi.JSON(strings.NewReader(abiJson))
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &abiObj, nil
// }

// // 准备abi中指定method的实际参数
// // 因为args是传递过来的string类型的
// // 要把他们转换成实际的值，比如*big.Int common.Address []byte 等等
// func abiArgs(abiObj *abi.ABI, methodName string, args ...string) (string, []string, []interface{}, error) {
// 	var (
// 		realArgs   []interface{}
// 		paramNames []string
// 		logger     = utils.GetLogger("abiArgs")
// 	)

// 	methodNum := len(abiObj.Methods)
// 	if methodNum == 0 {
// 		return "", nil, nil, fmt.Errorf("no method found in abi")
// 	}

// 	var method *abi.Method
// 	// 如果abi中只有一个method，那么忽略methodName
// 	if methodNum == 1 {
// 		for name, m := range abiObj.Methods {
// 			if methodName != "" {
// 				logger.Debug().Msgf("ignore method name")
// 			}
// 			logger.Debug().Msgf("use unique method: %v", name)
// 			method = &m
// 		}
// 	} else {
// 		if m, ok := abiObj.Methods[methodName]; ok {
// 			method = &m
// 		}
// 	}

// 	if method == nil {
// 		return "", nil, nil, fmt.Errorf("can not get abi method by name: %v", methodName)
// 	}

// 	if len(args) != len(method.Inputs) {
// 		return "", nil, nil, fmt.Errorf("arg count not match abi input count")
// 	}

// 	for i, m := range method.Inputs {
// 		arg := args[i]

// 		v, err := parseAbiType(m.Type, arg)
// 		if err != nil {
// 			return "", nil, nil, err
// 		}
// 		logger.Debug().Msgf("input type: %v, input value: %v", m.Type.String(), arg)

// 		realArgs = append(realArgs, v)
// 		paramNames = append(paramNames, m.Type.String())
// 	}

// 	return method.Name, paramNames, realArgs, nil
// }

// type NameValue struct {
// 	Name  string
// 	Value string
// }

// func parseOutput(abiObj *abi.ABI, methodName string, results []any) ([]NameValue, error) {
// 	methodNum := len(abiObj.Methods)
// 	if methodNum == 0 {
// 		return nil, fmt.Errorf("no method found in abi")
// 	}

// 	var method *abi.Method
// 	// 如果abi中只有一个method，那么忽略methodName
// 	if methodNum == 1 {
// 		for name, m := range abiObj.Methods {
// 			if methodName != "" {
// 				fmt.Printf("ignore method name\n")
// 			}
// 			fmt.Printf("use unique method: %v\n", name)
// 			method = &m
// 		}
// 	} else {
// 		if m, ok := abiObj.Methods[methodName]; ok {
// 			method = &m
// 		}
// 	}

// 	if method == nil {
// 		return nil, fmt.Errorf("can not get abi method by name: %v", methodName)
// 	}

// 	if len(results) != len(method.Outputs) {
// 		return nil, fmt.Errorf("result count not match abi output count")
// 	}

// 	var nameValues []NameValue

// 	for i, output := range method.Outputs {
// 		result := results[i]
// 		r, err := decodeOutput(output.Type, result)
// 		if err != nil {
// 			return nil, err
// 		}

// 		nameValues = append(nameValues, NameValue{
// 			Name:  output.Name,
// 			Value: r,
// 		})

// 	}

// 	return nameValues, nil
// }

func ReadContract(ctx context.Context, client *ethclient.Client, net *database.Network, contract, abiJson, methodName string, args ...string) ([]transaction.NameValue, error) {
	logger := utils.GetLogger("ReadContract")
	logger.Debug().Msgf("abi: %v", abiJson)

	logger.Info().Msg("parse abi")
	abiObj, err := transaction.ParseAbiJson(abiJson)
	if err != nil {
		return nil, fmt.Errorf("parse abi error: %w", err)
	}

	contractAddress := common.HexToAddress(contract)
	logger.Info().Msgf("contract address: %v", contractAddress)

	logger.Debug().Msgf("network info: %v", net)

	logger.Info().Msgf("network: %v", net.Name)
	logger.Info().Msgf("dial rpc: %v", net.Rpc)

	logger.Info().Msg("prepare abi args")
	methodName, _, realArgs, err := transaction.AbiArgs(abiObj, methodName, args...)
	if err != nil {
		return nil, err
	}

	logger.Info().Msgf("call method: %v with args: %v", methodName, realArgs)
	var results []any
	boundContract := bind.NewBoundContract(contractAddress, *abiObj, client, nil, nil)
	err = boundContract.Call(&bind.CallOpts{Context: ctx}, &results, methodName, realArgs...)
	if err != nil {
		return nil, fmt.Errorf("call contract error: %w", err)
	}

	logger.Debug().Msgf("raw results: %v", results)
	outputs, err := transaction.ParseOutput(abiObj, methodName, results)
	if err != nil {
		return nil, err
	}

	logger.Debug().Msgf("parsed output: %v", outputs)

	return outputs, nil
}

func WriteContract(ctx context.Context, client *ethclient.Client, net *database.Network, accountDetails *types.AccountDetails, contract, abiJson, methodName, accountName, nonce, value, gasLimitRatio, gasLimit, gasRatio, gasPrice, gasFeeCap, gasTipCap string, accountIndex uint, eip1559 bool, noconfirm bool, args ...string) error {
	logger := utils.GetLogger("WriteContract")

	logger.Debug().Msgf("network info: %v", net)

	logger.Info().Msgf("network: %v", net.Name)

	addressStr, err := accountDetails.Address()
	if err != nil {
		return fmt.Errorf("get account address error: %w", err)
	}

	logger.Info().Msgf("account info: name: %v address: %v account index: %v", accountDetails.Name, addressStr, accountDetails.CurrentIndex)

	privateKeyStr, err := accountDetails.PrivateKey()
	if err != nil {
		return fmt.Errorf("get account private key error: %w", err)
	}

	pk := strings.TrimPrefix(privateKeyStr, "0x")
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return fmt.Errorf("create private key error: %w", err)
	}

	logger.Info().Msg("query chain id")
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("get chain id error: %w", err)
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return err
	}

	logger.Info().Msgf("parse abi")
	abiObj, err := transaction.ParseAbiJson(abiJson)
	if err != nil {
		return fmt.Errorf("parse abi error: %w", err)
	}

	methodName, paramNames, realArgs, err := transaction.AbiArgs(abiObj, methodName, args...)
	if err != nil {
		return err
	}

	contractAddress := common.HexToAddress(contract)
	logger.Info().Msgf("contract address: %v", contractAddress)

	input, err := abiObj.Pack(methodName, realArgs...)
	if err != nil {
		return fmt.Errorf("abi pack error: %v", err)
	}

	txParams, err := transaction.GetTxParams(ctx, client, addressStr, contract, chainId.String(), nonce, value, gasLimitRatio, gasLimit, gasRatio, gasPrice, gasFeeCap, gasTipCap, eip1559, input)
	if err != nil {
		return fmt.Errorf("GetTxParams error: %w", err)
	}

	transactor.Nonce = txParams.Nonce
	transactor.GasLimit = txParams.GasLimit
	transactor.GasPrice = txParams.GasPrice
	transactor.GasFeeCap = txParams.GasFeeCap
	transactor.GasTipCap = txParams.GasTipCap
	transactor.Value = txParams.Value

	logger.Info().Msgf("From: %v", addressStr)
	logger.Info().Msgf("To: %v (contract)", contract)
	logger.Info().Msgf("Value: %v", value)
	logger.Info().Msgf("Nonce: %v", transactor.Nonce)
	logger.Info().Msgf("GasLimit: %v", transactor.GasLimit)
	logger.Info().Msgf("GasPrice: %v", transactor.GasPrice.String())
	logger.Info().Msgf("GasFeeCap: %v", transactor.GasFeeCap.String())
	logger.Info().Msgf("GasTipCap: %v", transactor.GasTipCap.String())
	logger.Info().Msgf("Method: %v", methodName)
	for i, param := range paramNames {
		logger.Info().Msgf("Arg%d: %v (%v)", i, realArgs[i], param)
	}
	logger.Info().Msgf("Data: %v", hex.EncodeToString(input))

	if !noconfirm {
		input, err := utils.ReadChar("Send ? [y/N] ")
		utils.ExitWhenErr(logger, err, "read input error: %s", err)

		if input != 'y' {
			os.Exit(0)
		}

	}
	boundContract := bind.NewBoundContract(contractAddress, *abiObj, client, client, nil)
	tx, err := boundContract.Transact(transactor, methodName, realArgs...)
	if err != nil {
		return fmt.Errorf("transact error: %v", err)
	}

	ctx2, cancel2 := utils.DefaultTimeoutContext()
	defer cancel2()

	logger.Info().Msgf("waiting for confirmation..")
	receipt, err := bind.WaitMined(ctx2, client, tx)
	if err != nil {
		return fmt.Errorf("get receipt for tx: %v error: %w", tx.Hash(), err)
	}

	utils.ShowReceipt(logger, receipt)

	logger.Info().Msgf("tx hash: %v", tx.Hash())
	logger.Info().Msgf("tx url: %v", fmt.Sprintf("%v/tx/%v", net.Explorer, tx.Hash()))

	return nil
}
