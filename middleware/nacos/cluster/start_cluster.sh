#!/bin/bash

# 定义默认配置路径和日志路径
DEFAULT_CFG_PATH="/Users/a123/Downloads/go/learn-go/learn-go/middleware/nacos/nacos.json"
DEFAULT_LOG_PATH="/Users/a123/Downloads/go/learn-go/learn-go/middleware/nacos/logs"

# 检查并创建日志目录
if [ ! -d "$DEFAULT_LOG_PATH" ]; then
    mkdir -p "$DEFAULT_LOG_PATH"
fi

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case "$1" in
        --cfgPath)
            CFG_PATH="$2"
            shift 2
            ;;
        --logPath)
            LOG_PATH="$2"
            shift 2
            ;;
        *)
            echo "未知参数: $1"
            exit 1
            ;;
    esac
done

# 使用用户提供的路径或默认路径
CFG_PATH=${CFG_PATH:-$DEFAULT_CFG_PATH}
LOG_PATH=${LOG_PATH:-$DEFAULT_LOG_PATH}

echo "使用配置路径: $CFG_PATH"
echo "使用日志路径: $LOG_PATH"

# 启动所有节点
./node/node --cfgPath="$CFG_PATH" --logPath="$LOG_PATH/node1.log" --serviceName="gse" --groupName="cmdb" --port=8090 &
./node/node --cfgPath="$CFG_PATH" --logPath="$LOG_PATH/node2.log" --serviceName="gse" --groupName="cmdb" --port=8091 &
./node/node --cfgPath="$CFG_PATH" --logPath="$LOG_PATH/node3.log" --serviceName="gse" --groupName="cmdb" --port=8092 &
./node/node --cfgPath="$CFG_PATH" --logPath="$LOG_PATH/node4.log" --serviceName="gse" --groupName="cmdb" --port=8093 &
./node/node --cfgPath="$CFG_PATH" --logPath="$LOG_PATH/node5.log" --serviceName="gse" --groupName="cmdb" --port=8094 &


echo "所有节点已启动"
echo "使用 'ps aux | grep node[1-3]' 查看进程"
echo "使用 'pkill -f node[1-3]' 停止所有节点"