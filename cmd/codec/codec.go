/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package codec

import (
	"my-ether-tool/cmd"

	"github.com/spf13/cobra"
)

// CodecCmd represents the codec command
var CodecCmd = &cobra.Command{
	Use:   "codec",
	Short: "encode and decode in ethereum",
	Long:  `abi encode and decode`,
	Run:   nil,
}

func init() {
	cmd.RootCmd.AddCommand(CodecCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// codecCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// codecCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
