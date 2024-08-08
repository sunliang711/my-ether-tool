package utils

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/shopspring/decimal"
)

const (
	eth  = "1_000_000_000_000_000_000"
	gwei = "1_000_000_000"
)

type Unit int8

const (
	UnitEth Unit = iota
	UnitGwei
)

// FormatUnits ("1000000000000000000","eth") ->  1
// FormatUnits ("1000000000","gwei") -> 1
func FormatUnits(value string, unit Unit) (string, error) {
	var base decimal.Decimal
	switch unit {
	case UnitEth:
		base = decimal.New(1, 18)
	case UnitGwei:
		base = decimal.New(1, 9)
	default:
		return "", errors.New("invalid unit")
	}

	valueDecimal, err := decimal.NewFromString(value)
	if err != nil {
		return "", fmt.Errorf("invalid value: %w", err)
	}

	result := valueDecimal.Div(base)
	return result.String(), nil

}

// ParseUnits ("1.2","eth") -> wei
// ParseUnits ("1.2","gwei") -> wei
func ParseUnits(value string, unit Unit) (ret *big.Int, err error) {
	var base decimal.Decimal
	switch unit {
	case UnitEth:
		base = decimal.New(1, 18)
	case UnitGwei:
		base = decimal.New(1, 9)
	default:
		return nil, errors.New("invalid unit")
	}

	valueDecimal, err := decimal.NewFromString(value)
	if err != nil {
		return nil, fmt.Errorf("invalid value: %w", err)
	}

	result := valueDecimal.Mul(base)

	ret = big.NewInt(0)
	ret, ok := ret.SetString(result.String(), 10)
	if !ok {
		return nil, fmt.Errorf("set result: %v failed", result.String())
	}

	return
}

// 存储在string中的大浮点数相乘并取整
func StringMul(a string, b string) (result *big.Int, err error) {
	af, ok := big.NewFloat(0).SetString(a)
	if !ok {
		err = errors.New("invalid float number")
		return
	}

	bf, ok := big.NewFloat(0).SetString(b)
	if !ok {
		err = errors.New("invalid float number")
		return
	}

	r := af.Mul(af, bf)

	resultStr := r.Text('f', 0)
	result, ok = new(big.Int).SetString(resultStr, 10)
	if !ok {
		err = errors.New("set result failed")
		return
	}

	return

}

func Wei2Gwei(wei string) (string, error) {
	r, err := StringDiv(wei, "1000000000")
	if err != nil {
		return "", err
	}

	return r.String(), nil
}

// 存储在string中的大浮点数相除
func StringDiv(a string, b string) (result *big.Float, err error) {
	af, ok := big.NewFloat(0).SetString(a)
	if !ok {
		err = errors.New("invalid float number")
		return
	}

	bf, ok := big.NewFloat(0).SetString(b)
	if !ok {
		err = errors.New("invalid float number")
		return
	}

	result = af.Quo(af, bf)

	return

}

func Erc20AmountToHuman(amount string, decimals string) (string, error) {
	d, err := strconv.ParseInt(decimals, 10, 32)
	if err != nil {
		return "", fmt.Errorf("parse decimals error: %w", err)
	}

	base := decimal.New(1, int32(d))

	amountDecimal, err := decimal.NewFromString(amount)
	if err != nil {
		return "", fmt.Errorf("invalid original amount: %w", err)
	}

	humanAmount := amountDecimal.Div(base)

	return humanAmount.String(), nil
}

func Erc20AmountFromHuman(humanAmount string, decimals string) (string, error) {
	d, err := strconv.ParseInt(decimals, 10, 32)
	if err != nil {
		return "", fmt.Errorf("parse decimals: %v error: %w", decimals, err)
	}

	base := decimal.New(1, int32(d))

	humanAmountDecimal, err := decimal.NewFromString(humanAmount)
	if err != nil {
		return "", fmt.Errorf("invalid human amount: %w", err)
	}

	amount := humanAmountDecimal.Mul(base)
	return amount.String(), nil

}
