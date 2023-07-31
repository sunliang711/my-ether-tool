package list

import (
	"my-ether-tool/cmd/account"
	"my-ether-tool/database"
	"my-ether-tool/utils"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"show"},
	Short:   "list account",
	Long:    "list account",
	Run:     listAccount,
}

var name *string
var insecure *bool

func init() {
	account.AccountCmd.AddCommand(listCmd)

	name = listCmd.Flags().String("name", "", "show specify account instead of all accounts")
	insecure = listCmd.Flags().Bool("insecure", false, "show private key or mnemonic")
}

func listAccount(cmd *cobra.Command, args []string) {
	if *name != "" {
		acc, err := database.QueryAccount(*name)
		utils.ExitWhenError(err, "query account: %s error: %s", *name, err)
		account.ShowAccount(acc, *insecure)
	} else {

		accounts, err := database.QueryAllAccounts()
		utils.ExitWhenError(err, "query accounts error: %s", err)

		for i := range accounts {
			account.ShowAccount(accounts[i], *insecure)
		}
	}

}
