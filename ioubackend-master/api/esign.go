package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"

	simplejson "github.com/bitly/go-simplejson"
	. "github.com/despard/ioubackend/config"
	db "github.com/despard/ioubackend/db"
	pbi "github.com/despard/iouproto/iou"
	log "github.com/despard/log"
	pb "github.com/golang/protobuf/proto"
)

var (
	HTTP_PROTO            = "https://"
	ESIGN_DOMAIN          = "test.esign.cn"
	URL_AUTH              = HTTP_PROTO + ESIGN_DOMAIN + "/v1/oauth2/access_token"
	URL_CREATE_ACCOUNT    = HTTP_PROTO + ESIGN_DOMAIN + "/v1/accounts"
	URL_CREATE_SEAL       = HTTP_PROTO + ESIGN_DOMAIN + "/v1/accounts/{accountId}/seals/personaltemplate"
	URL_CREATE_FLOW       = HTTP_PROTO + ESIGN_DOMAIN + "/v1/signflows"
	URL_QUERY_FLOW        = HTTP_PROTO + ESIGN_DOMAIN + "/v1/signflows"
	URL_GETUPLOAD         = HTTP_PROTO + ESIGN_DOMAIN + "/v1/files/getUploadUrl"
	URL_ATTACH_DOC        = HTTP_PROTO + ESIGN_DOMAIN + "/v1/signflows/%s/documents"
	URL_HAND_SIGN         = HTTP_PROTO + ESIGN_DOMAIN + "/v1/signflows/%s/signfields/handSign"
	URL_TEMPLATE_CONTRACT = HTTP_PROTO + ESIGN_DOMAIN + "/v1/files/createByTemplate"
	URL_TEMPLATE_INFO     = HTTP_PROTO + ESIGN_DOMAIN + "/v1/docTemplates/%s"
	URL_START_FLOW        = HTTP_PROTO + ESIGN_DOMAIN + "/v1/signflows/%s/start"
	URL_EXTRACT_URL       = HTTP_PROTO + ESIGN_DOMAIN + "/v1/signflows/%s/executeUrl?accountId=%s"
	URL_OCR               = HTTP_PROTO + ESIGN_DOMAIN + "/v2/identity/auth/api/ocr/idcard"
	URL_VERIFY_FACTOR     = HTTP_PROTO + ESIGN_DOMAIN + "/v2/identity/verify/individual/telecom3Factors"
	URL_AUTH_FACE         = HTTP_PROTO + ESIGN_DOMAIN + "/v2/identity/auth/api/individual/%s/face"
	URL_AUTH_TELECOM      = HTTP_PROTO + ESIGN_DOMAIN + "/v2/identity/auth/api/individual/%s/telecom3Factors"
	URL_AUTH_VCODE        = HTTP_PROTO + ESIGN_DOMAIN + "/v2/identity/auth/pub/individual/%s/telecom3Factors"
	URL_DOCS_DOWNLOAD     = HTTP_PROTO + ESIGN_DOMAIN + "/v1/signflows/%s/documents"
	URL_PREVIEW_URL       = HTTP_PROTO + ESIGN_DOMAIN + "/v1/signflows/%s/executeUrl?accountId=%s&urlType=1"
)

var AccessToken string = ""
var TokenExpire int64
var lock sync.Mutex
var HttpCli *http.Client = &http.Client{}

func GetAccessToken(appid string, appsecret string) (string, int64) {
	log.Debug("Start get new accesstoken from esign")
	url := fmt.Sprintf("%s?appId=%s&secret=%s&grantType=client_credentials", URL_AUTH, appid, appsecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", int64(0)
	}

	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", int64(0)
	}
	log.Debug("AccessToken res: %s", string(res))

	js, err := simplejson.NewJson(res)
	if err != nil {
		return "", int64(0)
	}

	accesstoken, _ := js.Get("data").Get("token").String()
	strExpire, _ := js.Get("data").Get("expiresIn").String()
	tokenexpire, _ := strconv.ParseInt(strExpire, 10, 64)
	log.Debug("Token expire in %d", tokenexpire/1000)
	return accesstoken, int64(tokenexpire / 1000)
}

func SetAuthHeader(req *http.Request) {
	now := time.Now().Unix()
	if TokenExpire <= now {
		lock.Lock()
		if TokenExpire <= (now - 600) {
			AccessToken, TokenExpire = GetAccessToken(Settings.Appid, Settings.AppSecret)
		}
		lock.Unlock()
	}

	req.Header.Add("X-Tsign-Open-Token", AccessToken)
	req.Header.Add("X-Tsign-Open-App-Id", Settings.Appid)
}

func EsignRequest(method string, url string, body []byte) (*simplejson.Json, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	SetAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := HttpCli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error("Esign Request status code : %d", resp.StatusCode)
		return nil, errors.New(fmt.Sprintf("Http Resp Code:%s", resp.Status))

	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	res, err := simplejson.NewJson(respBody)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func EsignRequestWithHeader(method string, uri string, body []byte, headers http.Header) (http.Header, []byte, error) {
	url := HTTP_PROTO + ESIGN_DOMAIN + uri
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	SetAuthHeader(req)
	for key, values := range headers {
		for _, v := range values {
			req.Header.Set(key, v)
		}
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := HttpCli.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	return resp.Header, res, err
}

//创建e签宝的账户ID
func CreateEsignAccount(userId string) (string, error) {
	user, err := db.QueryUser(userId)
	if err != nil {
		return "", err
	}

	ms := map[string]interface{}{}
	ms["idNumber"] = user.GetIDCardId()
	ms["idType"] = "CRED_PSN_CH_IDCARD"
	ms["mobile"] = user.GetMobilePhone()
	ms["name"] = user.GetName()
	//	ms["thirdPartyUserId"] = user.GetOpenId()
	body, _ := json.Marshal(ms)
	fmt.Println(string(body))

	res, err := EsignRequest("POST", URL_CREATE_ACCOUNT, body)
	if err != nil {
		return "", err
	}

	accountId, _ := res.Get("data").Get("accountId").String()

	return accountId, nil
}

//创建e签宝的个人签章
func CreateEsignSeal(userId string) (string, error) {

	return "", nil
}

//创建一个签署流程，由出借方创建
func CreateEsignFlow(iou *pbi.Iou, redirectUrl string) (string, error) {
	user, err := db.QueryUser(iou.GetFrom())
	if err != nil {
		return "", err
	}
	if user.GetAuth() != pbi.UserAuthStatus_AuthSucc {
		return "", errors.New("Unauthoried user.")
	}
	if user.GetESignId() == "" {
		return "", errors.New("No esign account.")
	}

	ms := map[string]interface{}{}
	ms["businessScene"] = "借款合同"
	ms["initiatorAccountId"] = user.GetESignId()
	ms["autoArchive"] = true //自动归档

	config := map[string]interface{}{}
	config["noticeType"] = ""
	config["signPlatform"] = "1"
	config["redirectUrl"] = redirectUrl
	config["noticeDeveloperUrl"] = fmt.Sprintf("%s/sign?iouid=%s", Settings.NoticeUrl, iou.GetId())

	ms["configInfo"] = config

	body, _ := json.Marshal(ms)

	res, err := EsignRequest("POST", URL_CREATE_FLOW, body)
	if err != nil {
		return "", err
	}

	if code, _ := res.Get("code").Int(); code != 0 {
		return "", errors.New(fmt.Sprintf("msg: %s", res.Get("message")))
	}

	flowId, _ := res.Get("data").Get("flowId").String()
	return flowId, nil
}

//查询一个签署流程
func QueryEsignFlow(flowId string) (string, error) {
	url := fmt.Sprintf("%s/%s", URL_QUERY_FLOW, flowId)
	res, err := EsignRequest("GET", url, []byte{})
	if err != nil {
		return "", err
	}
	bs, err := res.MarshalJSON()
	return string(bs), err
}

//创建一个合同文件
/*
func CreateEsignContract(contract *pbi.AddContractReq) (string, string, error) {
	ms := map[string]interface{}{}
	ms["contentMd5"] = contract.GetContentMD5()
	ms["fileName"] = contract.GetFileName()
	ms["fileSize"] = contract.GetFileSize()
	ms["contentType"] = "application/octet-stream"
	ms["convert2Pdf"] = true
	body, _ := json.Marshal(ms)

	res, err := EsignRequest("POST", URL_GETUPLOAD, body)
	if err != nil {
		return "", "", err
	}

	if code, _ := res.Get("code").Int(); code != 0 {
		return "", "", errors.New(fmt.Sprintf("msg: %s", res.Get("message")))
	}

	fileid, _ := res.Get("data").Get("fileId").String()
	uploadUrl, _ := res.Get("data").Get("uploadUrl").String()

	return fileid, uploadUrl, nil
}
*/

//上传合同文件
func UploadContract(iou *pbi.Iou, file []byte) error {
	fileId := iou.GetFileId()
	uploadUrl := iou.GetUploadUrl()
	filesize := iou.GetFileSize()
	md5 := iou.GetContentMD5()

	if fileId == "" || uploadUrl == "" || filesize == 0 || md5 == "" {
		return errors.New("Upload failed.Param invalid.")
	}

	req, err := http.NewRequest("PUT", uploadUrl, bytes.NewBuffer(file))
	SetAuthHeader(req)
	req.Header.Set("Content-MD5", md5)
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := HttpCli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Http Resp Code:%s", resp.Status))

	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	res, err := simplejson.NewJson(respBody)
	if err != nil {
		return err
	}

	if rmd5, _ := res.Get("Content-MD5").String(); rmd5 != md5 {
		return errors.New("Upload contract failed.")
	}

	return nil
}

//添加签署合同到签署流程
func AttachDocument(iouid string) error {
	iou, err := db.QueryIou(iouid)
	if err != nil {
		return err
	}

	ms := map[string]interface{}{}
	doc := map[string]interface{}{}
	doc["encryption"] = 0
	doc["fileId"] = iou.GetFileId()
	ms["docs"] = []map[string]interface{}{doc}

	body, _ := json.Marshal(ms)

	res, err := EsignRequest("POST", fmt.Sprintf(URL_ATTACH_DOC, iou.GetFlowId()), body)
	if err != nil {
		return err
	}

	if code, _ := res.Get("code").Int(); code != 0 {
		return errors.New(fmt.Sprintf("msg: %s", res.Get("message")))
	}

	return nil
}

//查询模板信息
func QueryTemplateInfo(tid string) (*simplejson.Json, error) {
	tinfo, err := EsignRequest("GET", fmt.Sprintf(URL_TEMPLATE_INFO, tid), []byte{})
	if err != nil {
		return nil, err
	}

	if code, _ := tinfo.Get("code").Int(); code != 0 {
		return nil, errors.New(fmt.Sprintf("msg: %s", tinfo.Get("message")))
	}

	return tinfo, err
}

//通过模版创建合同
func CreateTemplateContract(iouid string) (string, error) {
	iou, err := db.QueryIou(iouid)
	if err != nil {
		return "", err
	}
	fromUser, _ := db.QueryUser(iou.GetFrom())
	toUser, _ := db.QueryUser(iou.GetTo())
	if fromUser == nil || toUser == nil {
		return "", errors.New("缺乏借款双方人信息")
	}

	simpleFormFields := map[string]interface{}{}
	simpleFormFields["f1"] = fromUser.GetName()
	simpleFormFields["f3"] = toUser.GetName()
	simpleFormFields["f2"] = fromUser.GetIDCardId()
	simpleFormFields["f4"] = toUser.GetIDCardId()
	simpleFormFields["f5"] = Transfer(int(iou.GetAmount()))
	simpleFormFields["f6"] = iou.GetAmount()
	bTime := time.Unix(iou.GetBorrowAt(), 0)
	pTime := time.Unix(iou.GetPayBackAt(), 0)
	cTime := time.Unix(iou.GetCreateAt(), 0)
	simpleFormFields["f7"] = bTime.Year()
	simpleFormFields["f8"] = bTime.Month()
	simpleFormFields["f9"] = bTime.Day()
	simpleFormFields["f10"] = pTime.Year()
	simpleFormFields["f11"] = pTime.Month()
	simpleFormFields["f12"] = pTime.Day()

	pentaly := float32(iou.GetInterestRate()) * 1.5
	rate := fmt.Sprintf("万分之%.2f", float64(iou.GetInterestRate())/10000)
	pentalyRate := fmt.Sprintf("万分之%.2f", pentaly/10000)
	days := pTime.Sub(bTime).Hours() / 24
	interest := fmt.Sprintf("%.2f",
		days*float64(iou.GetInterestRate())/10000*float64(iou.GetAmount())/10000)

	simpleFormFields["f13"] = int(days)
	simpleFormFields["f14"] = rate
	simpleFormFields["f15"] = interest
	//simpleFormFields["f17"] = interestPerDay
	simpleFormFields["f16"] = pentalyRate

	simpleFormFields["f1-2"] = 3
	simpleFormFields["f2-2"] = cTime.Year()
	simpleFormFields["f3-2"] = cTime.Month()
	simpleFormFields["f4-2"] = cTime.Day()
	simpleFormFields["f5-2"] = cTime.Year()
	simpleFormFields["f6-2"] = cTime.Month()
	simpleFormFields["f7-2"] = cTime.Day()

	ms := map[string]interface{}{}
	ms["name"] = fmt.Sprintf("%s", "借条合同")
	ms["templateId"] = Settings.TemplateId
	ms["simpleFormFields"] = simpleFormFields

	body, _ := json.Marshal(ms)
	res, err := EsignRequest("POST", URL_TEMPLATE_CONTRACT, body)
	if err != nil {
		return "", err
	}

	if code, _ := res.Get("code").Int(); code != 0 {
		return "", errors.New(fmt.Sprintf("msg: %s", res.Get("message")))
	}

	fileId, err := res.Get("data").Get("fileId").String()
	log.Info("Create file from template succ.")

	return fileId, err
}

//开启签署流程
func StartEsignFlow(sign *pbi.SignReq) error {
	uid := sign.GetUserId()
	user, err := db.QueryUser(uid)
	if err != nil {
		return err
	}

	//该用户是否已实名
	if user.GetAuth() != pbi.UserAuthStatus_AuthSucc {
		return errors.New("用户未实名认证")
	}

	//生成合同文件
	fileId, err := CreateTemplateContract(sign.GetIouId())
	if err != nil {
		log.Error("Create iou %s contract from template failed.", sign.GetIouId())
		return err
	}

	//DB回写fileId
	if err := db.UpdateIou(&pbi.Iou{
		Id:     pb.String(sign.GetIouId()),
		FileId: pb.String(fileId),
	}); err != nil {
		log.Error("Update iou %s status failed", sign.GetIouId())
		return err
	}

	//添加流程文档
	if err := AttachDocument(sign.GetIouId()); err != nil {
		log.Error("Attach document to flow failed. iou %s", sign.GetIouId())
		return err
	}

	iou, err := db.QueryIou(sign.GetIouId())
	if err != nil {
		return err
	}

	//添加合同签署区
	if err := AddHandSign(iou); err != nil {
		log.Error("Add sign hand area of iou %s failed.", iou.GetId())
		return err
	}

	//开启签署流程
	res, err := EsignRequest("PUT", fmt.Sprintf(URL_START_FLOW, iou.GetFlowId()), []byte{})
	if err != nil {
		return err
	}

	if code, _ := res.Get("code").Int(); code != 0 {
		return errors.New(fmt.Sprintf("msg: %s", res.Get("message")))
	}

	return nil
}

//添加双方的签署区
func AddHandSign(iou *pbi.Iou) error {
	if iou == nil {
		return errors.New("Iou is nil.")
	}

	fromUser, _ := db.QueryUser(iou.GetFrom())
	toUser, _ := db.QueryUser(iou.GetTo())
	if fromUser == nil || toUser == nil {
		return errors.New("缺乏借款双方人信息")
	}

	//查询模版信息
	tinfo, err := QueryTemplateInfo(Settings.TemplateId)
	if err != nil {
		log.Error("Query iou %s contract template (%s) failed.", iou.GetId(), Settings.TemplateId)
		return err
	}

	log.Info("Template info :%v", tinfo)

	fromSignArea := map[string]interface{}{}
	toSignArea := map[string]interface{}{}
	fromSignArea["posPage"] = 4
	fromSignArea["posX"] = 179.0
	fromSignArea["posY"] = 590.0
	toSignArea["posPage"] = 4
	toSignArea["posX"] = 478.0
	toSignArea["posY"] = 590.0
	/*
		components, _ := tinfo.Get("structComponents").Array()
		for i, _ := range components {
			comp := tinfo.Get("structComponents").GetIndex(i)
			if ar, _ := comp.Get("key").String(); ar == "FromSignArea" {
				fromSignArea["posPage"], _ = comp.Get("context").Get("pos").Get("page").Int()
				fromSignArea["posX"], _ = comp.Get("context").Get("pos").Get("x").Float64()
				fromSignArea["posY"], _ = comp.Get("context").Get("pos").Get("y").Float64()
			}
			if ar, _ := comp.Get("key").String(); ar == "ToSignArea" {
				toSignArea["posPage"], _ = comp.Get("context").Get("pos").Get("page").Int()
				toSignArea["posX"], _ = comp.Get("context").Get("pos").Get("x").Float64()
				toSignArea["posY"], _ = comp.Get("context").Get("pos").Get("y").Float64()
			}
		}
	*/

	log.Debug("AreaSign:%v, %v", fromSignArea, toSignArea)

	signfields := []map[string]interface{}{}
	//甲方签署区
	signfield1 := map[string]interface{}{}
	signfield1["fileId"] = iou.GetFileId()
	signfield1["signerAccountId"] = fromUser.GetESignId()
	signfield1["posBean"] = fromSignArea
	//signfield1 =
	signfield1["sealType"] = "1" //个人签章
	signfield1["signType"] = 1

	//乙方签署区
	signfield2 := map[string]interface{}{}
	signfield2["fileId"] = iou.GetFileId()
	signfield2["signerAccountId"] = toUser.GetESignId()
	signfield2["posBean"] = toSignArea
	signfield2["sealType"] = "1" //个人签章
	signfield2["signType"] = 1

	signfields = append(signfields, signfield1)
	signfields = append(signfields, signfield2)

	ms := map[string]interface{}{}
	ms["signfields"] = signfields
	ms["flowId"] = iou.GetFlowId()
	body, _ := json.Marshal(ms)

	res, err := EsignRequest("POST", fmt.Sprintf(URL_HAND_SIGN, iou.GetFlowId()), body)
	if err != nil {
		return err
	}

	if code, _ := res.Get("code").Int(); code != 0 {
		return errors.New(fmt.Sprintf("msg: %s", res.Get("message")))
	}
	return nil
}

//获取签署地址
func GetSignUrl(sign *pbi.SignReq) (string, error) {
	if sign == nil {
		return "", errors.New("Request is nil")
	}

	iou, erri := db.QueryIou(sign.GetIouId())
	if erri != nil {
		log.Error("Query iou %s failed.err:%v", sign.GetIouId(), erri)
		return "", erri
	}

	var flowId string = iou.GetFlowId()
	var shortUrl string
	var url string
	var err error
	switch iou.GetStatus() {
	case pbi.IouStatus_Create, pbi.IouStatus_Confirm, pbi.IouStatus_SignCreate:
		//创建流程
		flowId, err = CreateEsignFlow(iou, sign.GetRedirectUrl())
		if err != nil {
			log.Error("Create esign flow of iou %s failed.err:%v", iou.GetId(), err)
			return "", err
		}
		//回写DB
		if flowId != "" {
			iou.FlowId = pb.String(flowId)
			iou.Status = pbi.IouStatus_SignCreate.Enum()
			if err := db.UpdateIou(iou); err != nil {
				log.Error("Update iou %s flowid %s failed.err:%v", iou.GetId(), flowId, err)
				return "", err
			}
		}
		//开启流程
		if err := StartEsignFlow(sign); err != nil {
			log.Error("Start sign flow %s of iou %s failed.err:%v", flowId, iou.GetId(), err)
			return "", err
		}
		//回写DB
		iou.FlowId = pb.String(flowId)
		iou.Status = pbi.IouStatus_SignPartA.Enum()
		if err := db.UpdateIou(iou); err != nil {
			log.Error("Update iou %s flowid %s failed.err:%v", iou.GetId(), flowId, err)
			return "", err
		}
		fallthrough
	case pbi.IouStatus_SignPartA, pbi.IouStatus_SignPartB:
		accountId := ""
		if iou.GetStatus() == pbi.IouStatus_SignPartA {
			user, err := db.QueryUser(iou.GetFrom())
			if err != nil {
				log.Error("Get iou %s signer accountId failed.", iou.GetId())
				return "", err
			}

			accountId = user.GetESignId()
		} else {
			user, err := db.QueryUser(iou.GetTo())
			if err != nil {
				log.Error("Get iou %s signer accountId failed.", iou.GetId())
				return "", err
			}

			accountId = user.GetESignId()
		}
		res, err := EsignRequest("GET", fmt.Sprintf(URL_EXTRACT_URL, flowId, accountId), []byte{})
		if err != nil {
			return "", err
		}

		if code, _ := res.Get("code").Int(); code != 0 {
			return "", errors.New(fmt.Sprintf("msg: %s", res.Get("message")))
		}

		shortUrl, _ = res.Get("data").Get("shortUrl").String()
		url, _ = res.Get("data").Get("url").String()
		log.Debug("Url: %s, ShortUrl: %s", url, shortUrl)
	default:
		log.Info("Iou %s Sign flow finished.", sign.GetIouId())
		return "", errors.New("Sign flow finished")
	}

	return url, nil
}

func EsignFaceAuth(rq *simplejson.Json) ([]byte, error) {
	if rq == nil {
		return []byte{}, errors.New("Request body nil")
	}
	openid, _ := rq.Get("openid").String()
	user, err := db.QueryUser(openid)
	if err != nil {
		return []byte{}, err
	}

	accountId := user.GetESignId()
	rq.Set("contextId", openid)
	rq.Set("notifyUrl", fmt.Sprintf("%s/auth", Settings.NoticeUrl))

	body, _ := rq.MarshalJSON()
	res, err := EsignRequest("POST", fmt.Sprintf(URL_AUTH_FACE, accountId), body)
	if err != nil {
		return []byte{}, err
	}

	return res.MarshalJSON()
}

func EsignTelecomAuth(rq *simplejson.Json) ([]byte, error) {
	if rq == nil {
		return []byte{}, errors.New("Request body nil")
	}
	openid, _ := rq.Get("openid").String()
	user, err := db.QueryUser(openid)
	if err != nil {
		return []byte{}, err
	}

	accountId := user.GetESignId()
	rq.Set("contextId", openid)
	rq.Set("notifyUrl", fmt.Sprintf("%s/auth", Settings.NoticeUrl))

	body, _ := rq.MarshalJSON()
	res, err := EsignRequest("POST", fmt.Sprintf(URL_AUTH_TELECOM, accountId), body)
	if err != nil {
		return []byte{}, err
	}

	if code, _ := res.Get("code").Int(); code == 30500153 || code == 30500152 {
		//设置实名
		db.UpdateUserAuthStatus(openid, 1)
	}

	return res.MarshalJSON()
}

//验证实名验证码
func EsignVeriCode(rq *simplejson.Json) ([]byte, error) {
	if rq == nil {
		return []byte{}, errors.New("Request body nil")
	}
	flowId, _ := rq.Get("flowId").String()
	openId, _ := rq.Get("openId").String()
	user, err := db.QueryUser(openId)
	if err != nil {
		log.Error("Auth target user %s is not found.flow %s,%v",
			openId, flowId, err)
		return []byte{}, err
	}
	rq.Del("flowId")
	body, _ := rq.MarshalJSON()
	res, err := EsignRequest("PUT", fmt.Sprintf(URL_AUTH_VCODE, flowId), body)
	if err != nil {
		return []byte{}, err
	}

	code, err := res.Get("code").Int()

	if err != nil || code != 0 {
		log.Info("Auth flow %s failed.", flowId)
		user.Auth = pbi.UserAuthStatus_AuthFail.Enum()
	} else {
		user.Auth = pbi.UserAuthStatus_AuthSucc.Enum()
	}

	db.AddUserToDB(user)

	return res.MarshalJSON()
}

func EsignGetDownloadUrl(iouid string) ([]byte, error) {
	iou, err := db.QueryIou(iouid)
	if err != nil {
		return []byte{}, err
	}

	flowId := iou.GetFlowId()

	resp, err := EsignRequest("GET", fmt.Sprintf(URL_DOCS_DOWNLOAD, flowId), []byte{})
	if err != nil {
		return []byte{}, err
	}

	return resp.MarshalJSON()
}

func EsignGetPreviewUrl(iouid string, openid string) ([]byte, error) {
	iou, err := db.QueryIou(iouid)
	if err != nil {
		return []byte{}, err
	}

	flowId := iou.GetFlowId()

	user, err1 := db.QueryUser(openid)
	if err1 != nil {
		return []byte{}, err1
	}

	accountId := user.GetESignId()

	res, err2 := EsignRequest("GET", fmt.Sprintf(URL_PREVIEW_URL, flowId, accountId), []byte{})
	if err2 != nil {
		return []byte{}, err2
	}

	return res.MarshalJSON()
}
