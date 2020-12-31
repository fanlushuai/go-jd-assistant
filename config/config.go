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

func InitConfig(configPath string) Jd {
	//相对工作空间的目录。或者使用"." 当前目录
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("config read error")
	}

	//todo 大坑：这个操作。结构体内的属性，必须大写开头。不然搞不进去！！！！
	var jd Jd
	viper.Unmarshal(&jd)
	return jd
}

func (ac *Account) GetCookie() string {
	return ac.Cookie
}

func (ac *Account) SetCookie(cookie string) {
	ac.Cookie = cookie
	//写入配置文件
	v := viper.New()
	v.Set("cookie", ac.Cookie)
	v.WriteConfig()
}
