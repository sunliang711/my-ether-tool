package rm

import (
	"met/cmd/account"
	database "met/database"
	utils "met/utils"

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
	logger := utils.GetLogger("switchAccount")
	utils.ExitWhen(logger, *name == "", "need name")

	err := database.SwitchAccount(*name, *accountIndex)
	utils.ExitWhenErr(logger, err, "switch account error: %s", err)

}
