package api

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/despard/ioubackend/db"
	pbi "github.com/despard/iouproto/iou"
	log "github.com/despard/log"
	pb "github.com/golang/protobuf/proto"
)

func init() {
	if err := db.InitDB("127.0.0.1", 27017); err != nil {
		fmt.Printf("Connect to Mongodb failed:%v\n", err)
		os.Exit(1)
	}

	log.SetLogLevel(log.LOGLEVEL_INFO)

}

/*
func TestCreateEsignAccount(t *testing.T) {
	if err := db.InitDB("127.0.0.1", 27017); err != nil {
		fmt.Printf("Connect to Mongodb failed:%v\n", err)
		return
	}

	if eid, err := CreateEsignAccount("or8kF5iYRtkTCRZ0ANCpi3KjYQQ0"); err != nil {
		t.Fatalf("err: %v", err)
		return
	} else {
		t.Logf("Esign id :%s\n", eid)
	}
	return
}

/*

func TestCreateEsignFlow(t *testing.T) {
	if err := db.InitDB("127.0.0.1", 27017); err != nil {
		fmt.Printf("Connect to Mongodb failed:%v\n", err)
		return
	}

	col, _ := db.Collection("ious")
	iou := &pbi.Iou{}
	col.Find(bson.M{"from": "or8kF5iYRtkTCRZ0ANCpi3KjYQQ0"}).One(iou)

	fmt.Println(*iou)
	if fid, err := CreateEsignFlow(iou); err != nil {
		t.Fatalf("err: %v", err)
		return
	} else {
		t.Logf("Esign flow id :%s\n", fid)
		if res, err := QueryEsignFlow(fid); err != nil {
			t.Fatalf("err: %v", err)
		} else {
			t.Logf("Esign flow :%s\n", res)
		}
	}

	return
}

func TestQueryEsignFlow(t *testing.T) {
	if err := db.InitDB("127.0.0.1", 27017); err != nil {
		fmt.Printf("Connect to Mongodb failed:%v\n", err)
		return
	}

	col, _ := db.Collection("ious")
	iou := &pbi.Iou{}
	col.Find(bson.M{"from": "or8kF5iYRtkTCRZ0ANCpi3KjYQQ0"}).One(iou)

	if fid, err := QueryEsignFlow(iou.GetToFlowId()); err != nil {
		t.Fatalf("err: %v", err)
		return
	} else {
		t.Logf("Esign flow id :%s\n", fid)
	}
	return
}

func TestSignReq(t *testing.T) {
	if err := db.InitDB("127.0.0.1", 27017); err != nil {
		fmt.Printf("Connect to Mongodb failed:%v\n", err)
		return
	}

	col, _ := db.Collection("ious")
	iou := &pbi.Iou{}
	col.Find(bson.M{"from": "or8kF5iYRtkTCRZ0ANCpi3KjYQQ0"}).One(iou)

	sign := &pbi.SignReq{
		IouId:       iou.GetId(),
		UserId:      iou.GetFrom(),
		RedirectUrl: "",
	}

	url, err := GetSignUrl(sign)
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	t.Log(url)
}
*/
func TestSetAuthHeader(t *testing.T) {
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetAuthHeader(tt.args.req)
		})
	}
}

func TestEsignRequest(t *testing.T) {
	type args struct {
		method string
		url    string
		body   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *simplejson.Json
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EsignRequest(tt.args.method, tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("EsignRequest(%v, %v, %v) error = %v, wantErr %v", tt.args.method, tt.args.url, tt.args.body, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EsignRequest(%v, %v, %v) = %v, want %v", tt.args.method, tt.args.url, tt.args.body, got, tt.want)
			}
		})
	}
}

func TestCreateEsignSeal(t *testing.T) {
	type args struct {
		userId string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateEsignSeal(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateEsignSeal(%v) error = %v, wantErr %v", tt.args.userId, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateEsignSeal(%v) = %v, want %v", tt.args.userId, got, tt.want)
			}
		})
	}
}

func TestUploadContract(t *testing.T) {
	type args struct {
		iou  *pbi.Iou
		file []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UploadContract(tt.args.iou, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("UploadContract(%v, %v) error = %v, wantErr %v", tt.args.iou, tt.args.file, err, tt.wantErr)
			}
		})
	}
}

func TestAttachDocument(t *testing.T) {
	type args struct {
		iouid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AttachDocument(tt.args.iouid); (err != nil) != tt.wantErr {
				t.Errorf("AttachDocument(%v) error = %v, wantErr %v", tt.args.iouid, err, tt.wantErr)
			}
		})
	}
}

func TestQueryTemplateInfo(t *testing.T) {
	type args struct {
		tid string
	}
	tests := []struct {
		name    string
		args    args
		want    *simplejson.Json
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryTemplateInfo(tt.args.tid)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryTemplateInfo(%v) error = %v, wantErr %v", tt.args.tid, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryTemplateInfo(%v) = %v, want %v", tt.args.tid, got, tt.want)
			}
		})
	}
}

func TestCreateTemplateContract(t *testing.T) {
	type args struct {
		iouid string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTemplateContract(tt.args.iouid)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTemplateContract(%v) error = %v, wantErr %v", tt.args.iouid, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateTemplateContract(%v) = %v, want %v", tt.args.iouid, got, tt.want)
			}
		})
	}
}

func TestStartEsignFlow(t *testing.T) {
	type args struct {
		sign *pbi.SignReq
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StartEsignFlow(tt.args.sign); (err != nil) != tt.wantErr {
				t.Errorf("StartEsignFlow(%v) error = %v, wantErr %v", tt.args.sign, err, tt.wantErr)
			}
		})
	}
}

func TestAddHandSign(t *testing.T) {
	type args struct {
		iou *pbi.Iou
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddHandSign(tt.args.iou); (err != nil) != tt.wantErr {
				t.Errorf("AddHandSign(%v) error = %v, wantErr %v", tt.args.iou, err, tt.wantErr)
			}
		})
	}
}

func TestGetSignUrl(t *testing.T) {
	type args struct {
		sign *pbi.SignReq
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "testNoIouUserId",
			args: args{sign: &pbi.SignReq{
				IouId:  pb.String("bebcfee2-80a8-439a-a5c7-e12c5e98d0a2"),
				UserId: pb.String("or8kF5iYRtkTCRZ0ANCpi3KjYQQ0"),
			}},
			wantErr: true,
			want:    "",
		},
		{
			name: "testNotAuthUserId",
			args: args{sign: &pbi.SignReq{
				IouId:  pb.String("bebcfee2-80a8-439a-a5c7-e12c5e98d0a2"),
				UserId: pb.String("or8kF5iYRtkTCRZ0ANdadadadaad"),
			}},
			wantErr: true,
			want:    "",
		},
		{
			name: "testAuthUserButNotFromIou",
			args: args{sign: &pbi.SignReq{
				IouId:  pb.String("bebcfee2-80a8-439a-a5c7-e12c5e98d0a3"),
				UserId: pb.String("or8kF5iYRtkTCRZ0ANCpi3KjYQQ0"),
			}},
			wantErr: true,
			want:    "",
		},
		{
			name: "testAuthUserAndRedirectUrl",
			args: args{sign: &pbi.SignReq{
				IouId:       pb.String("bebcfee2-80a8-439a-a5c7-e12c5e98d0a9"),
				UserId:      pb.String("or8kF5iYRtkTCRZ0ANCpi3KjYQQ0"),
				RedirectUrl: pb.String("www.baidu.com"),
			}},
			wantErr: false,
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSignUrl(tt.args.sign)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSignUrl(%v) error = %v, wantErr %v", tt.args.sign, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetSignUrl(%v) = %v, want %v", tt.args.sign, got, tt.want)
			}
		})
	}
}

func TestEsignFaceAuth(t *testing.T) {
	type args struct {
		rq *simplejson.Json
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EsignFaceAuth(tt.args.rq)
			if (err != nil) != tt.wantErr {
				t.Errorf("EsignFaceAuth(%v) error = %v, wantErr %v", tt.args.rq, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EsignFaceAuth(%v) = %v, want %v", tt.args.rq, got, tt.want)
			}
		})
	}
}
