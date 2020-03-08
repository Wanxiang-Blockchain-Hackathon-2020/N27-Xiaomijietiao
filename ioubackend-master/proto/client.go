package proto

// 查询用户信息
// GET URL Path: /api/v1/users/query?openid=xxx
type UserRsp struct {
	Error    *Error //错误信息，成功则空
	UserInfo *User  //用户信息
}

// 查询用户关联的借条
// POST URL Path: /api/v1/users/query_iou_of_user
type UserIouReq struct {
	Head    *ReqHead //请求头
	UserId  string   //查询的用户ID
	OpenId  string   //查询的用户OPENID
	LastCnt uint32   //拉最近几条, 不填则返回全部，客户端分页(?)
	From    bool     //拉取借出的
	To      bool     //拉取借入的
}

type IouWithUserInfo struct {
	Iou  *Iou  //借条
	From *User //出借人信息
	To   *User //借款人信息
}

type UserIouRsp struct {
	Error  *Error             //错误信息，成功则空
	Id     string             //对应用户的ID
	Openid string             //对应用户的OPENID
	Ious   []*IouWithUserInfo //关联的借条
}

// 更新状态请求
// POST URL Path: /api/v1/ious/update/status
type UpdateStatusReq struct {
	Head   *ReqHead  //请求头
	IouId  string    //借条ID
	Status IouStatus //借条状态
}

type UpdateStatusRsp struct {
	Error *Error //错误信息，成功则空
}

// 添加借条信息
// POST URL Path: /api/v1/ious/add
type UpsertIouReq struct {
	Head *ReqHead
	Iou  *Iou
}

type UpsertIouRsp struct {
	Error *Error //错误信息，成功则空
	Iou   *Iou   //借条信息
}

// 请求借条详情信息
// GET URL Path: /api/v1/ious/query?iouid=xxxxx
type QueryIouRsp struct {
	Error   *Error           //错误信息，成功则空
	IouInfo *IouWithUserInfo //借条详情，包括借款双方的信息
}

// 签署的流程
// 请求合同签署
// POST URL Path: /api/v1/iou
type SignReq struct {
	IouId       string //借条ID
	RedirectUrl string //签署完成回调url
	UserId      string //签署方的用户ID(即OpenId)_
}

// 返回签署流程地址，签署页面地址
type SignRsp struct {
	Error   *Error //错误信息，成功则空
	SignUrl string //签署URL
}

// 请求签署成功后的文件下载地址
// POST URL Path: /api/v1/ious/sign/getDowloadUrl
type ContractDownloadUrlReq struct {
	IouId string //借条ID
}

// 返回签署后的合同地址
type ContractDownloadUrlRsp struct {
	DownloadUrl string //下载地址
}
