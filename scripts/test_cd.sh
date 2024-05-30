#!/bin/bash

PATH=$PATH:/usr/local/go/bin
go mod tidy

mkdir ./build -p
PWD=`pwd`

# 创建apiserver系统进程文件apiserver.service
APISERVER_CONF_FILE="/etc/systemd/system/apiserver.service"
echo "[Service]" > $APISERVER_CONF_FILE
echo "User=root" >> $APISERVER_CONF_FILE
echo "ExecStart=${PWD}/build/apiserver" >> $APISERVER_CONF_FILE

# 创建shedular系统进程文件scheduler.service
SCHEDULER_CONF_FILE="/etc/systemd/system/scheduler.service"
echo "[Service]" > $SCHEDULER_CONF_FILE
echo "User=root" >> $SCHEDULER_CONF_FILE
echo "ExecStart=${PWD}/build/scheduler" >> $SCHEDULER_CONF_FILE

# 创建kubelet系统进程文件kubelet.service
KUBELET_CONF_FILE="/etc/systemd/system/kubelet.service"
echo "[Service]" > $KUBELET_CONF_FILE
echo "User=root" >> $KUBELET_CONF_FILE
echo "ExecStart=${PWD}/build/kubelet" >> $KUBELET_CONF_FILE

# 创建controller系统进程文件controller.service
CONTROLLER_CONF_FILE="/etc/systemd/system/controller.service"
echo "[Service]" > $CONTROLLER_CONF_FILE
echo "User=root" >> $CONTROLLER_CONF_FILE
echo "ExecStart=${PWD}/build/controller" >> $CONTROLLER_CONF_FILE

# 创建serveless系统进程文件serverless.service
SERVERLESS_CONF_FILE="/etc/systemd/system/serverless.service"
echo "[Service]" > $SERVERLESS_CONF_FILE
echo "User=root" >> $SERVERLESS_CONF_FILE
echo "ExecStart=${PWD}/build/serverless" >> $SERVERLESS_CONF_FILE

# 启动，并查看状态
systemctl daemon-reload
systemctl start apiserver
systemctl start scheduler
systemctl start kubelet
systemctl start controller
systemctl start serverless

systemctl status apiserver
systemctl status scheduler
systemctl status kubelet
systemctl status controller
systemctl status serverless