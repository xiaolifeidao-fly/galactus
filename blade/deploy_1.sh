
cluster_name="cluster1"
remote_path="/data/program/app/barry/bootstrap/barry-bootstrap-pro"
app_name="barry-bootstrap-pro-0.0.1-SNAPSHOT.jar"
# 建立SSH连接并执行远程命令
sshpass -p "$barry_bootstrap_pro_password" ssh -o StrictHostKeyChecking=no -T "$barry_bootstrap_pro_remote_server" << EOF
  mkdir -p $remote_path/$cluster_name
  cd $remote_path/$cluster_name
  rm -rf *.jar
EOF
sshpass -p "$barry_bootstrap_pro_password" scp -q ./target/$app_name "$barry_bootstrap_pro_remote_server:$remote_path/$cluster_name"
