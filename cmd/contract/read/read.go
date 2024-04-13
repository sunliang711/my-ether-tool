/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package read

import (
	"context"
	"met/cmd/contract"
	"met/consts"
	utils "met/utils"
	"time"

	"github.com/spf13/cobra"
)

// TxCmd represents the tx command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "read contract",
	Long:  `read contract`,
	Run:   readContract,
}

var (
// network *string

// contractAddress *string
// abi             *string
// method          *string
// abiArgs         *[]string
)

func init() {
	// cmd.RootCmd.AddCommand(writeCmd)
	contract.ContractCmd.AddCommand(readCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// TODO: 下面的一些参数可以移到上级的ContractCmd中
	// network = readCmd.Flags().String("network", "", "network name")

	// contractAddress = readCmd.Flags().String("contract", "", "contract address")
	// abi = readCmd.Flags().String("abi", "", "abi json")
	// method = readCmd.Flags().String("method", "", "method name")
	// abiArgs = readCmd.Flags().StringArray("args", nil, "arguments of abi (--args xx1 --args xx2 ...)")
}

func readContract(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("readContract")

	network := cmd.Flag("network").Value.String()
	contractAddress := cmd.Flag("contract").Value.String()
	abi := cmd.Flag("abi").Value.String()
	method := cmd.Flag("method").Value.String()
	abiArgs, err := cmd.Flags().GetStringArray("args")
	utils.ExitWhenErr(logger, err, "get args error: %v", err)

	utils.ExitWhen(logger, contractAddress == "", "missing contract")
	utils.ExitWhen(logger, abi == "", "missing abi")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*consts.DefaultTimeout)
	defer cancel()

	outputs, err := contract.ReadContract(ctx, network, contractAddress, abi, method, abiArgs...)
	utils.ExitWhenErr(logger, err, "read contract error: %v", err)

	for _, output := range outputs {
		logger.Info().Msgf("contract output: Name: %v Value: %v", output.Name, output.Value)
	}
}
