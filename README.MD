# 开发环境
系统版本：Ubuntu 24.04.2 LTS
Go版本：go1.21.3
# 方法介绍
参照这个开源项目：https://github.com/MuxueStudio/Rust-AndroidMemoryTools/blob/master/README.MD
# 获取帮助:
### 沐雪社区:
https://www.muxue.pro
### 哔哩哔哩:
https://space.bilibili.com/3546693165386233 
### 抖音:
https://v.douyin.com/I-xKYusrX5M/
# 示例代码
```go
package main

import "fmt"

func main() {
	pid := getPid("com.cyhxzhdzy.kz")
	h := getLibHead(pid, "libil2cpp.so")
	addr, err := readPoint(pid, h["Cd"]+0x4EF18, []int64{0x5C, 0xA0, 0xC})
	if err != nil {
		fmt.Println("readPoint:", err)
		return
	}
	val, err := readVal[DWORD](pid, addr)
	if err != nil {
		fmt.Println("readVal:", err)
	}
	err = writeVal[DWORD](pid, addr, 999)
	if err != nil {
		fmt.Println("writeVal:", err)
	}
	fmt.Println(val)
}

```