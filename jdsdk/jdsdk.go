package jdsdk

import (
	"encoding/json"
	"fmt"
	"github.com/asmcos/requests"
	"github.com/vdobler/ht/cookiejar"
	"go-jd-assistant/util"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var sessionReq = requests.Requests()
var jarP *cookiejar.Jar

func init() {
	sessionReq.Header.Set("User-Agent", userAgent)
	//不允许重定向
	sessionReq.Client.CheckRedirect =
		func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	//使用一个可以导出导入cookiejar的实现	github.com/vdobler/ht/cookiejar
	jar, _ := cookiejar.New(nil)
	jarP = jar
	sessionReq.Client.Jar = jar
}

//record log by proxy
func Proxy(proxy string) {
	sessionReq.Proxy(proxy)
}

func SaveCookies(filePath string) {
	util.SaveCookiesFromJar(jarP, filePath)
}

func ReLoadCookies(filePath string) {
	jar := util.LoadCookies(filePath)
	jarP = jar
	sessionReq.Client.Jar = jar
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

func GetKillInitInfo(skuId string, num string) InitData {
	url := "https://marathon.jd.com/seckillnew/orderService/pc/init.action"
	header := requests.Header{
		"user-agent": userAgent,
		"Host":       "marathon.jd.com",
	}

	data := requests.Datas{
		"sku":             skuId,
		"num":             num,
		"isModifyAddress": "false",
	}

	resp, err := sessionReq.Post(url, header, data)
	if err != nil {
		return InitData{}
	}

	var initdata InitData
	resp.Json(&initdata)
	return initdata
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

	sessionReq.Get(url, header)
}

func SubmitOrder(skuId string, num string, datas *map[string]string) bool {
	url := "https://marathon.jd.com/seckillnew/orderService/pc/submitOrder.action"
	//todo ？ 这个rid是Referer的链接。也就是说，不知道这个重要不，是否需要一个真正的值。还是按照格式来一个就行
	rid := genTime()
	header := requests.Header{
		"User-Agent": userAgent,
		"Host":       "marathon.jd.com",
		"Referer":    fmt.Sprintf("https://marathon.jd.com/seckill/seckill.action?skuId=%v&num=%v&rid=%v", skuId, num, rid),
	}

	param := requests.Params{
		"skuId": skuId,
	}

	resp, err := sessionReq.Post(url, header, param, &datas)
	if err != nil {
		return false
	}

	type Ret struct {
		success      bool //todo 类型等待验证
		errorMessage string
		orderId      string
		resultCode   string
	}

	var r Ret

	resp.Json(&r)

	fmt.Println(r)

	if r.success {
		return true
	}

	return false
}

func ValidCookie() bool {
	url := "https://order.jd.com/center/list.action"
	header := requests.Header{
		"dnt":                       "1",
		"sec-fetch-dest":            "document",
		"sec-fetch-mode":            "navigate",
		"sec-fetch-site":            "none",
		"upgrade-insecure-requests": "1",
		"user-agent":                userAgent,
	}

	param := requests.Params{
		"rid": genTime(),
	}

	resp, err := sessionReq.Get(url, header, param)
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
		ServerTime int //1609878734768
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

type InitData struct {
	addressList []address
	token       string
	invoiceInfo invoiceInfo
}

type invoiceInfo struct {
	invoiceTitle       string
	invoiceContentType string
	invoicePhone       string
	invoicePhoneKey    string
}

type address struct {
	id            string
	name          string
	provinceId    string
	cityId        string
	countyId      string
	townId        string
	addressDetail string
	email         string
}

func BuildSubmitOrderPostData(pw string, fp string, eid string, skuid string, num string, initData *InitData) *map[string]string {

	var data SubmitOrderPostData

	{
		data.num = num
		data.skuId = skuid
		data.fp = fp
		data.eid = eid
		data.password = pw
	}

	{
		defaultAddress := initData.addressList[0]
		invInfo := initData.invoiceInfo

		data.token = initData.token

		data.addressDetail = defaultAddress.addressDetail
		data.addressId = defaultAddress.id
		data.countyId = defaultAddress.countyId
		data.cityId = defaultAddress.cityId
		data.provinceId = defaultAddress.provinceId
		data.townId = defaultAddress.townId
		data.name = defaultAddress.name
		data.email = defaultAddress.email

		if invInfo != (invoiceInfo{}) {
			data.invoice = "true"
			data.invoiceTitle = invInfo.invoiceTitle
			data.invoiceContent = invInfo.invoiceContentType
			data.invoicePhoneKey = invInfo.invoicePhoneKey
			data.invoicePhone = invInfo.invoicePhone
		}
	}

	data.pru = ""
	data.phone = ""
	data.overseas = "0"
	data.areaCode = ""
	data.paymentType = "4"
	data.codTimeType = "3"
	data.invoiceEmail = ""
	data.invoiceTaxpayerNO = ""
	data.invoiceCompanyName = ""
	data.postCode = ""
	data.isModifyAddress = "false"
	data.yuShou = "" //todo

	return toMapstringstring(data)
}

func toMapstringstring(data SubmitOrderPostData) *map[string]string {
	m := make(map[string]string)
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	for i := 0; i < v.NumField(); i++ {
		m[t.Field(i).Name] = v.Field(i).Interface().(string)
	}
	return &m
}

type SubmitOrderPostData struct {
	skuId              string
	num                string
	addressId          string
	yuShou             string
	isModifyAddress    string
	name               string
	provinceId         string
	cityId             string
	countyId           string
	townId             string
	addressDetail      string
	mobile             string
	mobileKey          string
	email              string
	postCode           string
	invoiceTitle       string
	invoiceCompanyName string
	invoiceContent     string
	invoiceTaxpayerNO  string
	invoiceEmail       string
	invoicePhone       string
	invoicePhoneKey    string
	invoice            string
	password           string
	codTimeType        string
	paymentType        string
	areaCode           string
	overseas           string
	phone              string
	eid                string
	fp                 string
	token              string
	pru                string
}
