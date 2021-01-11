package jdsdk

import (
	"crypto/tls"
	"encoding/json"
	"errors"
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
	sessionReq.Client.Timeout = time.Duration(10) * time.Second
	sessionReq.Client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

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
		"Referer": "https://passport.jd.com/new/login.aspx",
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
		"Referer": "https://passport.jd.com/new/login.aspx",
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
		"Referer": "https://passport.jd.com/uc/login?ltype=logout",
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
		"Referer": "https://order.jd.com/center/list.action",
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

func GetKillInitInfo(skuId string, num string) (initData InitData, err error) {
	url := "https://marathon.jd.com/seckillnew/orderService/pc/init.action"
	header := requests.Header{
		"Host": "marathon.jd.com",
	}

	data := requests.Datas{
		"sku":             skuId,
		"num":             num,
		"isModifyAddress": "false",
	}

	//这个头按说不应该添加，奈何开源库的bug。没有对这个进行处理 https://github.com/asmcos/requests/issues/23
	header["Content-Type"] = "application/x-www-form-urlencoded"
	resp, err := sessionReq.Post(url, header, data)
	if err != nil {
		fmt.Println("fuck initinfo 获取失败了。好好思考一下")
		return InitData{}, errors.New("请求错误")
	}

	var initdata InitData
	resp.Json(&initdata)
	if len(initData.AddressList) == 0 {
		return InitData{}, errors.New("响应错误，估计被频率限制了")
	}

	return initdata, nil
}

func GetKillUrl(skuId string) string {
	url := "https://itemko.jd.com/itemShowBtn"
	header := requests.Header{
		"Host":    "itemko.jd.com",
		"Referer": fmt.Sprintf("https://item.jd.com/%v.html", skuId),
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
	type Ret struct {
		Url string
	}

	var r Ret
	json.Unmarshal([]byte(getJsonStr(resp.Text())), &r)

	if len(r.Url) > 3 {
		r.Url = strings.Replace(r.Url, "divide", "marathon", -1)
		killUrl := strings.Replace(r.Url, "user_routing", "captcha.html", -1)
		return "https:" + killUrl
	}

	return ""
}

func RequestKillUrl(skuId string, killUrl string) {
	url := killUrl
	header := requests.Header{
		"Host":    "marathon.jd.com",
		"Referer": fmt.Sprintf("https://item.jd.com/%v.html", skuId),
	}

	sessionReq.Get(url, header)
}

func SubmitOrder(skuId string, num string, datas *map[string]string) bool {
	url := "https://marathon.jd.com/seckillnew/orderService/pc/submitOrder.action"
	//todo ？ 这个rid是Referer的链接。也就是说，不知道这个重要不，是否需要一个真正的值。还是按照格式来一个就行
	rid := genTime()
	header := requests.Header{
		"Host":    "marathon.jd.com",
		"Referer": fmt.Sprintf("https://marathon.jd.com/seckill/seckill.action?skuId=%v&num=%v&rid=%v", skuId, num, rid),
	}

	param := requests.Params{
		"skuId": skuId,
	}

	header["Content-Type"] = "application/x-www-form-urlencoded"
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

	resp, err := requests.Get(url)
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
	AddressList  []address
	InvoiceInfo  invoiceInfo
	SeckillSkuVO seckillSkuVO
	Token        string
}

type address struct {
	Id            int
	Name          string
	ProvinceId    int
	CityId        int
	CountyId      int
	TownId        int
	AddressDetail string
	Email         string
}

type invoiceInfo struct {
	InvoiceTitle       int
	InvoiceContentType int
	InvoicePhone       string
	InvoicePhoneKey    string
}

type seckillSkuVO struct {
	ExtMap extMap
}

type extMap struct {
	YuShou string
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
		defaultAddress := initData.AddressList[0]
		invInfo := initData.InvoiceInfo

		data.token = initData.Token

		data.addressDetail = defaultAddress.AddressDetail
		data.addressId = strconv.Itoa(defaultAddress.Id)
		data.countyId = strconv.Itoa(defaultAddress.CountyId)
		data.cityId = strconv.Itoa(defaultAddress.CityId)
		data.provinceId = strconv.Itoa(defaultAddress.ProvinceId)
		data.townId = strconv.Itoa(defaultAddress.TownId)
		data.name = defaultAddress.Name
		data.email = defaultAddress.Email

		if invInfo != (invoiceInfo{}) {
			data.invoice = "true"
			data.invoiceTitle = strconv.Itoa(invInfo.InvoiceTitle)
			data.invoiceContent = strconv.Itoa(invInfo.InvoiceContentType)
			data.invoicePhoneKey = invInfo.InvoicePhoneKey
			data.invoicePhone = invInfo.InvoicePhone
		}

		//爱了。go语言对于null。完全不用判断。
		if initData.SeckillSkuVO.ExtMap.YuShou == "" {
			data.yuShou = "0"
		} else {
			data.yuShou = initData.SeckillSkuVO.ExtMap.YuShou
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

	return toMapstringstring(data)
}

func toMapstringstring(data SubmitOrderPostData) *map[string]string {
	m := make(map[string]string)
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	for i := 0; i < v.NumField(); i++ {
		m[t.Field(i).Name] = v.Field(i).String()
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
