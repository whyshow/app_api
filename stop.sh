#!/bin/bash
appName="app_api"
pid=`ps aux|grep $appName|grep -v grep|awk '{print $2}'`
kill -9 $pid