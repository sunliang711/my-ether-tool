package transaction

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	uint256Type abi.Type
	bytes32Type abi.Type
	addressType abi.Type
)

func init() {

	var err error
	uint256Type, err = abi.NewType("uint256", "", nil)
	if err != nil {
		panic("create uint256 failed")
	}
	bytes32Type, err = abi.NewType("bytes32", "", nil)
	if err != nil {
		panic("create bytes32 failed")
	}
	addressType, err = abi.NewType("address", "", nil)
	if err != nil {
		panic("create address failed")
	}
}

func AbiEncode(abiStr string, abiArgs []string) (result []byte, err error) {
	// selector is the first 4 bytes of hash of abiStr
	selector := crypto.Keccak256Hash([]byte(abiStr)).Bytes()[:4]

	result = selector

	fields := strings.FieldsFunc(abiStr, func(r rune) bool {
		return r == '(' || r == ')'
	})

	arguments := abi.Arguments{}
	argumentsValue := []any{}

	if len(fields) == 2 {
		argTypes := fields[1]

		types := strings.Split(argTypes, ",")

		if len(types) != len(abiArgs) {
			err = errors.New("number of arguments  not match with abi string")
			return

		}

		for i := range types {
			abiArg := (abiArgs)[i]
			fmt.Printf("> abi arg: %s\n", abiArg)

			switch types[i] {
			case "address":
				fmt.Printf("> address type\n")
				arguments = append(arguments, abi.Argument{Type: addressType})
				argumentsValue = append(argumentsValue, common.HexToAddress(abiArg))
			case "uint256":
				fmt.Printf("> uint256 type\n")
				arguments = append(arguments, abi.Argument{Type: uint256Type})
				v, ok := new(big.Int).SetString(abiArg, 10)
				if !ok {
					err = errors.New("invalid uint256 type argument")
					return
				}
				argumentsValue = append(argumentsValue, v)
			case "bytes32":
				var decoded []byte
				fmt.Printf("> bytes32 type\n")
				arguments = append(arguments, abi.Argument{Type: bytes32Type})
				decoded, err = hex.DecodeString(abiArg)
				if err != nil {
					err = fmt.Errorf("invalid bytes32 type argument: %s", err)
					return
				}
				argumentsValue = append(argumentsValue, decoded)
			// TODO: other types
			default:
				err = fmt.Errorf("not supprt type: %s", types[i])
				return
			}
		}
		var packed []byte
		packed, err = arguments.Pack(argumentsValue...)
		if err != nil {
			err = fmt.Errorf("pack arguments error: %s", err)
			return
		}
		result = append(result, packed...)

	}

	return
}
