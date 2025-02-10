import * as dotenv from 'dotenv';
dotenv.config(); // 加载 .env 文件中的环境变量

import { randomUUID } from "crypto";
import { DyEngine } from "@src/door/dy/dy.engine";
import { DyDeviceCollectMonitor } from "@src/door/dy/monitor/device/monitor";


async function sleep(ms: number){
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function collectDevice(num : number){
    const engine = new DyEngine<{}>("dy_device_collect");
    try{
        const monitor = new DyDeviceCollectMonitor(num);
        const page = await engine.init();
        if(!page){ 
            return;
        }
        //return await engine.openWaitMonitor(page, "https://www.browserscan.net/", monitor);
        return await engine.openWaitMonitor(page, "https://www.douyin.com/?recommend=1", monitor);
    }catch(e){  
        console.log("collectDevice error ", e);
    }finally{
        console.log("collectDevice finally ", num);
        try{
            await engine.closeContext();
            await sleep(2000);
        }catch(e){
            console.error("collectDevice close error", e);
        }
    }  
}


async function main(){
    const num = process.env.COLLECT_NUM;
    console.log("COLLECT_NUM is ", num);
    if(!num){
        console.log("COLLECT_NUM is not set");
        return;
    }
    for(let i = 0; i < Number(num); i++){
        console.log("collectDevice start ", i);
        await collectDevice(i);
        await sleep(2000);
        console.log("collectDevice end ", i);
    }
}

main();