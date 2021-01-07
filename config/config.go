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
	//也许配置文件和代码分离的设计。也是go特意考虑的，代码不管在哪里运行。都能做到代码和文件分离。
	//所以，你懂得，在哪里运行，请把config.yml文件拷贝到哪里
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("config read error，请理解一下go代码和配置分离的思想。在哪里运行，把config拷贝到相对执行文件的位置" +
			"凡是用到这个config.yml文件的代码，在执行的位置，拷贝一份")
	}

	//todo 大坑：这个操作。结构体内的属性，必须大写开头。不然搞不进去！！！！
	viper.Unmarshal(&Config)
}
