package util

import (
	"fmt"
	"testing"
	"time"
)

func TestOpen(t *testing.T) {
	qrfile := "../jdsdk/qrcode.png"
	Open(qrfile)
}

func TestSendEmail(t *testing.T) {
	SendEmailOutlook("../jdsdk/qrcode.png", "fanlushuai@outlook.com")
}

func TestTime(t *testing.T) {
	timeMs := "2009-01-01T01:02:01.1+08:00"
	//layout:="2006-01-02 03:04:05.000"
	layout := time.RFC3339
	parse, err := time.Parse(layout, timeMs)
	if err == nil {

	}
	fmt.Println(parse.UnixNano())
	fmt.Println(parse.UnixNano() / 1000000)

}
