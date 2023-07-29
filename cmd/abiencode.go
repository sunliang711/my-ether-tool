/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"my-ether-tool/utils"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var (
	abiStr  *string
	abiArgs *[]string
)

// abiencodeCmd represents the abiencode command
var abiencodeCmd = &cobra.Command{
	Use:   "abiencode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: abiEncode,
}

func init() {
	codecCmd.AddCommand(abiencodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// abiencodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// abiencodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	abiStr = abiencodeCmd.Flags().String("abi", "", "abi string, eg: transfer(address,uint256)")
	abiArgs = abiencodeCmd.Flags().StringArray("args", nil, "arguments of abi(--args 0x... --args 200)")

	var err error
	uint256Type, err = abi.NewType("uint256", "", nil)
	if err != nil {
		panic("create uint256 failed")
	}
	bytes32Type, err = abi.NewType("bytes32", "", nil)
	if err != nil {
		panic("create bytes32 failed")
	}
	addressType, err = abi.NewType("address", "", nil)
	if err != nil {
		panic("create address failed")
	}
}

var (
	uint256Type abi.Type
	bytes32Type abi.Type
	addressType abi.Type
)

func abiEncode(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*abiStr == "", "need --abi <abi_string>")

	// selector is the first 4 bytes of hash of abiStr
	selector := crypto.Keccak256Hash([]byte(*abiStr)).Bytes()[:4]

	result := selector

	fields := strings.FieldsFunc(*abiStr, func(r rune) bool {
		return r == '(' || r == ')'
	})

	fmt.Printf("fields: %v\n", fields)

	arguments := abi.Arguments{}
	argumentsValue := []any{}

	if len(fields) == 2 {
		argTypes := fields[1]

		types := strings.Split(argTypes, ",")

		utils.ExitWithMsgWhen(len(types) != len(*abiArgs), "number of arguments  not match with abi string\n")

		for i := range types {
			abiArg := (*abiArgs)[i]
			fmt.Printf("> abi arg: %s\n", abiArg)

			switch types[i] {
			case "address":
				fmt.Printf("> address type\n")
				arguments = append(arguments, abi.Argument{Type: addressType})
				argumentsValue = append(argumentsValue, common.HexToAddress(abiArg))
			case "uint256":
				fmt.Printf("> uint256 type\n")
				arguments = append(arguments, abi.Argument{Type: uint256Type})
				v, ok := new(big.Int).SetString(abiArg, 10)
				utils.ExitWithMsgWhen(!ok, "invalid uint256 type argument\n")
				argumentsValue = append(argumentsValue, v)
			case "bytes32":
				fmt.Printf("> bytes32 type\n")
				arguments = append(arguments, abi.Argument{Type: bytes32Type})
				decoded, err := hex.DecodeString(abiArg)
				utils.ExitWithMsgWhen(err != nil, "invalid bytes32 type argument\n")
				argumentsValue = append(argumentsValue, decoded)
			// TODO: other types
			default:
				fmt.Printf("donot support type: %s", types[i])
			}
		}
		packed, err := arguments.Pack(argumentsValue...)
		utils.ExitWithMsgWhen(err != nil, "pack arguments error: %s\n", err)
		result = append(result, packed...)

	}

	fmt.Printf("result: 0x%s\n", hex.EncodeToString(result))
}
