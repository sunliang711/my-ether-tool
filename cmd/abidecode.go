/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// abidecodeCmd represents the abidecode command
var abidecodeCmd = &cobra.Command{
	Use:   "abidecode",
	Short: "abi encode",
	Long:  `encode abi`,
	Run:   abiEncode,
}

func init() {
	codecCmd.AddCommand(abidecodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// abidecodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// abidecodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

