package contract

import (
	"encoding/hex"
	"fmt"
	"math/big"
	utils "met/utils"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// 根据abiType，把其string类型的值转换成go类型
func parseAbiType(abiType abi.Type, arg string) (any, error) {
	switch abiType.T {
	case abi.BoolTy:
		return parseBoolType(arg)

	case abi.IntTy:
		return parseIntType(abiType.String(), arg)

	case abi.UintTy:
		return parseUintType(abiType.String(), arg)

	case abi.AddressTy:
		return parseAddressType(arg)

	case abi.StringTy:
		return arg, nil

	case abi.BytesTy:
		return parseBytesType(arg)

	case abi.FixedBytesTy:
		panic("not support fixedBytes type")
	case abi.ArrayTy:
		panic("not support array type")
	case abi.TupleTy:
		panic("not support tuple type")
	default:
		panic("invalid abi type")

	}
}

func parseBoolType(arg string) (any, error) {
	switch arg {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return 0, fmt.Errorf("invalid bool type arg: %v", arg)
	}
}

// 解析intN类型
// 分为两种情况:
// 1. go中已经有的int类型(int8 int16 int32 int64)就用go已有的
// 2. go中没有的类型(比如:int22 int53 int256等等)，用*big.Int
func parseIntType(stringKind, arg string) (any, error) {
	if !strings.HasPrefix(stringKind, "int") {
		return 0, fmt.Errorf("invalid int type: %v", stringKind)
	}

	bitNumberStr := strings.TrimPrefix(stringKind, "int")
	bitNumber, err := strconv.ParseUint(bitNumberStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse bit number string: %v error: %w", bitNumberStr, err)
	}
	if bitNumber == 0 || bitNumber > 256 {
		return 0, fmt.Errorf("bitnumber range invalid: %v", bitNumber)
	}

	switch bitNumber {
	case 8:
		v, err := parseGoInt(arg)
		if err != nil {
			return 0, err
		}
		return int8(v), nil

	case 16:
		v, err := parseGoInt(arg)
		if err != nil {
			return 0, err
		}
		return int16(v), nil

	case 32:
		v, err := parseGoInt(arg)
		if err != nil {
			return 0, err
		}
		return int32(v), nil

	case 64:
		v, err := parseGoInt(arg)
		if err != nil {
			return 0, err
		}
		return int64(v), nil
	default: // *big.Int
		v, ok := big.NewInt(0).SetString(arg, 10)
		if !ok {
			return 0, fmt.Errorf("invalid intN value: %v", arg)
		}
		return v, nil
	}

}

func parseGoInt(arg string) (int64, error) {
	v, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse go int value: %v error: %w", arg, err)
	}
	return v, nil
}

// 解析uintN类型
// 分为两种情况:
// 1. go中已经有的uint类型(uint8 uint16 uint32 uint64)就用go已有的
// 2. go中没有的类型(比如:uint22 uint53 uint256等等)，用*big.Int
func parseUintType(stringKind, arg string) (any, error) {
	if !strings.HasPrefix(stringKind, "uint") {
		return 0, fmt.Errorf("invalid uint type: %v", stringKind)
	}

	bitNumberStr := strings.TrimPrefix(stringKind, "uint")
	bitNumber, err := strconv.ParseUint(bitNumberStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse bit number string: %v error: %w", bitNumberStr, err)
	}
	if bitNumber == 0 || bitNumber > 256 {
		return 0, fmt.Errorf("bitnumber range invalid: %v", bitNumber)
	}

	switch bitNumber {
	case 8:
		v, err := parseGoUint(arg)
		if err != nil {
			return 0, err
		}
		return uint8(v), nil

	case 16:
		v, err := parseGoUint(arg)
		if err != nil {
			return 0, err
		}
		return uint16(v), nil

	case 32:
		v, err := parseGoUint(arg)
		if err != nil {
			return 0, err
		}
		return uint32(v), nil

	case 64:
		v, err := parseGoUint(arg)
		if err != nil {
			return 0, err
		}
		return uint64(v), nil
	default: // *big.Int
		v, ok := big.NewInt(0).SetString(arg, 10)
		if !ok {
			return 0, fmt.Errorf("invalid uintN value: %v", arg)
		}
		return v, nil
	}

}

func parseGoUint(arg string) (uint64, error) {
	v, err := strconv.ParseUint(arg, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse go uint value: %v error: %w", arg, err)
	}
	return v, nil
}

func parseAddressType(arg string) (common.Address, error) {
	if !utils.IsValidAddress(arg) {
		return common.Address{}, fmt.Errorf("invalid address: '%v'", arg)
	}

	return common.HexToAddress(arg), nil

}

func parseBytesType(arg string) (any, error) {
	arg = strings.ToLower(arg)
	arg = strings.TrimPrefix(arg, "0x")

	bytes, err := hex.DecodeString(arg)
	if err != nil {
		return nil, fmt.Errorf("decode bytes: %v error: %w", arg, err)
	}

	return bytes, nil
}
