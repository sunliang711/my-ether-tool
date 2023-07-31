/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package tx

import (
	"my-ether-tool/cmd"

	"github.com/spf13/cobra"
)

// TxCmd represents the tx command
var TxCmd = &cobra.Command{
	Use:   "tx",
	Short: "transaction related",
	Long:  `send transaction (write or read)`,
	Run:   nil,
}

func init() {
	cmd.RootCmd.AddCommand(TxCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// txCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// txCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
