const express = require('express');
const { createProxyMiddleware } = require('http-proxy-middleware');
const path = require('path');
const PORT = process.env.PORT || 5000;
require('dotenv').config();
const app = express();
app.use(express.static('dist')); // 这里是存放静态文件的目录

const apiPrefix = process.env.REACT_APP_SERVER_API;

const preApi = apiPrefix;

app.use(express.json())
function init(){
    app.all("*", function (req, res, next) {
        // 设置是否运行客户端设置 withCredentials
        // 即在不同域名下发出的请求也可以携带 cookie
        res.header("Access-Control-Allow-Credentials", true)
        // 第二个参数表示允许跨域的域名，* 代表所有域名  
        res.header('Access-Control-Allow-Origin', '*')
        res.header('Access-Control-Allow-Methods', 'GET, PUT, POST, OPTIONS') // 允许的 http 请求的方法
        // 允许前台获得的除 Cache-Control、Content-Language、Content-Type、Expires、Last-Modified、Pragma 这几张基本响应头之外的响应头
        res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization, Content-Length, X-Requested-With')
        if (req.method == 'OPTIONS') {
            res.sendStatus(200)
        } else {
            next()
        }
    });
}

function use(path, func){
    app.use(path, func)
}


function start(port){
    //配置服务端口
    app.listen(port, () => {
        console.log(`localhost:${port} successfully`);
    });
}

module.exports = {
    use,
    start,
    init
}