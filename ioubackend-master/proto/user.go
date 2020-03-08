package proto

// 用户实名状态
type UserAuthStatus int32

const (
	UserAuthStatus_Unkown   UserAuthStatus = 0 //未实名
	UserAuthStatus_AuthSucc UserAuthStatus = 1 //实名成功
	UserAuthStatus_AuthFail UserAuthStatus = 2 //实名失败
)

// 用户ID类型
type UserIdType int32

const (
	UserIdType_OpenId UserIdType = 0 //采用微信openid
	UserIdType_Normal UserIdType = 1 //采用普通ID
)

type User struct {
	Id          string         //用户ID
	OpenId      string         //用户OPENID
	AvataUrl    string         //用户头像
	City        string         //城市
	Gender      int32          //性别
	NickName    string         //昵称
	Country     string         //国家
	Province    string         //省份
	Name        string         //真实名字
	IDCardId    string         //身份证ID
	MobilePhone string         //手机号码
	Auth        UserAuthStatus //实名状态
	EqianId     string         //对应e签宝的账户ID
	IDCardPic   string         //身份证正面照Base64编码后的字符串
	IdType      UserIdType     //用户ID类型
	ESignId     string         //此用户在e签宝的账户ID
}
