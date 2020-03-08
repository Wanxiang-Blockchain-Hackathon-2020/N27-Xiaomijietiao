package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	. "github.com/despard/ioubackend/config"

	simplejson "github.com/bitly/go-simplejson"
	db "github.com/despard/ioubackend/db"
	pbc "github.com/despard/iouproto/common"
	pbi "github.com/despard/iouproto/iou"
	"github.com/despard/log"
	pb "github.com/golang/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/oauth2.v3/server"
	stat "github.com/despard/stat"
)

type HandleFunc func (http.ResponseWriter, *http.Request, []byte) ([]byte, error) 
type HandleFunc2 func (http.ResponseWriter, *http.Request, httprouter.Params) ([]byte, error) 

func SetCommRespHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func HandleAccess(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr("ioubackend.http.access", 1)
	accessRsp := map[string]interface{}{}
	accessRsp["Version"] = Settings.Version
	resp, _ := json.Marshal(accessRsp)
	io.WriteString(w, string(resp))
}

func HandleIssue(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr("ioubackend.http.issue", 1)
	issueRsp := map[string]interface{}{}
	issueRsp["Code"] = 0
	issueRsp["Desc"] = "success."
	resp, _ := json.Marshal(issueRsp)
	io.WriteString(w, string(resp))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))
	issue := &db.Issue{}
	if err := json.Unmarshal(body, issue); err != nil {
		log.Error("Unknow format.%v", err)
		return
	}

	db.AddIssueToDB(issue)
}


func HandleBanners(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive banners get.")
	bannerList := []map[string]interface{}{}

	list, err := db.GetBannerList()
	if err != nil {
		w.WriteHeader(404)
		return 
	}

	for i := 0; i < len(list); i++ {
		banner := map[string]interface{}{}
		banner["path"] = list[i]
		banner["tag"] = "banner"
		bannerList = append(bannerList, banner)
	}

	resp, _ := json.Marshal(bannerList)
	io.WriteString(w, string(resp))
}

func HandleAddUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive user Post request")
	SetCommRespHeader(w)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))
	user := &pbi.User{}
	if err := json.Unmarshal(body, user); err != nil {
		log.Error("Unknow format.%v", err)
		return
	}
	if us,_ := db.QueryUser(user.GetOpenId()); us == nil {
		user.Auth = pbi.UserAuthStatus_Unkown.Enum()
	}

	if err := db.AddUserToDB(user); err != nil {
		log.Error("Insert user record err:%v", err)
		return
	}
	res, _ := json.Marshal(pbc.Error{Code: pb.Int32(int32(pbc.ErrCode_ECSuccess)), Desc: pb.String("success.")})

	io.WriteString(w, string(res))

}

func HandleQueryUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive user Get request")
	SetCommRespHeader(w)
	q := req.URL.Query()
	log.Debug("%v", q)
	openid := q.Get("openid")
	if openid == "" {
		res, _ := json.Marshal(pbc.Error{Code: pb.Int32(int32(pbc.ErrCode_ECParamErr)), Desc: pb.String("Invaild Params.")})
		w.WriteHeader(400)
		io.WriteString(w, string(res))
		return
	}

	user, err := db.QueryUser(openid)
	if err != nil {
		res, _ := json.Marshal(pbc.Error{Code: pb.Int32(int32(pbc.ErrCode_ECNotExistErr)), Desc: pb.String("User not exsit.")})
		w.WriteHeader(200)
		io.WriteString(w, string(res))
		return
	}

	rsp := &pbi.UserRsp{
		UserInfo: user,
	}

	res, _ := json.Marshal(rsp)

	io.WriteString(w, string(res))

}

func HandleQueryIouOfUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Verbose("Receive query IOU request")
	SetCommRespHeader(w)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))

	queryIou := &pbi.UserIouReq{}
	if err := json.Unmarshal(body, queryIou); err != nil {
		log.Error("Unknow format.%v", err)
		return
	}

	log.Verbose("Query user(openid:%s)'s relate IOUs", queryIou.GetOpenId())
	ious, err := db.QueryIousOfUser(queryIou)
	if err != nil {
		log.Error("Query user(openid:%s) IOUs error:%v", queryIou.GetOpenId(), err)
	}

	ius := []*pbi.IouWithUserInfo{}

	for _, iou := range ious {
		log.Verbose("ious of user(openid:%s): %+v", queryIou.GetOpenId(), iou)
		from, err1 := db.QueryUser(iou.GetFrom())
		if err1 != nil {
			log.Error("Query from person msg failed:%v", err)
		}
		to, err2 := db.QueryUser(iou.GetTo())
		if err2 != nil {
			log.Error("Query to person msg failed:%v", err)
		}
		iu := &pbi.IouWithUserInfo{
			Iou:  iou,
			From: from,
			To:   to,
		}
		ius = append(ius, iu)
	}

	IOURsp := &pbi.UserIouRsp{
		Openid: queryIou.OpenId,
		Ious:   ius,
	}

	res, _ := json.Marshal(IOURsp)
	io.WriteString(w, string(res))

}

func HandleAddIou(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive Add IOU request")
	SetCommRespHeader(w)

	res := &pbi.UpsertIouRsp{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))

	iouReq := &pbi.UpsertIouReq{}
	if err := json.Unmarshal(body, iouReq); err != nil {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(pbc.ErrCode_ECParamErr)),
			Desc: pb.String("Params invalid."),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return
	}

	//生成UUID
	id := uuid.NewV4()
	iou := iouReq.GetIou()
	iou.Id = pb.String(id.String())

	if err := db.CreateIou(iou); err != nil {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(3)),
			Desc: pb.String(fmt.Sprintf("DB error.%s", err.Error())),
		}
	}

	res.Iou = iou

	resp, _ := json.Marshal(res)
	io.WriteString(w, string(resp))

	return

}

func HandleUpdateIou(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive Update IOU request")
	SetCommRespHeader(w)
	res := &pbi.UpsertIouRsp{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))

	iouReq := &pbi.UpsertIouReq{}
	if err := json.Unmarshal(body, iouReq); err != nil {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(pbc.ErrCode_ECParamErr)),
			Desc: pb.String("Params invalid."),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return
	}

	if iouReq.GetIou().GetId() == "" {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(pbc.ErrCode_ECParamErr)),
			Desc: pb.String("Params invalid."),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return
	}

	if err := db.UpdateIou(iouReq.GetIou()); err != nil {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(pbc.ErrCode_ECDbErr)),
			Desc: pb.String(err.Error()),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return
	}

	iou, _ := db.QueryIou(iouReq.GetIou().GetId())
	res.Iou = iou

	resp, _ := json.Marshal(res)
	io.WriteString(w, string(resp))
}

func HandleAgreeIou(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive Update IOU request")
	SetCommRespHeader(w)
	res := &pbi.UpsertIouRsp{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))

	iouReq := &pbi.UpsertIouReq{}
	if err := json.Unmarshal(body, iouReq); err != nil {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(pbc.ErrCode_ECParamErr)),
			Desc: pb.String("Params invalid."),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return
	}

	if iouReq.GetIou().GetId() == "" {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(pbc.ErrCode_ECParamErr)),
			Desc: pb.String("Params invalid."),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return
	}

	if err := db.AgreeIou(iouReq.GetIou()); err != nil {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(pbc.ErrCode_ECDbErr)),
			Desc: pb.String(err.Error()),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return
	}

	iou, _ := db.QueryIou(iouReq.GetIou().GetId())
	if iou.GetFrom() != iouReq.GetIou().GetFrom() {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(pbc.ErrCode_ECFromExistErr)),
			Desc: pb.String("iou have been modified before."),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return
	}

	res.Iou = iou

	resp, _ := json.Marshal(res)
	io.WriteString(w, string(resp))
}

func HandleQueryIou(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive Query IOU request")
	SetCommRespHeader(w)

	res := &pbi.QueryIouRsp{}
	q := req.URL.Query()
	iouid := q.Get("iouid")
	iou, err := db.QueryIou(iouid)
	if err != nil {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(3)),
			Desc: pb.String(fmt.Sprintf("DB error.%s", err.Error())),
		}
	} else {
		iu := &pbi.IouWithUserInfo{}
		iu.Iou = iou
		user1, _ := db.QueryUser(iou.GetFrom())
		user2, _ := db.QueryUser(iou.GetTo())
		iu.From = user1
		iu.To = user2

		res.IouInfo = iu
	}

	resp, _ := json.Marshal(res)
	io.WriteString(w, string(resp))

	return
}

func HandleSign(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive sign request")
	SetCommRespHeader(w)
	res := &pbi.SignRsp{}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))

	signReq := &pbi.SignReq{}
	if err := json.Unmarshal(body, signReq); err != nil {
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(pbc.ErrCode_ECParamErr)),
			Desc: pb.String("Params invalid."),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return
	}

	signurl, err := GetSignUrl(signReq)
	if err != nil {
		log.Error("Get sign url of iou %s failed.", signReq.GetIouId())
		res.Error = &pbc.Error{
			Code: pb.Int32(int32(4)),
			Desc: pb.String(err.Error()),
		}
		resp, _ := json.Marshal(res)
		io.WriteString(w, string(resp))
		return 
	}

	res.SignUrl = pb.String(signurl)

	resp, _ := json.Marshal(res)
	io.WriteString(w, string(resp))
}

func HandleORC(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	//直接转发到esign
	SetCommRespHeader(w)
	body, _ := ioutil.ReadAll(req.Body) 
	req, err := http.NewRequest("POST", URL_OCR, bytes.NewBuffer(body))
	SetAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := HttpCli.Do(req)
	if err != nil {
		log.Error("Request OCR api err:%v", err)
		io.WriteString(w, string("forward error."))
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		io.WriteString(w, string("error"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	rBody, _ := ioutil.ReadAll(resp.Body)
	io.WriteString(w, string(rBody))
}

func HandleVerify3Factor(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	//直接转发到esign
	body, _ := ioutil.ReadAll(req.Body) 
	log.Debug("Receive verify 3 factor request. %s", string(body))
	req, err := http.NewRequest("POST", URL_VERIFY_FACTOR, bytes.NewBuffer(body))
	SetAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := HttpCli.Do(req)
	if err != nil {
		log.Error("Request Verify 3 fator api err:%v", err)
		io.WriteString(w, string("forward error."))
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		log.Error("Verify 3 factor status code: %d", resp.StatusCode)
		io.WriteString(w, string("error"))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	rBody, _ := ioutil.ReadAll(resp.Body)
	io.WriteString(w, string(rBody))
}

func HandleCreateEsignAccount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	SetCommRespHeader(w)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("Receive create esign account request: %s", string(body))

	rq, err := simplejson.NewJson(body)
	if err != nil {
		return
	}

	idcardno, _ := rq.Get("idNumber").String()
	mobile, _ := rq.Get("mobile").String()
	name, _ := rq.Get("name").String()
	userid, _ := rq.Get("openid").String()

	if userid == "" || name == "" || mobile == "" || idcardno == "" {
		stat.SetAttr("ioubackend.create.esign.err.-3", 1)
		io.WriteString(w, "{\"code\": -3, \"message\":\"Param invalid.\"}")
		return
	}
	user := &pbi.User{}
	user.OpenId = pb.String(userid)
	user.IDCardId = pb.String(idcardno)
	user.MobilePhone = pb.String(mobile)
	user.Name = pb.String(name)

	//更新用户信息
	db.AddUserToDB(user)
	accountId, err := CreateEsignAccount(userid)
	if err != nil {
		stat.SetAttr("ioubackend.create.esign.err.-1", 1)
		io.WriteString(w, fmt.Sprintf("{\"code\": -1, \"message\":%s}", err.Error()))
		return
	}

	if err := db.SetEsignAccount(userid, accountId); err != nil {
		stat.SetAttr("ioubackend.create.esign.err.-2", 1)
		io.WriteString(w, fmt.Sprintf("{\"code\": -2, \"message\":%s}", err.Error()))
		return
	}

	io.WriteString(w, "{\"code\": 0, \"message\":\"success\"}")

	return
}

func HandleAuthFace(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive face authority request")
	SetCommRespHeader(w)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))

	rq, err := simplejson.NewJson(body)
	if err != nil {
		io.WriteString(w, "{\"code\":1, \"message\": \"Param invalid.\"}")
		return
	}

	res, err := EsignFaceAuth(rq)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("{\"code\": -3, \"message\":%s}", err.Error()))
		return
	}

	io.WriteString(w, string(res))
}
func HandleAuthTelecom(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive telecom authority request")
	SetCommRespHeader(w)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))

	rq, err := simplejson.NewJson(body)
	if err != nil {
		stat.SetAttr("ioubackend.auth.telecom.err.1", 1)
		io.WriteString(w, "{\"code\":1, \"message\": \"Param invalid.\"}")
		return
	}

	res, err := EsignTelecomAuth(rq)
	if err != nil {
		stat.SetAttr("ioubackend.auth.telecom.err.-3", 1)
		io.WriteString(w, fmt.Sprintf("{\"code\": -3, \"message\":%s}", err.Error()))
		return
	}

	io.WriteString(w, string(res))
}

func HandleVeriCode(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	log.Debug("Receive verify code authority request")
	SetCommRespHeader(w)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))

	rq, err := simplejson.NewJson(body)
	if err != nil {
		stat.SetAttr("ioubackend.vericode.err.1", 1)
		io.WriteString(w, "{\"code\":1, \"message\": \"Param invalid.\"}")
		return
	}

	res, err := EsignVeriCode(rq)
	if err != nil {
		stat.SetAttr("ioubackend.vericode.err.-3", 1)
		io.WriteString(w, fmt.Sprintf("{\"code\": -3, \"message\":%s}", err.Error()))
		return
	}

	io.WriteString(w, string(res))
}

func HandleReqEsignDirect(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	fullUrl := req.URL.Path
	log.Debug("Receive esign request: %s", fullUrl)	
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.esgin.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)

	s := strings.Split(fullUrl,"esign/default")
	if len(s) < 2 {
		return 
	}

	url := s[1]

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	log.Debug("%s", string(body))

	rHdr, resp, err := EsignRequestWithHeader(req.Method, url, body, req.Header)
	if err != nil {
		io.WriteString(w, fmt.Sprintf("{\"code\": -3, \"message\":%s}", err.Error()))
		return
	}

	for key, values := range rHdr {
		for _, v := range values {
			w.Header().Set(key, v)
		}
	}	
	io.WriteString(w, string(resp))
}

func HandleReqContractDownload(w http.ResponseWriter, req *http.Request, body []byte) ([]byte, error) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	q := req.URL.Query()
	iouid := q.Get("iouid")

	if iouid == "" {
		stat.SetAttr("ioubackend.contract.err.1", 1)
		return []byte("{\"code\":1, \"message\": \"Param invalid.\"}"), nil
	}

	return EsignGetDownloadUrl(iouid)
}

func HandleReqContractPreview(w http.ResponseWriter, req *http.Request, body []byte) ([]byte, error) {
	stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
		strings.Replace(req.URL.Path, ".", "_", -1)), 1)
	q := req.URL.Query()
	iouid := q.Get("iouid")
	openid := q.Get("openid")

	if iouid == "" || openid == "" {
		return []byte("{\"code\":1, \"message\": \"Param invalid.\"}"), nil
	}

	return EsignGetPreviewUrl(iouid, openid)
}

func HandleTestAuth(w http.ResponseWriter, req *http.Request, body []byte) ([]byte, error) {
	return []byte("{\"code\": 0}"), nil
}

func HandleHttpRequest(f HandleFunc) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		log.Debug("Receive http request: %s",req.URL.Path)	
		stat.SetAttr(fmt.Sprintf("ioubackend.http.req.%s", 
			strings.Replace(req.URL.Path, ".", "_", -1)), 1)
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Error("Read request body error, %v", err)
			return
		}
		log.Debug("%s", string(body))
		resp, err := f(w, req, body)	
		if err != nil {
			io.WriteString(w, fmt.Sprintf("{\"code\": -3, \"message\":%s}", err.Error()))
			return 
		}
		SetCommRespHeader(w)
		io.WriteString(w, string(resp))
	}
}

func HandleHttpRequestWithAuth(f HandleFunc, oauth *server.Server) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		log.Debug("Receive http request: %s",req.URL.Path)	
		stat.SetAttr(fmt.Sprintf("ioubackend.http.req.auth.%s", 
			strings.Replace(req.URL.Path, ".", "_", -1)), 1)
		if _, err := oauth.ValidationBearerToken(req); err != nil {
	       http.Error(w, err.Error(), http.StatusBadRequest)
	       return
	    }
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Error("Read request body error, %v", err)
			return
		}
		log.Debug("%s", string(body))
		resp, err := f(w, req, body)	
		if err != nil {
			io.WriteString(w, fmt.Sprintf("{\"code\": -3, \"message\":%s}", err.Error()))
			return 
		}
		SetCommRespHeader(w)
		io.WriteString(w, string(resp))
	}
}