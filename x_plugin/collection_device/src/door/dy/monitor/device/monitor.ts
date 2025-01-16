import { DoorEntity } from "@src/door/entity";
import {  Response, Request } from "playwright";
import { getUrlParameter } from "@src/door/monitor/monitor";
import { DyMonitorResponse } from "@src/door/dy/monitor/dy.monitor";
import { instance } from "@src/utils/axios";


export class DyDeviceCollectMonitor extends DyMonitorResponse<{}>{

    private num: number;

    constructor(num: number){
        super();
        this.num = num;
    }

    getApiName(): string {
       return "web/hot/search/list/"
    }

    getKey(): string {
        return "dy_device_collect"
    }

    getTtwidAndOdinTtFromCookie(cookie: string): {odinTt: string, ttwid: string} {
        const cookieJson = cookie.split(';').map(item => item.trim());
        let odinTt = cookieJson.find(item => item.startsWith('odin_tt='));
        let ttwid = cookieJson.find(item => item.startsWith('ttwid='));
        if(odinTt){
            odinTt = odinTt.split('=')[1];
        }else{
            odinTt = '';
        }
        if(ttwid){
            ttwid = ttwid.split('=')[1];
        }else{
            ttwid = '';
        }
        return {odinTt, ttwid};
    }

    validateParams(params: {[key: string]: string}){
        if(!params){
            return false;
        }
        if(!params['ttwid']){
            return false;
        }
        if(!params['webid']){
            return false;
        }
        return true;
    }

    async doCallback(doorEntity: DoorEntity<{}>, request?: Request | undefined, response?: Response | undefined): Promise<void> {
            const data = doorEntity.data;
            if (!('status_code' in data) || data['status_code'] != 0){
                 console.log("dy device collect monitor callback error", data);
                 return;
            }
            const headers = await request?.allHeaders();
            if(!headers){
                console.log("dy device collect headers error", data);
                return;
            }
            const uifid = headers['uifid'];
            const userAgent = headers['user-agent'];
            const cookie = headers['cookie'];
            const url = request?.url();
            if(!url){
                return;
            }
            console.log("dy device collect url is ", url);
            console.log("dy device collect headers is ", headers);
            const urlParams = getUrlParameter(url);
            if(!urlParams){
                console.log("dy device collect urlParams error", data);
                return;
            }
            // urlParams 转成json
            const urlParamsJson: {[key: string]: string} = {};
            urlParams.forEach((value, key) => {
                urlParamsJson[key] = value;
            });
            urlParamsJson['uifid'] = uifid;
            urlParamsJson['user_agent'] = userAgent;
            urlParamsJson['verify_fp'] = urlParamsJson['verifyFp'];
            urlParamsJson['sec_ch_ua_platform'] = headers['sec-ch-ua-platform'];
            urlParamsJson['sec_ch_ua'] = headers['sec-ch-ua'];
            urlParamsJson['cookie'] = cookie;
            const {odinTt, ttwid} = this.getTtwidAndOdinTtFromCookie(cookie);
            urlParamsJson['odin_tt'] = odinTt;
            urlParamsJson['ttwid'] = ttwid;
            if(!this.validateParams(urlParamsJson)){
                console.log("dy device collect urlParams error", data);
                return;
            }
            await this.saveDevice(urlParamsJson);
    }

    async saveDevice(params: {[key: string]: string}){
        const response = await instance.post('/devices/save', params);
        console.log(this.num, " result is : ", response)
    }

    


}


