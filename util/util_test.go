package util

import (
	"testing"
)

func TestOpen(t *testing.T) {
	qrfile := "./qrcode.png"
	Open(qrfile)
}
