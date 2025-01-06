
cluster_name="cluster1"
remote_path="/data/program/app/galactus/x_plugin"
app_name="dy"
# 建立SSH连接并执行远程命令
sshpass -p "$blade_password" ssh -o StrictHostKeyChecking=no -T "$blade_remote_server" << EOF
  mkdir -p $remote_path/$app_name/$cluster_name
  cd $remote_path/$app_name/$cluster_name
  rm -rf dy
  rm -rf *.js
  rm -rf dy.tar.gz
  rm -rf *.json
EOF
sshpass -p "$blade_password" scp -q dy.tar.gz "$blade_remote_server:$remote_path/$app_name/$cluster_name"

sshpass -p "$blade_password" ssh -o StrictHostKeyChecking=no -T "$blade_remote_server" << EOF
  cd $remote_path/$app_name/$cluster_name
  tar -xzvf dy.tar.gz
EOF

