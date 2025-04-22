package v1

import (
	"strconv"

	"doovvvDP/dal/model"
	"doovvvDP/dal/redis"
	"doovvvDP/dto"
	"doovvvDP/utils"
)

func FollowUser(userId, followId uint64, isFollow bool) *dto.Result {
	result := &dto.Result{}
	follow := model.Follow{
		UserID:       userId,
		FollowUserID: followId,
	}
	if isFollow {
		// 关注
		err := model.FollowUser(follow)
		if err != nil {
			return result.Fail(err.Error())
		}
		key := utils.FLLOW_KEY + strconv.FormatUint(userId, 10)
		redis.RDB.SAdd(redis.RCtx, key, strconv.FormatUint(followId, 10))
	} else {
		// 取消关注
		err := model.UnFollowUser(follow)
		if err != nil {
			return result.Fail(err.Error())
		}
		redis.RDB.SRem(redis.RCtx, utils.FLLOW_KEY+strconv.FormatUint(userId, 10), strconv.FormatUint(followId, 10))
	}
	return result.Ok()
}

func QueryFollowByUserId(userId, followId uint64) *dto.Result {
	result := &dto.Result{}
	follow := model.Follow{
		UserID:       userId,
		FollowUserID: followId,
	}
	count, err := model.QueryFollowByUserId(follow)
	if err != nil {
		return result.Fail(err.Error())
	}
	return result.OkWithData(count > 0)
}

func FollowCommons(userId, otherId uint64) *dto.Result {
	result := &dto.Result{}
	userKey := utils.FLLOW_KEY + strconv.FormatUint(userId, 10)
	otherKey := utils.FLLOW_KEY + strconv.FormatUint(otherId, 10)
	// 交集
	commonIds, err := redis.RDB.SInter(redis.RCtx, userKey, otherKey).Result()
	if err != nil {
		return result.Fail(err.Error())
	}

	commonIdsInt := make([]uint64, len(commonIds))
	for i, id := range commonIds {
		idInt, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			return result.Fail(err.Error())
		}
		commonIdsInt[i] = idInt
	}
	users, err := model.GetUserByIds(commonIdsInt)
	if err != nil {
		return result.Fail(err.Error())
	}
	userVos := make([]*dto.UserVo, len(users))
	for i, user := range users {
		userVo := &dto.UserVo{
			Id:       user.ID,
			NickName: user.NickName,
			Icon:     user.Icon,
		}
		userVos[i] = userVo
	}
	return result.OkWithData(userVos)
}
