/*
	实时打印日志
*/
package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func main() {
	command := exec.Command("bash", "./printtime.sh")
	stdout, err := command.StdoutPipe()
	if err != nil {
		fmt.Println("不能标准输出")
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		fmt.Println("不能标准错误输出")
	}

	go logsync(stdout)
	go logsync(stderr)
	err = command.Start()
	if err != nil {
		fmt.Println("不能正确开始命令")
	}
	err = command.Wait()
	if err != nil {
		fmt.Println("不能正等待命令执行")
	}
}

func logsync(reader io.ReadCloser) (err error) {
	cache := ""
	buf := make([]byte, 4096, 4096)
	for {
		num, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF || strings.Contains(err.Error(), "closed") {
				err = nil
			}
			return err
		}
		if num > 0 {
			oByte := buf[:num]
			oSlice := strings.Split(string(oByte), "\n")
			line := strings.Join(oSlice[:len(oSlice)-1], "\n")
			fmt.Printf("%s%s\n", cache, line)
			cache = oSlice[len(oSlice)-1]
		}
	}
	return nil
}
