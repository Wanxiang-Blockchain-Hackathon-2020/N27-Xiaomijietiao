package proto

// 借条状态
type IouStatus int32

const (
	IouStatus_Create     IouStatus = 0 //新建
	IouStatus_Confirm    IouStatus = 1 //待确认出款人
	IouStatus_SignCreate IouStatus = 2 //待签署流程创建
	IouStatus_SignPartA  IouStatus = 3 //待出款人签名
	IouStatus_SignPartB  IouStatus = 4 //待借款人签名
	IouStatus_Pay        IouStatus = 5 //待出款
	IouStatus_Paid       IouStatus = 6 //待还款
	IouStatus_End        IouStatus = 7 //已结清
)

// 还款方式
type PayType int32

const (
	PayType_PTPrincipalAndInterest PayType = 0
)

// 借条结构
type Iou struct {
	Id           string    //借条ID, 采用UUID方案
	Amount       uint32    //借条对应的金额
	CreateAt     int64     //借条产生时间,时间戳
	BorrowAt     int64     //借条借款时间,时间戳
	PayBackAt    int64     //借条约定还款时间,时间戳
	InterestRate uint32    //借条利率，基数10000
	From         string    //借条对应的出借人
	To           string    //借条对应的借款人
	Status       IouStatus //借条状态
	PayBackType  PayType   //还款方式
	FlowId       string    //对应签署的流程ID
	FileId       string    //对应借条的pdf,实际文件存e签宝
	ContentMD5   string    //合同文件MD5，先去e签宝上传文件生成FileId，接下来调用上传文件接口上传
	UploadUrl    string    //合同上传URL
	FileSize     int32     //文件大小
}

// 借条状态流转流水
type IouFlow struct {
	FlowId     string    //流水ID
	IouId      string    //流水关联的借条ID
	PreStatus  IouStatus //流转前状态
	PostStatus IouStatus //流转后状态
	Timestamp  int64     //发生时间，服务器时间
}
