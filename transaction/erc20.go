package transaction

import (
	"context"
	"errors"
	"fmt"
	"met/database"
	utils "met/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Erc20ReadFuncType int8

const (
	Erc20Name Erc20ReadFuncType = iota
	Erc20Symbol
	Erc20Decimals
	Erc20TotalSupply
	Erc20BalanceOf
	Erc20Allowance
)

// 读erc20
func ReadErc20(ctx context.Context, contract string, client *ethclient.Client, net *database.Network, funcType Erc20ReadFuncType, arg1 string, arg2 string) (string, error) {
	logger := utils.GetLogger("ReadErc20")

	if net != nil {
		logger.Debug().Msgf("network info: %v", net)
	}

	contractAddress := common.HexToAddress(contract)
	erc20Instance, err := utils.NewErc20(contractAddress, client)
	if err != nil {
		return "", fmt.Errorf("NewErc20 error: %w", err)
	}

	switch funcType {
	case Erc20Name:
		logger.Info().Msg("call erc20 name")
		tokenName, err := erc20Instance.Name(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", nil
		}

		return tokenName, nil

	case Erc20Symbol:
		logger.Info().Msg("call erc20 symbol")
		symbol, err := erc20Instance.Symbol(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", err
		}

		return symbol, nil

	case Erc20Decimals:
		logger.Info().Msg("call erc20 decimals")
		decimals, err := erc20Instance.Decimals(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", decimals), nil

	case Erc20TotalSupply:
		logger.Info().Msg("call erc20 totalSupply")
		totalSupply, err := erc20Instance.TotalSupply(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", err
		}

		return totalSupply.String(), nil

	case Erc20BalanceOf:
		logger.Info().Msg("call erc20 balanceOf")
		if arg1 == "" {
			return "", errors.New("missing address")
		}
		address := common.HexToAddress(arg1)
		balance, err := erc20Instance.BalanceOf(&bind.CallOpts{Context: ctx}, address)
		if err != nil {
			return "", err
		}

		return balance.String(), nil

	case Erc20Allowance:
		logger.Info().Msg("call erc20 allowance")
		if arg1 == "" {
			return "", errors.New("missing owner address")
		}
		if arg2 == "" {
			return "", errors.New("missing spender address")
		}
		owner := common.HexToAddress(arg1)
		spender := common.HexToAddress(arg2)

		allowance, err := erc20Instance.Allowance(&bind.CallOpts{Context: ctx}, owner, spender)
		if err != nil {
			return "", err
		}
		return allowance.String(), nil

	default:
		return "", errors.New("invalid erc20 read func type")

	}
}
