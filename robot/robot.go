package robot

import (
	"errors"
	"fmt"
	"go-jd-assistant/config"
	"go-jd-assistant/jdsdk"
	"go-jd-assistant/util"
	"time"
)

func Run() {
	login()

	//基于时间校准，一个新的触发时间和抢购时间
	c := config.Config
	diffTimeMs := diffLocalServerTime()
	buytime, _ := time.Parse(time.RFC3339, c.Account.Sku.BuyTime)
	//时间格式：06-01-02 03:04:05.000  奇葩的go语言，奇葩的时间格式
	triggerTimeMs := int(buytime.UnixNano()/1000000) - diffTimeMs

	wartToTriggerTimeMs := triggerTimeMs - int(time.Now().UnixNano()/1000000)
	time.Sleep(time.Duration(wartToTriggerTimeMs) * time.Millisecond)

	//开启定时器，准备逼近抢购时间，调度一次
	kill(triggerTimeMs)
}

var qrTimeout = errors.New("验证码扫描超时")

func login() (jdConfig *config.Jd, err error) {
	c := config.Config

	if util.Exists(c.Account.CookieFilePath) {
		jdsdk.ReLoadCookies(c.Account.CookieFilePath)
		if jdsdk.ValidCookie() {
			return &c, nil
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
		jdsdk.SaveCookies(c.Account.CookieFilePath)
		return &c, nil
	}

	return nil, err
}

func diffLocalServerTime() int {
	t1 := time.Now()
	serverTimeMS := jdsdk.GetServerTime()
	elapsed := time.Since(t1)

	diff := int(time.Now().UnixNano()/1000000) - (int(elapsed.Milliseconds())/2 + serverTimeMS)
	return diff
}

func kill(triggerTimeMs int) {
	ac := config.Config.Account

	//应该什么时候出发呢？
	initInfo := jdsdk.GetKillInitInfo(ac.Sku.Id, ac.Sku.Count)
	submitOrderPostData := jdsdk.BuildSubmitOrderPostData(
		ac.Pwd,
		ac.Fp,
		ac.Eid,
		ac.Sku.Id,
		ac.Sku.Count,
		&initInfo,
	)

	//外层方法要提起一点调度。这个方法用来弥补，调度框架的精度问题
	nervousBlockWait(triggerTimeMs)

	//确保第一时间获取到killurl
	killUrl := getKillUrl(ac.Sku.Id)

	jdsdk.RequestKillUrl(ac.Sku.Id, killUrl)

	//多次并发抢购
	submitOrder(ac.Sku.Id, ac.Sku.Count, submitOrderPostData)
}

// 这种循环方式，可能比定时器，会准一点。
func nervousBlockWait(timeMs int) {
	timeNano := int64(timeMs * 1000000)
	fmt.Println("最后的等待")
	for {
		if time.Now().UnixNano() > timeNano {
			break
		}
	}
	fmt.Println("开始执行")
}

func getKillUrl(skuid string) string {
	ch := make(chan string, 200)

	goOn := GoOn{
		value: true,
	}
	fastGetKillUrl(skuid, ch, &goOn)
	killUrl := <-ch

	return killUrl
}

type GoOn struct {
	value bool
}

func fastGetKillUrl(skuId string, c chan string, goon *GoOn) {
	for {
		if !goon.value {
			break
		}

		go func() {
			url := jdsdk.GetKillUrl(skuId)
			if len(url) > 0 {
				c <- url
			}
		}()

		time.Sleep(20 * time.Millisecond)
	}
}

func submitOrder(skuid string, num string, data *map[string]string) {
	fastSubmitOrder(skuid, num, data)
}

func fastSubmitOrder(skuid string, num string, data *map[string]string) {
	tryTimes := 8

	for tryTimes > 0 {
		tryTimes--
		go jdsdk.SubmitOrder(skuid, num, data)
		go jdsdk.SubmitOrder(skuid, num, data)
		time.Sleep(20 * time.Millisecond)
	}
}
