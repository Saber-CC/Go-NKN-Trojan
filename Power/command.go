package Power

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"golang.org/x/text/encoding/simplifiedchinese"
)

func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}

func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			i++
			for j := 0; j < num-1; j++ {
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			return false
		}
	}
	return true
}

func NewCommand() Commander {
	var cmd Commander
	switch runtime.GOOS {
	case "windows":
		cmd = NewWindowsCommand()
	case "linux":
		cmd = NewLinuxCommand()
	}
	return cmd
}

const (
	GBK  string = "GBK"
	UTF8 string = "UTF8"
)

func GetStrCoding(data []byte) string {
	if isUtf8(data) == true {
		return UTF8
	} else {
		return GBK
	}
}

type Commander interface {
	Exec(args ...string) (int, string, error)
}

func ConvertFormat(byte []byte, format string) string {
	var str string
	switch format {
	case GBK:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}

type WindowsCommand struct {
}

func NewWindowsCommand() *WindowsCommand {
	return &WindowsCommand{}
}

func (lc *WindowsCommand) Exec(args ...string) (int, string, error) {
	args = append([]string{"/c"}, args...)
	cmd := exec.Command("cmd", args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{}

	outpip, err := cmd.StdoutPipe()
	defer outpip.Close()

	if err != nil {
		return 0, "", err
	}

	err = cmd.Start()
	if err != nil {
		return 0, "", err
	}

	out, err := ioutil.ReadAll(outpip)
	if err != nil {
		return 0, "", err
	}
	cmdout := ConvertFormat(out, GetStrCoding(out))

	return cmd.Process.Pid, cmdout, nil
}

type LinuxCommand struct {
}

func NewLinuxCommand() *LinuxCommand {
	return &LinuxCommand{}
}

func (lc *LinuxCommand) Exec(args ...string) (int, string, error) {
	args = append([]string{"-c"}, args...)
	cmd := exec.Command(os.Getenv("SHELL"), args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{}

	outpip, err := cmd.StdoutPipe()
	defer outpip.Close()

	if err != nil {
		return 0, "", err
	}

	err = cmd.Start()
	if err != nil {
		return 0, "", err
	}

	out, err := ioutil.ReadAll(outpip)
	if err != nil {
		return 0, "", err
	}
	cmdout := ConvertFormat(out, GetStrCoding(out))

	return cmd.Process.Pid, cmdout, nil
}
