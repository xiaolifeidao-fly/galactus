package dy

import (
	"galactus/blade/internal/service/dy/response"
	"strings"
)

type UserInfoEntity struct {
	*DyBaseEntity
	SecUid string
}

func GetUserInfoByWeb(userInfoEntity *UserInfoEntity) (map[string]interface{}, error) {
	url := "https://www-hj.douyin.com/aweme/v1/web/user/profile/other/?"
	userInfoEntity.Init(url)
	userInfoEntity.
		AppendUrlParams("land_to", "1").
		AppendUrlParams("sec_user_id", userInfoEntity.SecUid).
		AppendUrlParams("publish_video_strategy_type", "2").
		AppendUrlParams("personal_center_strategy", "1")
	return DoGet(userInfoEntity)
}

func GetUserInfo(userInfoEntity *UserInfoEntity) *response.ExtItemDTO {
	userInfoDTO := &response.ExtItemDTO{}
	userInfoDTO.DataStatus = response.ERROR
	userInfo, err := GetUserInfoByWeb(userInfoEntity)
	if err != nil {
		return userInfoDTO
	}
	if _, ok := userInfo["status_code"]; !ok {
		userInfoDTO.DataStatus = response.NOT_GET_DATA
		return userInfoDTO
	}
	statusCode := userInfo["status_code"].(float64)
	if statusCode != 0 {
		return userInfoDTO
	}
	
	userInfoDTO.BusinessId = userInfo["secUid"].(string)
	userInfoDTO.ExtParams["uid"] = userInfo["uid"].(string)
	userInfoDTO.DataStatus = response.SUCCESS
	return userInfoDTO
}

type UserFavoriteEntity struct {
	*DyBaseEntity
	SecUid    string
	MaxCursor int
	MinCursor int
	Count     int
}

func GetUserFavoriteByWeb(userFavoriteEntity *UserFavoriteEntity) (map[string]interface{}, error) {
	url := "https://www-hj.douyin.com/aweme/v1/web/aweme/favorite/?"
	userFavoriteEntity.Init(url)
	if userFavoriteEntity.Count == 0 {
		userFavoriteEntity.Count = 18
	}
	userFavoriteEntity.
		AppendUrlParams("sec_user_id", userFavoriteEntity.SecUid).
		AppendUrlParams("publish_video_strategy_type", "2").
		AppendUrlParams("cut_version", "1").
		AppendUrlParams("whale_cut_token", "").
		AppendUrlParams("count", userFavoriteEntity.Count).
		AppendUrlParams("max_cursor", userFavoriteEntity.MaxCursor).
		AppendUrlParams("min_cursor", userFavoriteEntity.MinCursor)
	return DoGet(userFavoriteEntity)
}

func GetUserFavorite(userFavoriteEntity *UserFavoriteEntity) *response.ExtItemListDTO {
	extItemListDTO := &response.ExtItemListDTO{}
	extItemListDTO.DataStatus = response.ERROR
	userFavoriteResult, err := GetUserFavoriteByWeb(userFavoriteEntity)
	if err != nil {
		return extItemListDTO
	}
	if _, ok := userFavoriteResult["status_code"]; !ok {
		extItemListDTO.DataStatus = response.NOT_GET_DATA
		return extItemListDTO
	}
	statusCode := userFavoriteResult["status_code"].(float64)
	if statusCode != 0 {
		return extItemListDTO
	}
	awemeList := userFavoriteResult["aweme_list"]
	if awemeList == nil {
		extItemListDTO.DataStatus = response.SUCCESS
		return extItemListDTO
	}
	//转成数组
	awemeListArray := awemeList.([]any)
	data := map[string]*response.ExtItemDTO{}
	for _, aweme := range awemeListArray {
		extItem := &response.ExtItemDTO{}
		awemeMap := aweme.(map[string]any)
		awemeId := awemeMap["aweme_id"].(string)
		extItem.BusinessId = awemeId
		data[awemeId] = extItem
	}
	extItemListDTO.Data = data
	// float64转bool
	extItemListDTO.HadMore = userFavoriteResult["has_more"].(float64) != 0
	extItemListDTO.NextIndex = int64(userFavoriteResult["max_cursor"].(float64))
	extItemListDTO.DataStatus = response.SUCCESS
	extItemListDTO.TotalNum = int(len(data))
	return extItemListDTO
}

type UserFollowingEntity struct {
	*DyBaseEntity
	UserId string
	SecUid string
	Offset int
	Count  int
}

func GetUserFollowingByWeb(userFollowEntity *UserFollowingEntity) (map[string]interface{}, error) {
	url := "https://www.douyin.com/aweme/v1/web/user/following/list/?"
	userFollowEntity.Init(url)
	userFollowEntity.Count = 20
	userFollowEntity.
		AppendUrlParams("user_id", userFollowEntity.UserId).
		AppendUrlParams("sec_user_id", userFollowEntity.SecUid).
		AppendUrlParams("offset", userFollowEntity.Offset).
		AppendUrlParams("count", userFollowEntity.Count).
		AppendUrlParams("source_type", "4").
		AppendUrlParams("is_top", "1").
		AppendUrlParams("min_time", "0").
		AppendUrlParams("max_time", "0").
		AppendUrlParams("gps_access", "0").
		AppendUrlParams("address_book_access", "0")
	return DoGet(userFollowEntity)
}

func getUid(uidType string, url string, userInfoEntity *UserInfoEntity) string {
	if strings.EqualFold(response.DY_UID_TYPE, uidType) {
		startIndex := strings.Index(url, "user/")
		endIndex := len(url)
		secUid := url[startIndex+5 : endIndex]
		userInfoEntity.SecUid = secUid
		userInfo, err := GetUserInfoByWeb(userInfoEntity)
		if err != nil {
			return ""
		}
		statusCode := userInfo["status_code"].(float64)
		if statusCode != 0 {
			return ""
		}
		user := userInfo["user"]
		if user == nil {
			return ""
		}
		userMap := user.(map[string]any)
		uid := userMap["uid"].(string)
		return uid
	}

	if strings.EqualFold(response.HS_UID_TYPE, uidType) {
		startIndex := strings.Index(url, "?to_user_id=")
		endIndex := strings.Index(url, "&")
		if endIndex == -1 {
			endIndex = len(url)
		}
		return url[startIndex+12 : endIndex]
	}
	return ""
}

func getUidType(url string) string {
	if strings.Contains(url, "www.douyin.com/user/") {
		return response.DY_UID_TYPE
	}
	if strings.Contains(url, "share.huoshan.com/pages/user/index.html/") {
		return response.HS_UID_TYPE
	}
	return ""
}

func GetUrlByUrl(url string, userInfoEntity *UserInfoEntity) *response.ConvertUrlItemDTO {
	convertUrlItemDTO := &response.ConvertUrlItemDTO{}
	convertUrlItemDTO.DataStatus = response.ERROR
	uidType := getUidType(url)
	if uidType == "" {
		return convertUrlItemDTO
	}
	convertUrlItemDTO.UidType = uidType
	uid := getUid(uidType, url, userInfoEntity)
	if uid == "" {
		return convertUrlItemDTO
	}
	convertUrlItemDTO.Uid = uid
	convertUrlItemDTO.DataStatus = response.SUCCESS
	return convertUrlItemDTO
}
