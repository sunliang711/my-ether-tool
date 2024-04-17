package add

import (
	"errors"
	"fmt"

	"met/cmd/account"
	database "met/database"
	hd "met/hd"
	types "met/types"
	utils "met/utils"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var importCmd = &cobra.Command{
	Use:        "import",
	ArgAliases: []string{"im"},
	Short:      "import account",
	Long:       "import account",
	Run:        importAccount,
}

var (
	name        *string
	accountType *string
	value       *string
	pathFormat  *string
	passphrase  *string
)

func init() {
	account.AccountCmd.AddCommand(importCmd)

	name = importCmd.Flags().String("name", "", "account name")
	accountType = importCmd.Flags().String("type", types.MnemonicType, "account type: 'mnemonic' or 'private key'")
	value = importCmd.Flags().String("value", "", "mnemonic or private key")
	pathFormat = importCmd.Flags().String("path-format", "", "bip32 path format,eg m/44'/60'/0'/0/x (x is placeholder)")
	passphrase = importCmd.Flags().String("passphrase", "", "bip32 passphrase")
}

func importAccount(cmd *cobra.Command, args []string) {
	var (
		err    error
		logger = utils.GetLogger("importAccount")
	)

	utils.ExitWhen(logger, *name == "", "need name")
	// utils.ExitWithMsgWhen(*value == "", "need value\n")

	_, err = database.QueryAccount(*name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ok
		} else {
			// other error
			utils.ExitWhenErr(logger, err, "query acocunt: %v error: %v", *name, err)
		}
	} else {
		utils.ExitWhen(logger, true, "account: %v already exist", *name)
	}

	if *accountType != types.MnemonicType && *accountType != types.PrivateKeyType {
		utils.ExitWhen(logger, true, "invalid account type, use 'mnemonic' or 'private key'")
	}

	if *value == "" {
		*value, err = utils.ReadSecret(fmt.Sprintf("Enter %s: ", *accountType))
		utils.ExitWhenErr(logger, err, "Read user input error: %s", err)
	}

	if *accountType == types.MnemonicType {
		if *pathFormat == "" {
			*pathFormat = types.DefaultHDPath
		}

		// check hd path
		err = hd.CheckHdPath(*pathFormat)
		utils.ExitWhenErr(logger, err, "invalid hd path: %s", err)

	}

	account := database.Account{
		Name:         *name,
		Type:         *accountType,
		Value:        *value,
		Encrypted:    false,
		PathFormat:   *pathFormat,
		Passphrase:   *passphrase,
		Current:      false,
		CurrentIndex: 0,
	}

	details, err := types.AccountToDetails(&account)
	utils.ExitWhenErr(logger, err, "invalid data")

	err = database.AddAccount(&account)
	utils.ExitWhenErr(logger, err, "Add account error: %s", err)

	logger.Info().Msgf("Account imported")
	logger.Info().Msgf("Account Name: %v", details.Name)
	logger.Info().Msgf("Account Type: %v", details.Type)

	addressStr, err := details.Address()
	utils.ExitWhenErr(logger, err, "get address error: %v", err)
	logger.Info().Msgf("Account Address: %v", addressStr)
}
