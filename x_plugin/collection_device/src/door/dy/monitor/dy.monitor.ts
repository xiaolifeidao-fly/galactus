import { Response } from "playwright";
import { Monitor, MonitorChain, MonitorRequest, MonitorResponse } from "@src/door/monitor/monitor";


export abstract class DyMonitorRequest<T> extends MonitorRequest<T> {



}

export abstract class DyMonitorResponse<T> extends MonitorResponse<T> {

    getType(): string{
        return "DY";
    }

    async filter(url: string, resourceType: string, method: string, headers: { [key: string]: string; }): Promise<boolean> {
        if (url.includes('png') || url.includes('.png') || url.includes('.jpg') || url.includes('.jpeg') || url.includes('.gif') || url.includes('.mp4') || url.includes('.webm')) {
            return true;
        } 
        return false;
    }

    public async isMatch(url : string, method: string, headers: { [key: string]: string; }, resourceType: string): Promise<boolean> {
        if(url.includes(".js")){
            return false;
        }
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