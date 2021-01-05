package jdsdk

import (
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
