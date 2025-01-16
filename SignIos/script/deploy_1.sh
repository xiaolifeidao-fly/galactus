
cluster_name="cluster1"
remote_path="/data/program/app/galactus/SignIos"
app_name="SignIos"
# 建立SSH连接并执行远程命令
sshpass -p "$blade_password" ssh -o StrictHostKeyChecking=no -T "$blade_remote_server" << EOF
  mkdir -p $remote_path/$cluster_name
  cd $remote_path/$cluster_name
  rm -rf SignIos
EOF
sshpass -p "$blade_password" scp -q SignIos "$blade_remote_server:$remote_path/$cluster_name"
