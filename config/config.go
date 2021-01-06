package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Jd struct {
	Account Account
}

type Account struct {
	Pwd            string
	CookieFilePath string
	Eid            string
	Fp             string
	TrackId        string
	RiskControl    string
	Sku            Sku
}

type Sku struct {
	Id      string
	Count   string
	BuyTime string
}

var Config Jd

func init() {
	//关于配置文件的路径问题。难搞。哈哈哈
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("config read error")
	}

	//todo 大坑：这个操作。结构体内的属性，必须大写开头。不然搞不进去！！！！
	viper.Unmarshal(&Config)
}
