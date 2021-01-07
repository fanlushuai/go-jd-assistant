package jdsdk

import (
	"encoding/json"
	"fmt"
	"go-jd-assistant/util"
	"testing"
	"time"
)

func init() {
	//Proxy("http://localhost:8888")
}

func TestGetLoginPage(t *testing.T) {
	GetLoginPage()
}

func TestGetQR(t *testing.T) {
	GetLoginPage()
	qrfile := "./qrcode.png"
	token := GetQR(qrfile)
	fmt.Println(token)
}

func TestGetQrTicket(t *testing.T) {
	GetLoginPage()
	qrfile := "./qrcode.png"
	token := GetQR(qrfile)
	util.Open(qrfile)
	fmt.Println(token)
	time.Sleep(8 * time.Second)
	ticket := GetQrTicket(token)
	fmt.Println("ticket->" + ticket)
}

func TestValidQRTicket(t *testing.T) {
	GetLoginPage()
	qrfile := "./qrcode.png"
	token := GetQR(qrfile)
	util.Open(qrfile)
	fmt.Println(token)
	time.Sleep(8 * time.Second)
	ticket := GetQrTicket(token)
	fmt.Println("ticket->" + ticket)
	ok := ValidQRTicket(ticket)
	fmt.Print("ret-->")
	fmt.Println(ok)
}

func TestGetUserInfo(t *testing.T) {
	GetLoginPage()
	qrfile := "./qrcode.png"
	token := GetQR(qrfile)
	util.Open(qrfile)
	fmt.Println(token)
	time.Sleep(8 * time.Second)
	ticket := GetQrTicket(token)
	fmt.Println("ticket->" + ticket)
	ok := ValidQRTicket(ticket)
	fmt.Print("ret-->")
	fmt.Println(ok)
	nickName := GetUserInfo()
	fmt.Println(nickName)
}

func TestSaveCookies(t *testing.T) {
	GetLoginPage()
	qrfile := "./qrcode.png"
	token := GetQR(qrfile)
	util.Open(qrfile)
	fmt.Println(token)
	time.Sleep(8 * time.Second)
	ticket := GetQrTicket(token)
	fmt.Println("ticket->" + ticket)
	ok := ValidQRTicket(ticket)
	fmt.Print("ret-->")
	fmt.Println(ok)

	SaveCookies("./cookie.cookies")
}

func TestLoadCookies(t *testing.T) {
	ReLoadCookies("./cookie.cookies")
	nickName := GetUserInfo()
	fmt.Println(nickName)
}

func TestValidCookie(t *testing.T) {
	GetLoginPage()
	qrfile := "./qrcode.png"
	token := GetQR(qrfile)
	util.Open(qrfile)
	fmt.Println(token)
	time.Sleep(8 * time.Second)
	ticket := GetQrTicket(token)
	fmt.Println("ticket->" + ticket)
	ok := ValidQRTicket(ticket)
	fmt.Print("ret-->")
	fmt.Println(ok)
	nickName := GetUserInfo()
	fmt.Println(nickName)

	validOk := ValidCookie()
	fmt.Println("validOk")
	fmt.Println(validOk)
}

func TestGetKillInitInfo(t *testing.T) {
	ReLoadCookies("../my.cookies")
	GetKillInitInfo("100012043978", "1")
}

func TestSubmitOrder(t *testing.T) {
	//ReLoadCookies("../my.cookies")
	//initdata:=GetKillInitInfo("100012043978",1)
	////
	////submitOrderPostDatas:=BuildSubmitOrderPostData(,,&initdata)
	////datasMapStringString:structs.Map(submitOrderPostDatas)
	////SubmitOrder(,,,datasMapStringString)
}

func TestGetServerTime(t *testing.T) {
	fmt.Println(GetServerTime())
}

func TestPareJson(t *testing.T) {
	str := "{\"addressList\":[{\"addressDetail\":\"沙*****号楼\",\"areaCode\":\"86\",\"cityId\":2901,\"cityName\":\"昌平区\",\"countyId\":55561,\"countyName\":\"沙河地区\",\"defaultAddress\":true,\"email\":\"\",\"id\":2357421508,\"mobile\":\"153**79\",\"mobileKey\":\"3b0e11dfdsfdf37c7bf5426e47e5\",\"name\":\"小天才\",\"overseas\":0,\"phone\":\"\",\"postCode\":\"\",\"provinceId\":1,\"provinceName\":\"北京\",\"townId\":0,\"townName\":\"\",\"yuyueAddress\":false}],\"buyNum\":2,\"code\":\"200\",\"freight\":0,\"invoiceInfo\":{\"invoiceCode\":\"\",\"invoiceCompany\":\"\",\"invoiceContentType\":1,\"invoiceEmail\":\"\",\"invoicePhone\":\"153****3379\",\"invoicePhoneKey\":\"3b0e11sdfdsfdfb61037c7bf5426e47e5\",\"invoiceTitle\":4,\"invoiceType\":3},\"paymentTypeList\":[{\"paymentId\":4,\"paymentName\":\"在线支付\"}],\"seckillSkuVO\":{\"color\":\"飞天 53%vol 500ml 贵州茅台酒\",\"extMap\":{\"YuShou\":\"1\",\"is7ToReturn\":\"0\",\"new7ToReturn\":\"8\",\"thwa\":\"0\",\"SoldOversea\":\"0\"},\"height\":0.0,\"jdPrice\":1499.00,\"length\":0.0,\"num\":1,\"rePrice\":0.00,\"size\":\"\",\"skuId\":100012043978,\"skuImgUrl\":\"jfs/t1/97097/12/15694/245806/5e7373e6Ec4d1b0ac/9d8c13728cc2544d.jpg\",\"skuName\":\"飞天 53%vol  500ml 贵州茅台酒（带杯）\",\"skuPrice\":1499.00,\"thirdCategoryId\":0.0,\"venderName\":\"京东自营\",\"venderType\":0,\"weight\":1.120,\"width\":0.0},\"shipmentParam\":{\"shipmentTimeType\":3,\"shipmentTimeTypeName\":\"工作日、双休日与假期均可送货\",\"shipmentType\":65,\"shipmentTypeName\":\"京东配送\"},\"token\":\"c68fcac59fdsffdbdf1c45ff2d\"}"
	var initdata InitData
	json.Unmarshal([]byte(str), &initdata)
	fmt.Println(initdata)

	m := BuildSubmitOrderPostData("pw", "fp", "eid", "skuid", "num", &initdata)
	fmt.Println(m)
}

func TestBuildSubmitOrderPostData(t *testing.T) {
	ReLoadCookies("../my.cookies")
	initdata := GetKillInitInfo("100012043978", "1")
	m := BuildSubmitOrderPostData("pw", "fp", "eid", "skuid", "num", &initdata)
	fmt.Println(m)
}
