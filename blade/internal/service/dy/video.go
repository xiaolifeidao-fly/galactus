package dy

import (
	"galactus/blade/internal/service/dy/response"
	dto "galactus/blade/internal/service/dy/response"
	"galactus/common/middleware/http"
	"log"
	"strconv"
	"strings"
)

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

func GetVideoItemInfo(videoInfo *VideoInfo) *dto.ExtItemDTO {
	extItem := &dto.ExtItemDTO{}
	extItem.DataStatus = response.ERROR
	videoResponse, err := GetVideoInfo(videoInfo)
	if err != nil {
		return extItem
	}
	statusCode := videoResponse["status_code"].(float64)
	if statusCode != 0 {
		return extItem
	}
	awemeDetail := videoResponse["aweme_detail"]
	if awemeDetail == nil {
		filterDetail := videoResponse["filter_detail"]
		if filterDetail != nil {
			filterDetailMap := filterDetail.(map[string]any)
			filterReason := filterDetailMap["filter_reason"].(string)
			if strings.Contains(filterReason, "status_deleted") {
				extItem.DataStatus = response.DELETE
				log.Println("video delete ", " videoId ", videoInfo.VideoId)
				return extItem
			}
			if strings.Contains(filterReason, "author_secret") {
				extItem.DataStatus = response.SECRET
				log.Println("video secret ", " videoId ", videoInfo.VideoId)
				return extItem
			}
		}
		return extItem
	}
	awemeDetailMap := awemeDetail.(map[string]any)
	statistics := awemeDetailMap["statistics"]
	if statistics == nil {
		return extItem
	}
	extItem.NowNum = int64(statistics.(map[string]any)["digg_count"].(float64))
	extItem.BusinessId = videoInfo.VideoId
	extItem.DataStatus = response.SUCCESS
	desc := awemeDetailMap["desc"].(string)
	extItem.Name = desc
	anchorInfo := awemeDetailMap["author"]
	if anchorInfo != nil {
		anchorInfoMap := anchorInfo.(map[string]any)
		extItem.Uid = anchorInfoMap["uid"].(string)
		extItem.ExtParams = map[string]interface{}{
			"secUid":   anchorInfoMap["sec_uid"].(string),
			"assistId": extItem.Uid,
			"hsFlag":   false,
		}
	}
	return extItem
}

func ConvertByVideoUrl(businessKey string) *response.ConvertItemDTO {
	convertItemDTO := &response.ConvertItemDTO{}
	convertItemDTO.DataStatus = response.ERROR
	if strings.HasPrefix(businessKey, "http") {
		if strings.Contains(businessKey, "v.douyin.com") {
			headers := map[string]string{
				"Referer":    "https://www.douyin.com",
				"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
			}
			response, err := http.GetToResponse(businessKey, "", headers)
			if err != nil {
				return convertItemDTO
			}
			defer response.Body.Close()
			businessKey = response.Request.URL.String()
		}
		if strings.Contains(businessKey, "www.douyin.com") {
			start := strings.Index(businessKey, "video/")
			end := strings.Index(businessKey, "?")
			if start == -1 {
				convertItemDTO.DataStatus = dto.DELETE
				return convertItemDTO
			}
			if end == -1 {
				end = len(businessKey)
			}
			businessKey = businessKey[start+6 : end]
		}
	}

	_, err := strconv.ParseUint(businessKey, 10, 64)
	if err != nil {
		convertItemDTO.DataStatus = dto.DELETE
		return convertItemDTO
	}
	convertItemDTO.ConvertValue = businessKey
	convertItemDTO.Property = map[string]interface{}{
		"url": "https://www.douyin.com/video/" + businessKey,
	}
	convertItemDTO.DataStatus = dto.SUCCESS
	return convertItemDTO
}
