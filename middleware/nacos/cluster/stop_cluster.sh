#!/bin/bash

# 定义要检查的端口列表
PORTS=(8090 8091 8092 8093 8094)

# 循环处理每个端口
for port in "${PORTS[@]}"; do
    # 查找监听指定端口的进程ID
    pid=$(lsof -ti :$port)

    if [ -z "$pid" ]; then
        echo "端口 $port 没有找到运行的进程"
    else
        echo "找到端口 $port 的进程: $pid"
        # 杀死进程
        kill -9 $pid
        # 验证进程是否被杀死
        if ps -p $pid > /dev/null; then
            echo "警告: 进程 $pid (端口 $port) 未能被终止"
        else
            echo "成功终止进程 $pid (端口 $port)"
        fi
    fi
done

echo "端口检查完成"