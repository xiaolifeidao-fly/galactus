#!/bin/bash

# 从环境变量获取配置
remote_server="root@$galactus_dy_remote_server_1"
remote_password="$galactus_dy_password_1"
remote_port="$galactus_dy_port_1"

# 执行构建脚本
echo "开始构建应用..."
bash build.sh

remote_path="/data/program/app/galactus/x_plugin"
app_name="dy"
cluster_name="cluster1"

# 检查环境变量是否设置
if [ -z "$remote_server" ] || [ -z "$remote_password" ] || [ -z "$remote_port" ]; then
    echo "错误: 请设置环境变量 galactus_dy_remote_server_1、galactus_dy_password_1 和 galactus_dy_port_1"
    exit 1
fi

# 建立SSH连接并执行远程命令
sshpass -p "$remote_password" ssh -p $remote_port -o StrictHostKeyChecking=no -T "$remote_server" << EOF
  mkdir -p $remote_path/$app_name/$cluster_name
  cd $remote_path/$app_name/$cluster_name
  rm -rf dy
  rm -rf *.js
  rm -rf dy.tar.gz
  rm -rf *.json
EOF

# 上传新的压缩包，显示上传进度
sshpass -p "$remote_password" scp -P $remote_port -q dy.tar.gz "$remote_server:$remote_path/$app_name/$cluster_name"

# 执行解压命令
sshpass -p "$remote_password" ssh -p $remote_port -o StrictHostKeyChecking=no -T "$remote_server" << EOF
  cd $remote_path/$app_name/$cluster_name
  tar -xzvf dy.tar.gz
EOF

