const  {start, init} = require("./application")
const dyRouter = require("./dy/router")


function register(){
    dyRouter.init();
}

function run(){
    init();
    register();
    start(9001);
}

run();

