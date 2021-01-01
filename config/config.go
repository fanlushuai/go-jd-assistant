package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Jd struct {
	Account Account
}

type Account struct {
	Name        string
	Pwd         string
	Cookie      string
	Eid         string
	Fp          string
	TrackId     string
	RiskControl string
	Sku         Sku
}

type Sku struct {
	Allow       bool
	Id          string
	Count       int
	ReserveTime string
	BuyTime     string
}

var Config Jd

func init() {
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

func writeCookie(cookie string) {
	viper.Set("account.cookie", cookie)
	viper.WriteConfig()
}

func (ac *Account) GetCookie() string {
	return ac.Cookie
}

func (ac *Account) SetCookie(cookie string) {
	ac.Cookie = cookie
	//写入配置文件
	writeCookie(cookie)
}
