package util

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/jordan-wright/email"
	"net/http"
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

func ToGobStr(obj []*http.Cookie) string {

	buf := bytes.Buffer{}

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(obj)
	if err != nil {
		fmt.Println("编码失败,错误原因: ", err)
		return "编码失败"
	}

	gobStr := string(buf.Bytes())
	return gobStr
}

func FromGobStr(gobStr string) []*http.Cookie {
	decoder := gob.NewDecoder(bytes.NewReader([]byte(gobStr)))
	var v []*http.Cookie
	decoder.Decode(&v)
	return v
}
