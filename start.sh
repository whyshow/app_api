#!/bin/sh
time=$(date "+%Y%m%d")
appName="app_api"
chmod 777 $appName
#创建文件夹

if [ -d "./logs/$time" ];then
      nohup ./$appName >>logs/$time/$time.log 2>&1 &
    else
        mkdir -p ./logs/$time/
        nohup ./$appName >>logs/$time/$time.log 2>&1 &
    fi