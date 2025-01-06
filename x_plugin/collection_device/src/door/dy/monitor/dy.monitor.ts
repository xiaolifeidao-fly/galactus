import { Response } from "playwright";
import { Monitor, MonitorChain, MonitorRequest, MonitorResponse } from "@src/door/monitor/monitor";


export abstract class DyMonitorRequest<T> extends MonitorRequest<T> {


}

export abstract class DyMonitorResponse<T> extends MonitorResponse<T> {

    getType(): string{
        return "DY";
    }

    public async isMatch(url : string, method: string, headers: { [key: string]: string; }): Promise<boolean> {
        if(url.includes(".js")){
            return false;
        }
        console.log("url", url);
        if(url.includes(this.getApiName())){
            return true;
        }
        return false;
    }

    abstract getApiName(): string;

}

export abstract class DyMonitorChain<T> extends MonitorChain<T>{

    getType(): string {
        return "DY"
    }
}