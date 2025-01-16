function isTypedArray(value) {
    return value instanceof Array ||
        value instanceof Int8Array ||
        value instanceof Uint8Array ||
        value instanceof Int16Array ||
        value instanceof Uint16Array ||
        value instanceof Int32Array ||
        value instanceof Uint32Array ||
        value instanceof Float32Array ||
        value instanceof Float64Array;
}

const arrayUtils = {
    Array: (value) => Array.isArray(value),
    Int8Array: (value) => value instanceof Int8Array,
    Uint8Array: (value) => value instanceof Uint8Array,
    Int16Array: (value) => value instanceof Int16Array,
    Uint16Array: (value) => value instanceof Uint16Array,
    Int32Array: (value) => value instanceof Int32Array,
    Uint32Array: (value) => value instanceof Uint32Array,
    Float32Array: (value) => value instanceof Float32Array,
    Float64Array: (value) => value instanceof Float64Array,
}

class XysLogSaver {
    constructor() {
        this.chromeLog = false; // 用于控制是否使用 Chrome 的日志功能
        this.switch = true; // 用于控制日志是否保存
        this.logs = [];  // 用于存储日志内容
        this.logCount = 1;  // 用于记录日志条数
        this.blob = null;  // 用于保存最新的日志 Blob
        this.maxLogCount = 10000;  // 设置最大日志条数，超过这个数量后清空
    }

    // 将多个参数的日志内容添加到日志数组
    addLog(...logArgs) {
        if (!this.switch) return;
        const logMessage = logArgs.map(arg => this.formatLog(arg)).join(' , ');
        const logEntry = `第 ${this.logCount} 条: ${logMessage}`;
        if (this.chromeLog) {
            console.log.apply(this, [`第 ${this.logCount} 条: `, ...logArgs]);
            this.logCount++;
            return;
        }
        this.logs.push(logEntry);
        this.logCount++;  // 每次添加日志时，增加计数器

        // 如果日志条数超过最大限制，清空日志数组并保存
        if (this.logs.length >= this.maxLogCount) {
            this.saveLogsToBlob();
            this.logs = [];  // 清空日志数组
        }
    }

    // 格式化日志，根据类型转换为合适的字符串
    formatLog(log) {
        if (log && typeof log === 'object') {

            for (const [type, checkFn] of Object.entries(arrayUtils)) {
                if (checkFn(log)) {
                    if (arrayUtils.Array(log) && log.some(item => typeof item === 'object')) {
                        return `[${log.map(item => this.formatLog(item)).join(',')}]`;
                    }
                    return `[${log}] (${type}类型，长度:${log.length})`;
                }
            }

            try {
                // 如果是对象，转换为 JSON 字符串
                return JSON.stringify(log);
            } catch (e) {
                // 如果对象有循环引用，返回特定信息
                return `${log} (循环引用)`;
            }
        } else if (typeof log === 'function') {
            // 如果是函数，转换为函数字符串
            return log.toString();
        } else if (log === null) {
            // 对 null 的特殊处理
            return 'null';
        } else {
            // 其他类型直接转为字符串
            return String(log);
        }
    }

    // 下载所有日志并生成新的 Blob
    saveLogsToBlob() {
        const logsContent = this.logs.join('\n');

        // 创建一个 Blob 对象
        const blob = new Blob([logsContent], { type: 'text/plain' });

        // 创建下载链接
        const link = document.createElement('a');

        // 为下载链接指定 Blob 对象的 URL
        link.href = URL.createObjectURL(blob);

        // 设置下载文件的名称
        link.download = `logs${this.logCount}.txt`;

        // 模拟点击链接，开始下载
        link.click();

        // 释放 URL 对象
        URL.revokeObjectURL(link.href);
    }
}

