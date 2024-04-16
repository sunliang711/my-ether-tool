/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package erc20

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	cmd "met/cmd"
	database "met/database"
	types "met/types"
	utils "met/utils"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
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
	log := log.With().Str("func", "ReadErc20").Logger()

	log.Info().Msgf("query network: %v", networkName)
	net, err := database.QueryNetworkOrCurrent(networkName)
	if err != nil {
		return "", fmt.Errorf("query network error: %w", err)
	}
	log.Debug().Msgf("network info: %v", net)

	log.Info().Msgf("dial rpc: %v", net.Rpc)
	client, err := ethclient.DialContext(ctx, net.Rpc)
	if err != nil {
		return "", fmt.Errorf("dial rpc error: %w", err)
	}

	contractAddress := common.HexToAddress(contract)
	erc20Instance, err := utils.NewErc20(contractAddress, client)

	switch funcType {
	case Erc20Name:
		log.Info().Msg("call erc20 name")
		tokenName, err := erc20Instance.Name(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", nil
		}

		return tokenName, nil

	case Erc20Symbol:
		log.Info().Msg("call erc20 symbol")
		symbol, err := erc20Instance.Symbol(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", err
		}

		return symbol, nil

	case Erc20Decimals:
		log.Info().Msg("call erc20 decimals")
		decimals, err := erc20Instance.Decimals(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v", decimals), nil

	case Erc20TotalSupply:
		log.Info().Msg("call erc20 totalSupply")
		totalSupply, err := erc20Instance.TotalSupply(&bind.CallOpts{Context: ctx})
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s", totalSupply.String()), nil

	case Erc20BalanceOf:
		log.Info().Msg("call erc20 balanceOf")
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
		log.Info().Msg("call erc20 allowance")
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
func WriteErc20(ctx context.Context, contract string, noconfirm bool, networkName string, accountName string, accountIndex uint, funcType Erc20WritFuncType, arg1, arg2, arg3 string) (string, error) {
	logger := utils.GetLogger("WriteErc20")

	logger.Info().Msgf("query network: %v", networkName)
	net, err := database.QueryNetworkOrCurrent(networkName)
	if err != nil {
		return "", fmt.Errorf("query network error: %w", err)
	}
	logger.Debug().Msgf("network info: %v", net)

	logger.Info().Msgf("dial rpc: %v", net.Rpc)
	client, err := ethclient.DialContext(ctx, net.Rpc)
	if err != nil {
		return "", fmt.Errorf("dial rpc error: %w", err)
	}

	logger.Info().Msgf("query account: %v with index: %v", accountName, accountIndex)
	account, err := database.QueryAccountOrCurrent(accountName, accountIndex)
	if err != nil {
		return "", fmt.Errorf("query account error: %w", err)
	}

	accountDetails, err := types.AccountToDetails(account)
	if err != nil {
		return "", fmt.Errorf("get account details error: %w", err)
	}
	addressStr, err := accountDetails.Address()
	if err != nil {
		return "", fmt.Errorf("get account address error: %w", err)
	}
	logger.Info().Msgf("account info: name: %v address: %v", accountDetails.Name, addressStr)

	privateKeyStr, err := accountDetails.PrivateKey()
	if err != nil {
		return "", fmt.Errorf("get account private key error: %w", err)
	}

	pk := strings.TrimPrefix(privateKeyStr, "0x")
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		return "", fmt.Errorf("create private key error: %w", err)
	}

	logger.Info().Msg("query chain id")
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
	if err != nil {
		return "", fmt.Errorf("NewErc20 error: %w", err)
	}

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
		if _, ok := amount.SetString(arg2, 10); !ok {
			return "", fmt.Errorf("invalid amount: %v", arg2)
		}

		logger.Info().Msgf("Call erc20 transfer")
		logger.Info().Msgf("From: %v", addressStr)
		logger.Info().Msgf("To: %v", to)
		logger.Info().Msgf("Amount: %v", amount)
		logger.Info().Msgf("Amount readable: %v", arg2)

		if !noconfirm {
			input, err := utils.ReadChar("Send ? [y/N] ")
			utils.ExitWhenErr(logger, err, "read input error: %s", err)

			if input != 'y' {
				os.Exit(0)
			}

		}
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
		if _, ok := amount.SetString(arg3, 10); !ok {
			return "", fmt.Errorf("invalid amount")
		}

		logger.Info().Msgf("Call erc20 transferFrom")
		logger.Info().Msgf("From: %v", from)
		logger.Info().Msgf("To: %v", to)
		logger.Info().Msgf("Amount: %v", amount)
		logger.Info().Msgf("Amount readable: %v", arg3)

		if !noconfirm {
			input, err := utils.ReadChar("Send ? [y/N] ")
			utils.ExitWhenErr(logger, err, "read input error: %s", err)

			if input != 'y' {
				os.Exit(0)
			}

		}

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
		if _, ok := amount.SetString(arg2, 10); !ok {
			return "", fmt.Errorf("invalid amount: %v", arg2)
		}

		logger.Info().Msgf("Call erc20 approve")
		logger.Info().Msgf("Owner: %v", addressStr)
		logger.Info().Msgf("Spender: %v", spender)
		logger.Info().Msgf("Amount: %v", amount)
		logger.Info().Msgf("Amount readable: %v", arg2)

		if !noconfirm {
			input, err := utils.ReadChar("Send ? [y/N] ")
			utils.ExitWhenErr(logger, err, "read input error: %s", err)

			if input != 'y' {
				os.Exit(0)
			}

		}

		tx, err := erc20Instance.Approve(transactor, spender, amount)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%v/tx/%v", net.Explorer, tx.Hash().Hex()), nil

	default:
		return "", errors.New("invalid erc20 write func type")

	}
}
