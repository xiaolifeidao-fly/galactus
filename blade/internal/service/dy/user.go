package dy

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

type UserFavoriteEntity struct {
	*DyBaseEntity
	SecUid    string
	MaxCursor int
	MinCursor int
	Count     int
}

func GetUserFavoriteByWeb(userFavoriteEntity *UserFavoriteEntity) (map[string]interface{}, error) {
	url := "https://www.douyin.com/aweme/v1/web/aweme/favorite/?"
	userFavoriteEntity.Init(url)
	userFavoriteEntity.Count = 18
	userFavoriteEntity.
		AppendUrlParams("sec_user_id", userFavoriteEntity.SecUid).
		AppendUrlParams("max_cursor", userFavoriteEntity.MaxCursor).
		AppendUrlParams("min_cursor", userFavoriteEntity.MinCursor).
		AppendUrlParams("publish_video_strategy_type", "2").
		AppendUrlParams("cut_version", "1").
		AppendUrlParams("pc_libra_divert", "Mac").
		AppendUrlParams("count", userFavoriteEntity.Count)
	return DoGet(userFavoriteEntity)
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
