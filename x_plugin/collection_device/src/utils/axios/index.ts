import axios, { AxiosInstance, AxiosResponse, AxiosError, InternalAxiosRequestConfig } from 'axios';
import { plainToInstance } from 'class-transformer';
// 定义一个 HttpError 类，扩展自 Error
class HttpError extends Error {
  code?: any;

  constructor(message: string, code?: any) {
    super(message);
    this.name = 'HttpError';
    if (code != null) {
      this.code = code;
    }
  }
}

// 抛出 http 异常
function rejectHttpError(message: string, code?: any): Promise<never> {
  const error = new HttpError(message, code);
  return Promise.reject(error);
}

function getBaseUrl() {
  return process.env.SERVER_TARGET;
}

const instance: AxiosInstance = axios.create({
  timeout: 60000,
  // baseURL: '',
  baseURL: getBaseUrl(),
  withCredentials: true,
  // 登录成功后，设置请求头 Authorization
  // headers: {
  //   Authorization: ""
  // },
});


// 前端请求拦截器
instance.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    return config;
  },
);

// 后端响应拦截器
instance.interceptors.response.use(
  (response: AxiosResponse) => {
    const responseData = response.data as { code: number, data: any, error: string };
    if (responseData.code != 0) {
      return rejectHttpError(
        responseData.error || '请求异常！',
        500
      );
    }
    return responseData.data;
  },
  (error: AxiosError) => {
    if (error.response) {
      const data = error.response.data as { error?: string; code?: any };
      if (data && data.error) {
        return rejectHttpError(data.error, data.code);
      }
      return rejectHttpError('请求异常：' + error.response.statusText);
    }

    if (error.request) {
      return rejectHttpError('请求异常：无返回结果');
    }
    return rejectHttpError(error.message);
  }
);
const getData = async <T>(clazz: new (...args: any[]) => T, url : string, params? : {}): Promise<T|null> =>{
    const data = await instance.get(url, {params})
    return plainToInstance(clazz, data)
}

const getDataList = async <T>(clazz: new (...args: any[]) => T, url : string, params? : {}) : Promise<T[]> => {
  const data : {}[] = await instance.get(url, {params})
  if (data == undefined || data.length  == 0) {
     return [];
  }
  const dataList : T[] = []
  data.forEach(item => {
     const item_instance : T = plainToInstance(clazz, item);
     dataList.push(item_instance);
  })
  return dataList;
}

class PageData<T> {
   data : T[];
   total : number;

   constructor(data : T[], total : number) {
      this.data = data;
      this.total = total;
   }

} 

const getPage = async <T>(clazz: new (...args: any[]) => T, url : string, params? : {}) : Promise<PageData<T>> => {
  const pageData : any = await instance.get(url, {params})
  if (pageData == undefined) {
     return new PageData<T>([], 0);
  }
  const dataList : T[] = []
  const responseDataList : {}[] = pageData.data;
  responseDataList.forEach(item=> {
     const item_instance : T = plainToInstance(clazz, item);
     dataList.push(item_instance);
  })
  const total : number = pageData.total;
  return new PageData<T>(dataList, total);
}


export {
  getData,
  getDataList,
  getPage,
  PageData,
  instance
} ;


