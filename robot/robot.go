package robot

import (
	"go-jd-assistant/config"
	"go-jd-assistant/jdsdk"
)

func Run() {
	jd := login()
	//todo 启动定时器
	reserve(jd)
	syncTime()
	//todo 启动定时器
	kill(jd)
}

func login() config.Jd {
	//load from file
	config := config.InitConfig("config")
	cookie := config.Account.GetCookie()

	//valid cookie
	jdsdk.ValidCookie(cookie)
	//login requests logic
	jdsdk.GetLoginPage()
	jdsdk.GetQR()

	jdsdk.GetQrTicket()

	jdsdk.ValidQRTicket("")
	jdsdk.GetUserInfo()
	config.Account.SetCookie("")
	return config
}

func reserve(config config.Jd) {

}

func kill(config config.Jd) {
	skuId := config.Account.Sku.Id
	num := config.Account.Sku.Count
	jdsdk.GetKillInitInfo(skuId, num)

	//todo 定时，并发
	jdsdk.GetKillUrl(skuId)
	killUrl := "sdfsdf"
	jdsdk.RequestKillUrl(skuId, killUrl)
	rid := "fdsfsdfsdfd"
	jdsdk.SubmitOrder(skuId, num, rid)
}

func syncTime() {
	jdsdk.GetServerTime()
}
