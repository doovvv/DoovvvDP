package v1

import (
	"fmt"
	"strconv"

	"doovvvDP/dal/model"
	"doovvvDP/dal/redis"
	"doovvvDP/dto"
	"doovvvDP/utils"

	"github.com/google/uuid"
)

func SendCode(phone string) (result *dto.Result) {
	result = &dto.Result{}
	// 检验手机号
	if !utils.IsValidPhoneNumber(phone) {
		return result.Fail("手机号格式错误")
	}
	// 生产验证码
	code := utils.GenerateRandomCode(6)
	// 保存验证码
	// session.Set("code", code)
	// if err := session.Save(); err!= nil {
	// 	return result.Fail("验证码保存失败")
	// }

	err := redis.RDB.Set(redis.RCtx, utils.LOGIN_CODE_KEY+phone, code, utils.LOGIN_CODE_TTL).Err()
	if err != nil {
		return result.Fail("验证码保存失败")
	}
	fmt.Println(code)

	// 发送验证码
	// fmt.Fprintf(os.Stdout,"code sent success! %s\n",phone)
	// 返回响应
	return result.Ok()
}

func Login(userdto dto.UserDTO) (result *dto.Result) {
	result = &dto.Result{}
	fmt.Println(userdto.Phone)
	// 检验手机号
	if !utils.IsValidPhoneNumber(userdto.Phone) {
		return result.Fail("手机号格式错误")
	}
	// 检验验证码
	code := redis.RDB.Get(redis.RCtx, utils.LOGIN_CODE_KEY+userdto.Phone).Val()
	fmt.Println("real", code)
	if code != userdto.Code {
		return result.Fail("验证码错误")
	}
	var user model.TbUser
	var err error
	var ok bool
	// 创建用户并保存到数据库
	if user, ok = model.CheckUserNotExist(userdto.Phone); ok {
		user, err = model.CreateUserWithPhone(userdto.Phone)
		if err != nil {
			return result.Fail("用户创建失败")
		}
	}
	// 保存用户信息到redis
	token := uuid.New().String()
	err = redis.RDB.HSet(redis.RCtx, utils.LOGIN_TOKEN_KEY+token,
		map[string]string{
			"id":       strconv.Itoa(int(user.ID)),
			"nickName": user.NickName,
			"icon":     user.Icon,
		}).Err()
	redis.RDB.Expire(redis.RCtx, utils.LOGIN_TOKEN_KEY+token, utils.LOGIN_TOKEN_TTL)
	if err != nil {
		return result.Fail(err.Error())
	}

	// session.Set("user", userVo)
	// if err := session.Save(); err!= nil {
	// 	return result.Fail(err.Error())
	// }
	// fmt.Println(session.Get("user"))
	// 返回响应
	return result.OkWithData(token)
}

func QueryUserById(id uint64) (result *dto.Result) {
	result = &dto.Result{}
	user, err := model.GetUserById(id)
	userVo := dto.UserVo{
		Id:       user.ID,
		NickName: user.NickName,
		Icon:     user.Icon,
	}
	if err != nil {
		return result.Fail(err.Error())
	}
	return result.OkWithData(userVo)
}
