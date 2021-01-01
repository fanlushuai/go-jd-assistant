package util

import (
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
