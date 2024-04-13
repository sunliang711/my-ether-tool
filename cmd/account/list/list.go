package list

import (
	"met/cmd/account"
	database "met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"show"},
	Short:   "list account",
	Long:    "list account",
	Run:     listAccount,
}

var (
	name    *string
	hdIndex *int
	count   *uint

	insecure *bool
)

func init() {
	account.AccountCmd.AddCommand(listCmd)

	// 指定名称时,只显示该账号信息,而不是所有账号信息
	name = listCmd.Flags().String("name", "", "show specify account instead of all accounts")
	// 如果指定名称,并且该账号是mnemonic类型,那么可以指定hd_index 和count,来显示多个子账号信息,否则使用账号内部的current_index
	// hd_index是-1是个非法的hd path index,这样可以用来表示不知道它,从而用current_index
	hdIndex = listCmd.Flags().Int("hd-index", -1, "hd index")
	count = listCmd.Flags().Uint("count", 1, "sub account count")

	insecure = listCmd.Flags().Bool("insecure", false, "show private key or mnemonic")
}

func listAccount(cmd *cobra.Command, args []string) {
	if *name != "" {
		acc, err := database.QueryAccount(*name)
		utils.ExitWhenError(err, "query account: %s error: %s", *name, err)

		startIndex := acc.CurrentIndex
		if *hdIndex >= 0 {
			startIndex = uint(*hdIndex)
		}

		i := startIndex
		for i < uint(startIndex)+*count {
			subAccount := acc.SwitchTo(i)
			account.ShowAccount(subAccount, *insecure)

			i += 1
		}

	} else {

		accounts, err := database.QueryAllAccounts()
		utils.ExitWhenError(err, "query accounts error: %s", err)

		for i := range accounts {
			account.ShowAccount(accounts[i], *insecure)
		}
	}

}
