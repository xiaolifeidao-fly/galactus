from random import choice
from random import randint
from random import random
from time import time
from urllib.parse import urlencode

from gmssl import sm3, func


class ABogus:
    __end_string = "cus"
    __str = {
        "s0": "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",
        "s1": "Dkdpgh4ZKsQB80/Mfvw36XI1R25+WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe=",
        "s2": "Dkdpgh4ZKsQB80/Mfvw36XI1R25-WUAlEi7NLboqYTOPuzmFjJnryx9HVGcaStCe=",
        "s3": "ckdp1h4ZKsUB80/Mfvw36XIgR25+WQAlEi7NLboqYTOPuzmFjJnryx9HVGDaStCe",
        "s4": "Dkdpgh2ZmsQB80/MfvV36XI1R45-WUAlEixNLwoqYTOPuzKFjJnry79HbGcaStCe",
    }

    def __init__(self, user_agent: str = "", platform: str = None):
        self.user_agent = user_agent
        self.ua_code = self.generate_ua_code(user_agent)
        self.browser = self.generate_browser_info(platform)
        self.browser_len = len(self.browser)
        self.browser_code = self.char_code_at(self.browser)

    def generate_ua_code(self, user_agent: str) -> list:
        numbers = [0.00390625, 1, 14]
        key_string = ''.join(chr(int(num)) for num in numbers)
        return self.sm3_to_array(self.generate_result(self.rc4_encrypt(user_agent, key_string), "s3"))

    def list_1(self, a=170, b=85, c=45) -> list:
        return self.random_list(a, b, 1, 2, 5, c & a)

    def list_2(self, a=170, b=85) -> list:
        return self.random_list(a, b, 1, 0, 0, 0)

    def list_3(self, a=170, b=85) -> list:
        return self.random_list(a, b, 1, 0, 5, 0)

    def random_list(self, b=170, c=85, d=0, e=0, f=0, g=0) -> list:
        r = random() * 10000
        v = [r, int(r) & 255, int(r) >> 8]
        s = v[1] & b | d
        v.append(s)
        s = v[1] & c | e
        v.append(s)
        s = v[2] & b | f
        v.append(s)
        s = v[2] & c | g
        v.append(s)
        return v[-4:]

    def from_char_code(self, *args):
        return "".join(chr(code) for code in args)

    def generate_string_1(self):
        return self.from_char_code(*self.list_1()) + self.from_char_code(
            *self.list_2()) + self.from_char_code(*self.list_3())

    def generate_string_2(self, url_params: str, method="GET") -> str:
        a = self.generate_string_2_list(url_params, method)
        e = self.end_check_num(a)
        a.extend(self.browser_code)
        a.append(e)
        return self.rc4_encrypt(self.from_char_code(*a), "y")

    def generate_string_2_list(self, url_params: str, method="GET") -> list:
        start_time = int(time() * 1000)
        end_time = start_time + randint(4, 8)
        params_array = self.generate_params_code(url_params)
        method_array = self.generate_method_code(method)
        return self.list_4(
            (end_time >> 24) & 255,
            params_array[21],
            self.ua_code[23],
            (end_time >> 16) & 255,
            params_array[22],
            self.ua_code[24],
            (end_time >> 8) & 255,
            (end_time >> 0) & 255,
            (start_time >> 24) & 255,
            (start_time >> 16) & 255,
            (start_time >> 8) & 255,
            (start_time >> 0) & 255,
            method_array[21],
            method_array[22],
            int(end_time / 256 / 256 / 256 / 256) >> 0,
            int(start_time / 256 / 256 / 256 / 256) >> 0,
            self.browser_len,
        )

    def list_4(self,
               a: int,
               b: int,
               c: int,
               d: int,
               e: int,
               f: int,
               g: int,
               h: int,
               i: int,
               j: int,
               k: int,
               m: int,
               n: int,
               o: int,
               p: int,
               q: int,
               r: int,
               ) -> list:
        return [
            44,
            a,
            0,
            0,
            0,
            0,
            24,
            b,
            n,
            0,
            c,
            d,
            0,
            0,
            0,
            1,
            0,
            239,
            e,
            o,
            f,
            g,
            0,
            0,
            0,
            0,
            h,
            0,
            0,
            14,
            i,
            j,
            0,
            k,
            m,
            3,
            p,
            1,
            q,
            1,
            r,
            0,
            0,
            0]

    def end_check_num(self, a: list):
        r = 0
        for i in a:
            r ^= i
        return r

    def convert_to_char_code(self, a):
        d = []
        for i in a:
            d.append(ord(i))
        return d

    def split_array(self, arr, chunk_size=64):
        result = []
        for i in range(0, len(arr), chunk_size):
            result.append(arr[i:i + chunk_size])
        return result

    def char_code_at(self, s):
        return [ord(char) for char in s]

    def generate_result(self, s, e="s4"):
        r = []

        for i in range(0, len(s), 3):
            if i + 2 < len(s):
                n = (
                        (ord(s[i]) << 16)
                        | (ord(s[i + 1]) << 8)
                        | ord(s[i + 2])
                )
            elif i + 1 < len(s):
                n = (ord(s[i]) << 16) | (
                        ord(s[i + 1]) << 8
                )
            else:
                n = ord(s[i]) << 16

            for j, k in zip(range(18, -1, -6), (0xFC0000, 0x03F000, 0x0FC0, 0x3F)):
                if j == 6 and i + 1 >= len(s):
                    break
                if j == 0 and i + 2 >= len(s):
                    break
                r.append(self.__str[e][(n & k) >> j])

        r.append("=" * ((4 - len(r) % 4) % 4))
        return "".join(r)

    def generate_method_code(self, method: str = "GET") -> list[int]:
        return self.sm3_to_array(self.sm3_to_array(method + self.__end_string))

    def generate_params_code(self, params: str) -> list[int]:
        return self.sm3_to_array(self.sm3_to_array(params + self.__end_string))

    def sm3_to_array(self, data: str | list) -> list[int]:
        if isinstance(data, str):
            b = data.encode("utf-8")
        else:
            b = bytes(data)  # 将 List[int] 转换为字节数组

        # 将字节数组转换为适合 sm3.sm3_hash 函数处理的列表格式
        h = sm3.sm3_hash(func.bytes_to_list(b))

        # 将十六进制字符串结果转换为十进制整数列表
        return [int(h[i: i + 2], 16) for i in range(0, len(h), 2)]

    def generate_browser_info(self, platform: str = "Win32") -> str:
        inner_width = randint(1280, 1920)
        inner_height = randint(720, 1080)
        outer_width = randint(inner_width, 1920)
        outer_height = randint(inner_height, 1080)
        screen_x = 0
        screen_y = choice((0, 30))
        value_list = [
            inner_width,
            inner_height,
            outer_width,
            outer_height,
            screen_x,
            screen_y,
            0,
            0,
            outer_width,
            outer_height,
            outer_width,
            outer_height,
            inner_width,
            inner_height,
            24,
            24,
            platform,
        ]
        return "|".join(str(i) for i in value_list)

    def rc4_encrypt(self, plaintext, key):
        s = list(range(256))
        j = 0

        for i in range(256):
            j = (j + s[i] + ord(key[i % len(key)])) % 256
            s[i], s[j] = s[j], s[i]

        i = 0
        j = 0
        cipher = []

        for k in range(len(plaintext)):
            i = (i + 1) % 256
            j = (j + s[i]) % 256
            s[i], s[j] = s[j], s[i]
            t = (s[i] + s[j]) % 256
            cipher.append(chr(s[t] ^ ord(plaintext[k])))

        return ''.join(cipher)

    def generate_a_bogus(self, url_params: dict | str) -> str:
        string_1 = self.generate_string_1()
        string_2 = self.generate_string_2(urlencode(url_params))
        string = string_1 + string_2
        return self.generate_result(string, "s4")


if __name__ == "__main__":
    import requests

    url = "https://www-hj.douyin.com/aweme/v1/web/aweme/favorite/?"

    headers = {
        'user-agent': "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
        'Cookie': "ttwid=1%7CMic1PT2poCovHTRr3WW_Fyljnl1fTkhFg_C7r4zdNhY%7C1736058777%7Cab1c6541b1d5bd54fc9ac1d5341a5067e28ec21a6e4ffd0074c22b887570f61e; UIFID_TEMP=0de8750d2b188f4235dbfd208e44abbb976428f0720eb983255afefa45d39c0c0b11f5a4a5274f3afe6c779377308d74b5c8cdd68d7e05d752d05bb9ebded043d8bf8fabed12cd53e158c0927ccfac21; hevc_supported=true; home_can_add_dy_2_desktop=%220%22; odin_tt=1a3604485e8412f08b3ab80ee84720b2e3a90a63f7a560f0465b68fb86745d39def020727cc0853c5346cb68486605b52de6c482665fd0cf1286a057c0404432402f1cec0a69cb1b744eb043fd6e73b4; FORCE_LOGIN=%7B%22videoConsumedRemainSeconds%22%3A180%7D; passport_csrf_token=b22cb0775107cc4e14fa80f732efe9a3; passport_csrf_token_default=b22cb0775107cc4e14fa80f732efe9a3; bd_ticket_guard_client_web_domain=2; WallpaperGuide=%7B%22showTime%22%3A1736058859772%2C%22closeTime%22%3A0%2C%22showCount%22%3A1%2C%22cursor1%22%3A10%2C%22cursor2%22%3A2%7D; stream_recommend_feed_params=%22%7B%5C%22cookie_enabled%5C%22%3Atrue%2C%5C%22screen_width%5C%22%3A1792%2C%5C%22screen_height%5C%22%3A1120%2C%5C%22browser_online%5C%22%3Atrue%2C%5C%22cpu_core_num%5C%22%3A6%2C%5C%22device_memory%5C%22%3A8%2C%5C%22downlink%5C%22%3A10%2C%5C%22effective_type%5C%22%3A%5C%224g%5C%22%2C%5C%22round_trip_time%5C%22%3A50%7D%22; strategyABtestKey=%221736243727.277%22; is_dash_user=1; volume_info=%7B%22isUserMute%22%3Afalse%2C%22isMute%22%3Atrue%2C%22volume%22%3A0.5%7D; stream_player_status_params=%22%7B%5C%22is_auto_play%5C%22%3A0%2C%5C%22is_full_screen%5C%22%3A0%2C%5C%22is_full_webscreen%5C%22%3A0%2C%5C%22is_mute%5C%22%3A0%2C%5C%22is_speed%5C%22%3A1%2C%5C%22is_visible%5C%22%3A0%7D%22; download_guide=%223%2F20250107%2F0%22; biz_trace_id=5dc9e443; bd_ticket_guard_client_data=eyJiZC10aWNrZXQtZ3VhcmQtdmVyc2lvbiI6MiwiYmQtdGlja2V0LWd1YXJkLWl0ZXJhdGlvbi12ZXJzaW9uIjoxLCJiZC10aWNrZXQtZ3VhcmQtcmVlLXB1YmxpYy1rZXkiOiJCR0hHU0o2Vk5KNDRCRjE2SnpXcnRyaWJhTVJMZWRnVGZVaHplcDEydTE4SFVsYkxBNnlhRGdGZzBZTHFPbURzVGJrTFJhYUE2QjlnWFByaUdlbktMQVE9IiwiYmQtdGlja2V0LWd1YXJkLXdlYi12ZXJzaW9uIjoyfQ%3D%3D; IsDouyinActive=true"
    }

    ## url 问号后边的参数转 params json
    from urllib.parse import parse_qs
    url_params = "device_platform=webapp&aid=6383&channel=channel_pc_web&sec_user_id=MS4wLjABAAAAAeXzDPimVN_Q71x0uCRfdmfO2Jfy0q3lSaL8TrB_a6Y&max_cursor=0&min_cursor=0&whale_cut_token=&cut_version=1&count=18&publish_video_strategy_type=2&update_version_code=170400&pc_client_type=1&pc_libra_divert=Mac&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1792&screen_height=1120&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Chrome&browser_version=131.0.0.0&browser_online=true&engine_name=Blink&engine_version=131.0.0.0&os_name=Mac+OS&os_version=10.15.7&cpu_core_num=6&device_memory=8&platform=PC&downlink=10&effective_type=4g&round_trip_time=50&webid=7456315560445429283&uifid=0de8750d2b188f4235dbfd208e44abbb976428f0720eb983255afefa45d39c0c88416d5bb97ea03cc5c8617e0fe45229836a56a3324dc84e4d75ae380941412565e0974d98c378bf3a277193418cc93c88206c465aec231551e7b83aa5eb736b8e572e0ca20c52ec8a30f92c2976ae619d03b276ca043cdda81da1e8d985f93b4751c619319c2630df323e49a8e33d08109754cc5f904052b2282768a79c0244&msToken=nhYETIaQxnmewdSLi24OqLkqJMTVfflyqRJHz4QWL5AATA3sSXM5Sunet5sGnkEoEqF5tpSvZwWLWXh56-XY7shsai9OxXf21rThIx45atU6z_uCXcnwesXPx4JAcpBDJte0_rNzhWZC0XAQrilXLyrG-_aa3k66irK19kNI0gdC1c3lcvWxL50%3D"
    params = {k: v[0] for k, v in parse_qs(url_params).items()}

    bogus = ABogus(headers['user-agent'])

    a_bogus = bogus.generate_a_bogus(params)
    print(a_bogus)

    params['a_bogus'] = a_bogus
    response = requests.get(url, params=params, headers=headers)

    # print(response.json())

    user_agent_11 = "Mozilla/5.0"
    print(bogus.generate_ua_code(user_agent_11))