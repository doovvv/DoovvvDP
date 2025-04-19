package v1

import (
	"doovvvDP/dal/model"
	"doovvvDP/dto"
)

func CreateBlog(blog model.Blog, userID uint64) *dto.Result {
	result := &dto.Result{}
	blog.UserID = userID
	err := model.CreateBlog(blog)
	if err != nil {
		return result.Fail(err.Error())
	}
	return result.Ok()
}
