package api

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/julienschmidt/httprouter"

	. "github.com/despard/ioubackend/config"
	db "github.com/despard/ioubackend/db"
	"github.com/despard/log"
)

var testurl string = "http://47.106.91.220:5000"

func init() {
	if err := db.InitDB(Settings.MongoIp, Settings.MongoPort); err != nil {
		fmt.Printf("Connect to Mongodb failed:%v\n", err)
		os.Exit(1)
	}

	log.SetLogLevel(log.LOGLEVEL_INFO)

}

func TestSetCommRespHeader(t *testing.T) {
	type args struct {
		w http.ResponseWriter
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCommRespHeader(tt.args.w)
		})
	}
}

func TestHandleAddUser(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleAddUser(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleQueryUser(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleQueryUser(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleQueryIouOfUser(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleQueryIouOfUser(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleAddIou(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleAddIou(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleUpdateIou(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleUpdateIou(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleQueryIou(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleQueryIou(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleSign(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleSign(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleORC(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleORC(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleVerify3Factor(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleVerify3Factor(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleCreateEsignAccount(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleCreateEsignAccount(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleAuthFace(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleAuthFace(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleReqEsignDirect(t *testing.T) {
	type args struct {
		w   http.ResponseWriter
		req *http.Request
		in2 httprouter.Params
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			HandleReqEsignDirect(tt.args.w, tt.args.req, tt.args.in2)
		})
	}
}

func TestHandleReqContractDownload(t *testing.T) {
	e := httpexpect.New(t, testurl)     //创建一个httpexpect实例
	postdata := map[string]interface{}{ //创建一个json变量
		"flag": 1,
		"msg":  "terrychow",
	}
	contentType := "application/json;charset=utf-8"

	e.POST("/api/v1/ious/getDownloadUrl"). //post 请求
						WithHeader("ContentType", contentType). //定义头信息
						WithJSON(postdata).                     //传入json body
						Expect().
						Status(http.StatusOK). //判断请求是否200
						JSON().
						Object().                      //json body实例化
						ContainsKey("msg").            //检验是否包括key
						ValueEqual("msg", "terryzhou") //对比key的value，value不匹配
}
