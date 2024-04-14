package current

import (
	"met/cmd/account"
	database "met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "show current account",
	Long:  "show current account",
	Run:   showCurrentAccount,
}

var (
	insecure *bool
)

func init() {
	account.AccountCmd.AddCommand(currentCmd)

	insecure = currentCmd.Flags().Bool("insecure", false, "show private key or mnemonic")
}

func showCurrentAccount(cmd *cobra.Command, args []string) {
	logger := utils.GetLogger("showCurrentAccount")

	current, err := database.CurrentAccount()
	utils.ExitWhenErr(logger, err, "show current account error: %s", err)

	account.ShowAccount(current, *insecure)
}
