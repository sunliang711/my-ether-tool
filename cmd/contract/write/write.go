/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package write

import (
	"context"
	"met/cmd/contract"
	utils "met/utils"
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

	nonce         *string
	value         *string
	gasLimitRatio *string
	gasLimit      *string
	gasRatio      *string
	gasPrice      *string
	gasFeeCap     *string
	gasTipCap     *string
	eip1559       *bool

	noconfirm *bool
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

	account = writeCmd.Flags().String("account", "", "account name(empty to use current account)")
	accountIndex = writeCmd.Flags().Uint("accountIndex", 0, "account index")

	nonce = writeCmd.Flags().String("nonce", "", "custom nonce")
	value = writeCmd.Flags().String("value", "0", "custom value")
	gasLimitRatio = writeCmd.Flags().String("gasLimitRatio", "", "gasLimitRatio")
	gasLimit = writeCmd.Flags().String("gasLimit", "", "custom gasLimit")
	gasRatio = writeCmd.Flags().String("gasRatio", "", "gasRatio")
	gasPrice = writeCmd.Flags().String("gasPrice", "", "custom gasPrice")
	gasFeeCap = writeCmd.Flags().String("gasFeeCap", "", "custom gasFeeCap")
	gasTipCap = writeCmd.Flags().String("gasTipCap", "", "custom gasTipCap")

	eip1559 = writeCmd.Flags().Bool("eip1559", true, "eip1559 (use --eip1559=false to disable)")
	noconfirm = writeCmd.Flags().Bool("noconfirm", false, "noconfirm")

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

	err = contract.WriteContract(ctx, network, contractAddress, abi, method, *account, *nonce, *value, *gasLimitRatio, *gasLimit, *gasRatio, *gasPrice, *gasFeeCap, *gasTipCap, *accountIndex, *eip1559, *noconfirm, abiArgs...)
	utils.ExitWhenErr(logger, err, "write contract error: %v", err)
}
