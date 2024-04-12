/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package contract

import (
	"context"
	"fmt"
	"my-ether-tool/cmd"
	"my-ether-tool/database"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
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
}

func parseAbi(abiJson string) (*abi.ABI, error) {
	abiObj, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		return nil, err
	}

	return &abiObj, nil
}

func abiArgs(abiObj *abi.ABI, methodName string, args ...string) (string, []interface{}, error) {
	methodNum := len(abiObj.Methods)
	if methodNum == 0 {
		return "", nil, fmt.Errorf("no method found in abi")
	}

	var method *abi.Method
	// 如果abi中只有一个method，那么忽略methodName
	if methodNum == 1 {
		for name, m := range abiObj.Methods {
			if methodName != "" {
				fmt.Printf("ignore method name\n")
			}
			fmt.Printf("use unique method: %v\n", name)
			method = &m
		}
	} else {
		if m, ok := abiObj.Methods[methodName]; ok {
			method = &m
		}
	}

	if method == nil {
		return "", nil, fmt.Errorf("can not get abi method by name: %v", methodName)
	}

	var params []interface{}
	for i, m := range method.Inputs {
		arg := args[i]
		// TODO
		_ = arg
		switch m.Type.T {
		case abi.IntTy:
			// params = append(params)
			// int1 int2 ..
		case abi.UintTy:
			// uint1 uint2 ..
		case abi.BoolTy:
		case abi.StringTy:
		case abi.AddressTy:
		case abi.BytesTy:

		case abi.ArrayTy:
			panic("not support array type")
		case abi.SliceTy:
			panic("not support slice type")
		case abi.TupleTy:
			panic("not support tuple type")
		case abi.FixedBytesTy:
			panic("not support fixedBytes type")
		case abi.HashTy:
			panic("not support hash type")
		case abi.FixedPointTy:
			panic("not support fixedPoint type")
		case abi.FunctionTy:
			panic("not support function type")
		}
	}

	return method.Name, params, nil
}

func ReadContract(ctx context.Context, networkName, contract, abiJson, methodName string, args ...string) (string, error) {
	abiObj, err := parseAbi(abiJson)
	if err != nil {
		return "", fmt.Errorf("parse abi error: %w", err)
	}
	contractAddress := common.HexToAddress(contract)

	log.Info().Msgf("query network: %v", networkName)
	net, err := database.QueryNetworkOrCurrent(networkName)
	if err != nil {
		return "", fmt.Errorf("query network error: %w", err)
	}
	log.Debug().Msgf("network info: %v", net)

	client, err := ethclient.Dial(net.Rpc)
	if err != nil {
		return "", fmt.Errorf("dial rpc error: %w", err)
	}
	defer client.Close()

	methodName, params, err := abiArgs(abiObj, methodName, args...)

	var results []interface{}
	boundContract := bind.NewBoundContract(contractAddress, *abiObj, client, nil, nil)
	boundContract.Call(&bind.CallOpts{Context: ctx}, &results, methodName, params)

	return "", nil
}

func WriteContract(contract, abiJson, methodName string, args ...string) error {
	abiObj, err := parseAbi(abiJson)
	if err != nil {
		return fmt.Errorf("parse abi error: %w", err)
	}
	_ = abiObj

	return nil
}
