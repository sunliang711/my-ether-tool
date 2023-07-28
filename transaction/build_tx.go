package transaction

import (
	"context"
	"encoding/hex"
	"errors"
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
// @param data  optional
// @param gasLimit optional 为0时从rpc获取
// @param nonce optional 为空时从rpc获取
// @param chainID optional 为0时从rpc获取
// @param gasPrice used when eip1559 is false unit: gwei
// @param gasTipCap used when eip1559 is true unit: gwei
// @param gasFeeCap used when eip1559 is true unit: gwei
// @param eip1559:
// eip1559为true时，当gasTipCap 和 gasFeeCap都不为空时使用它们，否则从rpc获取这两个值
// eip1559为false时，当gasPrice不为空时使用gasPrice，否则从rpc获取
func BuildTransaction(rpc string, from string, to string, value string, data string, gasLimit uint64, nonce string, chainID int, gasPrice string, gasTipCap string, gasFeeCap string, eip1559 bool) (tx *types.Transaction, err error) {
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

	value1, err := utils.ParseUnits(value, "eth")
	if err != nil {
		return
	}

	if data != "" {
		data0, err = hex.DecodeString(data)
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
		nonce0, err = client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return
		}
	} else {
		nonce0, err = strconv.ParseUint(nonce, 10, 64)
		if err != nil {
			return
		}
	}

	// if chainID == 0; get by rpc
	if chainID == 0 {
		chainID0, err = client.ChainID(context.Background())
		if err != nil {
			return
		}
	} else {
		chainID0 = big.NewInt(int64(chainID))
	}

	// if gasLimit == 0; estimateGas
	if gasLimit == 0 {
		gasLimit, err = client.EstimateGas(context.Background(), ethereum.CallMsg{
			From:  fromAddress,
			To:    &toAddress,
			Value: value1,
			Data:  data0,
		})
		if err != nil {
			return
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
			var gasPrice0 *big.Int
			gasPrice0, err = client.SuggestGasPrice(context.Background())
			if err != nil {
				return
			}

			tipCap, err = client.SuggestGasTipCap(context.Background())
			if err != nil {
				return
			}

			baseFee := new(big.Int).Sub(gasPrice0, tipCap)
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
			price, err = client.SuggestGasPrice(context.Background())
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
