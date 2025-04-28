#!/bin/bash
# 强制编译当前目录所有关联文件的版本

# 参数配置
SOURCE_DIR="../"  # 源码目录（根据实际情况调整）
BIN_NAME="AndroidMemoryTools"
DEVICE_PATH="/data/local/tmp"

# 获取设备架构
DEVICE_ABI=$(adb shell getprop ro.product.cpu.abi | tr -d '\r\n')
case "$DEVICE_ABI" in
    *arm64*) GOARCH="arm64" ;;
    *armeabi*|*armv7*) GOARCH="arm" ;;
    *x86_64*) GOARCH="amd64" ;;
    *x86*) GOARCH="386" ;;
    *) echo "不支持的架构: $DEVICE_ABI"; exit 1 ;;
esac

# 编译命令（关键修改！）
echo "▸ 编译 $SOURCE_DIR 下的所有关联文件 (GOARCH=$GOARCH)"
cd "$SOURCE_DIR" && \
CGO_ENABLED=0 GOOS=android GOARCH=$GOARCH \
go build -trimpath -ldflags="-s -w" -o "$BIN_NAME" . || {
    echo "编译失败！常见原因："
    echo "1. 函数名大小写错误（Go区分大小写）"
    echo "2. 文件头package声明不一致"
    echo "3. 函数未首字母大写（无法跨包调用）"
    exit 1
}

# 推送执行
adb push "$BIN_NAME" "$DEVICE_PATH/" && \
adb shell chmod 777 "$DEVICE_PATH/$BIN_NAME" && \
echo "▸ 执行结果:"
adb shell "su -c 'cd $DEVICE_PATH && ./$BIN_NAME'"