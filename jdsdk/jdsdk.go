package jdsdk

import (
	"fmt"
	"github.com/asmcos/requests"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GetLoginPage() {
	url := "https://passport.jd.com/new/login.aspx"
	header := requests.Header{
		"user-agent": userAgent,
	}
	requests.Get(url, header)
}

func GetQR() {
	url := "https://qr.m.jd.com/show"
	header := requests.Header{
		"user-agent": userAgent,
		"Referer":    "https://passport.jd.com/new/login.aspx",
	}
	param := requests.Params{
		"appid": "133",
		"size":  "147",
		"t":     genTime(),
	}

	resp, err := requests.Get(url, header, param)
	if err != nil {
		return
	}
	//resp.ResponseDebug()
	println(resp.Text())
}

func GetQrTicket() interface{} {
	url := "https://qr.m.jd.com/check"
	header := requests.Header{
		"user-agent": userAgent,
		"Referer":    "https://passport.jd.com/new/login.aspx",
	}

	//todo
	wlfstk_smdl := "wlfstk_smdl get from cookie"

	rand.Seed(time.Now().Unix())
	param := requests.Params{
		"appid":    "133",
		"callback": genCallback(),
		"token":    wlfstk_smdl,
		"_":        genTime(),
	}

	resp, err := requests.Get(url, header, param)
	if err != nil {
		return ""
	}
	var json map[string]interface{}
	resp.Json(&json)
	return json["ticket"]
}

func ValidQRTicket(ticket string) bool {
	url := "https://passport.jd.com/uc/qrCodeTicketValidation"
	header := requests.Header{
		"user-agent": userAgent,
		"Referer":    "https://passport.jd.com/uc/login?ltype=logout",
	}

	param := requests.Params{
		"t": ticket,
	}

	resp, err := requests.Get(url, header, param)
	if err != nil {
		return false
	}
	var json map[string]interface{}
	resp.Json(&json)
	return json["returnCode"] == 0
}

func GetUserInfo() interface{} {
	url := "https://passport.jd.com/user/petName/getUserInfoForMiniJd.action"
	header := requests.Header{
		"user-agent": userAgent,
		"Referer":    "https://order.jd.com/center/list.action",
	}

	param := requests.Params{
		"callback": genCallback(),
		"_":        genTime(),
	}

	resp, err := requests.Get(url, header, param)
	if err != nil {
		return false
	}
	var json map[string]interface{}
	resp.Json(&json)
	return json["nickName"]
}

func GetKillInitInfo(skuId string, num int) interface{} {
	url := "https://marathon.jd.com/seckillnew/orderService/pc/init.action"
	header := requests.Header{
		"user-agent": userAgent,
		"Host":       "marathon.jd.com",
	}

	data := requests.Datas{
		"sku":             skuId,
		"num":             strconv.Itoa(num),
		"isModifyAddress": "false",
	}

	resp, err := requests.Post(url, header, data)
	if err != nil {
		return false
	}
	var json map[string]interface{}
	resp.Json(&json)
	return json
}

func GetKillUrl(skuId string) string {
	url := "https://itemko.jd.com/itemShowBtn"
	header := requests.Header{
		"user-agent": userAgent,
		"Host":       "itemko.jd.com",
		"Referer":    fmt.Sprint("https://item.jd.com/%v.html", skuId),
	}

	param := requests.Params{
		"callback": genCallback(),
		"skuId":    skuId,
		"from":     "pc",
		"_":        genTime(),
	}

	resp, err := requests.Get(url, header, param)
	if err != nil {
		return ""
	}
	var json map[string]interface{}
	resp.Json(&json)

	url = json["url"].(string)
	url = strings.Replace(url, "divide", "marathon", -1)
	killUrl := strings.Replace(url, "user_routing", "captcha.html", -1)
	return "https:" + killUrl
}

func RequestKillUrl(skuId string, killUrl string) {

	url := killUrl
	header := requests.Header{
		"user-agent": userAgent,
		"Host":       "marathon.jd.com",
		"Referer":    fmt.Sprint("https://item.jd.com/%v.html", skuId),
	}

	//todo allow_redirects=False
	requests.Get(url, header)
}

func SubmitOrder(skuId string, num int, rid string) bool {
	url := "https://marathon.jd.com/seckillnew/orderService/pc/submitOrder.action"
	header := requests.Header{
		"User-Agent": userAgent,
		"Host":       "marathon.jd.com",
		"Referer":
		fmt.Sprint("https://marathon.jd.com/seckill/seckill.action?skuId=%v&num=%v&rid=%v", skuId, num, rid),
	}

	param := requests.Params{
		"skuId": skuId,
	}

	data := requests.Datas{}
	//todo 获取起初，初始化的订单data
	resp, err := requests.Post(url, header, param, data)
	if err != nil {
		return false
	}
	var json map[string]interface{}
	resp.Json(&json)
	//todo 怎么判断成功
	//json["success"] == 0
	return true
	//json["orderId"]
	//json["totalMoney"]
	//json["pcUrl"]
}

func ValidCookie(cookie string) {
	url := "https://order.jd.com/center/list.action"
	header := requests.Header{
		"dnt":                       "1",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"upgrade-insecure-requests": "1",
		"user-agent":                userAgent,
		"Cookie":                    cookie,
	}

	param := requests.Params{
		"rid": strconv.Itoa(time.Now().Second() * 1000),
	}

	resp, err := requests.Get(url, header, param)
	if err != nil {
		return
	}
	//resp.ResponseDebug()
	println(resp.Text())
}

func GetServerTime() interface{} {
	url := "https://a.jd.com//ajax/queryServerData.html"
	header := requests.Header{
		"User-Agent": userAgent,
	}

	resp, err := requests.Get(url, header)
	if err != nil {
		return false
	}
	var json map[string]interface{}
	resp.Json(&json)
	return json["serverTime"]
}

func genCallback() string {
	return "jQuery{}" + strconv.Itoa(int(1000000+rand.Int31n(8999999)))
}

func genTime() string {
	return strconv.Itoa(time.Now().Second() * 1000)
}

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"

func genUserAgent() string {
	return userAgent
}
