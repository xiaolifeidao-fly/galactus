package main

import (
	"fmt"
	"galactus/blade/a"
	webDeviceService "galactus/blade/internal/service/device"
	"galactus/blade/internal/service/device/support/dy"
)

func main() {
	a.Init()
	webDeviceDTO, _ := webDeviceService.FindWebDeviceById(6)
	result, _ := dy.GetUserInfoByWeb("MS4wLjABAAAAkeDX3Vsj5KqMVJnwzOEZ8U9VLA7DUiOSqMdkZ16gpS8", webDeviceDTO)
	fmt.Println(result)
	// result, _ := http.Get("https://www-hj.douyin.com/aweme/v1/web/user/profile/other/?channel=channel_pc_web&os_name=Mac%20OS&downlink=10&device_platform=webapp&update_version_code=170400&version_name=17.4.0&browser_online=true&engine_name=Blink&os_version=10.15.7&screen_height=1120&verify_fp=verify_m2baxsc2_0gnoyeB1_ALLJ_4ieq_9FVN_hCFnHcXCvf6I&cpu_core_num=12&browser_language=zh-CN&browser_name=Chrome&engine_version=130.0.0.0&platform=PC&round_trip_time=50&source=channel_pc_web&cookie_enabled=1&effective_type=4&webid=7204670535381141049&fp=verify_m2baxsc2_0gnoyeB1_ALLJ_4ieq_9FVN_hCFnHcXCvf6I&msToken=t83bXziTG2hsAxVAKQLh6EtfidNkccPxrSnmQwfUPS9k8HmgZtdZmRehbMZ5HRPRPTxN5WD5lO9N7issUPFiOzf2bItwMiDig9uUgCwd9yt&pc_client_type=1&browser_platform=MacIntel&browser_version=130.0.0.0&aid=6383&version_code=170400&screen_width=1792&device_memory=8&uifid=96cd3b166f3029d7c1cc3f64582454ab8a83ff1f9e6d6689076dd47ef1dca5f838478297937fcb88316a6ee9493aba5f8d54cd67ea1d1784b59d9556552d13f105e3f5242002dfd9e2235682f0ed8bc2bca11f5182cadab8b57e511e045ec0bb4de6fe3074b07c589e096bbf2cb03bc38cceaa7ff28af647ba70fe054ab39588c42e7c7050c93bdcdfd2dc7daffe47f32717bda207d1a5d22c8e4bc7c1efafc9&&land_to=1&sec_user_id=MS4wLjABAAAAkeDX3Vsj5KqMVJnwzOEZ8U9VLA7DUiOSqMdkZ16gpS8&publish_video_strategy_type=2&personal_center_strategy=1&a_bogus=dXWZBfw6difkgfSX512LfY3qV5H3Y0S50SVkMDheJnVYC639HMOG9exYqDUvbUueLG/dIbDjy4hbTrOgrQ2G0Zwf9Skw/2A2mESkKl5Q5xSSs1XyeykgJUhimktRSeo2RkBlrOfkqJKGKuRplnl60fAAPn6=", "", nil)
	// fmt.Println(result)

}
