package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

type (
	DWORD  int32
	FLOAT  float32
	WORD   int16
	DOUBLE float64
	QWORD  int64
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

func getLibHead(pid int, soName string) map[string]int64 {
	m := make(map[string]int64)
	command := exec.Command("sh", "-c", fmt.Sprintf("cat /proc/%d/maps | grep -F '%s'", pid, soName))
	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Printf("未成功获取so文件: %v\n", err)
	}
	for line := range strings.Lines(string(output)) {
		addrStr := strings.Split(line, "-")[0]
		if strings.Contains(line, "rw-p") {
			m["Cd"], _ = strconv.ParseInt(addrStr, 16, 64)
		} else {
			m["Xa"], _ = strconv.ParseInt(addrStr, 16, 64)
		}
	}
	return m
}

func readVal[T any](pid int, addr int64) (T, error) {
	var zero T
	path := fmt.Sprintf("/proc/%d/mem", pid)
	fd, err := syscall.Open(path, syscall.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("readVal error: %v\n", err)
	}
	defer func(fd int) {
		err := syscall.Close(fd)
		if err != nil {
			fmt.Printf("close file error: %v\n", err)
		}
	}(fd)
	if fd == -1 {
		err := fmt.Errorf("open file error")
		return zero, err
	}
	size := unsafe.Sizeof(zero)
	buf := make([]byte, size)
	// 2. 调用 pread 读取数据（等效于 pread64）
	n, err := syscall.Pread(fd, buf, addr)
	if err != nil {
		return zero, fmt.Errorf("读取失败: %v", err)
	}
	if n != int(size) {
		return zero, fmt.Errorf("读取字节数不足: %d < %d", n, size)
	}
	return *(*T)(unsafe.Pointer(&buf[0])), nil
}
func rpoint(pid int, address int64) (int64, error) {
	// 打开 /proc/<pid>/mem
	memPath := fmt.Sprintf("/proc/%d/mem", pid)
	fd, err := syscall.Open(memPath, syscall.O_RDONLY, 0)
	if err != nil {
		fmt.Printf("failed to open %s: %v", memPath, err)
	}
	defer func(fd int) {
		err := syscall.Close(fd)
		if err != nil {
			fmt.Printf("close file error: %v\n", err)
		}
	}(fd)

	// 读取4字节u32值
	var val uint32
	_, err = syscall.Pread(fd, (*(*[4]byte)(unsafe.Pointer(&val)))[:], address)
	if err != nil {
		fmt.Println(err)
	}
	return int64(val), nil
}

func readPoint(pid int, address int64, offsets []int64) (int64, error) {
	p1, err := rpoint(pid, address)
	if err != nil {
		return 0, err
	}

	size := len(offsets)
	for i := 0; i < size-1; i++ {
		p1, err = rpoint(pid, p1+offsets[i])
		if err != nil {
			return 0, err
		}
	}

	return p1 + offsets[size-1], nil
}

func writeVal[T any](pid int, address int64, val T) error {
	memPath := fmt.Sprintf("/proc/%d/mem", pid)
	fd, err := syscall.Open(memPath, syscall.O_RDWR, 0)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", memPath, err)
	}
	defer func(fd int) {
		err := syscall.Close(fd)
		if err != nil {
			fmt.Printf("failed to close fd %d: %v\n", fd, err)
		}
	}(fd)

	// 获取类型大小
	size := unsafe.Sizeof(val)

	// 写入数据
	_, err = syscall.Pwrite(fd, (*(*[1 << 30]byte)(unsafe.Pointer(&val)))[:size], address)
	if err != nil {
		return fmt.Errorf("write failed at 0x%x: %v", address, err)
	}

	return nil
}
