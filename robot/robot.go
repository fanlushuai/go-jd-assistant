package robot

import (
	"errors"
	"fmt"
	"go-jd-assistant/config"
	"go-jd-assistant/jdsdk"
	"go-jd-assistant/util"
	"time"
)

var c config.Jd
var ac config.Account

func init() {
	c = config.Config
	ac = config.Config.Account
}

func Run() {
	login()

	//基于时间校准，一个新的触发时间和抢购时间
	triggerTimeMs := getTriggerTime()

	//不用sleep准点触发，对其精度表示怀疑。提前5s
	waitToTriggerTimeMs := triggerTimeMs - int(time.Now().UnixNano()/1000000)
	time.Sleep(time.Duration(waitToTriggerTimeMs-5*1000) * time.Millisecond)

	//先把初始化数据搞下，抢的时候，不浪费时间
	submitOrderPostData := getSubmitOrderPostData(ac.Sku)
	//消耗一下cpu，触发,这种方式也许会准点
	nervousBlockWait(triggerTimeMs)

	kill(submitOrderPostData)
}

func login() (err error) {
	if util.Exists(c.Account.CookieFilePath) {
		jdsdk.ReLoadCookies(c.Account.CookieFilePath)
		if jdsdk.ValidCookie() {
			fmt.Println("本地cookie登录成功")
			return nil
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
			return errors.New("验证码扫描超时")
		}
	}

	if jdsdk.ValidQRTicket(ticket) {
		jdsdk.SaveCookies(c.Account.CookieFilePath)
		return nil
	}

	return errors.New("预期流程未能正确登录，请检查代码")
}

func getTriggerTime() int {
	diffTimeMs := diffLocalServerTime()
	buytime, _ := time.Parse(time.RFC3339, c.Account.Sku.BuyTime)
	//时间格式：06-01-02 03:04:05.000  奇葩的go语言，奇葩的时间格式
	triggerTimeMs := int(buytime.UnixNano()/1000000) - diffTimeMs
	return triggerTimeMs
}

func diffLocalServerTime() int {
	t1 := time.Now()
	serverTimeMS := jdsdk.GetServerTime()
	elapsed := time.Since(t1)

	diff := int(time.Now().UnixNano()/1000000) - (int(elapsed.Milliseconds())/2 + serverTimeMS)
	return diff
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

func kill(submitOrderPostData *map[string]string) {
	//确保第一时间获取到killurl
	killUrl := getKillUrl(ac.Sku.Id)

	jdsdk.RequestKillUrl(ac.Sku.Id, killUrl)

	//多次并发抢购
	submitOrder(ac.Sku.Id, ac.Sku.Count, submitOrderPostData)
}

func getSubmitOrderPostData(sku config.Sku) *map[string]string {
	initInfo := jdsdk.GetKillInitInfo(sku.Id, sku.Count)
	submitOrderPostData := jdsdk.BuildSubmitOrderPostData(
		ac.Pwd,
		ac.Fp,
		ac.Eid,
		ac.Sku.Id,
		ac.Sku.Count,
		&initInfo,
	)
	return submitOrderPostData
}

func getKillUrl(skuId string) string {
	ch := make(chan string, 1314)

	go func(ch chan string) {

		type Goon struct {
			value bool
		}

		giveUpLoopLeftTimes := 100
		goon := Goon{value: true}
		for {

			if !goon.value || giveUpLoopLeftTimes == 0 {
				break
			}

			go func(g *Goon, ch chan string) {
				url := jdsdk.GetKillUrl(skuId)
				if len(url) > 0 {
					ch <- url
					g.value = false
				}
			}(&goon, ch)

			time.Sleep(20 * time.Millisecond)
			giveUpLoopLeftTimes--
		}

	}(ch)

	killUrl := <-ch

	return killUrl
}

func submitOrder(skuId string, num string, data *map[string]string) {
	tryTimes := 8

	for tryTimes > 0 {
		tryTimes--
		go jdsdk.SubmitOrder(skuId, num, data)
		go jdsdk.SubmitOrder(skuId, num, data)
		time.Sleep(20 * time.Millisecond)
	}
}
