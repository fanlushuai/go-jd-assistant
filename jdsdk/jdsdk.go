package jdsdk

import (
	"encoding/json"
	"fmt"
	"github.com/asmcos/requests"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var sessionReq = requests.Requests()

func init() {
	sessionReq.Header.Set("User-Agent", userAgent)
	//不允许重定向
	sessionReq.Client.CheckRedirect =
		func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
}

func GetSessionReqCookies() {
	sessionReq.Client.Jar = nil
	json.Marshal(sessionReq.Client.Jar.Cookies())

}

func Proxy(proxy string) {
	sessionReq.Proxy(proxy)
}

func GetLoginPage() {
	url := "https://passport.jd.com/new/login.aspx"
	sessionReq.Get(url)
}

func GetQR(qrPath string) (token string) {
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

	resp, err := sessionReq.Get(url, header, param)
	if err != nil {
		return
	}
	resp.SaveFile(qrPath)

	resp.Cookies()
	for _, c := range resp.Cookies() {
		if c.Name == "wlfstk_smdl" {
			return c.Value
		}
	}

	return
}

func GetQrTicket(token string) string {
	url := "https://qr.m.jd.com/check"
	header := requests.Header{
		"user-agent": userAgent,
		"Referer":    "https://passport.jd.com/new/login.aspx",
	}

	rand.Seed(time.Now().Unix())
	param := requests.Params{
		"appid":    "133",
		"callback": genCallback(),
		//token=wlfstk_smdl := "wlfstk_smdl get from cookie"
		"token": token,
		"_":     genTime(),
	}

	resp, err := sessionReq.Get(url, header, param)
	if err != nil {
		return ""
	}
	type Ret struct {
		Ticket string
	}
	var r Ret
	json.Unmarshal([]byte(getJsonStr(resp.Text())), &r)
	return r.Ticket
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

	resp, err := sessionReq.Get(url, header, param)
	if err != nil {
		return false
	}
	type Ret struct {
		ReturnCode int
		Url        string
	}
	var r Ret
	resp.Json(&r)
	return r.ReturnCode == 0
}

func GetUserInfo() string {
	url := "https://passport.jd.com/user/petName/getUserInfoForMiniJd.action"
	header := requests.Header{
		"user-agent": userAgent,
		"Referer":    "https://order.jd.com/center/list.action",
	}

	param := requests.Params{
		"callback": genCallback(),
		"_":        genTime(),
	}

	resp, err := sessionReq.Get(url, header, param)
	if err != nil {
		return "获取nick称失败"
	}

	type Ret struct {
		NickName string
	}
	var r Ret
	json.Unmarshal([]byte(getJsonStr(resp.Text())), &r)
	return r.NickName
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

	resp, err := sessionReq.Post(url, header, data)
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
		"Referer":    fmt.Sprintf("https://item.jd.com/%v.html", skuId),
	}

	param := requests.Params{
		"callback": genCallback(),
		"skuId":    skuId,
		"from":     "pc",
		"_":        genTime(),
	}

	resp, err := sessionReq.Get(url, header, param)
	if err != nil {
		return ""
	}
	var json map[string]interface{}
	resp.Json(&json)

	url = json["url"].(string)
	if len(url) < 1 {
		return url
	}
	url = strings.Replace(url, "divide", "marathon", -1)
	killUrl := strings.Replace(url, "user_routing", "captcha.html", -1)
	return "https:" + killUrl
}

func RequestKillUrl(skuId string, killUrl string) {
	url := killUrl
	header := requests.Header{
		"user-agent": userAgent,
		"Host":       "marathon.jd.com",
		"Referer":    fmt.Sprintf("https://item.jd.com/%v.html", skuId),
	}

	requests.Get(url, header)
}

func SubmitOrder(skuId string, num int, rid string) bool {
	url := "https://marathon.jd.com/seckillnew/orderService/pc/submitOrder.action"
	header := requests.Header{
		"User-Agent": userAgent,
		"Host":       "marathon.jd.com",
		"Referer":    fmt.Sprintf("https://marathon.jd.com/seckill/seckill.action?skuId=%v&num=%v&rid=%v", skuId, num, rid),
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

func ValidCookie(cookie string) bool {
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
		"rid": genTime(),
	}

	resp, err := requests.Get(url, header, param)
	if err != nil {
		return false
	}

	return resp.R.StatusCode == 200
}

func GetServerTime() int {
	url := "https://a.jd.com//ajax/queryServerData.html"
	header := requests.Header{
		"User-Agent": userAgent,
	}

	resp, err := requests.Get(url, header)
	if err != nil {
		return -1
	}
	type Ret struct {
		ServerTime int
	}
	var r Ret
	resp.Json(&r)
	return r.ServerTime
}

func genCallback() string {
	return "jQuery" + strconv.Itoa(int(1000000+rand.Int31n(8999999)))
}

func genTime() string {
	return strconv.Itoa(time.Now().Second() * 1000)
}

func getJsonStr(text string) string {
	fromIndex := strings.Index(text, "{")
	endIndex := strings.LastIndex(text, "}")
	return text[fromIndex : endIndex+1]
}

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"

func genUserAgent() string {
	return userAgent
}
