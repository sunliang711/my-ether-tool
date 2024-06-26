package database

import (
	"fmt"
	"met/utils"

	"gorm.io/gorm"
)

type Account struct {
	Name string `gorm:"unique;"`

	Type string

	Value string

	// Value 字段是否已经加密
	Encrypted bool

	PathFormat string
	Passphrase string

	// 是否是当前账号
	Current bool
	// 当Current为true 并且 Type 是MnemonicType时，所对应的助记词的index
	CurrentIndex uint
}

func (account Account) SwitchTo(newIndex uint) Account {
	newAccount := account
	newAccount.CurrentIndex = newIndex
	return newAccount
}

const (
	AccountTableName = "accounts"
)

func (Account) TableName() string {
	return AccountTableName
}

// op

func QueryAccount(name string) (account Account, err error) {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	err = Conn.WithContext(ctx).Model(&Account{}).First(&account, "name = ?", name).Error
	return
}

func SwitchAccount(name string, index int) (err error) {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	err = Conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Account{}).Where("current = true").Update("current", false).Error
		if err != nil {
			return err
		}

		err = tx.Model(&Account{}).Where("name = ?", name).Updates(map[string]any{"current": true, "current_index": index}).Error
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func AddAccount(account *Account) error {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	_, err := QueryAccount(account.Name)
	if err == gorm.ErrRecordNotFound {
		result := Conn.WithContext(ctx).Create(account)
		err = result.Error
		if err != nil {
			return err
		}

		// switch account
		err = SwitchAccount(account.Name, 0)
		return err
	} else {
		return fmt.Errorf("account: %s already exists", account.Name)
	}
}

func QueryAllAccounts() (accounts []Account, err error) {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	err = Conn.WithContext(ctx).Model(&Account{}).Find(&accounts).Error

	return
}

func RemoveAccount(name string) error {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	_, err := QueryAccount(name)
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("account: %s not exist", name)
	} else {
		result := Conn.WithContext(ctx).Delete(&Account{}, "name = ?", name)
		err = result.Error
		if err != nil {
			return err
		}

		// switch account
		allAccounts, err := QueryAllAccounts()
		if err != nil {
			return err
		}

		if len(allAccounts) > 0 {
			err = SwitchAccount(allAccounts[0].Name, 0)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("no remaining account to switch to after remove\n")
		}

		return nil
	}
}

func CurrentAccount() (account Account, err error) {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	err = Conn.WithContext(ctx).Model(&Account{}).First(&account, "current = true").Error
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("no such account")
	}
	return
}

func QueryAccountOrCurrent(name string, index uint) (*Account, error) {
	var (
		acc    Account
		err    error
		logger = utils.GetLogger("QueryAccountOrCurrent")
	)

	if name != "" {
		logger.Info().Msgf("Query account: %v", name)
		acc, err = QueryAccount(name)
		acc2 := acc.SwitchTo(index)
		return &acc2, err
	}

	logger.Info().Msgf("Query current account")
	acc, err = CurrentAccount()
	if err != nil {
		return nil, err
	}

	return &acc, nil
}

// 那么为空时表示所有
func LockAccount(name string, password string) error {
	var (
		accountList []Account
		logger      = utils.GetLogger("LockAccount")
	)
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	err := Conn.WithContext(ctx).Model(&Account{}).Where(&Account{Name: name}).Find(&accountList).Error
	if err != nil {
		return fmt.Errorf("query account by name: %s error: %w", name, err)
	}

	for _, acc := range accountList {
		if acc.Encrypted {
			logger.Info().Msgf("account: %v already locked,skip", acc.Name)
			continue
		}
		// encrypt
		logger.Info().Msgf("lock account: %v", acc.Name)
		encrypted := utils.Encrypt(password, acc.Value)
		err = Conn.WithContext(ctx).Model(&Account{}).Where(&Account{Name: acc.Name}).Updates(map[string]any{"encrypted": true, "value": encrypted}).Error
		if err != nil {
			return fmt.Errorf("lock account: %v error: %w", acc.Name, err)
		}
	}

	return nil
}

func UnlockAccount(name string, password string) error {

	var (
		accountList []Account
		logger      = utils.GetLogger("UnlockAccount")
	)

	err := Conn.Model(&Account{}).Where(&Account{Name: name}).Find(&accountList).Error
	if err != nil {
		return fmt.Errorf("query account by name: %s error: %w", name, err)
	}

	for _, acc := range accountList {
		if !acc.Encrypted {
			logger.Info().Msgf("account: %v already unlocked,skip", acc.Name)
			continue
		}
		// encrypt
		logger.Info().Msgf("unlock account: %v", acc.Name)
		decrypted := utils.Decrypt(password, acc.Value)
		if decrypted == "" {
			return fmt.Errorf("wrong password")
		}
		err = Conn.Model(&Account{}).Where(&Account{Name: acc.Name}).Updates(map[string]any{"encrypted": false, "value": decrypted}).Error
		if err != nil {
			return fmt.Errorf("lock account: %v error: %w", acc.Name, err)
		}
	}

	return nil
}
