import * as dotenv from 'dotenv';
dotenv.config(); // 加载 .env 文件中的环境变量

import { randomUUID } from "crypto";
import { DyEngine } from "@src/door/dy/dy.engine";
import { DyDeviceCollectMonitor } from "@src/door/dy/monitor/device/monitor";




async function main(){

    const engine = new DyEngine(randomUUID().toString());
    try{
        const monitor = new DyDeviceCollectMonitor();
        const page = await engine.init();
        if(!page){ 
            return;
        }
        return await engine.openWaitMonitor(page, "https://www.douyin.com/?recommend=1", monitor);
    }finally{
        await engine.closePage();
    }  

}

main();