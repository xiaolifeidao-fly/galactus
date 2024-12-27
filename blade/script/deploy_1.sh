
cluster_name="cluster1"
remote_path="/data/program/app/galactus/blade"
app_name="blade"
# 建立SSH连接并执行远程命令
sshpass -p "$blade_password" ssh -o StrictHostKeyChecking=no -T "$blade_remote_server" << EOF
  mkdir -p $remote_path/$cluster_name
  cd $remote_path/$cluster_name
  rm -rf blade
EOF
sshpass -p "$blade_password" scp -q blade "$blade_remote_server:$remote_path/$cluster_name"
