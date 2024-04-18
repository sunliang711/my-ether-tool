package transaction

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// 把output中的any类型，断言成实际的类型(*big.Int,int8,int16, common.Address等等)，然后转换成字符串
func decodeOutput(abiType abi.Type, output any) (string, error) {
	switch abiType.T {
	case abi.BoolTy:
		return decodeBoolType(output)

	case abi.IntTy:
		return decodeIntType(abiType.String(), output)

	case abi.UintTy:
		return decodeUintType(abiType.String(), output)

	case abi.AddressTy:
		return decodeAddressType(output)

	case abi.StringTy:
		return decodeStringType(output)

	case abi.BytesTy:
		return decodeBytesType(output)

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

func decodeBoolType(output any) (string, error) {
	if v, ok := output.(bool); ok {
		return fmt.Sprintf("%v", v), nil
	} else {
		return "", fmt.Errorf("output value: %v is not bool type", output)
	}
}

func decodeAddressType(output any) (string, error) {
	if v, ok := output.(common.Address); ok {
		return v.Hex(), nil
	} else {
		return "", fmt.Errorf("output value: %v is not address type", output)
	}
}

func decodeStringType(output any) (string, error) {
	if v, ok := output.(string); ok {
		return v, nil
	} else {
		return "", fmt.Errorf("output value: %v is not string type", output)
	}
}

func decodeBytesType(output any) (string, error) {
	if v, ok := output.([]byte); ok {
		return "0x" + hex.EncodeToString(v), nil
	} else {
		return "", fmt.Errorf("output value: %v is not bytes type", output)
	}
}

// 解析intN类型
// 分为两种情况:
// 1. go中已经有的int类型(int8 int16 int32 int64)就用go已有的
// 2. go中没有的类型(比如:int22 int53 int256等等)，用*big.Int
func decodeIntType(stringKind string, output any) (string, error) {
	if !strings.HasPrefix(stringKind, "int") {
		return "", fmt.Errorf("invalid int type: %v", stringKind)
	}

	bitNumberStr := strings.TrimPrefix(stringKind, "int")
	bitNumber, err := strconv.ParseUint(bitNumberStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse bit number string: %v error: %w", bitNumberStr, err)
	}
	if bitNumber == 0 || bitNumber > 256 {
		return "", fmt.Errorf("bitnumber range invalid: %v", bitNumber)
	}

	switch bitNumber {
	case 8:
		if v, ok := output.(int8); ok {
			return fmt.Sprintf("%v", v), nil
		} else {
			return "", fmt.Errorf("output value: %v is not int8 type", output)
		}

	case 16:
		if v, ok := output.(int16); ok {
			return fmt.Sprintf("%v", v), nil
		} else {
			return "", fmt.Errorf("output value: %v is not int16 type", output)
		}

	case 32:
		if v, ok := output.(int32); ok {
			return fmt.Sprintf("%v", v), nil
		} else {
			return "", fmt.Errorf("output value: %v is not int32 type", output)
		}

	case 64:
		if v, ok := output.(int64); ok {
			return fmt.Sprintf("%v", v), nil
		} else {
			return "", fmt.Errorf("output value: %v is not int64 type", output)
		}

	default: // *big.Int
		if v, ok := output.(*big.Int); ok {
			return v.String(), nil
		} else {
			return "", fmt.Errorf("output value: %v is not *big.Int type", output)
		}
	}
}

// 解析uintN类型
// 分为两种情况:
// 1. go中已经有的uint类型(uint8 uint16 uint32 uint64)就用go已有的
// 2. go中没有的类型(比如:uint22 uint53 uint256等等)，用*big.Int
func decodeUintType(stringKind string, output any) (string, error) {
	if !strings.HasPrefix(stringKind, "uint") {
		return "", fmt.Errorf("invalid uint type: %v", stringKind)
	}

	bitNumberStr := strings.TrimPrefix(stringKind, "uint")
	bitNumber, err := strconv.ParseUint(bitNumberStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("parse bit number string: %v error: %w", bitNumberStr, err)
	}
	if bitNumber == 0 || bitNumber > 256 {
		return "", fmt.Errorf("bitnumber range invalid: %v", bitNumber)
	}

	switch bitNumber {
	case 8:
		if v, ok := output.(uint8); ok {
			return fmt.Sprintf("%v", v), nil
		} else {
			return "", fmt.Errorf("output value: %v is not uint8 type", output)
		}

	case 16:
		if v, ok := output.(uint16); ok {
			return fmt.Sprintf("%v", v), nil
		} else {
			return "", fmt.Errorf("output value: %v is not uint16 type", output)
		}

	case 32:
		if v, ok := output.(uint32); ok {
			return fmt.Sprintf("%v", v), nil
		} else {
			return "", fmt.Errorf("output value: %v is not uint32 type", output)
		}

	case 64:
		if v, ok := output.(uint64); ok {
			return fmt.Sprintf("%v", v), nil
		} else {
			return "", fmt.Errorf("output value: %v is not uint64 type", output)
		}

	default: // *big.Int
		if v, ok := output.(*big.Int); ok {
			return v.String(), nil
		} else {
			return "", fmt.Errorf("output value: %v is not *big.Int type", output)
		}
	}
}
