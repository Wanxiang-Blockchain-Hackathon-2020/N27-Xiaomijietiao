package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"

	//. "github.com/despard/ioubackend/config"
	"github.com/despard/log"
	"github.com/julienschmidt/httprouter"

	simplejson "github.com/bitly/go-simplejson"
	db "github.com/despard/ioubackend/db"
	pbi "github.com/despard/iouproto/iou"
	stat "github.com/despard/stat"
)

func HandleEsignNotice(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr("ioubackend.esign.notice", 1)
	log.Info("Receive Esign notice message.")
	appid := req.Header.Get("X-Tsign-Open-App-Id")
	signature := req.Header.Get("X-Tsign-Open-SIGNATURE")
	timestamp := req.Header.Get("X-Tsign-Open-TIMESTAMP ")
	alg := req.Header.Get("X-Tsign-Open-SIGNATURE-ALGORITHM")
	log.Debug("Notice params: appid %s, sinature %s, timestamp %s, alg %s",
		appid, signature, timestamp, alg)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}
	/*
		if appid == "" || signature == "" || timestamp == "" || alg == "" {
			return
		}
		if appid != Settings.Appid {
			return
		}
		//获取Params，并且按字典序排序
		var keys []string
		q := req.URL.Query()
		for k, _ := range q {
			log.Debug("Params key: %s", k)
			keys = append(keys, k)
		}
		sort.Strings(keys)

		vals := ""
		for _, k := range keys {
			v := q.Get(k)
			vals += v
		}

		bvals := []byte(vals)


		var buffer bytes.Buffer
		buffer.Write([]byte(timestamp))
		buffer.Write(bvals)
		buffer.Write(body)

		allByte := buffer.Bytes()

		sign := hmacSha256(allByte, Settings.AppSecret)
		if sign != signature {
			log.Error("Signature illegal.")
			return
		}
	*/
	log.Info("Esign notice message: %s", string(body))

	iouid := req.URL.Query().Get("iouid")
	if iouid == "" {
		log.Debug("IouId is nil")
		return
	}

	res, err := simplejson.NewJson(body)
	if err != nil {
		return
	}

	act, err := res.Get("action").String()
	if err != nil {
		return
	}

	switch act {
	case "SIGN_FLOW_UPDATE":
		HandleSignUpdate(iouid, res)
	case "SIGN_FLOW_FINISH":
		HandleSignFinish(res)
	}

	io.WriteString(w, "{\"code\":0, \"message\": \"success.\"}")
	return
}

func HandleSignUpdate(iouid string, su *simplejson.Json) {
	stat.SetAttr("ioubackend.esign.update", 1)
	flowId, _ := su.Get("flowId").String()
	iou, err := db.QueryIou(iouid)
	if err != nil || iou == nil {
		stat.SetAttr("ioubackend.esign.err.update", 1)
		log.Error("Receive sign flow %s update, but iou %s is not found ", flowId, iouid)
		return
	}

	if iou.GetFlowId() != flowId {
		stat.SetAttr("ioubackend.esign.err.update", 1)
		log.Error("Receive sign flow %s update, but flowId is not right for iou %s?", flowId, iouid)
		return
	}

	result, _ := su.Get("signResult").Int()
	if result != 2 {
		log.Error("Flow id %s , iouid %s sign result not succ.", flowId, iouid)
		return
	}

	if iou.GetStatus() == pbi.IouStatus_SignPartA {
		iou.Status = pbi.IouStatus_SignPartB.Enum()
	} else if iou.GetStatus() == pbi.IouStatus_SignPartB {
		iou.Status = pbi.IouStatus_Pay.Enum()
	}
	//更新状态
	if err := db.UpdateIou(iou); err != nil {
		stat.SetAttr("ioubackend.esign.err.update", 1)
		log.Error("Update iou %s status failed.", iouid)
	}

	log.Info("Update iou %s sign status succ.", iouid)

	return
}

func HandleSignFinish(su *simplejson.Json) {
	stat.SetAttr("ioubackend.esign.finished", 1)
	flowId, _ := su.Get("flowId").String()

	log.Info("Esign flow %s finished.", flowId)
	return
}

func hmacSha256(data []byte, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func HandleAuthNotice(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	stat.SetAttr("ioubackend.esign.notice.auth", 1)
	log.Info("Receive Esign face auth notice message.")
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("Read request body error, %v", err)
		return
	}

	log.Info("Face auth notice message:%v", string(body))
	res, err := simplejson.NewJson(body)
	if err != nil {
		return
	}

	flowId, _ := res.Get("flowId").String()
	accountId, _ := res.Get("accountId").String()
	succ, _ := res.Get("success").Bool()
	cId, _ := res.Get("contextId").String()

	user, err := db.QueryUser(cId)
	if err != nil {
		log.Error("Receive face auth notice msg ,but user %s is not found.flow %s, account %d,%v",
			cId, flowId, accountId, err)
		io.WriteString(w, "")
		return
	}

	if succ == true {
		user.Auth = pbi.UserAuthStatus_AuthSucc.Enum()
	} else {
		user.Auth = pbi.UserAuthStatus_AuthFail.Enum()
	}

	db.AddUserToDB(user)

	io.WriteString(w, "ok.")
}
