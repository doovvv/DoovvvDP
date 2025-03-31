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
	fmt.Println(code)
	//保存验证码
	//todo:修改key为手机号
	session.Set("code", code)
	if err := session.Save(); err!= nil {
		return result.Fail("验证码保存失败")
	}
	//发送验证码
	// fmt.Fprintf(os.Stdout,"code sent success! %s\n",phone)
	//返回响应
	return result.Ok()


}
func Login(userdto dto.UserDTO,session sessions.Session)(result *dto.Result){
	result = &dto.Result{}
	// 检验手机号
	if(!utils.IsValidPhoneNumber(userdto.Phone)){
		return result.Fail("手机号格式错误")
	}
	//检验验证码
	code := session.Get("code")
	if code != userdto.Code {
		return result.Fail("验证码错误")
	}
	//创建用户并保存到数据库
	if !model.CheckUserExist(userdto.Phone){
		err := model.CreateUserWithPhone(userdto.Phone)
		if err!= nil {
			return result.Fail("用户创建失败")
		}
	}
	//保存用户信息到session
	session.Set("user", userdto.Phone)
	if err := session.Save(); err!= nil {
		return result.Fail("用户信息保存失败")
	}

	//返回响应
	return result.Ok()

}