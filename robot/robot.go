package robot

import (
	"errors"
	"go-jd-assistant/config"
	"go-jd-assistant/jdsdk"
	"go-jd-assistant/util"
	"time"
)

func Run() {
	//设计为，定时登录。并且开机验证。登录cookie并且设置cookie登录时间，
	//计算抢购时间延后一小时，保证30分钟前有一次验证。一次验证周期为30分钟。逆向推算首次验证的周期
	//保证抢购时间延后一个小时，一定距离现在小于24小时。
	jd, err := login()
	if err != nil {
		//send 登录失败，重新登录。
	}
	//todo 启动定时器
	reserve(jd)
	syncTime()
	//todo 启动定时器
	kill(jd)
}

var qrTimeout = errors.New("验证码扫描超时")

func login() (jdConfig *config.Jd, err error) {
	// load config
	jd := config.Config

	if len(jd.Account.Cookie) > 0 && jdsdk.ValidCookie(jd.Account.Cookie) {
		//valid exist cookie suc as logined
		return &jd, nil
	}

	//reLogin
	jdsdk.GetLoginPage()

	QRFilePath := "./qrcode.png"
	token := jdsdk.GetQR(QRFilePath)

	//todo && send email
	util.Open(QRFilePath)

	checkLeftTimes := 5
	var ticket string
	for {
		time.Sleep(4 * time.Second)
		ticket = jdsdk.GetQrTicket(token)
		if len(ticket) > 0 {
			break
		}

		if checkLeftTimes--; checkLeftTimes == 0 {
			return nil, qrTimeout
		}
	}

	if jdsdk.ValidQRTicket(ticket) {
		jdsdk.GetUserInfo()
		//todo 获取cookie
		cookie := "parse from logined request"
		jd.Account.SetCookie(cookie)
		return &jd, nil
	}

	return nil, err
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
