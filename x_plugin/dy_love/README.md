# dy_love

这是一个使用Playwright自动化打开抖音网页版并处理登录流程的项目。

## 功能

- 自动打开浏览器并访问抖音网页版
- 等待用户手动完成登录操作
- 检测登录成功状态

## 环境要求

- Node.js 14+
- TypeScript
- Playwright

## 安装

1. 克隆项目
2. 安装依赖

```bash
cd x_plugin/dy_love
npm install
# 或者使用 pnpm
pnpm install
```

3. 配置环境变量

在项目根目录创建`.env`文件，内容如下：

```
CHROME_PATH=/Applications/Google Chrome.app/Contents/MacOS/Google Chrome
HEADLESS=false
```

注意：
- `CHROME_PATH`是Chrome浏览器的可执行文件路径，根据你的系统进行调整
- `HEADLESS=false`表示以有界面模式运行浏览器，设置为`true`则为无界面模式

## 运行

```bash
npm run dev
# 或者使用 pnpm
pnpm run dev
```

## 使用说明

1. 运行程序后，会自动打开Chrome浏览器并访问抖音网页版
2. 在打开的浏览器中手动完成登录操作
3. 程序会自动检测登录状态，登录成功后会在控制台输出提示信息
4. 登录成功后，浏览器会保持打开一段时间，然后自动关闭 

npx ts-node -r tsconfig-paths/register src/main.ts