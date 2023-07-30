package transaction

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"my-ether-tool/utils"
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

	addressSliceType abi.Type
	bytes32SliceType abi.Type
	uint256SliceType abi.Type

	bytesSliceType  abi.Type
	stringSliceType abi.Type
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

	addressSliceType, err = abi.NewType("address[]", "", nil)
	utils.ExitWhenError(err, "create address[] failed")
	uint256SliceType, err = abi.NewType("uint256[]", "", nil)
	utils.ExitWhenError(err, "create uint256[] failed")
	bytes32SliceType, err = abi.NewType("bytes32[]", "", nil)
	utils.ExitWhenError(err, "create bytes32[] failed")
	bytesSliceType, err = abi.NewType("bytes[]", "", nil)
	utils.ExitWhenError(err, "create bytes[] failed")
	stringSliceType, err = abi.NewType("string[]", "", nil)
	utils.ExitWhenError(err, "create string[] failed")

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

	var items []string

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
			// Basic Types Begin
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
			// Basic Types End

			// Slice Begin
			// slice类型对应的abiArg是个json数组,数组中每个元素都是一个string类型，不管它原来是什么类型
			case "address[]":
				fmt.Printf("> address[] type\n")
				items, err = DecodeJsonArray(abiArg)
				utils.ExitWhenError(err, "decode address[] arguments error: %s\n", err)

				addresses := []common.Address{}
				for _, item := range items {
					addresses = append(addresses, common.HexToAddress(item))
				}

				arg = abi.Argument{Type: addressSliceType}
				argValue = addresses

			case "uint256[]":
				fmt.Printf("> uint256[] type\n")
				items, err = DecodeJsonArray(abiArg)
				utils.ExitWhenError(err, "decode address[] arguments error: %s\n", err)

				uint256s := []*big.Int{}
				for _, item := range items {
					r, ok := new(big.Int).SetString(item, 10)
					utils.ExitWithMsgWhen(!ok, "set uint256 error: %s", item)

					uint256s = append(uint256s, r)
				}
				// fmt.Printf("uint256s: %v\n", uint256s)

				arg = abi.Argument{Type: uint256SliceType}
				argValue = uint256s

			case "bytes32[]":
				fmt.Printf("> bytes32[] type\n")
				items, err = DecodeJsonArray(abiArg)
				utils.ExitWhenError(err, "decode bytes32[] arguments error: %s\n", err)

				bytes32s := []common.Hash{}

				for _, item := range items {
					bytes32s = append(bytes32s, common.HexToHash(item))
				}
				arg = abi.Argument{Type: bytes32SliceType}
				argValue = bytes32s
			case "bytes[]":
				fmt.Printf("> bytes[] type\n")
				items, err = DecodeJsonArray(abiArg)
				utils.ExitWhenError(err, "decode bytes[] arguments error: %s\n", err)

				var decoded []byte
				byteses := [][]byte{}

				for _, item := range items {
					decoded, err = hex.DecodeString(item)
					utils.ExitWhenError(err, "decode bytes[] error: %s", err)
					byteses = append(byteses, decoded)
				}

				arg = abi.Argument{Type: bytesSliceType}
				argValue = byteses
			case "string[]":
				fmt.Printf("> string[] type\n")
				items, err = DecodeJsonArray(abiArg)
				utils.ExitWhenError(err, "decode string[] arguments error: %s\n", err)

				arg = abi.Argument{Type: stringSliceType}
				argValue = items

			// Slice End

			// Tuple Begin
			// Tuple End
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

func DecodeJsonArray(data string) (items []string, err error) {
	err = json.NewDecoder(strings.NewReader(data)).Decode(&items)
	return
}
