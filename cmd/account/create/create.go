package create

import (
	"fmt"
	"my-ether-tool/cmd/account"
	"my-ether-tool/database"
	"my-ether-tool/hd"
	"my-ether-tool/types"
	"my-ether-tool/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
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
	account.AccountCmd.AddCommand(createCmd)

	name = createCmd.Flags().String("name", "", "account name, leave it empty for temp use")
	accountType = createCmd.Flags().String("type", types.MnemonicType, "account type, available type: 'mnemonic' or 'private key' ")
	words = createCmd.Flags().Uint8("words", 12, "mnemonic words count when type is mnemonic")
	passphrase = createCmd.Flags().String("passphrase", "", "passphrase when type is mnemonic")

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

	fullAccount, err := types.AccountToFullAccount(newAccount)
	utils.ExitWhenError(err, "calculate address error: %s", err)

	if *name != "" {
		// save
		err := database.AddAccount(newAccount)
		utils.ExitWhenError(err, "add account to db error: %s\n", err)

		// show address
		fmt.Printf("Created Account Info:\n")
		fmt.Printf("Account Type: %s\n", *accountType)
		fmt.Printf("Address: %s\n", fullAccount.Address)

	} else {
		// show with mnemonic or private key
		fmt.Printf("Created Account Info:\n")
		fmt.Printf("Account Type: %s\n", *accountType)
		if *accountType == types.MnemonicType {
			fmt.Printf("Mnemonic: %s\n", newAccount.Value)
			fmt.Printf("Path: %s\n", fullAccount.Path)
		} else {
			fmt.Printf("Private Key: %s\n", newAccount.Value)
		}
		fmt.Printf("Address: %s\n", fullAccount.Address)
	}

}
