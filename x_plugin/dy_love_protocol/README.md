# 抖音点赞协议

这是一个实现抖音点赞功能的协议库，通过请求签名服务获取签名，然后替换到请求中实现点赞功能。

## 文件说明

- `dy_love.py`: 主要实现抖音点赞功能的Python模块

## 使用方法

### 基本用法

```python
from dy_love import DYLoveProtocol

# 创建点赞协议实例
dy_love = DYLoveProtocol()

# 点赞视频
result = dy_love.like_video("视频ID")
print(result)
```

### 自定义设备信息和token

```python
from dy_love import DYLoveProtocol

# 创建点赞协议实例
dy_love = DYLoveProtocol()

# 自定义设备信息
device_info = {
    "device_id": "你的设备ID",
    "iid": "安装ID",
    "cdid": "设备CD ID",
    "app_version": "32.8.0",
    "build_number": "328013",
    "os_version": "13.5",
    "device_type": "iPhone12,8",
    "aid": "1128",
    "channel": "App Store"
}

# 自定义token
token = "你的token"

# 点赞视频
result = dy_love.like_video("视频ID", token=token, device_info=device_info)
print(result)
```

## 参数说明

### DYLoveProtocol类

初始化参数:
- 无需参数，使用默认配置

### like_video方法

参数:
- `aweme_id`: 必填，视频ID
- `token`: 可选，用户token，默认使用初始化时的token
- `device_info`: 可选，设备信息，默认使用初始化时的设备信息

返回值:
- 字典类型，点赞请求的响应结果

## 实现原理

1. 构建原始点赞请求参数
2. 请求签名服务获取签名
3. 将签名添加到请求头中
4. 发送最终的点赞请求

## 注意事项

- 此协议依赖于签名服务，默认使用`http://47.122.2.184:19823/frfaf4343daf/sign`作为签名服务器
- 默认使用的token和设备信息来自测试文件，实际使用时建议替换为自己的参数

python -m venv myenv  # myenv为自定义环境名称
source myenv/bin/activate 
deactivate 
