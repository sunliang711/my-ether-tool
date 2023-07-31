package add

import (
	"my-ether-tool/cmd/account"
	"my-ether-tool/database"
	"my-ether-tool/types"
	"my-ether-tool/utils"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add account",
	Long:  "add account",
	Run:   addAccount,
}

var (
	name        *string
	accountType *string
	value       *string
	pathFormat  *string
	passphrase  *string
)

func init() {
	account.AccountCmd.AddCommand(addCmd)

	name = addCmd.Flags().String("name", "", "account name")
	accountType = addCmd.Flags().String("type", types.MnemonicType, "account type: 'mnemonic' or 'private key'")
	value = addCmd.Flags().String("value", "", "mnemonic or private key")
	pathFormat = addCmd.Flags().String("path-format", "", "bip32 path format,eg m/44'/60'/0'/0/x (x is placeholder)")
	passphrase = addCmd.Flags().String("passphrase", "", "bip32 passphrase")
}

func addAccount(cmd *cobra.Command, args []string) {
	var err error

	utils.ExitWithMsgWhen(*name == "", "need name\n")
	utils.ExitWithMsgWhen(*value == "", "need value\n")

	if *accountType != types.MnemonicType && *accountType != types.PrivateKeyType {
		utils.ExitWithMsgWhen(true, "invalid account type, use 'mnemonic' or 'private key'\n")
	}

	if *accountType == types.MnemonicType {
		if *pathFormat == "" {
			*pathFormat = types.DefaultHDPath
		}

		// TODO: check hd path
		// accounts.ParseDerivationPath(*pathFormat)

	}

	account := database.Account{
		Name:         *name,
		Type:         *accountType,
		Value:        *value,
		PathFormat:   *pathFormat,
		Passphrase:   *passphrase,
		Current:      false,
		CurrentIndex: 0,
	}

	err = database.AddAccount(&account)
	utils.ExitWhenError(err, "Add account error: %s\n", err)

}
