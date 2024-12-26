#!/bin/bash

# 设置远程服务器和路径
remote_path="/data/program/app/barry/bootstrap/barry-bootstrap-pro"
# 建立SSH连接并执行远程命令
sshpass -p "$barry_bootstrap_pro_password" ssh -o StrictHostKeyChecking=no -T "$barry_bootstrap_pro_remote_server" << EOF
  mkdir -p $remote_path
  rm -rf $remote_path/*.sh
EOF

sshpass -p "$barry_bootstrap_pro_password" scp -p ./start.sh "$barry_bootstrap_pro_remote_server:$remote_path"
sshpass -p "$barry_bootstrap_pro_password" scp -p ./stop.sh "$barry_bootstrap_pro_remote_server:$remote_path"
