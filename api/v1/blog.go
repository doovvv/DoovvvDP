package v1

import (
	"strconv"
	"time"

	"doovvvDP/dal/model"
	"doovvvDP/dal/redis"
	"doovvvDP/dto"
	"doovvvDP/utils"

	goredis "github.com/redis/go-redis/v9"
)

var idWorker, _ = utils.NewIdWorker(1)

func CreateBlog(blog model.Blog, userID uint64) *dto.Result {
	result := &dto.Result{}
	blog.UserID = userID
	err := model.CreateBlog(&blog)
	if err != nil {
		return result.Fail(err.Error())
	}
	fans, err := model.QueryFans(userID)
	if err != nil {
		return result.Fail(err.Error())
	}

	for _, fan := range fans {
		// 推送消息
		fanId := fan.UserID
		key := utils.FEED_KEY + strconv.FormatUint(fanId, 10)
		score, err := idWorker.Generate()
		if err != nil {
			return result.Fail(err.Error())
		}
		redis.RDB.ZAdd(redis.RCtx, key, goredis.Z{
			Score:  float64(score),
			Member: blog.ID,
		})
	}
	return result.OkWithData(blog.ID)
}

func QueryBlogByID(id uint64, selfId uint64) *dto.Result {
	result := &dto.Result{}
	blog, err := model.GetBlogById(id)
	if err != nil {
		return result.Fail(err.Error())
	}
	user, err := model.GetUserById(blog.UserID)
	if err != nil {
		return result.Fail(err.Error())
	}
	isLike := isBlogLiked(blog.ID, selfId)
	blogDTO := dto.BlogDTO{
		ID:         blog.ID,
		ShopID:     blog.ShopID,
		UserID:     blog.UserID,
		Title:      blog.Title,
		Images:     blog.Images,
		Content:    blog.Content,
		Liked:      blog.Liked,
		Comments:   blog.Comments,
		CreateTime: blog.CreateTime,
		UpdateTime: blog.UpdateTime,
		NickName:   user.NickName,
		Icon:       user.Icon,
		IsLike:     isLike,
	}
	return result.OkWithData(blogDTO)
}

func QueryHotBlogs(current int, selfId uint64) *dto.Result {
	result := &dto.Result{}
	blogs, err := model.QueryHotBlogs(current)
	if err != nil {
		return result.Fail(err.Error())
	}
	blogsDTO := make([]dto.BlogDTO, len(blogs))
	for i, blog := range blogs {
		// 获取博主信息
		user, err := model.GetUserById(blog.UserID)
		if err != nil {
			return result.Fail(err.Error())
		}
		var isLike bool
		if selfId != 0 {
			// 判断是否点赞
			isLike = isBlogLiked(blog.ID, selfId)
		} else {
			isLike = false
		}
		blogsDTO[i] = dto.BlogDTO{
			ID:         blog.ID,
			ShopID:     blog.ShopID,
			UserID:     blog.UserID,
			Title:      blog.Title,
			Images:     blog.Images,
			Content:    blog.Content,
			Liked:      blog.Liked,
			Comments:   blog.Comments,
			CreateTime: blog.CreateTime,
			UpdateTime: blog.UpdateTime,
			NickName:   user.NickName,
			Icon:       user.Icon,
			IsLike:     isLike,
		}
	}
	return result.OkWithData(blogsDTO)
}

func isBlogLiked(id uint64, userId uint64) bool {
	key := utils.LIKE_BLOG_KEY + strconv.FormatUint(id, 10)
	// 判断用户是否已点赞
	member := strconv.FormatUint(userId, 10) // Convert userId to string for ZMScore
	scores, err := redis.RDB.ZMScore(redis.RCtx, key, member).Result()
	if err != nil {
		return false
	}
	return scores[0] != 0
}

func LikeBlog(id uint64, userId uint64) *dto.Result {
	result := &dto.Result{}
	key := utils.LIKE_BLOG_KEY + strconv.FormatUint(id, 10)
	// 判断用户是否已点赞
	member := strconv.FormatUint(userId, 10) // Convert userId to string for ZMScore
	scores, err := redis.RDB.ZMScore(redis.RCtx, key, member).Result()
	if err != nil {
		return result.Fail(err.Error())
	}
	var isLiked bool
	isLiked = true
	if scores[0] == 0 {
		isLiked = false
	}
	if !isLiked {
		// 未点赞，点赞
		err = model.LikeBlog(id)
		if err != nil {
			return result.Fail(err.Error())
		}
		// 将用户ID添加到集合中
		member := goredis.Z{
			Score:  float64(time.Now().Unix()),
			Member: userId,
		}
		redis.RDB.ZAdd(redis.RCtx, key, member).Err()
	} else {
		// 已点赞，取消点赞
		err = model.UnLikeBlog(id)
		if err != nil {
			return result.Fail(err.Error())
		}
		// 将用户ID从集合中移除
		member := strconv.FormatUint(userId, 10)
		redis.RDB.ZRem(redis.RCtx, key, member).Err()
	}
	return result.Ok()
}

func QueryBlogLikes(id uint64) *dto.Result {
	result := &dto.Result{}
	key := utils.LIKE_BLOG_KEY + strconv.FormatUint(id, 10)
	// 获取所有点赞用户的ID
	members, err := redis.RDB.ZRange(redis.RCtx, key, 0, 4).Result()
	if err != nil {
		return result.Fail(err.Error())
	}
	// 将用户ID转换为uint64类型
	userIds := make([]uint64, 0)
	for _, userId := range members {
		if userId != "" {
			userIdUint, err := strconv.ParseUint(userId, 10, 64)
			if err == nil {
				userIds = append(userIds, userIdUint)
			}
		}
	}
	users, err := model.GetUserByIds(userIds)
	if err != nil {
		return result.Fail(err.Error())
	}
	userVos := make([]dto.UserVo, len(users))
	for i, user := range users {
		userVos[i] = dto.UserVo{
			Id:       user.ID,
			NickName: user.NickName,
			Icon:     user.Icon,
		}
	}
	return result.OkWithData(userVos)
}

func QueryBlogByUserId(userId uint64, current int) *dto.Result {
	result := &dto.Result{}
	blogs, err := model.QueryBlogOfUser(userId, current)
	if err != nil {
		return result.Fail(err.Error())
	}
	blogsDTO := make([]dto.BlogDTO, len(blogs))
	for i, blog := range blogs {
		// 获取博主信息
		user, err := model.GetUserById(blog.UserID)
		if err != nil {
			return result.Fail(err.Error())
		}
		blogsDTO[i] = dto.BlogDTO{
			ID:         blog.ID,
			ShopID:     blog.ShopID,
			UserID:     blog.UserID,
			Title:      blog.Title,
			Images:     blog.Images,
			Content:    blog.Content,
			Liked:      blog.Liked,
			Comments:   blog.Comments,
			CreateTime: blog.CreateTime,
			UpdateTime: blog.UpdateTime,
			NickName:   user.NickName,
			Icon:       user.Icon,
		}
	}
	return result.OkWithData(blogsDTO)
}

func QueryBlogOfFollow(userId uint64, maxId uint64, offset int) *dto.Result {
	result := &dto.Result{}

	key := utils.FEED_KEY + strconv.FormatUint(userId, 10)
	tubles, err := redis.RDB.ZRevRangeByScoreWithScores(redis.RCtx, key, &goredis.ZRangeBy{
		Min:    "0",
		Max:    strconv.FormatUint(maxId, 10),
		Offset: int64(offset),
	}).Result()
	if err != nil {
		return result.Fail(err.Error())
	}
	var minId int64
	ids := make([]uint64, len(tubles))
	for i, tuble := range tubles {
		minId = int64(tuble.Score)
		ids[i], err = strconv.ParseUint(tuble.Member.(string), 10, 64)
		if err != nil {
			return result.Fail(err.Error())
		}
	}
	blogs, err := model.QueryBlogsByIDs(ids)
	if err != nil {
		return result.Fail(err.Error())
	}
	blogsDTO := make([]any, len(ids))
	for i, blog := range blogs {
		// 获取博主信息
		user, err := model.GetUserById(blog.UserID)
		if err != nil {
			return result.Fail(err.Error())
		}
		var isLike bool
		blogsDTO[i] = dto.BlogDTO{
			ID:         blog.ID,
			ShopID:     blog.ShopID,
			UserID:     blog.UserID,
			Title:      blog.Title,
			Images:     blog.Images,
			Content:    blog.Content,
			Liked:      blog.Liked,
			Comments:   blog.Comments,
			CreateTime: blog.CreateTime,
			UpdateTime: blog.UpdateTime,
			NickName:   user.NickName,
			Icon:       user.Icon,
			IsLike:     isLike,
		}
	}

	scrollVo := dto.ScrollVo{
		List:    blogsDTO,
		MinTime: minId,
		Offset:  1,
	}
	return result.OkWithData(scrollVo)
}
