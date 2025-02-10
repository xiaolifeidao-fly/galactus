import { DoorEntity } from "@src/door/entity";
import {  Response, Request } from "playwright";
import { getUrlParameter } from "@src/door/monitor/monitor";
import { DyMonitorResponse } from "@src/door/dy/monitor/dy.monitor";
import { instance } from "@src/utils/axios";




export class DyDeviceCollectMonitor extends DyMonitorResponse<{}>{

    private num: number;

    private collectFlag : boolean;
    constructor(num: number, collectFlag : boolean = false){
        super();
        this.num = num;
        this.collectFlag = collectFlag;
    }


    public async isMatch(url : string, method: string, headers: { [key: string]: string; }, resourceType: string): Promise<boolean> {
        if(this.collectFlag){
            return false;
        }
        if(url.includes(".js") || url.includes(".css")){
            return false;
        }
        if(!resourceType.includes("xhr")){
            return false;
        }
        if(!headers){
            return false;
        }
        const cookie = headers['cookie'];
        if(!url){
            return false;
        }
        const urlParams = getUrlParameter(url);
        if(!urlParams){
            return false;
        }
        // urlParams 转成json
        const urlParamsJson: {[key: string]: string} = {};
        urlParams.forEach((value, key) => {
            urlParamsJson[key] = value;
        });
     
        const {odinTt, ttwid} = this.getTtwidAndOdinTtFromCookie(cookie);
        urlParamsJson['ttwid'] = ttwid;

        if(!this.validateParams(urlParamsJson)){
            return false;
        }
        this.collectFlag = true;
        return true;
    }
    
    getApiName(): string {
       return "web/hot/search/list/"
    }

    getKey(): string {
        return "dy_device_collect"
    }

    getTtwidAndOdinTtFromCookie(cookie: string): {odinTt: string, ttwid: string} {
        if (!cookie || cookie == ''){
            return {odinTt : '',ttwid : ''}
        }
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
        const url = request?.url();
        const headers = await request?.allHeaders();
        if (!('status_code' in data) || data['status_code'] != 0){
            console.log(data);
            console.log("dy device collect monitor callback error", data);
            return;
        }
        if(!headers){
            console.log("dy device collect headers error", data);
            return;
        }
        const uifid = headers['uifid'];
        const userAgent = headers['user-agent'];
        const cookie = headers['cookie'];
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

        //webapp	6383	channel_pc_web	6
        if(!('device_platform' in urlParamsJson)){
            urlParamsJson['device_platform'] = 'webapp';
        }
        if(!('aid' in urlParamsJson)){
            urlParamsJson['aid'] = '6383';
        }
        if(!('channel' in urlParamsJson)){
            urlParamsJson['channel'] = 'channel_pc_web';
        }
        if(!('source' in urlParamsJson)){
            urlParamsJson['source'] = '6';
        }
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
        this.collectFlag = true;
        await this.saveDevice(urlParamsJson);
    }

    async saveDevice(params: {[key: string]: string}){
        const response = await instance.post('/devices/save', params);
        console.log(this.num, " result is : ", response)
    }

    


}


