#!/usr/bin/env python
# -*- coding: utf-8 -*-

import requests
import json
import logging

# 配置日志
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(name)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

class DYLoveProtocol:
    """抖音点赞协议实现"""
    
    def __init__(self):
        # 默认签名服务器地址
        self.sign_server = "http://47.122.2.184:19823/frfaf4343daf/sign"
        
        # 默认设备信息和token，来自签名测试文件
        self.default_device_info = {
            "device_id": "2091719857561258",
            "iid": "2091719857565354",
            "cdid": "E2953BAA-7DF0-497A-8877-91A877844A0F",
            "app_version": "32.8.0",
            "build_number": "328013",
            "os_version": "13.5",
            "device_type": "iPhone12,8",
            "aid": "1128",
            "channel": "App Store"
        }
        
        # 默认请求头信息
        self.default_headers = {
            "Host": "api3-normal-m-hj.amemv.com",
            "Cookie": "passport_csrf_token=9826f8124bac5292b11493fb3d3a9f5b; passport_csrf_token_default=9826f8124bac5292b11493fb3d3a9f5b; ticket_guard_has_set_public_key=1; d_ticket=a096bf60a72af434b37ad491497ebfe43c3b4; passport_mfa_token=CjyFTJikb6VqOkL5PtLHcRuRL%2BC31IUu9hFESBxmTNVOjtF%2BxVYjfTYhJy19HNZwUWuAQAgh9DSbEGG0ZCsaSgo8AAAAAAAAAAAAAE6%2FYTjeobm%2FgqUKF53dsARDIpWVQhzUdP3vX4igFTTC7BAetzip6TWinuM5MX0vP6DyEK3p6w0Y9rHRbCACIgED1BrUkQ%3D%3D; is_staff_user=false; multi_sids=104508686021%3A67f13ff5401e6cd77afcef7d6eb6c675%7C4094224458067131%3A7614b500b207ad9a6836e026379f67ec; n_mh=9dLb6CZOvasdw437-JSA3bubTSMYFtEpceFKBP2bmsA; odin_tt=21990540d9494d57202b3fd0b3eae2068004b5e7402447f43ac35598cd9f11ad97e7fc4224e9262b9c53c591fd55cda27aba092328e0c926fff54d2fcb31c9d660fe3adff34abdc1c0ed8f58f878b20c; passport_assist_user=CkGmQpATBDDht4Umy95cJsUfR7QuJw8Qph46hn2prBGmN27o-2MzJ7I9TqvqWwHdyVY2d310kynb5AjRVofY68oFjRpKCjwAAAAAAAAAAAAATr_VfxYSECaGcRIWR7kpaKoy9Dkee77Gcj_6yVAsDjaJS7sGAFoDWgEYia8yTVtGrGQQpenrDRiJr9ZUIAEiAQOEWSo_; sessionid=7614b500b207ad9a6836e026379f67ec; sessionid_ss=7614b500b207ad9a6836e026379f67ec; sid_guard=7614b500b207ad9a6836e026379f67ec%7C1741783830%7C5184000%7CSun%2C+11-May-2025+12%3A50%3A30+GMT; sid_tt=7614b500b207ad9a6836e026379f67ec; uid_tt=cc74674db38b0f89e71a2ed827c3a037; uid_tt_ss=cc74674db38b0f89e71a2ed827c3a037; install_id=2091719857565354; ttreq=1$ca0fc9f55ebf3b51e4d358291db79ed0c1beb321; store-region=cn-sh; store-region-src=uid",
            "bd-ticket-guard-version": "3",
            "bd-ticket-guard-client-cert": "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUI3RENDQVpLZ0F3SUJBZ0lVVy9HY0xBdGxhVTY0a3k2SEdCcmZQSlZXemxVd0NnWUlLb1pJemowRUF3SXcKTVRFTE1Ba0dBMVVFQmhNQ1EwNHhJakFnQmdOVkJBTU1HWFJwWTJ0bGRGOW5kV0Z5WkY5allWOWxZMlJ6WVY4eQpOVFl3SGhjTk1qVXdNVEUxTVRNd01EVTNXaGNOTXpVd01URTFNakV3TURVM1dqQWFNUmd3RmdZRFZRUURFdzlpClpDMTBhV05yWlhRdFozVmhjbVF3V1RBVEJnY3Foa2pPUFFJQkJnZ3Foa2pPUFFNQkJ3TkNBQVFlZURpVUpob0wKWU9RVXpqRmZrbWZyeWorQkxRRlRxMnhWZ0pzWGxzblBxR0xZMXAwNDNBWWlHYmF6YndzQVhWLzNXNGRGcEpwSwpVNEQ3RjNWbUxQaWJvNEdlTUlHYk1BNEdBMVVkRHdFQi93UUVBd0lGb0RBeEJnTlZIU1VFS2pBb0JnZ3JCZ0VGCkJRY0RBUVlJS3dZQkJRVUhBd0lHQ0NzR0FRVUZCd01EQmdnckJnRUZCUWNEQkRBcEJnTlZIUTRFSWdRZ3l3bkQKQml4TEM3Ylk1V2lYeG1QR2JaR1JiQW9ncXVGQkNpOE1lMGlOQXNvd0t3WURWUjBqQkNRd0lvQWdNcVZuNm81awpTQktOekU1TlFIdHpGSnRIYlZONnBOR0ExM21VbDNzaVI0TXdDZ1lJS29aSXpqMEVBd0lEU0FBd1JRSWhBTEs4ClJaVXp4V09qNDBEYnA4czgrWEZ4TDg3cEVQa0tNNURCKzV1bXRZcWFBaUFaa3VYS0s3UWNtVWJTUTJCbzNjZncKcWYyMkw4Z2cwc25STmpzODhoYkxXUT09Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K",
            "bd-ticket-guard-iteration-version": "3",
            "bd-ticket-guard-client-data": "eyJyZXFfY29udGVudCI6InRpY2tldCxwYXRoLHRpbWVzdGFtcCIsInJlcV9zaWduX3JlZSI6IjBQUUhLM2dOMWZtc3JMdjJDd0VxOHZiXC9LeHpzMVNFSG1kXC90WXh0ZzZqdz0iLCJ0aW1lc3RhbXAiOjE3NDE3ODY3MTUsInJlcV9zaWduIjoiTUVRQ0lIYWVuN25hcnBmSHN0VVhUMDRvc1ZQcW96WDExdFVCQzFGaUdGRTFYQjhuQWlBZGtvK1A3TGJmWHM0RDFWd1lQbFYraG5FanVCVTNNT1NzK3RHSFRmY3JmZz09IiwidHNfc2lnbiI6InRzLjEuYjI3YTNiMDIxOTY0ZWJjYmJkMTdjMTk1NzU1NWM3OGUzNWMxNzllNTNjYmRlMzU3OTQ1OTY3OTE2YzBlMTU2YWM0ZmJlODdkMjMxOWNmMDUzMTg2MjRjZWRhMTQ5MTFjYTQwNmRlZGJlYmVkZGIyZTMwZmNlOGQ0ZmEwMjU3NWQiLCJ0c19zaWduX3JlZSI6InRzLjEuMTJiNjY2NTNjMWFkNjlkOGY0ZjdkYWFmMTA5MTFhYzdiYmU3NjdjOTMwNGRlMDU1NjQ4NTA3MTA1MGQzNjNlMGM0ZmJlODdkMjMxOWNmMDUzMTg2MjRjZWRhMTQ5MTFjYTQwNmRlZGJlYmVkZGIyZTMwZmNlOGQ0ZmEwMjU3NWQifQ==",
            "sdk-version": "2",
            "content-type": "application/x-www-form-urlencoded",
            "x-tt-token": "007614b500b207ad9a6836e026379f67ec03d6afbda17d0b9b2d4594c2358ed50dc640776c5df381d5854e2a56037c89ed085c7cf6e0c6a72d3f77dbb4379c2cf1cd281887a0d9c9b895b884ee403a4d60242db11d0b361f4e38b0438dbfac3c4fa5f--0a490a209826978c5e740c969560206b53ba60cf7acd2244e8bd5939dab2a93ab5153b4c122041e5e5bc49feb571db73ceaddb5125ebf540e3f785ce97b413c813ee68c786c718f6b4d309-3.0.0",
            "user-agent": "Aweme 32.8.0 rv:328013 (iPhone; iOS 13.5; zh_CN) Cronet",
            "x-vc-bdturing-sdk-version": "3.7.2",
            "passport-sdk-version": "6.0.0-alpha.55",
            "x-tt-dt": "AAAVIX2DUYH6APH76KNLZ57PM7XMIJWDV2LAA6ERAQOLJOODVSLCEQZQMCOAV6VQIJGHIIA75KHLKCQ4RIZMF4TVRF6XMLJPY2KSKCTQXIHGGSYPBPQBRPD5NFOSHN26IUU4JZOTQFXZWGAZ6OYQN7I",
            "bd-ticket-guard-ree-public-key": "BMGHWxhIDsci558bioA2kLADlZQhPqylmQlE318Qz6ZS/krasCDzMKM8aGOd/8yS5iVRt3nHUmw2SBtFga09sh0=",
            "x-tt-passport-mfa-token": "CjyFTJikb6VqOkL5PtLHcRuRL+C31IUu9hFESBxmTNVOjtF+xVYjfTYhJy19HNZwUWuAQAgh9DSbEGG0ZCsaSgo8AAAAAAAAAAAAAE6/YTjeobm/gqUKF53dsARDIpWVQhzUdP3vX4igFTTC7BAetzip6TWinuM5MX0vP6DyEK3p6w0Y9rHRbCACIgED1BrUkQ==",
            "x-tt-token-supplement": "030e8703c950706ee3747198f486676be4b9f6bf77775835304c9abbf35c447cecbbc3e3d51ea773eb9a2ff94d8e02ca17ddd0c3fd84d3db42d5b3076859ac9fe38",
            "x-ss-stub": "D98588B0A401ECBF32AD0EA93D9A380A",
            "x-tt-store-region": "cn-sh",
            "x-tt-store-region-src": "uid",
            "x-tt-request-tag": "s=1",
            "x-ss-dp": "1128",
            "x-tt-trace-id": "00-8a9417b80d76e687b145eaad1cc80468-8a9417b80d76e68701"
        }
        
        # 默认token信息
        self.default_token = "007614b500b207ad9a6836e026379f67ec03d6afbda17d0b9b2d4594c2358ed50dc640776c5df381d5854e2a56037c89ed085c7cf6e0c6a72d3f77dbb4379c2cf1cd281887a0d9c9b895b884ee403a4d60242db11d0b361f4e38b0438dbfac3c4fa5f--0a490a209826978c5e740c969560206b53ba60cf7acd2244e8bd5939dab2a93ab5153b4c122041e5e5bc49feb571db73ceaddb5125ebf540e3f785ce97b413c813ee68c786c718f6b4d309-3.0.0"
    
    def get_signature(self, aweme_id, device_info=None):
        """
        获取签名信息
        
        Args:
            aweme_id: 视频ID
            device_info: 设备信息，默认使用初始化时的设备信息
        
        Returns:
            dict: 包含签名信息的字典
        """
        if device_info is None:
            device_info = self.default_device_info
        
        # 构建URL，使用与签名测试文件中相同的格式
        base_url = "https://api3-normal-m-hj.amemv.com/aweme/v1/commit/item/digg/"
        url_params = {
            "package": "com.ss.iphone.ugc.Aweme",
            "need_personal_recommend": "1",
            "version_code": device_info.get("app_version", "32.8.0"),
            "js_sdk_version": "3.55.0.11",
            "tma_jssdk_version": "3.55.0.11",
            "app_name": "aweme",
            "app_version": device_info.get("app_version", "32.8.0"),
            "device_id": device_info.get("device_id", "2091719857561258"),
            "channel": device_info.get("channel", "App Store"),
            "mcc_mnc": "",
            "aid": device_info.get("aid", "1128"),
            "minor_status": "0",
            "screen_width": "750",
            "klink_egdi": "AAKMSLLBDWo7Ac0Rrd7QWBlLvghhARqyBT9nUMJ9sdZ9pR1mnnHFH-GV",
            "cdid": device_info.get("cdid", "E2953BAA-7DF0-497A-8877-91A877844A0F"),
            "os_api": "18",
            "ac": "wifi",
            "os_version": device_info.get("os_version", "13.5"),
            "appTheme": "light",
            "is_guest_mode": "0",
            "device_platform": "iphone",
            "build_number": device_info.get("build_number", "328013"),
            "iid": device_info.get("iid", "2091719857565354"),
            "is_vcd": "1",
            "device_type": device_info.get("device_type", "iPhone12,8")
        }
        
        # 构造URL
        url = base_url + "?" + "&".join([f"{k}={v}" for k, v in url_params.items()])
        
        # 构建签名请求数据
        sign_data = {
            "url": url,
            "stub": "D98588B0A401ECBF32AD0EA93D9A380A",  # 这个值在测试文件中是固定的
            "did": device_info.get("device_id", "2091719857561258"),
            "sdiToken": "",
            "lanusk": "",
            "lanusv": "",
            "app_ver": device_info.get("app_version", "32.8.0"),
            "code": "D53gAP",  # 这个值在测试文件中是固定的
            "os_ver": device_info.get("os_version", "13.5"),
            "sdk_ver": "v04.06.00-ml-iOS",
            "sdk_ver_code": "67502081"
        }
        
        try:
            logger.info(f"获取签名，请求数据: {sign_data}")
            response = requests.post(self.sign_server, json=sign_data)
            response.raise_for_status()
            
            sign_result = response.json()
            logger.info(f"获取签名成功: {sign_result}")
            
            return sign_result
        except Exception as e:
            logger.error(f"获取签名失败: {str(e)}")
            raise
    
    def like_video(self, aweme_id, token=None, device_info=None):
        """
        点赞视频
        
        Args:
            aweme_id: 视频ID
            token: 用户token，默认使用初始化时的token
            device_info: 设备信息，默认使用初始化时的设备信息
            
        Returns:
            dict: 点赞结果
        """
        if token is None:
            token = self.default_token
            
        if device_info is None:
            device_info = self.default_device_info
        
        # 获取签名
        sign_result = self.get_signature(aweme_id, device_info)
        
        # 构建请求头，复制默认头并更新签名信息
        headers = self.default_headers.copy()
        headers.update({
            "x-tt-token": token,
            "x-argus": sign_result.get("a", ""),
            "x-gorgon": sign_result.get("g", ""),
            "x-helios": sign_result.get("h", ""),
            "x-khronos": sign_result.get("k", ""),
            "x-ladon": sign_result.get("l", ""),
            "x-medusa": sign_result.get("m", "")
        })
        
        # 构建URL
        base_url = "https://api3-normal-m-hj.amemv.com/aweme/v1/commit/item/digg/"
        url_params = {
            "package": "com.ss.iphone.ugc.Aweme",
            "need_personal_recommend": "1",
            "version_code": device_info.get("app_version", "32.8.0"),
            "js_sdk_version": "3.55.0.11",
            "tma_jssdk_version": "3.55.0.11",
            "app_name": "aweme",
            "app_version": device_info.get("app_version", "32.8.0"),
            "device_id": device_info.get("device_id", "2091719857561258"),
            "channel": device_info.get("channel", "App Store"),
            "mcc_mnc": "",
            "aid": device_info.get("aid", "1128"),
            "minor_status": "0",
            "screen_width": "750",
            "klink_egdi": "AAKMSLLBDWo7Ac0Rrd7QWBlLvghhARqyBT9nUMJ9sdZ9pR1mnnHFH-GV",
            "cdid": device_info.get("cdid", "E2953BAA-7DF0-497A-8877-91A877844A0F"),
            "os_api": "18",
            "ac": "wifi",
            "os_version": device_info.get("os_version", "13.5"),
            "appTheme": "light",
            "is_guest_mode": "0",
            "device_platform": "iphone",
            "build_number": device_info.get("build_number", "328013"),
            "iid": device_info.get("iid", "2091719857565354"),
            "is_vcd": "1",
            "device_type": device_info.get("device_type", "iPhone12,8")
        }
        
        url = base_url + "?" + "&".join([f"{k}={v}" for k, v in url_params.items()])
        
        # 构建请求体数据
        data = {
            "aweme_id": aweme_id,
            "channel_id": "0",
            "enter_from": "homepage_hot",
            "friend_recommend": "0",
            "is_commerce": "0",
            "nearby_level": "0",
            "previous_page": "homepage_hot",
            "type": "1"
        }
        
        try:
            logger.info(f"发送点赞请求，视频ID: {aweme_id}")
            response = requests.post(url, headers=headers, data=data)
            response.raise_for_status()
            
            result = response.json()
            logger.info(f"点赞结果: {result}")
            
            return result
        except Exception as e:
            logger.error(f"点赞请求失败: {str(e)}")
            raise

# 示例用法
if __name__ == "__main__":
    # 创建点赞协议实例
    dy_love = DYLoveProtocol()
    
    # 点赞示例视频
    try:
        # result = dy_love.like_video("7478180544006737204")
        result = dy_love.like_video("7453475040403705147")
        print(f"点赞结果: {json.dumps(result, ensure_ascii=False, indent=2)}")
    except Exception as e:
        print(f"点赞失败: {str(e)}") 