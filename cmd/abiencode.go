/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "encoding/hex"
	"fmt"
	"my-ether-tool/transaction"
	"my-ether-tool/utils"
	// "math/big"
	// "my-ether-tool/transaction"
	// "my-ether-tool/utils"
	"os"
	// "strings"
	// "github.com/ethereum/go-ethereum/accounts/abi"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/crypto"
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

}

func abiEncode(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*abiStr == "", "need --abi <abi_string>")
	encoded, err := transaction.AbiEncode(*abiStr, *abiArgs)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("result: %s\n", encoded)
}
