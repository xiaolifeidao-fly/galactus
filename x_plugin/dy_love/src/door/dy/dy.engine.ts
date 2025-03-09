import { Page } from "playwright";
import { DoorEngine } from "../engine";

export class DyEngine<T> extends DoorEngine<T> {
    public getNamespace(): string {
        return "dy";
    }

    // 等待登录完成
    public async waitForLogin(page: Page): Promise<boolean> {
        try {
            console.log("开始检测登录状态...");
            
            // 尝试多种可能的登录成功标志
            const possibleSelectors = [
                '.avatar-img', // 用户头像
                '.login-success-element', // 登录成功元素
                '.nickname', // 用户昵称
                '.personal-tab', // 个人中心标签
                '[data-e2e="user-icon"]', // 用户图标
                '.dy-account-icon', // 抖音账号图标
                '.user-info', // 用户信息区域
                '.account-name', // 账号名称
                '.user-profile', // 用户资料
                '.login-btn.login-btn-disabled', // 登录按钮变为禁用状态
                '.user-avatar', // 用户头像
                '.user-name', // 用户名
                // 新增更多可能的选择器
                '.author-card', // 作者卡片
                '.profile-card', // 个人资料卡片
                '.user-card', // 用户卡片
                '.user-profile-card', // 用户资料卡片
                '.user-info-card', // 用户信息卡片
                '.user-header', // 用户头部
                '.user-panel', // 用户面板
                '.user-center', // 用户中心
                '.user-dropdown', // 用户下拉菜单
                '.user-menu', // 用户菜单
                '.user-action', // 用户操作
                '.user-operation', // 用户操作
                '.user-setting', // 用户设置
                '.user-config', // 用户配置
                '.user-control', // 用户控制
                '.user-manage', // 用户管理
                '.user-edit', // 用户编辑
                '.user-view', // 用户查看
                '.user-detail', // 用户详情
                '.user-info-detail', // 用户信息详情
                '.user-profile-detail', // 用户资料详情
                '.user-account-detail', // 用户账号详情
                '.user-account-info', // 用户账号信息
                '.user-account-profile', // 用户账号资料
                '.user-account-setting', // 用户账号设置
                '.user-account-config', // 用户账号配置
                '.user-account-control', // 用户账号控制
                '.user-account-manage', // 用户账号管理
                '.user-account-edit', // 用户账号编辑
                '.user-account-view', // 用户账号查看
                '.user-account-detail', // 用户账号详情
            ];
            
            console.log("等待登录成功标志出现...");
            
            // 创建一个Promise，当任何一个选择器匹配时就解析
            const loginPromise = new Promise<boolean>(async (resolve) => {
                let resolved = false;
                let checkUrlInterval: NodeJS.Timeout | null = null;
                let checkLoginButtonInterval: NodeJS.Timeout | null = null;
                let checkPageContentInterval: NodeJS.Timeout | null = null;
                let timeoutId: NodeJS.Timeout | null = null;
                
                // 定义一个函数来处理成功登录
                const handleSuccess = (reason: string) => {
                    if (!resolved) {
                        resolved = true;
                        console.log(`登录成功: ${reason}`);
                        
                        // 清理所有定时器
                        if (checkUrlInterval) clearInterval(checkUrlInterval);
                        if (checkLoginButtonInterval) clearInterval(checkLoginButtonInterval);
                        if (checkPageContentInterval) clearInterval(checkPageContentInterval);
                        if (timeoutId) clearTimeout(timeoutId);
                        
                        resolve(true);
                    }
                };
                
                // 定义一个函数来处理失败
                const handleFailure = (reason: string) => {
                    if (!resolved) {
                        resolved = true;
                        console.log(`登录失败: ${reason}`);
                        
                        // 清理所有定时器
                        if (checkUrlInterval) clearInterval(checkUrlInterval);
                        if (checkLoginButtonInterval) clearInterval(checkLoginButtonInterval);
                        if (checkPageContentInterval) clearInterval(checkPageContentInterval);
                        if (timeoutId) clearTimeout(timeoutId);
                        
                        resolve(false);
                    }
                };
                
                // 首先检查是否已经登录
                try {
                    for (const selector of possibleSelectors) {
                        try {
                            // 检查元素是否已经存在
                            const exists = await page.$(selector);
                            if (exists) {
                                handleSuccess(`检测到登录成功标志: ${selector}`);
                                return;
                            }
                        } catch (e) {
                            // 忽略错误，继续检查下一个选择器
                        }
                    }
                } catch (error) {
                    console.error("检查登录状态时出错:", error);
                }
                
                // 如果没有立即找到，设置监听器等待元素出现
                try {
                    const promises = possibleSelectors.map(selector => {
                        return page.waitForSelector(selector, { timeout: 60000 })
                            .then(() => {
                                handleSuccess(`检测到登录成功标志: ${selector}`);
                                return true;
                            })
                            .catch(() => false);
                    });
                    
                    // 任何一个Promise成功就算成功
                    Promise.any(promises)
                        .then(() => {
                            if (!resolved) {
                                handleSuccess("通过元素选择器检测到登录成功");
                            }
                        })
                        .catch(() => {
                            // 忽略错误
                        });
                } catch (error) {
                    console.error("设置元素监听器时出错:", error);
                }
                
                // 另一种检测方法：检查URL变化
                try {
                    checkUrlInterval = setInterval(async () => {
                        try {
                            if (resolved) return;
                            
                            const currentUrl = page.url();
                            // 如果URL包含用户相关的路径，可能表示已登录
                            if (currentUrl.includes('/user/') || 
                                currentUrl.includes('/profile/') || 
                                currentUrl.includes('/account/')) {
                                handleSuccess(`通过URL变化检测到登录成功: ${currentUrl}`);
                            }
                        } catch (e) {
                            // 忽略错误
                        }
                    }, 5000); // 每5秒检查一次
                } catch (error) {
                    console.error("设置URL检查定时器时出错:", error);
                }
                
                // 检查登录按钮是否消失
                try {
                    checkLoginButtonInterval = setInterval(async () => {
                        try {
                            if (resolved) return;
                            
                            // 检查登录按钮是否不存在
                            const loginButton = await page.$('.login-button, .login-btn, [data-e2e="login-button"], button:has-text("登录"), a:has-text("登录")');
                            if (!loginButton) {
                                // 再次确认是否有用户相关元素
                                const userElement = await page.$('.avatar-img, .nickname, .user-info, .user-avatar, .user-name');
                                if (userElement) {
                                    handleSuccess('检测到登录按钮消失且存在用户元素，确认登录成功');
                                }
                            }
                        } catch (e) {
                            // 忽略错误
                        }
                    }, 5000); // 每5秒检查一次
                } catch (error) {
                    console.error("设置登录按钮检查定时器时出错:", error);
                }
                
                // 检查页面内容变化
                try {
                    checkPageContentInterval = setInterval(async () => {
                        try {
                            if (resolved) return;
                            
                            // 检查页面内容是否包含用户相关信息
                            const pageContent = await page.content();
                            const userRelatedKeywords = [
                                '我的主页', '个人中心', '我的收藏', '我的喜欢', '我的关注',
                                '我的粉丝', '我的消息', '我的设置', '退出登录', '账号设置',
                                '个人资料', '账号安全', '隐私设置', '消息通知', '黑名单管理'
                            ];
                            
                            for (const keyword of userRelatedKeywords) {
                                if (pageContent.includes(keyword)) {
                                    handleSuccess(`检测到页面内容包含用户相关关键词: ${keyword}`);
                                    break;
                                }
                            }
                            
                            // 尝试执行JavaScript检查登录状态
                            const isLoggedIn = await page.evaluate(() => {
                                // 检查是否存在用户相关的全局变量或localStorage
                                return !!(
                                    localStorage.getItem('user') || 
                                    localStorage.getItem('userInfo') || 
                                    localStorage.getItem('token') || 
                                    localStorage.getItem('auth') ||
                                    document.cookie.includes('login') ||
                                    document.cookie.includes('user') ||
                                    document.cookie.includes('token') ||
                                    document.cookie.includes('auth')
                                );
                            });
                            
                            if (isLoggedIn) {
                                handleSuccess('通过JavaScript检测到登录状态');
                            }
                        } catch (e) {
                            // 忽略错误
                        }
                    }, 5000); // 每5秒检查一次
                } catch (error) {
                    console.error("设置页面内容检查定时器时出错:", error);
                }
                
                // 设置清理定时器的超时
                try {
                    timeoutId = setTimeout(() => {
                        handleFailure("登录检测超时");
                    }, 120000); // 2分钟超时
                } catch (error) {
                    console.error("设置超时定时器时出错:", error);
                }
            });
            
            return loginPromise;
        } catch (error) {
            console.log("等待登录过程中出错", error);
            return false;
        }
    }
} 