#!/usr/bin/env python
# -*- coding: utf-8 -*-

"""
抖音点赞协议使用示例
"""

import sys
import json
from dy_love import DYLoveProtocol

def main():
    """主函数"""
    # 检查是否提供了视频ID参数
    if len(sys.argv) < 2:
        print("使用方法: python example.py 视频ID")
        return
    
    # 获取视频ID
    aweme_id = sys.argv[1]
    
    print(f"准备给视频 {aweme_id} 点赞...")
    
    # 创建点赞协议实例
    dy_love = DYLoveProtocol()
    
    try:
        # 执行点赞
        result = dy_love.like_video(aweme_id)
        
        # 打印结果
        print("点赞结果:")
        print(json.dumps(result, ensure_ascii=False, indent=2))
        
        # 判断是否成功
        if result.get('status_code') == 0:
            print("点赞成功!")
        else:
            print("点赞失败!")
    except Exception as e:
        print(f"点赞过程中出现错误: {str(e)}")

if __name__ == "__main__":
    main() 