package rm

import (
	"my-ether-tool/cmd/account"
	"my-ether-tool/database"
	"my-ether-tool/utils"

	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "switch account",
	Long:  "switch account",
	Run:   switchAccount,
}

var name *string
var accountIndex *int

func init() {
	account.AccountCmd.AddCommand(switchCmd)

	name = switchCmd.Flags().String("name", "", "account name")
	accountIndex = switchCmd.Flags().Int("account-index", 0, "account index when mnemonic type")
}

func switchAccount(cmd *cobra.Command, args []string) {
	utils.ExitWithMsgWhen(*name == "", "need name\n")

	err := database.SwitchAccount(*name, *accountIndex)
	utils.ExitWhenError(err, "switch account error: %s", err)

}
