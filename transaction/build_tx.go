package transaction

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"met/utils"
	"strconv"
	"strings"

	mTypes "met/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

func BuildTx(client *ethclient.Client, from string, to string, value *string, data []byte, ledger bool, gasMode mTypes.GasMode, nonce, chainId, gasLimit, gasLimitRatio, gasRatio, gasPrice, gasTipCap, gasFeeCap string, sendAll bool) (tx *types.Transaction, err error) {
	var (
		nonce0     uint64
		gasLimit0  uint64
		chainId0   *big.Int
		value0     *big.Int
		gasFeeCap0 *big.Int
		gasTipCap0 *big.Int
		gasPrice0  *big.Int
		ok         bool
		logger     = utils.GetLogger("BuildTx")
	)
	ctx, cancel := utils.DefaultTimeoutContext()
	defer cancel()

	// 1. check flags
	if gasLimit != "" && gasLimitRatio != "" {
		return nil, errors.New("gasLimit conflicts with gasLimitRatio")
	}

	switch gasMode {
	case mTypes.GasModeAuto:
		if gasTipCap != "" {
			return nil, errors.New("no need for gasTipCap in auto gas mode")
		}
		if gasFeeCap != "" {
			return nil, errors.New("no need for gasFeeCap in auto gas mode")
		}
		if gasPrice != "" {
			return nil, errors.New("no need for gasPrice in auto gas mode")
		}
	case mTypes.GasModeLegacy:
		if gasTipCap != "" {
			return nil, errors.New("no need for gasTipCap in legacy gas mode")
		}
		if gasFeeCap != "" {
			return nil, errors.New("no need for gasFeeCap in legacy gas mode")
		}
	case mTypes.GasModeEip1559:
		if gasPrice != "" {
			return nil, errors.New("no need for gasPrice in eip1559 gas mode")
		}
	default:
		return nil, fmt.Errorf("invalid gasMode: %v", gasMode)
	}

	if sendAll {
		if gasMode != mTypes.GasModeLegacy {
			return nil, fmt.Errorf("sendAll only works in legacy gas mode")
		}
		if len(data) > 0 {
			return nil, fmt.Errorf("sendAll do not support input data")
		}
	}

	// 2. init vars
	from = strings.TrimPrefix(from, "0x")
	to = strings.TrimPrefix(to, "0x")
	fromAddress := common.HexToAddress(from)

	// 如果to为空，则是合约部署交易，此时toAddress为nil
	var toAddress *common.Address
	if to != "" {
		tt := common.HexToAddress(to)
		toAddress = &tt
	}

	// value
	if *value != "" {
		logger.Debug().Msgf("parse value: %v", *value)
		value0, err = utils.ParseUnits(*value, utils.UnitEth)
		if err != nil {
			return nil, fmt.Errorf("parse value: %v error: %w", *value, err)
		}
	}
	logger.Debug().Msgf("value: %v", value0.String())

	// nonce
	if nonce != "" {
		logger.Debug().Msgf("parse nonce: %v", nonce)
		nonce0, err = strconv.ParseUint(nonce, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse nonce: %v error: %w", nonce, err)
		}
	} else {
		logger.Debug().Msgf("query nonce..")
		nonce0, err = client.PendingNonceAt(ctx, fromAddress)
		if err != nil {
			return nil, fmt.Errorf("query nonce error: %w", err)
		}
	}
	logger.Debug().Msgf("nonce: %v", nonce0)

	// chainId
	if chainId != "" {
		logger.Debug().Msgf("parse chainId: %v", chainId)
		if chainId0, ok = big.NewInt(0).SetString(chainId, 10); !ok {
			return nil, fmt.Errorf("invalid chainId: %v", chainId)
		}

	} else {
		logger.Debug().Msgf("query chainId..")
		chainId0, err = client.ChainID(ctx)
		if err != nil {
			return nil, fmt.Errorf("query chainId error: %w", err)
		}
	}
	logger.Debug().Msgf("chainId: %v", chainId0.String())

	// gasLimit
	if gasLimit != "" {
		logger.Debug().Msgf("parse gasLimit: %v", gasLimit)
		gasLimit0, err = strconv.ParseUint(gasLimit, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse gasLimit: %v error: %w", gasLimit, err)
		}
	} else {
		gasLimit0, err = client.EstimateGas(ctx, ethereum.CallMsg{
			From:  fromAddress,
			To:    toAddress,
			Value: value0,
			Data:  data,
		})
		if err != nil {
			return nil, fmt.Errorf("estimate gas error: %w", err)
		}
	}
	logger.Debug().Msgf("gasLimit: %v", gasLimit0)

	// gasLimitRatio
	if gasLimitRatio != "" {
		logger.Debug().Msgf("parse gasLimitRatio: %v", gasLimitRatio)
		gasLimitRatio0, err := decimal.NewFromString(gasLimitRatio)
		if err != nil {
			return nil, fmt.Errorf("parse gasLimitRatio: %v error: %w", gasLimitRatio, err)
		}
		logger.Debug().Msgf("scale gasLimit by ratio: %v", gasLimitRatio)
		newGasLimit := decimal.NewFromInt(int64(gasLimit0)).Mul(gasLimitRatio0).IntPart()
		gasLimit0 = uint64(newGasLimit)
		logger.Debug().Msgf("after scale, gasLimit: %v", gasLimit0)
	}

	// ledger 只支持 legacy gas mode
	if ledger {
		gasMode = mTypes.GasModeLegacy
	}

	// gas
	if gasMode == mTypes.GasModeAuto {
		logger.Info().Msgf("Gas mode: auto")
		logger.Debug().Msgf("query latest block header..")
		header, err := client.HeaderByNumber(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("get latest block header error: %w", err)
		}

		baseFee := header.BaseFee
		if baseFee != nil {
			// support 1559
			logger.Debug().Msgf("support 1559")
			logger.Debug().Msgf("query gasTipCap")
			gasTipCap0, err = client.SuggestGasTipCap(ctx)
			if err != nil {
				return nil, fmt.Errorf("query gasTipCap error: %w", err)
			}

			gasFeeCap0 = new(big.Int).Add(gasTipCap0, new(big.Int).Mul(baseFee, big.NewInt(2)))

			logger.Debug().Msgf("gasFeeCap: %v", gasFeeCap0.String())
			logger.Debug().Msgf("gasTipCap: %v", gasTipCap0.String())

			gasMode = mTypes.GasModeEip1559

		} else {
			// not support 1559
			logger.Debug().Msgf("not support 1559")
			logger.Debug().Msgf("query gasPrice")
			gasPrice0, err = client.SuggestGasPrice(ctx)
			if err != nil {
				return nil, fmt.Errorf("query gasPrice error: %w", err)
			}

			logger.Debug().Msgf("gasPrice: %v", gasPrice0)

			gasMode = mTypes.GasModeLegacy
		}

	} else if gasMode == mTypes.GasModeLegacy {
		logger.Info().Msgf("Gas mode: legacy")
		if gasPrice != "" {
			logger.Debug().Msgf("parse gasPrice: %v", gasPrice)
			gasPrice0, err = utils.ParseUnits(gasPrice, utils.UnitGwei)
			if err != nil {
				return nil, fmt.Errorf("parse gasPrice: %v error: %w", gasPrice, err)
			}
		} else {
			gasPrice0, err = client.SuggestGasPrice(ctx)
			if err != nil {
				return nil, fmt.Errorf("query gasPrice error: %w", err)
			}
		}

		logger.Debug().Msgf("gasPrice: %v", gasPrice0)

		if sendAll {
			logger.Info().Msgf("sendAll mode")

			logger.Debug().Msgf("query balance for address: %v", from)
			currentBalance, err := client.BalanceAt(ctx, fromAddress, nil)
			if err != nil {
				return nil, fmt.Errorf("query balance for address: %v error: %w", from, err)
			}

			txFee := new(big.Int).Mul(big.NewInt(OnlyTransferGas), gasPrice0)
			logger.Info().Msgf("sendAll tx fee: %v", txFee.String())

			value0 = new(big.Int).Sub(currentBalance, txFee)
			logger.Info().Msgf("sendAll value: %v", value0.String())

			vv, err := utils.FormatUnits(value0.String(), utils.UnitEth)
			if err != nil {
				return nil, fmt.Errorf("format value: %v error: %w", value0.String(), err)
			}
			*value = vv
		}

	} else if gasMode == mTypes.GasModeEip1559 {
		logger.Info().Msgf("Gas mode: eip1559")

		if gasTipCap != "" {
			logger.Debug().Msgf("parse gasTipCap: %v", gasTipCap)
			gasTipCap0, err = utils.ParseUnits(gasTipCap, utils.UnitGwei)
			if err != nil {
				return nil, fmt.Errorf("parse gasTipCap: %v error: %w", gasTipCap, err)
			}

		} else {
			logger.Debug().Msgf("query gasTipCap..")
			gasTipCap0, err = client.SuggestGasTipCap(ctx)
			if err != nil {
				return nil, fmt.Errorf("query gasTipCap error: %w", err)
			}

		}
		logger.Debug().Msgf("gasTipCap: %v", gasTipCap0)

		if gasFeeCap != "" {
			logger.Debug().Msgf("parse gasFeeCap: %v", gasFeeCap)
			gasFeeCap0, err = utils.ParseUnits(gasFeeCap, utils.UnitGwei)
			if err != nil {
				return nil, fmt.Errorf("parse gasFeeCap: %v error: %w", gasFeeCap, err)
			}
		} else {
			logger.Debug().Msgf("calculate gasFeeCap")
			logger.Debug().Msgf("query latest block header..")
			header, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("get latest block header error: %w", err)
			}

			baseFee := header.BaseFee
			if baseFee == nil {
				return nil, fmt.Errorf("cannot get baseFee in eip1559 gas mode")
			}
			logger.Debug().Msgf("query gasTipCap..")

			gasFeeCap0 = new(big.Int).Add(gasTipCap0, new(big.Int).Mul(baseFee, big.NewInt(2)))

		}
		logger.Debug().Msgf("gasFeeCap: %v", gasFeeCap0.String())

	}

	if gasRatio != "" {
		logger.Debug().Msgf("parse gasRatio: %v", gasRatio)
		gasRatio0, err := decimal.NewFromString(gasRatio)
		if err != nil {
			return nil, fmt.Errorf("parse gasRatio: %v error: %w", gasRatio, err)
		}

		if gasPrice0 != nil {
			gasPrice0 = decimal.NewFromBigInt(gasPrice0, 0).Mul(gasRatio0).BigInt()
		}
		if gasFeeCap0 != nil {
			gasFeeCap0 = decimal.NewFromBigInt(gasFeeCap0, 0).Mul(gasRatio0).BigInt()
		}
		if gasTipCap0 != nil {
			gasTipCap0 = decimal.NewFromBigInt(gasTipCap0, 0).Mul(gasRatio0).BigInt()
		}

		logger.Debug().Msgf("after gasRatio, gasPrice: %v", gasPrice0.String())
		logger.Debug().Msgf("after gasRatio, gasFeeCap: %v", gasFeeCap0.String())
		logger.Debug().Msgf("after gasRatio, gasTipCap: %v", gasTipCap0.String())
	}

	// ledger 只支持 legacy tx
	if ledger {
		logger.Info().Msgf("Transaction type: legacy")
		tx = types.NewTx(&types.LegacyTx{
			Nonce:    nonce0,
			GasPrice: gasPrice0,
			Gas:      gasLimit0,
			To:       toAddress,
			Value:    value0,
			Data:     data,
		})
	} else {
		if gasMode == mTypes.GasModeLegacy {
			logger.Info().Msgf("Transaction type: accessList")
			tx = types.NewTx(&types.AccessListTx{
				ChainID:    chainId0,
				Nonce:      nonce0,
				GasPrice:   gasPrice0,
				Gas:        gasLimit0,
				To:         toAddress,
				Value:      value0,
				Data:       data,
				AccessList: []types.AccessTuple{},
			})
		} else if gasMode == mTypes.GasModeEip1559 {
			logger.Info().Msgf("Transaction type: dynamicFee (eip1559)")
			tx = types.NewTx(&types.DynamicFeeTx{
				ChainID:   chainId0,
				Nonce:     nonce0,
				GasTipCap: gasTipCap0,
				GasFeeCap: gasFeeCap0,
				Gas:       gasLimit0,
				To:        toAddress,
				Value:     value0,
				Data:      data,
			})
		}
	}

	return
}

// @param from  from address of tx
// @param to    to address of tx
// @param value amount of eth to be transfered, unit: ETH
// @param data  optional   conflict with abi
// @param abi  optional conflict with data
// @param args  optional
// @param gasLimit optional 为0时从rpc获取
// @param nonce optional 为空时从rpc获取
// @param chainID optional 为0时从rpc获取
// @param gasRatio
// @param gasPrice used when eip1559 is false unit: gwei
// @param gasTipCap used when eip1559 is true unit: gwei
// @param gasFeeCap used when eip1559 is true unit: gwei
// @param eip1559:
// eip1559为true时，当gasTipCap 和 gasFeeCap都不为空时使用它们，否则从rpc获取这两个值
// eip1559为false时，当gasPrice不为空时使用gasPrice，否则从rpc获取
func BuildTransaction(ctx context.Context, client *ethclient.Client, from string, to string, value *string, data string, abi string, args []string, gasLimit uint64, nonce string, chainID int, gasRatio, gasPrice string, gasTipCap string, gasFeeCap string, eip1559 bool, sendAll bool) (tx *types.Transaction, err error) {
	// check params
	var (
		nonce0   uint64
		chainID0 *big.Int = big.NewInt(0)
		data0    []byte
	)
	from = strings.TrimPrefix(from, "0x")
	to = strings.TrimPrefix(to, "0x")
	data = strings.TrimPrefix(data, "0x")

	logger := utils.GetLogger("BuildTransaction")
	fromAddress := common.HexToAddress(from)
	toAddress := common.HexToAddress(to)

	value1, err := utils.ParseUnits(*value, utils.UnitEth)
	if err != nil {
		return
	}

	if data != "" && abi != "" {
		err = errors.New("data conflict with abi,specify one")
		return
	}

	if sendAll {
		eip1559 = false
		logger.Info().Msgf("sendAll conflict with eip1559, disable eip1559")
	}

	if data != "" {
		data0, err = hex.DecodeString(data)
		if err != nil {
			return
		}
	}

	if abi != "" {
		data0, err = AbiEncode(abi, args)
		if err != nil {
			return
		}
	}

	// if nonce == ""; get by rpc
	if nonce == "" {
		nonce0, err = client.PendingNonceAt(ctx, fromAddress)
		if err != nil {
			return nil, fmt.Errorf("query nonce error: %w", err)
		}
	} else {
		nonce0, err = strconv.ParseUint(nonce, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse nonce error: %w", err)
		}
	}

	// if chainID == 0; get by rpc
	if chainID == 0 {
		chainID0, err = client.ChainID(ctx)
		if err != nil {
			return nil, fmt.Errorf("get chain id error: %w", err)
		}
	} else {
		chainID0 = big.NewInt(int64(chainID))
	}

	// if gasLimit == 0; estimateGas
	if gasLimit == 0 {
		gasLimit, err = client.EstimateGas(ctx, ethereum.CallMsg{
			From:  fromAddress,
			To:    &toAddress,
			Value: value1,
			Data:  data0,
		})
		if err != nil {
			return nil, fmt.Errorf("estimate gas error: %w", err)
		}
	}

	// signer := types.NewCancunSigner(chainID0)

	// gWei := big.NewInt(1_000_000_000)

	if eip1559 {
		// use gasTipCap gasFeeCap
		var tipCap *big.Int
		var feeCap *big.Int

		if gasTipCap != "" && gasFeeCap != "" {
			// baseFee=eth_gasPrice - maxPriorityFeePerGas
			// 最大费用，maxPriorityFeePerGas + 2 * base_fee

			tipCap, err = utils.ParseUnits(gasTipCap, utils.UnitGwei)
			if err != nil {
				return nil, fmt.Errorf("parse gasTipCap: %v error: %w", gasTipCap, err)
			}

			feeCap, err = utils.ParseUnits(gasFeeCap, utils.UnitGwei)
			if err != nil {
				return nil, fmt.Errorf("parse gasFeeCap: %v error: %w", gasFeeCap, err)
			}

		} else {

			tipCap, err = client.SuggestGasTipCap(ctx)
			if err != nil {
				return
			}

			header, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("get latest block header error: %w", err)
			}

			baseFee := header.BaseFee

			// baseFee := new(big.Int).Sub(gasPrice0, tipCap)
			feeCap = new(big.Int).Add(tipCap, new(big.Int).Mul(baseFee, big.NewInt(2)))
		}

		if gasRatio != "" {
			gRatio, err := decimal.NewFromString(gasRatio)
			if err != nil {
				return nil, fmt.Errorf("parse gasRatio: %v error: %w", gasRatio, err)
			}

			tipCap = decimal.NewFromBigInt(tipCap, 0).Mul(gRatio).BigInt()
			feeCap = decimal.NewFromBigInt(feeCap, 0).Mul(gRatio).BigInt()
			logger.Debug().Msgf("after gasRatio: %v, tipCap: %v", gasRatio, tipCap.String())
			logger.Debug().Msgf("after gasRatio: %v, feeCap: %v", gasRatio, feeCap.String())

		}
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   chainID0,
			Nonce:     nonce0,
			GasTipCap: tipCap,
			GasFeeCap: feeCap,
			Gas:       gasLimit,
			To:        &toAddress,
			Value:     value1,
			Data:      data0,
		})

	} else {
		// use gasPrice
		var price *big.Int
		if gasPrice == "" {
			price, err = client.SuggestGasPrice(ctx)
			if err != nil {
				return
			}
		} else {
			price, err = utils.ParseUnits(gasPrice, utils.UnitGwei)
			if err != nil {
				return nil, fmt.Errorf("parse gasPrice: %v error: %w", gasPrice, err)
			}

		}

		if gasRatio != "" {
			gRatio, err := decimal.NewFromString(gasRatio)
			if err != nil {
				return nil, fmt.Errorf("parse gasRatio: %v error: %w", gasRatio, err)
			}

			price = decimal.NewFromBigInt(price, 0).Mul(gRatio).BigInt()
			logger.Debug().Msgf("after gasRatio: %v, gasPrice: %v", gasRatio, price.String())
		}

		// 重新计算value
		if sendAll {
			logger.Info().Msgf("sendAll mode")
			// 查询当前balance
			currentBalance, err := client.BalanceAt(ctx, common.HexToAddress(from), nil)
			if err != nil {
				return nil, fmt.Errorf("query balance error: %w", err)
			}

			// 计算手续费 = 21000 * gasPrice
			txFee := big.NewInt(OnlyTransferGas)
			txFee.Mul(txFee, price)
			logger.Info().Msgf("sendAll tx fee: %v", txFee.String())

			// 剩下的value为所有待发送value
			value1 = currentBalance.Sub(currentBalance, txFee)
			logger.Info().Msgf("sendAll value: %v", value1.String())

			// newValue = value1
			vv, err := utils.FormatUnits(value1.String(), utils.UnitEth)
			if err != nil {
				return nil, fmt.Errorf("format value error: %w", err)
			}
			*value = vv
		}

		tx = types.NewTx(&types.AccessListTx{
			ChainID:  chainID0,
			Nonce:    nonce0,
			GasPrice: price,
			Gas:      gasLimit,
			To:       &toAddress,
			Value:    value1,
			Data:     data0,
		})

	}

	return
}

const (
	OnlyTransferGas = 21000
)

type TxParams struct {
	ChainId *big.Int
	Nonce   *big.Int

	GasLimit  uint64
	GasPrice  *big.Int
	GasFeeCap *big.Int
	GasTipCap *big.Int

	Value *big.Int
}

// chainId: 为空时从rpc读取
// nonce：为空时从rpc读取
// value: 单位是ETH，不为空时转换成wei，

// gasLimitRatio: 不为空时，把estimateGas的结果乘以这个比例
// gasLimit：不为空时，使用这个作为gasLimit
// gasLimitRatio和gasLimit 不能同时使用

// eip1559: 是否使用eip1559。为true时使用gasFeeCap gasTipCap，为false时使用gasPrice，前提是它们也都不为空，如果它们为空，则从rpc读取，并根据gasRatio是否为空进行缩放
// gasRatio: 不为空时，使用它把gasPrice(或gasTipCap gasFeeCap)乘以这个比例
// gasPrice: 不为空时，使用它作为gasPrice
// gasFeeCap： 不为空时使用
// gasTipCap： 不为空时使用
// gasRatio 和gasPrice(或 gasFeeCap gasTip Cap) 不能同时使用
func GetTxParams(ctx context.Context, client *ethclient.Client, fromAddress, contractAddress, chainId, nonce, value, gasLimitRatio, gasLimit, gasRatio, gasPrice, gasFeeCap, gasTipCap string, eip1559 bool, input []byte) (*TxParams, error) {
	logger := utils.GetLogger("GetTxParams")

	if gasLimitRatio != "" && gasLimit != "" {
		return nil, fmt.Errorf("gasLimitRatio conflicts with gasLimit")
	}

	if eip1559 {
		if gasRatio != "" && (gasFeeCap != "" || gasTipCap != "") {
			return nil, fmt.Errorf("gasRatio conflicts with gasFeeCap or gasTipCap")
		}

	} else {
		if gasRatio != "" && gasPrice != "" {
			return nil, fmt.Errorf("gasRatio conflicts with gasPrice")
		}

	}

	txParam := new(TxParams)

	from := common.HexToAddress(fromAddress)
	contract := common.HexToAddress(contractAddress)

	// Value
	if value != "" {
		logger.Debug().Msgf("parse value: %v", value)
		v, err := utils.ParseUnits(value, utils.UnitEth)
		if err != nil {
			return nil, fmt.Errorf("parse value:%v error: %w", value, err)
		}
		logger.Debug().Msgf("value: %v", v.String())
		txParam.Value = v
	}

	// ChainId
	if chainId == "" {
		logger.Debug().Msgf("query chain id")
		id, err := client.ChainID(ctx)
		if err != nil {
			return nil, fmt.Errorf("get chain id error: %w", err)
		}
		logger.Debug().Msgf("chain id: %v", id.String())
		txParam.ChainId = id
	} else {
		logger.Debug().Msgf("parse chain id: %v", chainId)
		if id, ok := big.NewInt(0).SetString(chainId, 10); ok {
			txParam.ChainId = id
		} else {
			return nil, fmt.Errorf("invalid chain id: %v", chainId)
		}
	}

	// Nonce
	if nonce == "" {
		from := common.HexToAddress(fromAddress)
		logger.Debug().Msgf("query nonce for address: %v", from.Hex())
		n, err := client.PendingNonceAt(ctx, from)
		if err != nil {
			return nil, fmt.Errorf("query nonce of address: %v error: %w", fromAddress, err)
		}
		logger.Debug().Msgf("nonce: %v", n)
		txParam.Nonce = big.NewInt(int64(n))
	} else {
		logger.Debug().Msgf("parse nonce: %v", nonce)
		n, err := strconv.ParseUint(nonce, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse nonce: %v error: %w", nonce, err)
		}
		logger.Debug().Msgf("nonce: %v", n)
		txParam.Nonce = big.NewInt(int64(n))
	}

	if gasLimit != "" {
		logger.Debug().Msgf("parse gasLimit: %v", gasLimit)
		l, err := strconv.ParseUint(gasLimit, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse gasLimit: %v error: %w", gasLimit, err)
		}

		txParam.GasLimit = l
	} else {
		logger.Debug().Msgf("estimate gas")
		msg := ethereum.CallMsg{
			From: from,
			To:   &contract,
			// GasPrice: gasPrice,
			Value: txParam.Value,
			Data:  input,
		}
		estimatedGas, err := client.EstimateGas(ctx, msg)
		if err != nil {
			return nil, fmt.Errorf("estimate gas error: %v", err)
		}

		txParam.GasLimit = estimatedGas
	}

	if gasLimitRatio != "" {
		logger.Debug().Msgf("parse gasLimitRatio")
		limitRatio, err := decimal.NewFromString(gasLimitRatio)
		if err != nil {
			return nil, fmt.Errorf("parse gasLimitRatio error: %w", err)
		}
		logger.Debug().Msgf("gasLimitRatio: %v", limitRatio.String())

		logger.Debug().Msgf("before gas limit ratio: %v, gas limit: %v", gasLimitRatio, txParam.GasLimit)
		gasLimit := decimal.NewFromInt(int64(txParam.GasLimit)).Mul(limitRatio)
		txParam.GasLimit = gasLimit.BigInt().Uint64()
		logger.Debug().Msgf("after gas limit ratio: %v, gas limit: %v", gasLimitRatio, txParam.GasLimit)
	}

	// Gas
	if eip1559 {
		logger.Debug().Msgf("eip1559 on, ignore gasPrice")

		if gasTipCap != "" {
			logger.Debug().Msgf("parse gasTipCap")

			tip, err := utils.ParseUnits(gasTipCap, utils.UnitGwei)
			if err != nil {
				return nil, fmt.Errorf("parse gasTipCap: %v error: %w", gasTipCap, err)
			}
			txParam.GasTipCap = tip

		} else {
			logger.Debug().Msgf("query gasTipCap")
			tipCap, err := client.SuggestGasTipCap(ctx)
			if err != nil {
				return nil, fmt.Errorf("query gas tip cap error: %w", err)
			}
			txParam.GasTipCap = tipCap
		}

		if gasFeeCap != "" {
			logger.Debug().Msgf("parse gasFeeCap")

			fee, err := utils.ParseUnits(gasFeeCap, utils.UnitGwei)
			if err != nil {
				return nil, fmt.Errorf("parse gasFeeCap: %v error: %w", gasFeeCap, err)
			}
			txParam.GasFeeCap = fee
		} else {
			logger.Debug().Msgf("query gasFeeCap")
			tipCap, err := client.SuggestGasTipCap(ctx)
			if err != nil {
				return nil, fmt.Errorf("query gas tip cap error: %w", err)
			}

			header, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				return nil, fmt.Errorf("get latest block header error: %w", err)
			}

			baseFee := header.BaseFee
			// feeCap = tipCap + 2 * baseFee
			feeCap := new(big.Int).Add(tipCap, new(big.Int).Mul(baseFee, big.NewInt(2)))

			txParam.GasFeeCap = feeCap
		}

		if gasRatio != "" {
			logger.Debug().Msgf("parse gasRatio")
			gRatio, err := decimal.NewFromString(gasRatio)
			if err != nil {
				return nil, fmt.Errorf("parse gasRatio: %v error: %w", gasRatio, err)
			}

			logger.Debug().Msgf("before gas ratio, tip cap: %v fee cap: %v", txParam.GasTipCap.String(), txParam.GasFeeCap.String())

			tip := decimal.NewFromBigInt(txParam.GasTipCap, 0).Mul(gRatio)
			fee := decimal.NewFromBigInt(txParam.GasFeeCap, 0).Mul(gRatio)

			logger.Debug().Msgf("after gas ratio: %v, tip cap: %v fee cap: %v", gasRatio, tip.String(), fee.String())

			txParam.GasTipCap = tip.BigInt()
			txParam.GasFeeCap = fee.BigInt()

		}

	} else {
		// legacy
		logger.Debug().Msgf("eip1559 off, ignore gasFeeCap and gasTipCap")

		if gasPrice != "" {
			logger.Debug().Msgf("parse gasPrice")
			gp, err := utils.ParseUnits(gasPrice, utils.UnitGwei)
			if err != nil {
				return nil, fmt.Errorf("parse gasPrice: %v error: %w", gasPrice, err)
			}
			txParam.GasPrice = gp
		} else {
			logger.Debug().Msgf("query gasPrice")
			gp, err := client.SuggestGasPrice(ctx)
			if err != nil {
				return nil, fmt.Errorf("query gas price error: %v", err)
			}
			txParam.GasPrice = gp
		}
		logger.Debug().Msgf("gasPrice: %v", txParam.GasPrice.String())

		if gasRatio != "" {
			logger.Debug().Msgf("parse gasRatio")
			gRatio, err := decimal.NewFromString(gasRatio)
			if err != nil {
				return nil, fmt.Errorf("parse gasRatio: %v error: %w", gasRatio, err)
			}
			logger.Debug().Msgf("gas ratio: %v", gRatio.String())

			logger.Debug().Msgf("before gas ratio: %v, gas price: %v", gasRatio, txParam.GasPrice.String())
			gasPrice := decimal.NewFromBigInt(txParam.GasPrice, 0).Mul(gRatio)
			logger.Debug().Msgf("after gas ratio: %v, gas price: %v", gasRatio, gasPrice.BigInt().String())

			txParam.GasPrice = gasPrice.BigInt()

		}

	}

	return txParam, nil
}
