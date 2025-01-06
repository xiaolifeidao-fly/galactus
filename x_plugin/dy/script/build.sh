rm -rf dy.tar.gz
cd ../
tar -czvf ./script/dy.tar.gz --exclude=./node_modules --exclude=./build --exclude=./nohup.out --exclude=./script *
