/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package erc20

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"my-ether-tool/cmd"
	"my-ether-tool/database"
	"my-ether-tool/types"
	"my-ether-tool/utils"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

// Erc20Cmd represents the tx command
var Erc20Cmd = &cobra.Command{
	Use:   "erc20",
	Short: "erc20 command",
	Long:  `erc20 command`,
	Run:   nil,
}

func init() {
	cmd.RootCmd.AddCommand(Erc20Cmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// Erc20Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// Erc20Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

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
func ReadErc20(ctx context.Context, contract string, networkName string, funcType Erc20ReadFuncType, arg1 string, arg2 string) (string, error) {
	net, err := database.QueryNetworkOrCurrent(networkName)
	if err != nil {
		return "", fmt.Errorf("query network error: %w", err)
	}

	client, err := ethclient.DialContext(ctx, net.Rpc)
	if err != nil {
		return "", fmt.Errorf("dial rpc error: %w", err)
	}

	contractAddress := common.HexToAddress(contract)
	erc20Instance, err := utils.NewErc20(contractAddress, client)

	switch funcType {
	case Erc20Name:
		tokenName, err := erc20Instance.Name(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", nil
		}

		return tokenName, nil

	case Erc20Symbol:
		symbol, err := erc20Instance.Symbol(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", err
		}

		return symbol, nil

	case Erc20Decimals:
		decimals, err := erc20Instance.Decimals(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", decimals), nil

	case Erc20TotalSupply:
		totalSupply, err := erc20Instance.TotalSupply(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s", totalSupply.String()), nil

	case Erc20BalanceOf:
		if arg1 == "" {
			return "", errors.New("missing address")
		}
		address := common.HexToAddress(arg1)
		balance, err := erc20Instance.BalanceOf(&bind.CallOpts{Context: ctx}, address)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s", balance.String()), nil

	case Erc20Allowance:
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
		return fmt.Sprintf("%s", allowance.String()), nil

	default:
		return "", errors.New("invalid erc20 read func type")

	}
}

type Erc20WritFuncType int8

const (
	Erc20Transfer Erc20WritFuncType = iota
	Erc20TransferFrom
	Erc20Approve
)

// 写erc20
func WriteErc20(ctx context.Context, contract string, networkName string, accountName string, accountIndex uint, funcType Erc20WritFuncType, arg1, arg2, arg3 string) (string, error) {

	net, err := database.QueryNetworkOrCurrent(networkName)
	if err != nil {
		return "", fmt.Errorf("query network error: %w", err)
	}

	client, err := ethclient.DialContext(ctx, net.Rpc)
	if err != nil {
		return "", fmt.Errorf("dial rpc error: %w", err)
	}

	account, err := database.QueryAccountOrCurrent(accountName, accountIndex)
	if err != nil {
		return "", fmt.Errorf("query account error: %w", err)
	}

	accountDetails, err := types.AccountToDetails(account)
	if err != nil {
		return "", fmt.Errorf("get account details error: %w", err)
	}

	pk := strings.TrimPrefix(accountDetails.PrivateKey, "0x")
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return "", fmt.Errorf("create private key error: %w", err)
	}

	chainId, err := client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("get chain id error: %w", err)
	}

	transactor, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return "", err
	}

	contractAddress := common.HexToAddress(contract)
	erc20Instance, err := utils.NewErc20(contractAddress, client)

	switch funcType {
	case Erc20Transfer:
		if arg1 == "" {
			return "", errors.New("missing to address")
		}
		if arg2 == "" {
			return "", errors.New("missing amount")
		}

		to := common.HexToAddress(arg1)
		amount := big.NewInt(0)
		amount.SetString(arg2, 10)

		tx, err := erc20Instance.Transfer(transactor, to, amount)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v/tx/%v", net.Explorer, tx.Hash().Hex()), nil

	case Erc20TransferFrom:
		if arg1 == "" {
			return "", errors.New("missing from address")
		}
		if arg2 == "" {
			return "", errors.New("missing to address")
		}
		if arg3 == "" {
			return "", errors.New("missing amount")
		}

		from := common.HexToAddress(arg1)
		to := common.HexToAddress(arg2)
		amount := big.NewInt(0)
		amount.SetString(arg3, 10)

		tx, err := erc20Instance.TransferFrom(transactor, from, to, amount)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v/tx/%v", net.Explorer, tx.Hash().Hex()), nil

	case Erc20Approve:
		if arg1 == "" {
			return "", errors.New("missing spender address")
		}
		if arg2 == "" {
			return "", errors.New("missing amount")
		}

		spender := common.HexToAddress(arg1)
		amount := big.NewInt(0)
		amount.SetString(arg2, 10)

		tx, err := erc20Instance.Approve(transactor, spender, amount)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v/tx/%v", net.Explorer, tx.Hash().Hex()), nil

	default:
		return "", errors.New("invalid erc20 write func type")

	}
}
