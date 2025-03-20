import * as dotenv from 'dotenv';
dotenv.config(); // 加载 .env 文件中的环境变量

import { chromium, Browser, BrowserContext, Page } from 'playwright';

async function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// 全局变量，用于存储浏览器、上下文和页面对象，以便后续操作
let browser: Browser | null = null;
let context: BrowserContext | null = null;
let page: Page | null = null;

async function startDouyinLogin() {
    console.log("开始启动抖音登录流程...");
    
    try {
        // 直接使用playwright打开浏览器
        console.log("正在启动浏览器...");
        
        // 设置浏览器启动选项 - 减少不必要的参数
        browser = await chromium.launch({
            headless: false,
            executablePath: process.env.CHROME_PATH,
            args: [
                '--disable-accelerated-2d-canvas', '--disable-webgl', '--disable-software-rasterizer',
                '--no-sandbox', // 取消沙箱，某些网站可能会检测到沙箱模式
                '--disable-setuid-sandbox',
                '--disable-webrtc-encryption',
                '--disable-webrtc-hw-decoding',
                '--disable-webrtc-hw-encoding',
                '--disable-extensions-file-access-check',
                '--disable-blink-features=AutomationControlled',  // 禁用浏览器自动化控制特性
              ]
        });
        console.log("浏览器启动成功");
        
        // 创建上下文 - 不设置viewport，使用窗口默认大小
        console.log("正在创建浏览器上下文...");
        context = await browser.newContext({
            viewport: null, // 不设置视口大小，使用窗口大小
            ignoreHTTPSErrors: true, // 忽略HTTPS错误
        });
        console.log("浏览器上下文创建成功");
        
        // 设置超时
        context.setDefaultTimeout(60000);
        
        // 创建页面
        console.log("正在创建新页面...");
        page = await context.newPage();
        console.log("新页面创建成功");
        
        // 减少监听器，只保留关键错误监听
        page.on('pageerror', err => console.error(`页面错误: ${err.message}`));
        
        // 尝试最大化窗口
        console.log("尝试最大化窗口...");
        try {
            // 在访问网站前先尝试最大化窗口
            await page.evaluate(() => {
                window.moveTo(0, 0);
                window.resizeTo(screen.width, screen.height);
            });
            console.log("窗口最大化完成");
        } catch (error) {
            console.error("窗口最大化失败:", error);
        }
        
        // 访问抖音网页版
        console.log("正在访问抖音网页版...");
        try {
            await page.goto("https://www.douyin.com/", { 
                timeout: 60000,
                waitUntil: 'domcontentloaded' // 使用domcontentloaded而不是networkidle，提高加载速度
            });
            console.log("抖音网页版加载成功");
        } catch (error) {
            console.error("抖音网页版加载失败:", error);
            console.log("尝试继续执行...");
        }
        
        // 等待页面加载 - 减少等待时间
        await sleep(5000);
        
        console.log("已打开抖音网页版，请在浏览器中完成登录操作");
        console.log("提示: 如果页面加载失败，请尝试刷新页面");
        console.log("      您也可以在浏览器中手动访问 https://www.douyin.com/");
        
        // 等待用户手动登录
        console.log("等待用户登录中...");
        console.log("请在浏览器中点击页面右上角的「点击登录」按钮完成登录流程");
        
        // 保持浏览器打开，直到用户手动关闭或程序结束
        console.log("浏览器将保持打开状态，直到您手动关闭或按Ctrl+C结束程序");
        console.log("登录成功后，您可以继续使用程序进行其他操作");
        
        // 等待用户输入命令
        await waitForUserCommand();
        
    } catch (error) {
        console.error("运行过程中出现错误:", error);
    }
}

// 等待用户输入命令
async function waitForUserCommand() {
    console.log("\n请输入命令:");
    console.log("1. 点赞视频 (输入视频链接)");
    console.log("2. 退出程序");
    
    // 这里可以添加读取用户输入的代码
    // 由于Node.js中读取用户输入需要额外的库，这里简化为等待一段时间
    
    // 模拟用户输入等待
    await sleep(60000);
    
    // 如果需要关闭浏览器
    if (browser) {
        console.log("正在关闭浏览器...");
        await browser.close();
        console.log("浏览器已关闭");
    }
}

// 点赞视频的函数
async function likeVideo(videoUrl: string) {
    if (!page) {
        console.error("浏览器未初始化，无法点赞视频");
        return;
    }
    
    try {
        console.log(`正在访问视频: ${videoUrl}`);
        await page.goto(videoUrl, { timeout: 60000, waitUntil: 'domcontentloaded' });
        
        // 等待页面加载 - 使用固定等待时间而不是networkidle
        await sleep(5000);
        
        // 查找点赞按钮
        console.log("正在查找点赞按钮...");
        const likeButtonSelectors = [
            '.like-icon',
            '.like-button',
            '[data-e2e="like-icon"]',
            'button:has-text("点赞")',
            'svg[aria-label="点赞"]',
            '.like',
            '.thumbs-up',
            '.heart-icon'
        ];
        
        let likeButtonFound = false;
        for (const selector of likeButtonSelectors) {
            try {
                const likeButton = await page.$(selector);
                if (likeButton) {
                    console.log(`找到点赞按钮: ${selector}`);
                    
                    // 滚动到按钮位置
                    await likeButton.scrollIntoViewIfNeeded();
                    console.log("已滚动到点赞按钮位置");
                    
                    // 等待一下确保按钮可见
                    await sleep(1000);
                    
                    // 点击按钮
                    await likeButton.click({ force: true });
                    console.log("已点击点赞按钮");
                    likeButtonFound = true;
                    break;
                }
            } catch (e) {
                console.error(`尝试点击 ${selector} 失败:`, e);
            }
        }
        
        if (!likeButtonFound) {
            console.log("未找到点赞按钮，尝试使用JavaScript查找并点击");
            
            // 使用JavaScript查找并点击点赞按钮
            const result = await page.evaluate(() => {
                // 查找所有可能的点赞按钮
                const allElements = document.querySelectorAll('button, svg, div, span, i');
                
                // 过滤出可能是点赞按钮的元素
                const likeElements = Array.from(allElements).filter(el => {
                    const text = el.textContent || '';
                    const className = el.className || '';
                    const id = el.id || '';
                    
                    return text.includes('赞') || 
                           text.includes('like') || 
                           text.includes('Like') || 
                           className.includes('like') || 
                           className.includes('Like') || 
                           className.includes('heart') || 
                           className.includes('Heart') || 
                           id.includes('like') || 
                           id.includes('Like');
                });
                
                if (likeElements.length > 0) {
                    // 尝试点击第一个找到的元素
                    try {
                        (likeElements[0] as HTMLElement).click();
                        return { success: true, message: "成功点击点赞按钮" };
                    } catch (e) {
                        return { success: false, message: `点击失败: ${e}` };
                    }
                }
                
                return { success: false, message: "未找到点赞按钮" };
            });
            
            console.log("JavaScript点击结果:", result);
        }
        
    } catch (error) {
        console.error("点赞视频过程中出现错误:", error);
    }
}

// 启动主函数
async function main() {
    await startDouyinLogin();
}

main().catch(error => {
    console.error("程序执行出错:", error);
    
    // 确保在出错时关闭浏览器
    if (browser) {
        browser.close().catch(() => {});
    }
}); 