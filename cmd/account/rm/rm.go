package rm

import (
	"my-ether-tool/cmd/account"
	"my-ether-tool/database"
	"my-ether-tool/utils"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:     "rm",
	Aliases: []string{"remove", "delete", "del"},
	Short:   "rm account",
	Long:    "rm account",
	Run:     removeAccount,
}

var name *string

func init() {
	account.AccountCmd.AddCommand(rmCmd)

	name = rmCmd.Flags().String("name", "", "account name")
}

func removeAccount(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*name == "", "need name\n")

	err := database.RemoveAccount(*name)
	utils.ExitWhenError(err, "remove account error: %s\n", err)

}
