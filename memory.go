package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func getPid(packageName string) int {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("ps -A | grep %s | awk '{print $2}'", packageName))
	output, err := cmd.CombinedOutput()
	if err != nil {
		// 修正3：grep在找不到进程时会返回错误，这是正常情况
		if strings.Contains(err.Error(), "exit status 1") {
			return 0
		}
		fmt.Printf("命令执行错误: %v\n", err)
		return 0
	}
	// 去除空白字符后解析PID
	pidStr := strings.TrimSpace(string(output))
	if pidStr == "" {
		return 0
	}
	var pid int
	_, err = fmt.Sscanf(pidStr, "%d", &pid)
	if err != nil {
		return 0
	}
	return pid
}
