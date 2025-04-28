#!/bin/bash
# 参数配置区（您只需要改这里）
SOURCE_PATH="../main.go"  # 您的Go源码绝对/相对路径
BIN_NAME="AndroidMemoryTools"                     # 生成的二进制文件名
DEVICE_PATH="/data/local/tmp"             # 设备存储路径

# 自动检测设备架构
DEVICE_ABI=$(adb shell getprop ro.product.cpu.abi | tr -d '\r\n')
case $DEVICE_ABI in
    *arm64*) GOARCH="arm64" ;;
    *armv7*) GOARCH="arm" ;;
    *x86_64*) GOARCH="amd64" ;;
    *x86*) GOARCH="386" ;;
    *) echo "不支持的架构: $DEVICE_ABI"; exit 1 ;;
esac

# 编译指定源码
echo "正在编译 $SOURCE_PATH -> $BIN_NAME (GOARCH=$GOARCH)"
CGO_ENABLED=0 GOOS=android GOARCH=$GOARCH \
go build -trimpath -ldflags="-s -w" -o "$BIN_NAME" "$SOURCE_PATH"

# 推送到设备
adb push "$BIN_NAME" "$DEVICE_PATH/"
adb shell chmod +x "$DEVICE_PATH/$BIN_NAME"
echo "程序已推送至设备:"
adb shell ls -lh "$DEVICE_PATH/$BIN_NAME"
adb shell su -c "$DEVICE_PATH/$BIN_NAME"