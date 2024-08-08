package transaction

import (
	"encoding/hex"
	"fmt"
	"met/consts"
	"met/utils"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
)

// ParseInput 解析输入参数, data和abi二选一，用于构造交易input
func ParseInput(data string, abi string, method string, abiArgs ...string) ([]byte, error) {
	var (
		input []byte
		err   error
	)
	if data != "" {
		// 使用data
		// 删除0x
		input, err = hex.DecodeString(strings.TrimPrefix(data, "0x"))
		if err != nil {
			return nil, fmt.Errorf("decode data: %v error: %v", data, err)
		}

	} else {
		// 使用abi
		input, err = ParseAbi(abi, method, abiArgs...)
		if err != nil {
			return nil, fmt.Errorf("parse abi error: %v", err)
		}
	}

	return input, nil
}

// ParseAbi 解析erc20 abi, 用于构造交易input
// 1. symbol: 代币符号,如果为空，则通过rpc查询;如果为了快速构造input，可以传入
// 2. decimals: 代币精度,如果为空，则通过rpc查询;如果为了快速构造input，可以传入
// 3. method: erc20方法
// 4. abiArgs: erc20方法参数,如果是amount类的参数，其类型为人类可读的数字，需要转换为区块链底层数字
func ParseErc20Input(client *ethclient.Client, contractAddress, sender, symbol, decimals, method string, abiArgs ...string) ([]byte, error) {
	logger := utils.GetLogger("ParseErc20Input")

	switch method {
	case consts.Erc20Transfer:
		if len(abiArgs) != 2 {
			return nil, fmt.Errorf("erc20 transfer method need 2 args")
		}
		to := abiArgs[0]
		humanAmount := abiArgs[1]
		symbol, decimals, err := getErc20SymbolAndDecimals(client, contractAddress, symbol, decimals)
		if err != nil {
			return nil, fmt.Errorf("getErc20SymbolAndDecimals error: %v", err)
		}

		// 转换amount
		amount, err := utils.Erc20AmountFromHuman(humanAmount, decimals)
		if err != nil {
			return nil, fmt.Errorf("erc20 amount from human amount: %v error: %v", humanAmount, err)
		}

		input, err := ParseInput("", consts.Erc20, method, to, amount)
		if err != nil {
			return nil, fmt.Errorf("parse input error: %v", err)
		}

		transferInfo := fmt.Sprintf(`
Erc20 Transfer Info:
Contract Address:      %s
From:                  %s
To:                    %s
Symbol:                %s
Amount:                %s
`,
			contractAddress,
			sender,
			to,
			symbol,
			humanAmount)
		logger.Info().Msgf(transferInfo)

		return input, nil

	case consts.Erc20Approve:
		if len(abiArgs) != 2 {
			return nil, fmt.Errorf("erc20 approve method need 2 args")
		}
		spender := abiArgs[0]
		humanAmount := abiArgs[1]
		symbol, decimals, err := getErc20SymbolAndDecimals(client, contractAddress, symbol, decimals)
		if err != nil {
			return nil, fmt.Errorf("getErc20SymbolAndDecimals error: %v", err)
		}

		// 转换amount
		amount, err := utils.Erc20AmountFromHuman(humanAmount, decimals)
		if err != nil {
			return nil, fmt.Errorf("erc20 amount from human amount: %v error: %v", humanAmount, err)
		}

		input, err := ParseInput("", consts.Erc20, method, spender, amount)
		if err != nil {
			return nil, fmt.Errorf("parse input error: %v", err)
		}

		approveInfo := fmt.Sprintf(`
Erc20 Approve Info:
Contract Address:      %s
Spender:               %s
Symbol:                %s
Amount:                %s
`,
			contractAddress,
			spender,
			symbol,
			humanAmount)
		logger.Info().Msgf(approveInfo)

		return input, nil

	case consts.Erc20TransferFrom:
		if len(abiArgs) != 3 {
			return nil, fmt.Errorf("erc20 transferFrom method need 3 args")
		}
		owner := abiArgs[0]
		to := abiArgs[1]
		humanAmount := abiArgs[2]
		symbol, decimals, err := getErc20SymbolAndDecimals(client, contractAddress, symbol, decimals)
		if err != nil {
			return nil, fmt.Errorf("getErc20SymbolAndDecimals error: %v", err)
		}

		// 转换amount
		amount, err := utils.Erc20AmountFromHuman(humanAmount, decimals)
		if err != nil {
			return nil, fmt.Errorf("erc20 amount from human amount: %v error: %v", humanAmount, err)
		}

		input, err := ParseInput("", consts.Erc20, method, owner, to, amount)
		if err != nil {
			return nil, fmt.Errorf("parse input error: %v", err)
		}

		transferFromInfo := fmt.Sprintf(`
Erc20 transferFrom Info:
Contract Address:      %s
Spender:               %s
From:                  %s
To:                    %s
Symbol:                %s
Amount:                %s
`,
			contractAddress,
			sender,
			owner,
			to,
			symbol,
			humanAmount)
		logger.Info().Msgf(transferFromInfo)

		return input, nil

	default:
		return nil, fmt.Errorf("unsupported erc20 method: %s", method)
	}

	return nil, nil
}

// getErc20SymbolAndDecimals 获取erc20代币的symbol和decimals, 如果symbol或decimals为空，则通过rpc查询
func getErc20SymbolAndDecimals(client *ethclient.Client, contractAddress, symbol string, decimals string) (string, string, error) {
	logger := utils.GetLogger("getErc20SymbolAndDecimals")
	var (
		err error
	)

	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	if symbol == "" {
		logger.Debug().Msgf("symbol is empty, try to get from contract: %s", contractAddress)
		// 通过rpc查询
		symbol, err = ReadErc20(ctx, contractAddress, client, nil, Erc20Symbol, "", "")
		if err != nil {
			return "", "", fmt.Errorf("get symbol error: %w", err)
		}
	}

	if decimals == "" {
		logger.Debug().Msgf("decimals is empty, try to get from contract: %s", contractAddress)
		// 通过rpc查询
		decimals, err = ReadErc20(ctx, contractAddress, client, nil, Erc20Decimals, "", "")
		if err != nil {
			return "", "", fmt.Errorf("get decimals error: %w", err)
		}
	}

	return symbol, decimals, nil
}
