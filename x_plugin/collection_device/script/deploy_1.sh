
cluster_name="cluster1"
remote_path="/data/program/app/galactus/x_plugin"
app_name="device_collect"
# 建立SSH连接并执行远程命令
sshpass -p "$blade_password" ssh -o StrictHostKeyChecking=no -T "$blade_remote_server" << EOF
  mkdir -p $remote_path/$app_name/$cluster_name
  cd $remote_path/$app_name/$cluster_name
  rm -rf device_collect
  rm -rf *.js
  rm -rf device_collect.tar.gz
  rm -rf *.json
EOF
sshpass -p "$blade_password" scp -q device_collect.tar.gz "$blade_remote_server:$remote_path/$app_name/$cluster_name"

sshpass -p "$blade_password" ssh -o StrictHostKeyChecking=no -T "$blade_remote_server" << EOF
  cd $remote_path/$app_name/$cluster_name
  tar -xzvf device_collect.tar.gz
EOF

