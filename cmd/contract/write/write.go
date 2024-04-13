/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package write

import (
	"context"
	"my-ether-tool/cmd/contract"
	"my-ether-tool/utils"
	"time"

	"github.com/spf13/cobra"
)

// TxCmd represents the tx command
var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "write contract",
	Long:  `write contract`,
	Run:   writeContract,
}

var (
	account      *string
	accountIndex *uint
)

func init() {
	// cmd.RootCmd.AddCommand(writeCmd)
	contract.ContractCmd.AddCommand(writeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	account = writeCmd.Flags().String("account", "", "account name")
	accountIndex = writeCmd.Flags().Uint("accountIndex", 0, "account index")

}

func writeContract(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("readContract")

	network := cmd.Flag("network").Value.String()
	contractAddress := cmd.Flag("contract").Value.String()
	abi := cmd.Flag("abi").Value.String()
	method := cmd.Flag("method").Value.String()
	abiArgs, err := cmd.Flags().GetStringArray("args")
	utils.ExitWhenErr(logger, err, "get args error: %v", err)

	utils.ExitWhen(logger, contractAddress == "", "missing contract")
	utils.ExitWhen(logger, abi == "", "missing abi")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = contract.WriteContract(ctx, network, contractAddress, abi, method, *account, *accountIndex, abiArgs...)
	utils.ExitWhenErr(logger, err, "write contract error: %v", err)
}
