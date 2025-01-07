import path from 'path';
import fs from 'fs'
import { Browser, chromium, devices,firefox, BrowserContext, Page, Route ,Request, Response} from 'playwright';
import { Monitor, MonitorChain, MonitorRequest, MonitorResponse } from './monitor/monitor';
import { DoorEntity } from './entity';
import { ActionChain, ActionResult } from './element/element';

const browserMap = new Map<string, Browser>();

export abstract class DoorEngine<T = any> {

    private chromePath: string | undefined;

    browser: Browser | undefined;

    context: BrowserContext | undefined;

    public resourceId : string;

    public headless: boolean = true;

    monitors : Monitor<T>[] = [];

    monitorsChain : MonitorChain<T>[] = [];

    page : Page | undefined;

    constructor(resourceId : string, headless: boolean = true, chromePath: string = ""){
        this.resourceId = resourceId;
        if(chromePath){
            this.chromePath = chromePath;
        }else{
            this.chromePath = this.getChromePath();
        }
        this.headless = headless;
    }

    getChromePath() : string | undefined{
        return process.env.CHROME_PATH;
    }

    addMonitor(monitor: Monitor){
        this.monitors.push(monitor);
    }

    addMonitorChain(monitorChain: MonitorChain<T>){
        this.monitorsChain.push(monitorChain);
        this.monitors.push(...monitorChain.getMonitors());
    }

    public async init(url : string|undefined = undefined) : Promise<Page | undefined> {
        if(this.browser){
            return undefined;
        }
        this.browser = await this.createBrowser();
        this.context = await this.createContext();
        if(!this.context){
            return undefined;
        }
        const page = await this.context.newPage();
        if(url){
            await page.goto(url);
        }
        console.log("new context and new page ")
        this.onRequest(page);
        this.onResponse(page);
        this.page = page;
        return page;
    }

    public async release(){
        await this.closePage();
        await this.closeContext();
        await this.closeBrowser();
    }

    public async closePage(){
        try{
            if(this.page && !this.page.isClosed()){
                await this.page.close();
            }
        }catch(e){
            console.error("closePage error", e);
        }
    }


    public async doBeforeRequest(request: Request, headers: { [key: string]: string; }){
        for(const monitor of this.monitors){
            if(monitor.finishTag){
                continue;
            }
            if(!(monitor instanceof MonitorRequest)){
                continue;
            }
            if(!await monitor.isMatch(request.url(), request.method(), headers)){
                continue;
            }
            const requestMonitor = monitor as MonitorRequest<T>;
            let data;
            if(requestMonitor.handler){
                data = await requestMonitor.handler(request, undefined);
            }
            monitor._doCallback(new DoorEntity(data ? true : false, data));
            monitor.setFinishTag(true);
        }
    }

    public async onRequest(page : Page){
        page.route("*/**", async (router : Route) => {
            // 获取请求对象
            const request = router.request();
            const headers = await request.allHeaders();
            await this.doBeforeRequest(request, headers);
            router.continue();
        });
    }

    public async doAfterResponse(response: Response){
        for(const monitor of this.monitors){
            if(monitor.finishTag){
                continue;
            }
            if(!(monitor instanceof MonitorResponse)){
                continue;
            }
            const responseMonitor = monitor as MonitorResponse<T>;
            if(!await monitor.doMatchResponse(response)){
                continue;
            }
            const data = await responseMonitor.getResponseData(response);
            const doorEntity = new DoorEntity<T>(data ? true : false, data);
            responseMonitor._doCallback(doorEntity, response.request(), response);
            responseMonitor.setFinishTag(true);
        }
    }

    public async onResponse(page : Page){
        page.on('response', async (response) => {
            await this.doAfterResponse(response);
        });
    }

    public async openWaitMonitor(page : Page,  url: string, monitor : Monitor<T | any>, headers: Record<string, string> = {}){
        try{
            this.addMonitor(monitor);
            await this.startMonitor();
            await page.goto(url, { timeout: 5000});
        }catch(e){
            console.error("openWaitMonitor error", e);
        }
        const doorEntity = await monitor.waitForAction();
        return doorEntity;
    }



    public async openWaitMonitorChain(page : Page,  url: string, monitorChain: MonitorChain<T | any>, headers: Record<string, string> = {}){
        try{
            this.addMonitorChain(monitorChain);
            await this.startMonitor();
            page.goto(url, { timeout: 5000});
        }catch(e){
            console.error("openWaitMonitorChain error", e);
        }
        const doorEntity = await monitorChain.waitForAction();
        return doorEntity;
    }

    public async startMonitor(){
        for(const monitor of this.monitors){
            monitor.start();
        }
    }

    public async doFillWaitForElement(page : Page, version: string, doorType: string, data? : any) {
        // const actionCommands :ActionCommand[]= [];//await getDoorList(version, doorType);
        // let prevResult : DoorEntity<T> | undefined = undefined;
        // for (const actionCommand of actionCommands) {
        //     const monitorChain = this.getMonitorChainFromChain(actionCommand.key);
        //     if(monitorChain){
        //         await monitorChain.start();
        //     }else{
        //         const monitor = this.getMonitor(actionCommand.key);
        //         if(monitor){
        //             await monitor.start();
        //         }
        //     }
        //     const dynamicFunction = new Function(actionCommand.code)();
        //     await dynamicFunction(page, prevResult, data);
        //     prevResult = await monitorChain?.waitForAction();
        //     if(!prevResult){
        //         continue;
        //     }
        //     if(!prevResult.code){
        //         return prevResult;
        //     }
        // }
        // return prevResult;
    }

    getMonitorChainFromChain(key : string) : MonitorChain<T> | undefined{
        if(!this.monitorsChain || this.monitorsChain.length == 0){
            return undefined;
        }
        for(const monitorChain of this.monitorsChain){
            if(monitorChain.getKey() == key){
                return monitorChain;
            }
        }
        return undefined;
    }

    getMonitor(key : string) : Monitor<T> | undefined{
        if(!this.monitors || this.monitors.length == 0){
            return undefined;
        }
        for(const monitor of this.monitors){
            if(monitor.getKey() == key){
                return monitor;
            }
        }
        return undefined;
    }

    public async closeContext(){
        try{
            if(this.context){
                await this.context.close();
            }
        }catch(e){
            console.error("closeContext error", e);
        }
    }

    public async closeBrowser(){
        try{
            if(this.browser){
                await this.browser.close();
            }
        }catch(e){
            console.error("closeBrowser error",e);
        }
    }

    getKey(){
        return `door_engine_${this.getNamespace()}_${this.resourceId}`;
    }



    abstract getNamespace(): string;

    async createContext(){
        if(!this.browser){
            return;
        }
        const macbook = devices['MacBook Pro 16'];
        const context = await this.browser?.newContext({
            ...macbook,
            bypassCSP: true,
            userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36', 
            extraHTTPHeaders: {
                'sec-ch-ua': '"Google Chrome";v="131", "Chromium";v="131", "Not_A Brand";v="24"',
                'sec-ch-ua-mobile': '?0', // 设置为移动设备
                'sec-ch-ua-platform': '"macOS"',
              },
              locale: 'zh-CN', 
            // proxy: {
            //     server: 'http://127.0.0.1:8888',
            // },
            viewport: { width: 1920, height: 1080 }
        });

        // 注入脚本来修改 navigator 信息
        await context.addInitScript(() => {
            // 模拟不同的 navigator 信息
            //@ts-ignore    
            Object.defineProperty(navigator, 'userAgent', {
            get: () => 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36',
            });
            //@ts-ignore
            Object.defineProperty(navigator, 'platform', {
            get: () => 'MacIntel',
            });
            //@ts-ignore
            Object.defineProperty(navigator, 'appVersion', {
            get: () => '5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36',
            });
        });
        return context;
    }

    async createBrowser(){
        let key = this.headless.toString();
        if (this.chromePath) {
            key += "_" + this.chromePath;
        }
        if(browserMap.has(key)){
            return browserMap.get(key);
        }
        const browser = await chromium.launch({
            headless: this.headless,
            executablePath: this.chromePath,
            args: [
                '--no-sandbox', // 取消沙箱，某些网站可能会检测到沙箱模式
                '--disable-setuid-sandbox',
                '--disable-blink-features=AutomationControlled',  // 禁用浏览器自动化控制特性
              ]
        });
        browserMap.set(key, browser);
        return browser;
    }

}

