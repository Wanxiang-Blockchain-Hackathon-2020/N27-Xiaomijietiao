package proto

type ErrCode int32

const (
	ErrCode_ECSuccess       ErrCode = 0
	ErrCode_ECParamErr      ErrCode = 1
	ErrCode_ECNotExistErr   ErrCode = 2
	ErrCode_ECUuidCreateErr ErrCode = 3
	ErrCode_ECDbErr         ErrCode = 4
)

type Error struct {
	Code ErrCode
	Desc string
}

type ReqHead struct {
	InvokeSource uint32 //请求来源
}
