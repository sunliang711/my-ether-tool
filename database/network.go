package database

import (
	"fmt"
	utils "met/utils"

	"gorm.io/gorm"
)

type Network struct {
	Name string `gorm:"unique;"`
	Rpc  string
	// native token symbol, eg: ETH BNB
	Symbol string

	Explorer string

	Current bool
}

const (
	NetworkTableName = "networks"
)

func (Network) TableName() string {
	return NetworkTableName
}

// op
func QueryNetwork(name string) (network Network, err error) {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	err = Conn.WithContext(ctx).Model(&Network{}).First(&network, "name = ?", name).Error

	return
}

func SwitchNetwork(name string) error {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	logger := utils.GetLogger("SwitchNetwork")
	logger.Info().Msgf("switch to network: %v", name)

	err := Conn.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&Network{}).Where("current = true").Update("current", false).Error
		if err != nil {
			return err
		}

		err = tx.Model(&Network{}).Where("name = ?", name).Update("current", true).Error
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

func AddNetwork(network *Network) error {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	logger := utils.GetLogger("AddNetwork")
	logger.Info().Msgf("add network: %v", network.Name)

	logger.Debug().Msgf("query network: %v", network.Name)
	_, err := QueryNetwork(network.Name)
	if err == gorm.ErrRecordNotFound {
		result := Conn.WithContext(ctx).Create(network)
		err = result.Error
		if err != nil {
			return err
		}
		logger.Info().Msgf("network: %v added", network.Name)
		//  switch network
		err = SwitchNetwork(network.Name)
		return err
	} else {
		return fmt.Errorf("network: %s already exists", network.Name)
	}

}

func QueryAllNetworks() (networks []Network, err error) {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	err = Conn.WithContext(ctx).Model(&Network{}).Find(&networks).Error

	return
}
func RemoveNetwork(name string) error {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	_, err := QueryNetwork(name)
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("network: %s not exist", name)
	} else {
		result := Conn.WithContext(ctx).Delete(&Network{}, "name = ?", name)
		err = result.Error
		if err != nil {
			return err
		}
		// switch network
		allNetworks, err := QueryAllNetworks()
		if err != nil {
			return err
		}

		if len(allNetworks) > 0 {
			err = SwitchNetwork(allNetworks[0].Name)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("no remaining network to switch to after remove\n")
		}

		return nil
	}
}

func CurrentNetwork() (network Network, err error) {
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	err = Conn.WithContext(ctx).Model(&Network{}).First(&network, "current = true").Error
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("no current network")
	}
	return
}

func QueryNetworkOrCurrent(name string) (*Network, error) {
	var (
		net    Network
		err    error
		logger = utils.GetLogger("QueryNetworkOrCurrent")
	)

	if name != "" {
		logger.Info().Msgf("Query network: %v", name)
		net, err = QueryNetwork(name)
		return &net, err
	}

	logger.Info().Msgf("Query current network")
	net, err = CurrentNetwork()
	if err != nil {
		return nil, err
	}

	return &net, nil
}
