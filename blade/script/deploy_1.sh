#!/bin/bash

# 从环境变量获取配置
remote_server="root@$galactus_blade_remote_server_1"
remote_password="$galactus_blade_password_1"
remote_port="$galactus_blade_port_1"

# 执行构建脚本
echo "开始构建应用..."
bash build.sh

remote_path="/data/program/app/blade"
app_name="blade"

# 检查环境变量是否设置
if [ -z "$remote_server" ] || [ -z "$remote_password" ] || [ -z "$remote_port" ]; then
    echo "错误: 请设置环境变量 galactus_blade_remote_server_1、galactus_blade_password_1 和 galactus_blade_port_1"
    exit 1
fi

# 建立SSH连接并执行远程命令
sshpass -p "$remote_password" ssh -p $remote_port -o StrictHostKeyChecking=no -T "$remote_server" << EOF
  mkdir -p $remote_path
  cd $remote_path
  rm -rf $app_name
EOF

# 上传新的二进制文件，显示上传进度
sshpass -p "$remote_password" scp -P $remote_port ./$app_name "$remote_server:$remote_path"

# 执行重启命令
sshpass -p "$remote_password" ssh -p $remote_port -o StrictHostKeyChecking=no -T "$remote_server" << EOF
  cd $remote_path
  
  # 查找并杀死已存在的进程
  pid=\$(ps -ef | grep ${app_name} | grep -v grep | awk '{print \$2}')
  if [ -n "\$pid" ]; then
      echo "killing pid: \$pid"
      kill -9 \$pid
  fi
  
  # 等待进程停止
  sleep 3
  
  # 启动应用
  nohup ./$app_name > server.log 2>&1 &
  
  # 等待应用启动
  sleep 5
  
  # 检查是否成功启动
  new_pid=\$(ps -ef | grep ${app_name} | grep -v grep | awk '{print \$2}')
  if [ -n "\$new_pid" ]; then
      echo "应用启动成功，进程ID: \$new_pid"
  else
      echo "应用启动失败"
      exit 1
  fi
EOF
