package robot

import (
	"errors"
	"fmt"
	"go-jd-assistant/config"
	"go-jd-assistant/jdsdk"
	"go-jd-assistant/util"
	"os"
	"sync"
	"time"
)

var c config.Jd
var ac config.Account

func init() {
	c = config.Config
	ac = config.Config.Account
	//jdsdk.Proxy("http://localhost:8888")
}

func Run() {
	login()
	diffTimeMs := diffLocalServerTime()
	fmt.Println("服务器本地时差", diffTimeMs, "毫秒")

	var wg sync.WaitGroup
	wg.Add(len(ac.Skus))

	for index := range ac.Skus {
		//注意，用索引的方式，可以获取到值，而不是值copy.可以配合动态监听配置文件，修改一些内容
		sku := ac.Skus[index]
		fmt.Println("开启doSku id=", sku.Id)
		go func() {
			doSku(sku, diffTimeMs)
			wg.Done()
		}()
	}

	wg.Wait()
}

func login() (err error) {
	if util.Exists(c.Account.CookieFilePath) {
		jdsdk.ReLoadCookies(c.Account.CookieFilePath)
		if jdsdk.ValidCookie() {
			fmt.Println("本地cookie 登录成功")
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
			fmt.Println("验证码扫描超时")
			os.Exit(-1)
			return errors.New("验证码扫描超时")
		}
	}

	if jdsdk.ValidQRTicket(ticket) {
		jdsdk.SaveCookies(c.Account.CookieFilePath)
		fmt.Println("QR 登录成功")
		return nil
	}

	return errors.New("预期流程未能正确登录，请检查代码")
}

//this can do for many sku if you need,code like: for sku skus {doSku(sku)}
func doSku(sku config.Sku, diffTimeMs int) {

	//基于时间校准，一个新的触发时间和抢购时间
	triggerTimeMs := getTriggerTime(sku, diffTimeMs)
	fmt.Println("校准后的触发时间", triggerTimeMs, "毫秒")

	//不用sleep准点触发，对其精度表示怀疑。提前5s
	waitToTriggerTimeMs := triggerTimeMs - int(time.Now().UnixNano()/1000000)
	fmt.Println("触发时间提前25秒获取一下基本信息。基本信息接口有频率限制，防止一次获取不成功")
	time.Sleep(time.Duration(waitToTriggerTimeMs-25*1000) * time.Millisecond)

	//先把初始化数据搞下，抢的时候，不浪费时间.这接口很危险。不知道会处理多久
	submitOrderPostData := getSubmitOrderPostData(sku)
	fmt.Println("获取秒杀需要的基本信息", submitOrderPostData)

	if triggerTimeMs-int(time.Now().UnixNano()/1000000) > 5000 {
		fmt.Println("触发时间提前5秒进入等待……")
		time.Sleep(time.Duration(waitToTriggerTimeMs-5*1000) * time.Millisecond)
	}

	//消耗一下cpu，触发,这种方式也许会准点
	fmt.Println("进入cpu紧张等待……")
	nervousBlockWait(triggerTimeMs)

	fmt.Println("进入秒杀！！！！！！！")
	kill(sku, submitOrderPostData)
}

func getTriggerTime(sku config.Sku, diffTimeMs int) int {
	buytime, _ := time.Parse(time.RFC3339, sku.BuyTime)
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

func kill(sku config.Sku, submitOrderPostData *map[string]string) {
	//确保第一时间获取到killurl
	killUrl := getKillUrl(sku.Id)
	fmt.Println("获取到killurl", killUrl)

	jdsdk.RequestKillUrl(sku.Id, killUrl)

	//多次并发抢购
	fmt.Println("提交")
	submitOrder(sku.Id, sku.Count, submitOrderPostData)
}

// 这个接口容易被频率限制，大概是5s。持续5s，返回null数据
func getSubmitOrderPostData(sku config.Sku) *map[string]string {
	for {
		initInfo, err := jdsdk.GetKillInitInfo(sku.Id, sku.Count)

		if err != nil {
			continue
		}

		time.Sleep(1 * time.Second)

		submitOrderPostData := jdsdk.BuildSubmitOrderPostData(
			ac.Pwd,
			ac.Fp,
			ac.Eid,
			sku.Id,
			sku.Count,
			&initInfo,
		)
		return submitOrderPostData
	}
}

//这种设计主要考虑，第一个请求，并没有成功。
//为了尽快的获取，所以，悲观的认为第一个失败了，但是其实我们并不知道第一个的情况

//间隔20毫秒，假设第一个请求失败了。那么可能第二个请求成功了，那么，我们相当于，一个正常的请求时间+20毫秒就能获取到。
//如果，第二个，第三个，第n个都失败了,直到n+1成功了。
//公式为：获取时间=n*20+正常的一个请求的耗时
//可以看到，如果n=0.那么我们就是一个请求的耗时。
//如果0=1，那么我们就是需要一个请求耗时+20ms.，

//对比其他两种方式：
//1.单纯for循环重试的耗时为，ReqT1+Tsleep+ReqT2
//2.直接上来并发5个，来进行健壮性增强。认为，总会有一个ok的。这种方式，大多数的耗时为ReqT1.
//  但是，这个只适合是单次健壮性增强的场景。此处还存在时间的问题。可能5个并发，全部命中了，killurl未刷新的情况。
//  所以，可能还是要继续轮训。和第一个就一样了。
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
