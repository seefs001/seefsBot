package util

import (
	"fmt"
	"testing"
)

func TestGetCoinPrice(t *testing.T) {
	price, err := GetCoinPrice("FIL")
	if err != nil {
		panic(err)
	}
	fmt.Println(price)
}
