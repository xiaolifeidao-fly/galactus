
cluster_name="cluster1"
remote_path="/data/program/app/galactus/x_plugin"
app_name="server.js"
# 建立SSH连接并执行远程命令
sshpass -p "$blade_password" ssh -o StrictHostKeyChecking=no -T "$blade_remote_server" << EOF
  mkdir -p $remote_path/$cluster_name
  cd $remote_path/$cluster_name
  rm -rf dy
  rm -rf *.js
  rm -rf x_plugin.tar.gz
  rm -rf *.json
EOF
sshpass -p "$blade_password" scp -q x_plugin.tar.gz "$blade_remote_server:$remote_path/$cluster_name"

