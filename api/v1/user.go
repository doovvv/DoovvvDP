package v1

import (
	"doovvvDP/dal/model"
	"doovvvDP/dto"
	"doovvvDP/utils"
	"fmt"

	"github.com/gin-contrib/sessions"
)
func SendCode(phone string,session sessions.Session)(result *dto.Result){
	result = &dto.Result{}
	// 检验手机号
	if(!utils.IsValidPhoneNumber(phone)){
		return result.Fail("手机号格式错误")
	}
	//生产验证码
	code := utils.GenerateRandomCode(6)
	//保存验证码
	//todo:修改key为手机号
	session.Set("code", code)
	if err := session.Save(); err!= nil {
		return result.Fail("验证码保存失败")
	}
	fmt.Println(session.Get("code"))
	//发送验证码
	// fmt.Fprintf(os.Stdout,"code sent success! %s\n",phone)
	//返回响应
	return result.Ok()


}
func Login(userdto dto.UserDTO,session sessions.Session)(result *dto.Result){
	result = &dto.Result{}
	fmt.Println(userdto.Phone)
	// 检验手机号
	if(!utils.IsValidPhoneNumber(userdto.Phone)){
		return result.Fail("手机号格式错误")
	}
	//检验验证码
	code := session.Get("code")
	fmt.Println("real",code)
	if code != userdto.Code {
		return result.Fail("验证码错误")
	}
	var user model.TbUser
	var err error
	var ok bool
	//创建用户并保存到数据库
	if user,ok = model.CheckUserNotExist(userdto.Phone);ok{
		user,err = model.CreateUserWithPhone(userdto.Phone)
		if err!= nil {
			return result.Fail("用户创建失败")
		}
	}
	userVo := dto.UserVo{
		Id: user.ID,
		NickName: user.NickName,
		Icon: user.Icon,
	}
	//保存用户信息到session
	session.Set("user", userVo)
	if err := session.Save(); err!= nil {
		return result.Fail(err.Error())
	}
	// fmt.Println(session.Get("user"))
	//返回响应
	return result.OkWithData(user)

}