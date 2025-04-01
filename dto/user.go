package dto
type UserDTO struct {
	Phone string `json:"phone"`
	Code string `json:"code"`
	Password string `json:"password"`
}
type UserVo struct{
	Id uint64 `json:"id"`
	NickName string `json:"nickName"`
	Icon string `json:"icon"`

}