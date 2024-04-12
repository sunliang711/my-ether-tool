package transaction

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"my-ether-tool/utils"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// @param from  from address of tx
// @param to    to address of tx
// @param value amount of eth to be transfered, unit: ETH
// @param data  optional   conflict with abi
// @param abi  optional conflict with data
// @param args  optional
// @param gasLimit optional 为0时从rpc获取
// @param nonce optional 为空时从rpc获取
// @param chainID optional 为0时从rpc获取
// @param gasPrice used when eip1559 is false unit: gwei
// @param gasTipCap used when eip1559 is true unit: gwei
// @param gasFeeCap used when eip1559 is true unit: gwei
// @param eip1559:
// eip1559为true时，当gasTipCap 和 gasFeeCap都不为空时使用它们，否则从rpc获取这两个值
// eip1559为false时，当gasPrice不为空时使用gasPrice，否则从rpc获取
func BuildTransaction(ctx context.Context, rpc string, from string, to string, value *string, data string, abi string, args []string, gasLimit uint64, nonce string, chainID int, gasPrice string, gasTipCap string, gasFeeCap string, eip1559 bool, sendAll bool) (tx *types.Transaction /*newValue *big.Int,*/, err error) {
	// check params

	var (
		nonce0   uint64
		chainID0 *big.Int = big.NewInt(0)
		// value0   *big.Float = big.NewFloat(0)
		data0 []byte
		ok    bool
	)
	from = strings.TrimPrefix(from, "0x")
	to = strings.TrimPrefix(to, "0x")
	data = strings.TrimPrefix(data, "0x")

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
		fmt.Printf("sendAll conflict with eip1559, disable eip1559\n")
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

	client, err := ethclient.Dial(rpc)
	if err != nil {
		return
	}
	defer client.Close()

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

	gWei := big.NewInt(1_000_000_000)

	if eip1559 {
		// use gasTipCap gasFeeCap
		var tipCap *big.Int
		var feeCap *big.Int

		if gasTipCap != "" && gasFeeCap != "" {
			// baseFee=eth_gasPrice - maxPriorityFeePerGas
			// 最大费用，maxPriorityFeePerGas + 2 * base_fee
			tipCap, ok = new(big.Int).SetString(gasTipCap, 10)
			if !ok {
				err = errors.New("set gas tip cap failed")
				return
			}
			tipCap = tipCap.Mul(tipCap, gWei)
			feeCap, ok = new(big.Int).SetString(gasFeeCap, 10)
			if !ok {
				err = errors.New("set gas fee cap failed")
				return
			}
			feeCap = feeCap.Mul(feeCap, gWei)

		} else {
			// get by rpc
			// var gasPrice0 *big.Int
			// gasPrice0, err = client.SuggestGasPrice(ctx)
			// if err != nil {
			// 	return
			// }

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
			price, ok = new(big.Int).SetString(gasPrice, 10)
			if !ok {
				err = errors.New("set gasPrice failed")
				return
			}
			price = price.Mul(price, gWei)
		}

		// 重新计算value
		if sendAll {
			// 查询当前balance
			currentBalance, err := client.BalanceAt(ctx, common.HexToAddress(from), nil)
			if err != nil {
				return nil, fmt.Errorf("query balance error: %w", err)
			}

			// 计算手续费 = 21000 * gasPrice
			txFee := big.NewInt(OnlyTransferGas)
			txFee.Mul(txFee, price)
			fmt.Printf("sendAll txFee: %s\n", txFee.String())

			// 剩下的value为所有待发送value
			value1 = currentBalance.Sub(currentBalance, txFee)
			fmt.Printf("sendAll value: %s\n", value1.String())

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
