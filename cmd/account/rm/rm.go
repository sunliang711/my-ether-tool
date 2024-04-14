package rm

import (
	"met/cmd/account"
	database "met/database"
	utils "met/utils"

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
	logger := utils.GetLogger("removeAccount")

	utils.ExitWhen(logger, *name == "", "need name")

	err := database.RemoveAccount(*name)
	utils.ExitWhenErr(logger, err, "remove account error: %s", err)

}
