#!/usr/bin/env python
# -*- coding: utf-8 -*-

import requests
import json
import logging
import re
import time
import random

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
            
            # 先记录原始响应内容，无论是否成功
            logger.info(f"签名服务器响应状态码: {response.status_code}")
            logger.info(f"签名服务器响应头: {dict(response.headers)}")
            logger.info(f"签名服务器原始响应内容: {response.text}")
            
            response.raise_for_status()
            
            sign_result = response.json()
            logger.info(f"获取签名成功: {sign_result}")
            
            return sign_result
        except requests.exceptions.RequestException as e:
            # 增强异常日志，记录完整响应内容
            error_msg = f"获取签名失败: {str(e)}"
            try:
                if hasattr(e, 'response') and e.response:
                    error_msg += f"\n响应状态码: {e.response.status_code}"
                    error_msg += f"\n响应内容: {e.response.text}"
            except Exception:
                pass
            logger.error(error_msg)
            raise
        except Exception as e:
            logger.error(f"获取签名失败（非请求异常）: {str(e)}")
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
            
            # 先记录原始响应内容，无论是否成功
            logger.info(f"点赞请求响应状态码: {response.status_code}")
            logger.info(f"点赞请求响应头: {dict(response.headers)}")
            logger.info(f"点赞请求原始响应内容: {response.text}")
            
            response.raise_for_status()
            
            # 尝试解析JSON前先检查响应内容是否为空
            if not response.text.strip():
                logger.warning("点赞请求返回空响应")
                return {"status_code": -1, "message": "空响应"}
            
            result = response.json()
            logger.info(f"点赞结果: {result}")
            
            return result
        except requests.exceptions.JSONDecodeError as e:
            # JSON解析错误，记录原始响应内容
            error_msg = f"点赞请求返回数据解析失败: {str(e)}"
            error_msg += f"\n响应状态码: {response.status_code}"
            error_msg += f"\n响应头信息: {dict(response.headers)}"
            error_msg += f"\n原始响应内容: {response.text}"
            logger.error(error_msg)
            # 返回错误信息而不是抛出异常，让程序继续运行
            return {"status_code": -1, "message": f"JSON解析失败: {str(e)}", "raw_response": response.text}
        except requests.exceptions.RequestException as e:
            # 请求异常，记录完整响应内容
            error_msg = f"点赞请求网络异常: {str(e)}"
            try:
                if hasattr(e, 'response') and e.response:
                    error_msg += f"\n响应状态码: {e.response.status_code}"
                    error_msg += f"\n响应头信息: {dict(e.response.headers)}"
                    error_msg += f"\n响应内容: {e.response.text}"
            except Exception:
                pass
            logger.error(error_msg)
            raise
        except Exception as e:
            logger.error(f"点赞请求其他异常: {str(e)}")
            raise
    
    def get_task(self, uid, sec_uid, code="MY_LOVE", uid_type="DY"):
        """
        获取点赞任务
        
        Args:
            uid: 用户ID
            sec_uid: 安全用户ID
            code: 任务代码，默认为MY_LOVE
            uid_type: 用户类型，默认为DY
            
        Returns:
            dict: 任务数据
        """
        url = f"http://111.180.188.251:9999/batch/tasks/get?uid={uid}&uidType={uid_type}&code={code}&secUid={sec_uid}"
        headers = {
            "PUB_TOKEN": "b07a1a53-8c8a-4b15-b975-c023c67d6b8a"
        }
        
        try:
            logger.info(f"获取任务，用户ID: {uid}")
            response = requests.get(url, headers=headers)
            
            # 先记录原始响应内容，无论是否成功
            logger.info(f"获取任务响应状态码: {response.status_code}")
            logger.info(f"获取任务响应头: {dict(response.headers)}")
            logger.info(f"获取任务原始响应内容: {response.text}")
            
            response.raise_for_status()
            
            # 尝试解析JSON前先检查响应内容是否为空
            if not response.text.strip():
                logger.warning("获取任务返回空响应")
                return {"status": -1, "message": "空响应"}
            
            result = response.json()
            logger.info(f"获取任务结果: {result}")
            
            return result
        except requests.exceptions.JSONDecodeError as e:
            # JSON解析错误，记录原始响应内容
            error_msg = f"获取任务返回数据解析失败: {str(e)}"
            error_msg += f"\n响应状态码: {response.status_code}"
            error_msg += f"\n响应头信息: {dict(response.headers)}"
            error_msg += f"\n原始响应内容: {response.text}"
            logger.error(error_msg)
            # 返回错误信息而不是抛出异常，让程序继续运行
            return {"status": -1, "message": f"JSON解析失败: {str(e)}", "raw_response": response.text}
        except requests.exceptions.RequestException as e:
            # 请求异常，记录完整响应内容
            error_msg = f"获取任务网络异常: {str(e)}"
            try:
                if hasattr(e, 'response') and e.response:
                    error_msg += f"\n响应状态码: {e.response.status_code}"
                    error_msg += f"\n响应头信息: {dict(e.response.headers)}"
                    error_msg += f"\n响应内容: {e.response.text}"
            except Exception:
                pass
            logger.error(error_msg)
            raise
        except Exception as e:
            logger.error(f"获取任务其他异常: {str(e)}")
            raise
    
    def extract_video_id(self, task_url):
        """
        从任务URL中提取视频ID
        
        Args:
            task_url: 任务URL
            
        Returns:
            str: 视频ID
        """
        # 提取形如 https://www.douyin.com/video/7454424455373262139 中的ID
        match = re.search(r'/video/(\d+)', task_url)
        if match:
            return match.group(1)
        return None
    
    def run_tasks(self, uid, sec_uid, max_tasks=100):
        """
        运行点赞任务
        
        Args:
            uid: 用户ID
            sec_uid: 安全用户ID
            max_tasks: 最大任务数，默认100
        """
        # 配置文件日志
        file_handler = logging.FileHandler('dy_love.log', encoding='utf-8')
        file_handler.setFormatter(logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s'))
        logger.addHandler(file_handler)
        
        task_count = 0
        
        logger.info(f"开始任务循环，用户ID: {uid}, 安全用户ID: {sec_uid}")
        
        while True:
            try:
                # 获取任务
                task_result = self.get_task(uid, sec_uid)
                
                if task_result.get("status") != 0 or not task_result.get("data"):
                    logger.warning("获取任务失败或没有可用任务")
                    time.sleep(60)  # 如果没有任务，等待60秒
                    continue
                
                task_data = task_result.get("data")
                task_url = task_data.get("taskUrl")
                
                # 提取视频ID
                video_id = self.extract_video_id(task_url)
                if not video_id:
                    logger.warning(f"无法从URL提取视频ID: {task_url}")
                    time.sleep(10)
                    continue
                
                # 执行点赞
                logger.info(f"开始点赞视频: {video_id}, URL: {task_url}")
                like_result = self.like_video(video_id)
                
                # 输出点赞结果
                if like_result.get("status_code") == 0:
                    logger.info(f"点赞成功，视频ID: {video_id}")
                else:
                    logger.warning(f"点赞失败，视频ID: {video_id}, 结果: {like_result}")
                
                # 计数器增加
                task_count += 1
                logger.info(f"已完成 {task_count}/{max_tasks} 个任务")
                
                # 达到最大任务数，重新开始循环
                if task_count >= max_tasks:
                    logger.info(f"已达到最大任务数 {max_tasks}，重新开始循环")
                    task_count = 0
                
                # 随机等待时间，以对抗风控
                sleep_time = random.randint(45, 120)
                logger.info(f"随机等待 {sleep_time} 秒...")
                time.sleep(sleep_time)
                
            except Exception as e:
                logger.error(f"任务执行出错: {str(e)}")
                # 打印完整的异常堆栈信息，便于调试
                import traceback
                logger.error(f"异常堆栈: {traceback.format_exc()}")
                time.sleep(30)  # 出错后等待30秒后继续

# 示例用法
if __name__ == "__main__":
    # 创建点赞协议实例
    dy_love = DYLoveProtocol()
    
    # 运行任务循环
    dy_love.run_tasks(
        uid="177336881", 
        sec_uid="MS4wLjABAAAA7Zd098eJ-fKLpj169UBS8PPebs0TqIzz0LXD34hbYHb"
    ) 