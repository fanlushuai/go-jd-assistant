package util

import (
	"testing"
)

func TestOpen(t *testing.T) {
	qrfile := "../jdsdk/qrcode.png"
	Open(qrfile)
}

func TestSendEmail(t *testing.T) {
	SendEmailOutlook("../jdsdk/qrcode.png", "fanlushuai@outlook.com")
}
