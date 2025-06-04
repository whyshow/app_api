#!/bin/bash
# 编译Linux可执行文件
GOOS=linux GOARCH=amd64 go build -o app_api

# 复制配置文件（如果有）
cp config.ini app_api.ini

echo "编译完成，生成文件: app_api"
