package dy

type VideoInfo struct {
	*DyBaseEntity
	VideoId string
}

func GetVideoInfo(videoInfo *VideoInfo) (map[string]any, error) {
	url := "https://www.douyin.com/aweme/v1/web/aweme/detail/?"
	videoInfo.Init(url)
	videoInfo.
		AppendUrlParams("aweme_id", videoInfo.VideoId).
		AppendUrlParams("pc_libra_divert", "Mac")
	return DoGet(videoInfo)
}

func PlayerVideo(videoInfo *VideoInfo) (map[string]any, error) {
	params := map[string]interface{}{
		"aweme_type": 68,
		"item_id":    videoInfo.VideoId,
		"play_delta": 1,
		"source":     0,
	}
	url := "https://www-hj.douyin.com/aweme/v2/web/aweme/stats/?"
	videoInfo.Init(url)
	videoInfo.
		AppendUrlParams("pc_libra_divert", "Mac")
	return DoPost(videoInfo, params, "application/x-www-form-urlencoded; charset=UTF-8")
}

func LoveVideo(videoInfo *VideoInfo) (map[string]any, error) {
	url := "https://www-hj.douyin.com/aweme/v1/web/commit/item/digg/?"
	params := map[string]interface{}{
		"aweme_id":  videoInfo.VideoId,
		"item_type": 0,
		"type":      1,
	}
	videoInfo.Init(url)
	videoInfo.
		AppendUrlParams("pc_libra_divert", "Mac")
	return DoPost(videoInfo, params, "application/x-www-form-urlencoded; charset=UTF-8")
}
