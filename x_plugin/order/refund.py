#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import requests
import time
import logging
from typing import Optional

# 配置日志
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler('refund.log'),
        logging.StreamHandler()
    ]
)

class RefundClient:
    def __init__(self, base_url: str = "http://111.180.188.251:50001"):
        self.base_url = base_url.rstrip('/')
        self.session = requests.Session()
        # 设置默认headers
        self.session.headers.update({
            'Accept': 'application/json, text/plain, */*',
            'Accept-Language': 'zh-CN,zh;q=0.9,en;q=0.8',
            'Connection': 'keep-alive',
            'Content-Length': '0',
            'Origin': 'http://111.180.188.251',
            'Referer': 'http://111.180.188.251/',
            'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36',
            'x-forwarded-for': '4.2.2.2'
        })
        self.token: Optional[str] = None
        self.login_user: Optional[str] = None

    def set_auth(self, token: str, login_user: str):
        """设置认证信息"""
        self.token = token
        self.login_user = login_user
        self.session.headers.update({
            'X-Token': token,
            'loginuser': login_user
        })

    def refund_order(self, order_id: int) -> bool:
        """退款单个订单"""
        if not self.token or not self.login_user:
            raise ValueError("Please set authentication first using set_auth()")

        url = f"{self.base_url}/kakrolot_web/orders/{order_id}/refund"
        try:
            response = self.session.post(url, verify=False)
            response.raise_for_status()
            logging.info(f"Order {order_id} refund success")
            return True
        except requests.exceptions.RequestException as e:
            logging.error(f"Order {order_id} refund failed: {str(e)}")
            return False

    def batch_refund(self, start_id: int, end_id: int, delay: float = 1.0):
        """批量退款订单"""
        success_count = 0
        fail_count = 0
        
        for order_id in range(start_id, end_id + 1):
            if self.refund_order(order_id):
                success_count += 1
            else:
                fail_count += 1
            
            if order_id < end_id:  # 不是最后一个订单时才延迟
                time.sleep(delay)
        
        logging.info(f"Batch refund completed. Success: {success_count}, Failed: {fail_count}")
        return success_count, fail_count

def main():
    # 使用示例
    client = RefundClient()
    
    # 设置认证信息（这里需要替换为实际的token和login_user）
    client.set_auth(
        token="809cd0f7-a099-453f-9d5f-fc8a6abf0112",
        login_user="oRgrw5FwQIx63ieHwe6EmXHdgWhI"
    )
    
    # 批量退款
    start_id = 569
    end_id = 2069
    success, failed = client.batch_refund(start_id, end_id, delay=1.0)
    
    print(f"退款完成，成功：{success}，失败：{failed}")

if __name__ == "__main__":
    main()
