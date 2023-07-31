/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package account

import (
	"fmt"
	"my-ether-tool/cmd"
	"my-ether-tool/database"

	"github.com/spf13/cobra"
)

// AccountCmd represents the account command
var AccountCmd = &cobra.Command{
	Use:   "account",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("account called")
	// },
}

func init() {
	cmd.RootCmd.AddCommand(AccountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// accountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// accountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ShowAccount(account database.Account, insecure bool) {
	fmt.Printf("Account: %s\n", account.Name)
	fmt.Printf("Address: %s\n", "TODO")
	fmt.Printf("Type: %s\n", account.Type)
	fmt.Printf("Private key: %s\n", account.Value)
	fmt.Printf("Mnemonic: %s\n", account.Value)
	fmt.Printf("PathFormat: %s\n", account.PathFormat)
	fmt.Printf("Passphrase: %s\n", account.Passphrase)
	fmt.Printf("Current: %v\n", account.Current)
	fmt.Printf("Current Index: %v\n", account.CurrentIndex)
	fmt.Println()

}
