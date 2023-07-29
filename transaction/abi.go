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

	bytesType  abi.Type
	stringType abi.Type
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

	bytesType, err = abi.NewType("bytes", "", nil)
	if err != nil {
		panic("create bytes failed")
	}

	stringType, err = abi.NewType("string", "", nil)
	if err != nil {
		panic("create string failed")
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
			var arg abi.Argument
			var argValue any

			switch types[i] {
			case "address":
				fmt.Printf("> address type\n")
				arg = abi.Argument{Type: addressType}
				argValue = common.HexToAddress(abiArg)
			case "uint256":
				fmt.Printf("> uint256 type\n")
				arg = abi.Argument{Type: uint256Type}
				v, ok := new(big.Int).SetString(abiArg, 10)
				if !ok {
					err = errors.New("invalid uint256 type argument")
					return
				}
				argValue = v
			case "bytes32":
				fmt.Printf("> bytes32 type\n")
				arg = abi.Argument{Type: bytes32Type}
				argValue = common.HexToHash(abiArg)
			case "bytes":
				fmt.Printf("> bytes type\n")
				arg = abi.Argument{Type: bytesType}
				var decoded []byte
				decoded, err = hex.DecodeString(abiArg)
				if err != nil {
					err = fmt.Errorf("invalid bytes type argument: %s", err)
					return
				}
				argValue = decoded
			case "string":
				fmt.Printf("> string type\n")
				arg = abi.Argument{Type: stringType}
				argValue = abiArg
			// TODO: other types
			default:
				err = fmt.Errorf("not supprt type: %s", types[i])
				return
			}
			arguments = append(arguments, arg)
			argumentsValue = append(argumentsValue, argValue)

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
