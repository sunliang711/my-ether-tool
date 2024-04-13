package new

import (
	"fmt"
	"met/cmd/account"
	database "met/database"
	hd "met/hd"
	types "met/types"
	utils "met/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create"},
	Short:   "create a new account",
	Long:    "create a new account",
	Run:     createWallet,
}

var (
	name        *string
	accountType *string
	words       *uint8
	passphrase  *string
)

func init() {
	account.AccountCmd.AddCommand(newCmd)

	name = newCmd.Flags().String("name", "", "account name, leave it empty for temp use")
	accountType = newCmd.Flags().String("type", types.MnemonicType, "account type, available type: 'mnemonic' or 'private key' ")
	words = newCmd.Flags().Uint8("words", 12, "mnemonic words count when type is mnemonic")
	passphrase = newCmd.Flags().String("passphrase", "", "passphrase when type is mnemonic")

}

func createWallet(cmd *cobra.Command, args []string) {
	var newAccount *database.Account

	switch *accountType {
	case types.MnemonicType:
		mnemonic, err := hd.CreateMnemonic(*words)
		utils.ExitWhenError(err, "create mnemonic error: %s", err)

		newAccount = &database.Account{
			Name:       *name,
			Type:       *accountType,
			Value:      mnemonic,
			PathFormat: types.DefaultHDPath,
			Passphrase: *passphrase,
		}

	case types.PrivateKeyType:
		privateKey, err := crypto.GenerateKey()
		utils.ExitWhenError(err, "generate private key error: %s\n", err)
		privateKeyBytes := crypto.FromECDSA(privateKey)

		newAccount = &database.Account{
			Name:  *name,
			Type:  *accountType,
			Value: hexutil.Encode(privateKeyBytes),
		}
	default:
		utils.ExitWithMsgWhen(true, "invalid account type, use 'mnemonic' or 'private key'\n")
	}

	fullAccount, err := types.AccountToDetails(newAccount)
	utils.ExitWhenError(err, "calculate address error: %s", err)

	if *name != "" {
		// save
		err := database.AddAccount(newAccount)
		utils.ExitWhenError(err, "add account to db error: %s\n", err)

		// query
		newAccount, err := database.QueryAccount(*name)
		utils.ExitWhenError(err, "query account by name: %s error: %s\n", *name, err)

		// format
		fullAccount, err = types.AccountToDetails(&newAccount)
		utils.ExitWhenError(err, "calculate address error: %s", err)
	}

	fmt.Printf("%s\n", fullAccount.AsString(*name == ""))

}
