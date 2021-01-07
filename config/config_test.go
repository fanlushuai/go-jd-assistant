package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Println(Config)

	for index := range Config.Account.Skus {
		sku := Config.Account.Skus[index]
		fmt.Println(sku)
	}
}
