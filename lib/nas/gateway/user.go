package gateway

type User struct {
	OpenId       string `json:"open_id"`  //用户唯一id
	UnionId      string `json:"union_id"` //微信union_id
	Channel      string `json:"channel"`  //登录类型
	Nickname     string `json:"nickname"` //昵称
	Gender       string `json:"gender"`   //0=>未知 1=>男 2=>女   twitter和line不会返回性别，所以这里是0，Facebook根据你的权限，可能也不会返回，所以也可能是0
	Avatar       string `json:"avatar"`   //头像
	Type         string `json:"type"`     //授权类型
}
