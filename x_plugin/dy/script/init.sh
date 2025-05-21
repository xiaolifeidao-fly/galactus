#!/bin/bash

# 从环境变量获取配置
remote_server="root@$galactus_dy_remote_server_1"
remote_password="$galactus_dy_password_1"
remote_port="$galactus_dy_port_1"

remote_path="/data/program/app/galactus/x_plugin/dy"

# 检查环境变量是否设置
if [ -z "$remote_server" ] || [ -z "$remote_password" ] || [ -z "$remote_port" ]; then
    echo "错误: 请设置环境变量 galactus_dy_remote_server_1、galactus_dy_password_1 和 galactus_dy_port_1"
    exit 1
fi

# 建立SSH连接并执行远程命令
sshpass -p "$remote_password" ssh -p $remote_port -o StrictHostKeyChecking=no -T "$remote_server" << EOF
  mkdir -p $remote_path
  rm -rf $remote_path/*.sh
EOF

# 上传脚本文件
sshpass -p "$remote_password" scp -P $remote_port -p ./start.sh "$remote_server:$remote_path"
sshpass -p "$remote_password" scp -P $remote_port -p ./stop.sh "$remote_server:$remote_path"
