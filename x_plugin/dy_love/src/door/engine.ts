import path from 'path';
import fs from 'fs';
import { Browser, chromium, devices, BrowserContext, Page, Route, Request, Response } from 'playwright';

export abstract class DoorEngine<T = any> {
    private chromePath: string | undefined;
    browser: Browser | undefined;
    context: BrowserContext | undefined;
    public resourceId: string;
    public headless: boolean = true;
    page: Page | undefined;

    constructor(resourceId: string, headless: boolean = true, chromePath: string = "") {
        this.resourceId = resourceId;
        if (chromePath) {
            this.chromePath = chromePath;
        } else {
            this.chromePath = this.getChromePath();
        }
        this.headless = headless;
    }

    getChromePath(): string | undefined {
        return process.env.CHROME_PATH;
    }

    public async init(url: string | undefined = undefined): Promise<Page | undefined> {
        if (this.browser) {
            return undefined;
        }
        
        try {
            this.browser = await this.createBrowser();
            this.context = await this.createContext();
            if (!this.context) {
                return undefined;
            }
            
            const page = await this.context.newPage();
            
            // 不设置视口大小，使用窗口大小
            console.log("正在设置页面为全屏模式...");
            
            // 设置额外的页面选项，使用try-catch包裹以避免错误
            try {
                await page.addInitScript(() => {
                    // 禁用某些可能导致性能问题的功能
                    window.onbeforeunload = null;
                    // 禁用右键菜单
                    document.addEventListener('contextmenu', e => e.preventDefault(), false);
                    
                    // 尝试设置全屏
                    try {
                        const requestFullscreen = document.documentElement.requestFullscreen 
                            || (document.documentElement as any).webkitRequestFullscreen 
                            || (document.documentElement as any).mozRequestFullScreen 
                            || (document.documentElement as any).msRequestFullscreen;
                        
                        if (requestFullscreen) {
                            requestFullscreen.call(document.documentElement);
                        }
                    } catch (e) {
                        console.error("设置全屏失败:", e);
                    }
                });
            } catch (error) {
                console.error("设置页面初始化脚本失败:", error);
            }
            
            if (url) {
                try {
                    console.log(`正在访问: ${url}`);
                    await page.goto(url, { 
                        timeout: 60000, 
                        waitUntil: 'domcontentloaded' 
                    });
                    console.log(`成功加载页面: ${url}`);
                    
                    // 等待页面稳定，使用try-catch包裹以避免错误
                    try {
                        await page.waitForLoadState('networkidle', { timeout: 30000 }).catch(() => {
                            console.log('等待网络稳定超时，继续执行');
                        });
                    } catch (error) {
                        console.error("等待页面稳定失败:", error);
                    }
                    
                    // 尝试调整窗口大小，确保内容可见
                    try {
                        // 使用page而不是context来获取窗口大小
                        const windowSize = await page.evaluate(() => ({
                            outerWidth: window.outerWidth,
                            outerHeight: window.outerHeight,
                            devicePixelRatio: window.devicePixelRatio
                        }));
                        
                        console.log(`当前窗口大小: ${JSON.stringify(windowSize)}`);
                        
                        // 尝试最大化窗口
                        await page.evaluate(() => {
                            window.moveTo(0, 0);
                            window.resizeTo(screen.width, screen.height);
                            
                            // 尝试进入全屏模式
                            try {
                                const requestFullscreen = document.documentElement.requestFullscreen 
                                    || (document.documentElement as any).webkitRequestFullscreen 
                                    || (document.documentElement as any).mozRequestFullScreen 
                                    || (document.documentElement as any).msRequestFullscreen;
                                
                                if (requestFullscreen) {
                                    requestFullscreen.call(document.documentElement);
                                }
                            } catch (e) {
                                console.error("设置全屏失败:", e);
                            }
                        });
                        console.log("已尝试调整窗口大小为全屏");
                    } catch (error) {
                        console.error("调整窗口大小失败:", error);
                    }
                } catch (error) {
                    console.error(`页面加载失败: ${error}`);
                    // 即使加载失败也返回页面对象，让用户可以继续操作
                }
            }
            
            this.page = page;
            return page;
        } catch (error) {
            console.error("浏览器初始化失败:", error);
            // 清理资源
            await this.release();
            return undefined;
        }
    }

    public async release() {
        await this.closePage();
        await this.closeContext();
        await this.closeBrowser();
    }

    public async closePage() {
        try {
            if (this.page && !this.page.isClosed()) {
                await this.page.close();
            }
        } catch (e) {
            console.error("closePage error", e);
        }
    }

    public async closeContext() {
        try {
            if (this.context) {
                await this.context.close();
                this.context = undefined;
            }
        } catch (e) {
            console.error("closeContext error", e);
        }
    }

    public async closeBrowser() {
        try {
            if (this.browser) {
                await this.browser.close();
                this.browser = undefined;
            }
        } catch (e) {
            console.error("closeBrowser error", e);
        }
    }

    getKey() {
        return `${this.getNamespace()}_${this.resourceId}`;
    }

    abstract getNamespace(): string;

    async createContext() {
        if (!this.browser) {
            return undefined;
        }
        
        // 设置更大的视口尺寸，避免界面被覆盖
        const context = await this.browser.newContext({
            viewport: { width: 1920, height: 1080 }, // 使用固定视口大小
            userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36',
            deviceScaleFactor: 1,
            javaScriptEnabled: true,
            bypassCSP: true,
            ignoreHTTPSErrors: true,
            screen: { width: 1920, height: 1080 },
            // 添加额外的权限
            permissions: ['geolocation', 'notifications'],
            // 设置颜色方案
            colorScheme: 'light',
            // 设置语言
            locale: 'zh-CN',
            // 设置时区
            timezoneId: 'Asia/Shanghai',
            // 设置地理位置
            geolocation: { latitude: 39.9042, longitude: 116.4074 },
            // 设置减少动画
            reducedMotion: 'no-preference',
            // 设置强制颜色
            forcedColors: 'none',
        });
        
        // 设置默认超时
        context.setDefaultTimeout(60000);
        
        return context;
    }

    async createBrowser() {
        const options: any = {
            headless: this.headless,
            args: [
                '--disable-gpu',
                '--disable-dev-shm-usage',
                '--disable-setuid-sandbox',
                '--no-sandbox',
                '--window-size=1920,1080',
                '--start-maximized', // 最大化窗口
                '--disable-extensions', // 禁用扩展
                '--disable-popup-blocking', // 禁用弹窗拦截
                '--disable-infobars', // 禁用信息栏
                '--window-position=0,0', // 设置窗口位置
                '--disable-notifications', // 禁用通知
                '--disable-translate', // 禁用翻译
                '--disable-features=TranslateUI', // 禁用翻译UI
                '--disable-default-apps', // 禁用默认应用
                '--no-first-run', // 禁用首次运行
                '--no-default-browser-check', // 禁用默认浏览器检查
                '--disable-background-networking', // 禁用后台网络
                '--disable-sync', // 禁用同步
                '--disable-background-timer-throttling', // 禁用后台计时器节流
                '--disable-backgrounding-occluded-windows', // 禁用背景遮挡窗口
                '--disable-breakpad', // 禁用崩溃报告
                '--disable-component-extensions-with-background-pages', // 禁用带有后台页面的组件扩展
                '--disable-features=BackForwardCache', // 禁用前进后退缓存
                '--disable-ipc-flooding-protection', // 禁用IPC洪水保护
                '--enable-automation', // 启用自动化
                '--password-store=basic', // 密码存储基本
                '--use-mock-keychain', // 使用模拟钥匙链
            ],
            ignoreDefaultArgs: ['--enable-automation'], // 忽略默认参数，避免浏览器被识别为自动化工具
        };

        if (this.chromePath) {
            options.executablePath = this.chromePath;
        }

        console.log("正在启动浏览器，使用以下选项:", JSON.stringify(options, null, 2));
        const browser = await chromium.launch(options);
        
        return browser;
    }
} 