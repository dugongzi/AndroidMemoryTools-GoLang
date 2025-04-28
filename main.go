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
