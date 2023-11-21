import axios, { AxiosRequestConfig, Method } from 'axios';
import { ElLoading } from 'element-plus';

// 创建请求实例  
const instance = axios.create({
    baseURL: import.meta.env.VITE_BASE_URL,
    // 指定请求超时的毫秒数  
    timeout: 10000,
    // 表示跨域请求时是否需要使用凭证  
    withCredentials: false,
});

// 设置请求头
instance.defaults.headers.post['Content-Type'] = 'application/json;charset=UTF-8';
instance.defaults.headers.put['Content-Type'] = 'application/x-www-form-urlencoded';
// instance.defaults.headers.put['Content-Type'] = 'application/json';

// 加载实例
let loadingInstance = null;
// 开始loading
const startLoading = () => {
    if (!loadingInstance) {
        loadingInstance = ElLoading.service({});
    }
};
// 结束loading
const endLoading = () => {
    if (loadingInstance) {
        loadingInstance.close();
        loadingInstance = null
    }
};

// 取消重复请求
const pending = [];

// 定义接口
interface PendingType {
    url?: string;
    method?: Method;
    params: any;
    data: any;
    cancel: any;
}

// 移除重复请求
const removePending = (config: AxiosRequestConfig) => {
    for (const key in pending) {
        const item: number = +key;
        const list: PendingType = pending[key];
        // 当前请求在数组中存在时执行函数体
        if (list.url === config.url && list.method === config.method && JSON.stringify(list.params) === JSON.stringify(config.params) && JSON.stringify(list.data) === JSON.stringify(config.data)) {
            // 执行取消操作
            list.cancel('操作太频繁，请稍后再试');
            // 从数组中移除记录
            pending.splice(item, 1);
        }
    }
};


// 请求拦截器（发起请求之前的拦截）  
instance.interceptors.request.use(
    (config): AxiosRequestConfig<any> => {
        // 不同基地址请求
        // if (config.requestBase === UserEnum.adminBaseUrl) {
        //     config.baseURL = window.configList.adminBaseUrl;
        // }
        removePending(config);
        config.cancelToken = new axios.CancelToken(c => {
            pending.push({ url: config.url, method: config.method, params: config.params, data: config.data, cancel: c });
        });
        /**  
         * 在这里一般会携带前台的参数发送给后台，比如下面这段代码：  
         * const token = sessionStorage.getItem('token')  
         * if (token) {  
         *  config.headers.Authorization = `Basic ${token}`  
         * }  
         */
        return config;
    },
    (error) => {
        return Promise.reject(error);
    },
);

// 响应拦截器（获取到响应时的拦截）  
instance.interceptors.response.use(
    (response) => {
        removePending(response.config);
        /**  
         * 根据你的项目实际情况来对 response 和 error 做处理  
         * 这里对 response 和 error 不做任何处理，直接返回  
         */
        return response;
    },
    (error) => {
        return Promise.reject(error);
    },
);

interface ResType<T> {
    code: number;
    data?: T;
    msg?: string;
    message?: string;
    err?: string;
}

interface IOptions {
    loading?: boolean;
    isFormUrlencoded?: boolean;
    requestBase?: string; // 修改基地址
}

interface Http {
    post<T>(url: string, data?: unknown,  options?: IOptions): Promise<ResType<T>>;
    get<T>(url: string, options?: IOptions): Promise<ResType<T>>;
    put<T>(url: string, data?: unknown, options?: IOptions): Promise<ResType<T>>;
    upload<T>(url: string, file?: unknown): Promise<ResType<T>>;
    _delete<T>(url: string, options?: IOptions): Promise<ResType<T>>;
}


 /**
 * 是否是x-www-form-urlencoded格式请求
 * @param {} 
 */
function isFormUrlencoded(url, params, options) {
    if (options?.isFormUrlencoded) {
        const list: any[] = [];
        for (const key in params) {
            if (params[key] !== null) {
                list.push(`${key}=${encodeURIComponent(params[key])}`);
            }
        }
        const newParams = list.join('&');
        const newOptions = {
            ...options,
            headers: {
                ...options.headers,
                'Content-Type': 'application/x-www-form-urlencoded',
            },
        };

        return { newUrl: url + '?' + newParams, newParams: {}, newOptions };
    }
    return { newUrl: url, newParams: params, newOptions: options };
}

// 导出常用函数  
const http: Http = {
    post(url, data, options?) {
        return new Promise((resolve, reject) => {
            options?.loading && startLoading();
            const { newUrl, newData, newOptions } = isFormUrlencoded(data, options);
            instance
                .post(newUrl, JSON.stringify(newData), newOptions)
                .then((res) => {
                    options?.loading && endLoading();
                    resolve(res.data);
                })
                .catch((err) => {
                    options?.loading && endLoading();
                    reject(err.data);
                });
        });
    },
    get(url, params, options?) {
        return new Promise((resolve, reject) => {
            options?.loading && startLoading();
            instance
                .get(url, { params, ...options })
                .then((res) => {
                    options?.loading && endLoading();
                    resolve(res.data);
                })
                .catch((err) => {
                    options?.loading && endLoading();
                    reject(err.data);
                });
        });
    },
    put(url, data, options?) {
        return new Promise((resolve, reject) => {
            options?.loading && startLoading();
            const { newUrl, newData, newOptions } = isFormUrlencoded(data, options);
            instance
                .put(newUrl, newData, newOptions)
                .then((res) => {
                    options?.loading && endLoading();
                    resolve(res.data);
                })
                .catch((err) => {
                    options?.loading && endLoading();
                    reject(err.data);
                });
        });
    },
     upload(url, file) {
        return new Promise((resolve, reject) => {
            // options?.loading && startLoading()
            instance
                .post(url, file, {
                    headers: { 'Content-Type': 'multipart/form-data' },
                })
                .then((res) => {
                    // options?.loading && endLoading()
                    resolve(res.data);
                })
                .catch((err) => {
                    // options?.loading && endLoading()
                    reject(err.data);
                });
        });
    },
    downFile(url) {
        const iframe = document.createElement('iframe');
        iframe.style.display = 'none';
        iframe.src = url;
        iframe.onload = function () {
            document.body.removeChild(iframe);
        };
        document.body.appendChild(iframe);
    },
    downBlob(url,fileName) {
        fetch(url, {
        method: 'get',
    })
        .then((res) => res.blob())
        .then((blob) => {
            const a = document.createElement('a');
            // 获取 blob 本地文件连接 (blob 为纯二进制对象，不能够直接保存到磁盘上)
            const url = window.URL.createObjectURL(
                new Blob([blob], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' }),
            );
            a.href = url;
            //定义导出的文件名
            a.download = `${fileName}.xls`;
            a.click();
            window.URL.revokeObjectURL(url);
        });
    },
    _delete(url, options?) {
        return new Promise((resolve, reject) => {
            options?.loading && startLoading();
            instance
                .delete(url, options)
                .then((res) => {
                    options?.loading && endLoading();
                    resolve(res.data);
                })
                .catch((err) => {
                    options?.loading && endLoading();
                    reject(err.data);
                });
        });
    }
}

export default http;
