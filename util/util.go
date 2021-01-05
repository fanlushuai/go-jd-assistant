package util

import (
	"encoding/json"
	"github.com/jordan-wright/email"
	"github.com/vdobler/ht/cookiejar"
	"io/ioutil"
	"log"
	"net/textproto"
	"os"
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

func SaveCookiesFromJar(jar *cookiejar.Jar, filename string) error {
	if jar == nil {
		return nil
	}

	cookies := make(map[string]cookiejar.Entry)
	for _, tld := range jar.ETLDsPlus1(nil) {
		for _, cookie := range jar.Entries(tld, nil) {
			id := cookie.ID()
			cookies[id] = cookie
		}
	}
	return saveCookies(cookies, filename)
}

func saveCookies(cookies map[string]cookiejar.Entry, filename string) error {
	b, err := json.MarshalIndent(cookies, "    ", "")
	if err != nil {
		return nil
	}
	return ioutil.WriteFile(filename, b, 0666)
}

func LoadCookies(filename string) *cookiejar.Jar {
	if filename == "" {
		return nil
	}
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("Cannot read cookie file: %s", err)
	}

	cookies := make(map[string]cookiejar.Entry)
	err = json.Unmarshal(buf, &cookies)
	if err != nil {
		log.Panicf("Cannot decode cookie file: %s", err)
	}
	cs := make([]cookiejar.Entry, 0, len(cookies))
	for _, c := range cookies {
		cs = append(cs, c)
	}

	jar, _ := cookiejar.New(nil)
	jar.LoadEntries(cs)
	return jar
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}
