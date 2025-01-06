rm -rf device_collect.tar.gz
cd ../
tar -czvf ./script/device_collect.tar.gz --exclude=.env --exclude=./node_modules --exclude=./dist --exclude=./nohup.out --exclude=./script *
