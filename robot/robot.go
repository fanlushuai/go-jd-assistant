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
	login()

	//todo 启动定时器
	reserve()
	//syncTime()
	////todo 启动定时器
	kill()
}

var qrTimeout = errors.New("验证码扫描超时")

func login() (jdConfig *config.Jd, err error) {
	// load config
	jd := config.Config

	if util.Exists(jd.Account.CookieFilePath) {
		jdsdk.ReLoadCookies(jd.Account.CookieFilePath)
		if jdsdk.ValidCookie() {
			return &jd, nil
		}
	}

	//reLogin
	jdsdk.GetLoginPage()

	QRFilePath := "./qrcode.png"
	token := jdsdk.GetQR(QRFilePath)

	//you can do this for server deploy
	//-> util.SendEmailOutlook()
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
		jdsdk.SaveCookies(jd.Account.CookieFilePath)
		return &jd, nil
	}

	return nil, err
}

func reserve() {

}

func kill() {
	ac := config.Config.Account

	initInfo := jdsdk.GetKillInitInfo(ac.Sku.Id, ac.Sku.Count)
	submitOrderPostData := jdsdk.BuildSubmitOrderPostData(
		ac.Pwd,
		ac.Fp,
		ac.Eid,
		ac.Sku.Id,
		ac.Sku.Count,
		&initInfo,
	)

	killUrl := jdsdk.GetKillUrl(ac.Sku.Id)
	jdsdk.RequestKillUrl(ac.Sku.Id, killUrl)

	jdsdk.SubmitOrder(ac.Sku.Id, ac.Sku.Count, submitOrderPostData)
}

func diffLocalServerTime() int {
	t1 := time.Now()
	serverTimeMS := jdsdk.GetServerTime()
	elapsed := time.Since(t1)

	diff := time.Now().Nanosecond()/1000 - (int(elapsed.Milliseconds())/2 + serverTimeMS)
	return diff
}
