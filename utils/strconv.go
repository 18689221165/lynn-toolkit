package utils

import (
	"github.com/shopspring/decimal"
	"strconv"
)

func Atoi(s string) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return 0
}

func NewFromString(value string) decimal.Decimal {
	if v, err := decimal.NewFromString(value); err == nil {
		return v
	}
	return decimal.Zero
}
