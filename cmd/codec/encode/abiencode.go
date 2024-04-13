/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package codec

import (
	"encoding/hex"
	"fmt"
	"os"

	"met/cmd/codec"
	transaction "met/transaction"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var (
	abiStr20  *string
	abiArgs20 *[]string
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
	codec.CodecCmd.AddCommand(abiencodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// abiencodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// abiencodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	abiStr20 = abiencodeCmd.Flags().String("abi", "", "abi string, eg: transfer(address,uint256)")
	abiArgs20 = abiencodeCmd.Flags().StringArray("args", nil, "arguments of abi(--args 0x... --args 200)")

}

func abiEncode(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*abiStr20 == "", "need --abi <abi_string>\n")
	encoded, err := transaction.AbiEncode(*abiStr20, *abiArgs20)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("result: 0x%s\n", hex.EncodeToString(encoded))
}
