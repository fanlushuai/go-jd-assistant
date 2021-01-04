package util

import (
	"fmt"
	"net/http"
	"testing"
)

func TestOpen(t *testing.T) {
	qrfile := "../jdsdk/qrcode.png"
	Open(qrfile)
}

func TestSendEmail(t *testing.T) {
	SendEmailOutlook("../jdsdk/qrcode.png", "fanlushuai@outlook.com")
}

func TestGob(t *testing.T) {
	var a = []*http.Cookie{
		{Name: "ChocolateChip", Value: "tasty"},
		{Name: "First", Value: "Hit"},
		{Name: "Second", Value: "Hit"},
	}
	fmt.Println("编码之前 = ", a)

	s := ToGobStr(a)
	fmt.Println("编码 字符串 = ", s)

	FromGobStr(s)

	fmt.Println("解码 = ", a)
}
