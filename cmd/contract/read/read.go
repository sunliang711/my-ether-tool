/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package read

import (
	"fmt"
	"my-ether-tool/cmd/contract"

	"github.com/spf13/cobra"
)

// TxCmd represents the tx command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "read contract",
	Long:  `read contract`,
	Run:   readContract,
}

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
}

func readContract(cmd *cobra.Command, args []string) {
	fmt.Printf("read contract\n")
}
