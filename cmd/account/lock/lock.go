package lock

import (
	"met/cmd/account"
	"met/database"
	utils "met/utils"

	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:     "lock",
	Aliases: []string{"l"},
	Short:   "lock or unlock account",
	Long:    "lock or unlock account",
	Run:     lockAccount,
}

var (
	name     *string
	unlock   *bool
	password *string
)

func init() {
	account.AccountCmd.AddCommand(lockCmd)

	// 指定名称时,只lock该账号信息,而不是所有账号信息
	name = lockCmd.Flags().String("name", "", "lock specify account instead of all accounts")

	unlock = lockCmd.Flags().BoolP("unlock", "u", false, "unlock account")
	password = lockCmd.Flags().String("password", "", "password")
}

func lockAccount(cmd *cobra.Command, args []string) {
	var (
		logger = utils.GetLogger("listAccount")
		err    error
	)

	if *password == "" {
		// read from stdin
		*password, err = utils.ReadSecret("Enter password:")
		utils.ExitWhenErr(logger, err, "read password error: %s", err)
	}

	if *unlock {
		err = database.UnlockAccount(*name, *password)
	} else {
		err = database.LockAccount(*name, *password)
	}
	utils.ExitWhenErr(logger, err, "(un)lock account error: %s", err)

}
