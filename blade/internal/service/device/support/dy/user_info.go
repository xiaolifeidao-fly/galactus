package dy

import (
	"galactus/blade/internal/service/device/dto"
)

type UserInfoEntity struct {
	*DyBaseEntity
}

type GetUserInfoBy struct {
	DyRequest[*UserInfoEntity]
}

func GetUserInfoByWeb(secUid string, webDevice *dto.WebDeviceDTO) (map[string]interface{}, error) {
	url := "https://www-hj.douyin.com/aweme/v1/web/user/profile/other/?"
	userInfoEntity := &UserInfoEntity{
		DyBaseEntity: newDyBaseEntity(url, "GET", webDevice),
	}
	userInfoEntity.AppendUrlParams("land_to", "1")
	userInfoEntity.AppendUrlParams("sec_user_id", secUid)
	userInfoEntity.AppendUrlParams("publish_video_strategy_type", "2")
	userInfoEntity.AppendUrlParams("personal_center_strategy", "1")
	abogus := userInfoEntity.GetAbogus(userInfoEntity.GetParams(), webDevice.UserAgent)
	userInfoEntity.AppendUrlParams("a_bogus", abogus)
	return DoRequest(userInfoEntity)
}
