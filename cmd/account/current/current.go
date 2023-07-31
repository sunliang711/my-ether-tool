package current

import (
	"my-ether-tool/cmd/account"
	"my-ether-tool/database"
	"my-ether-tool/utils"

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
	current, err := database.CurrentAccount()
	utils.ExitWhenError(err, "show current account error: %s", err)

	account.ShowAccount(current, *insecure)
}
