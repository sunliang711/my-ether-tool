/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package network

import (
	"fmt"
	cmd "met/cmd"
	database "met/database"

	"github.com/spf13/cobra"
)

// NetworkCmd represents the network command
var NetworkCmd = &cobra.Command{
	Use:   "network",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func init() {
	cmd.RootCmd.AddCommand(NetworkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ShowNetwork(network database.Network) {
	fmt.Printf("Network %s\n", network.Name)
	fmt.Printf("Rpc: %s\n", network.Rpc)
	fmt.Printf("Symbol: %s\n", network.Symbol)
	fmt.Printf("Explorer: %s\n", network.Explorer)
	fmt.Printf("Current: %v\n", network.Current)
	fmt.Println()
}
