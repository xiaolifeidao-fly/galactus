rm -rf x_plugin.tar.gz
cd ../
tar -czvf ./x_plugin/x_plugin.tar.gz --exclude=./x_plugin/node_modules --exclude=./x_plugin/build --exclude=./x_plugin/nohup.out ./x_plugin
