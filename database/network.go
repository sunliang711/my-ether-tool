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
	err = Conn.Model(&Network{}).First(&network, "name = ?", name).Error

	return
}

func SwitchNetwork(name string) error {
	logger := utils.GetLogger("SwitchNetwork")
	logger.Info().Msgf("switch to network: %v", name)

	// 清空老的current
	err := Conn.Model(&Network{}).Where("current = true").Update("current", false).Error
	if err != nil {
		return err
	}
	// 设置新的current
	err = Conn.Model(&Network{}).Where("name = ?", name).Update("current", true).Error
	return err
}

func AddNetwork(network *Network) error {
	logger := utils.GetLogger("AddNetwork")
	logger.Info().Msgf("add network: %v", network.Name)

	logger.Debug().Msgf("query network: %v", network.Name)

	_, err := QueryNetwork(network.Name)
	if err == gorm.ErrRecordNotFound {
		result := Conn.Create(network)
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
	err = Conn.Model(&Network{}).Find(&networks).Error

	return
}
func RemoveNetwork(name string) error {
	_, err := QueryNetwork(name)
	if err == gorm.ErrRecordNotFound {
		return fmt.Errorf("network: %s not exist", name)
	} else {
		result := Conn.Delete(&Network{}, "name = ?", name)
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
	err = Conn.Model(&Network{}).First(&network, "current = true").Error
	if err == gorm.ErrRecordNotFound {
		err = fmt.Errorf("no current network")
	}
	return
}

func QueryNetworkOrCurrent(name string) (*Network, error) {
	var (
		net Network
		err error
	)
	if name != "" {
		net, err = QueryNetwork(name)
		return &net, err
	}
	net, err = CurrentNetwork()

	return &net, err
}
