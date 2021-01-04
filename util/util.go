package util

import (
	"github.com/jordan-wright/email"
	"net/textproto"
	"os/exec"
	"runtime"
)

func Open(uri string) error {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd.exe", "/c", "start "+uri)
		return cmd.Start()
	}
	//其他平台，自己支持，本人只有windows环境
	var PicOpenImpErr error
	return PicOpenImpErr
}

func SendEmailOutlook(filePath string, to string) {
	e := &email.Email{
		To:      []string{to},
		From:    "JD ROBOT <fanlushuai@outlook.com>",
		Subject: "QR LOGIN by jd robot",
		Headers: textproto.MIMEHeader{},
	}
	e.AttachFile(filePath)
	auth := OutLookEmailLoginAuth("fanlushuai@outlook.com", "password !!!!")
	e.Send("smtp.office365.com:587", auth)
}
