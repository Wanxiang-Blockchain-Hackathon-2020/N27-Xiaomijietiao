package main

import (
	"fmt"
	"net/http"

	api "github.com/despard/ioubackend/api"
	. "github.com/despard/ioubackend/config"
	db "github.com/despard/ioubackend/db"
	serv "github.com/despard/ioubackend/server"
	"github.com/despard/log"
	"github.com/julienschmidt/httprouter"
	uuid "github.com/satori/go.uuid"

	"encoding/json"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	stat "github.com/despard/stat"	
)

func init() {
	log.SetLogFileName(Settings.LogFileName)
	log.SetLogLevel(Settings.LogLevel)
	fmt.Println("Init succ.")
}

func main() {
	fmt.Println("Starting ioubackend proccess...")

	log.Info("stat serverTag: aliyun")
	if err := stat.AttrInit(Settings.StatHost, "aliyun", Settings.StatPrefix); err != nil {
		log.Error("stat.AttrInit err:%v", err)
		return
	}

	stat.SetCpuStat("ioubackend.cpu")
	stat.SetMemoryStat("ioubackend.mem")
	stat.SetHeapStat("ioubackend.heap")

	stat.SetAttrStatus("ioubackend.start", 1)

	fmt.Println("Connecting Mongodb.")
	if err := db.InitDB(Settings.MongoIp, Settings.MongoPort, Settings.MongoDbName); err != nil {
		fmt.Printf("Connect to Mongodb failed:%v\n", err)
		return
	}
	fmt.Println("Connect Mongodb succ.")
	fmt.Println("Starting weyom backend HTTP Server")
	router := httprouter.New()

	//设置OAuth2.0
	authMgr := manage.NewDefaultManager()
	authMgr.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	authMgr.MustTokenStorage(store.NewFileTokenStore(Settings.TokenStore))

	clientStore := store.NewClientStore()
	authMgr.MapClientStorage(clientStore)

	srv := server.NewDefaultServer(authMgr)
	srv.SetAllowGetAccessRequest(true)
	srv.SetClientInfoHandler(server.ClientFormHandler)
	authMgr.SetRefreshTokenCfg(manage.DefaultRefreshTokenCfg)

	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		fmt.Println("Internal Error:", err.Error())
		return
	})

	srv.SetResponseErrorHandler(func(re *errors.Response) {
		fmt.Println("Response Error:", re.Error.Error())
	})

	router.GET("/api/auth/token", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		srv.HandleTokenRequest(w, r)
	})

	router.GET("/auth/credentials", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		clientId := uuid.NewV4().String()[:8]
		clientSecret := uuid.NewV4().String()
		err := clientStore.Set(clientId, &models.Client{
			ID:     clientId,
			Secret: clientSecret,
			Domain: fmt.Sprintf("http://%s", Settings.HttpAddr),
		})
		if err != nil {
			fmt.Println(err.Error())
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"CLIENT_ID": clientId, "CLIENT_SECRET": clientSecret})
	})

	router.POST("/api/v1/users/add", api.HandleAddUser)
	router.POST("/api/v1/users/query_iou_of_user", api.HandleQueryIouOfUser)
	router.POST("/api/v1/users/create_esign_account", api.HandleCreateEsignAccount)
	router.POST("/api/v1/ious/add", api.HandleAddIou)
	router.POST("/api/v1/ious/update", api.HandleUpdateIou)
	router.POST("/api/v1/ious/agree", api.HandleAgreeIou)

	router.POST("/api/server/esign/notice/sign", serv.HandleEsignNotice)
	router.POST("/api/server/esign/notice/auth", serv.HandleAuthNotice)
	router.POST("/api/v1/tools/orc", api.HandleORC)
	router.POST("/api/v1/tools/verify3factor", api.HandleVerify3Factor)
	router.POST("/api/v1/auth/face", api.HandleAuthFace)
	router.POST("/api/v1/auth/telecom3factor", api.HandleAuthTelecom)
	router.POST("/api/v1/auth/telecom3factor/", api.HandleAuthTelecom)
	router.POST("/api/v1/auth/verify_code", api.HandleVeriCode)
	router.POST("/api/v1/auth/verify_code/", api.HandleVeriCode)
	router.POST("/api/v1/ious/sign/getSignUrl", api.HandleSign)
	router.POST("/api/v1/issue", api.HandleIssue)

	router.GET("/api/v1/users/query", api.HandleQueryUser)
	router.GET("/api/v1/access", api.HandleAccess)
	router.GET("/api/v1/ious/query", api.HandleQueryIou)
	router.GET("/api/v1/web/images/banners", api.HandleBanners)
	router.GET("/api/v1/ious/getDownloadUrl", api.HandleHttpRequest(api.HandleReqContractDownload))
	router.GET("/api/v1/ious/getPreviewUrl", api.HandleHttpRequest(api.HandleReqContractPreview))
	router.GET("/oauth", api.HandleHttpRequestWithAuth(api.HandleTestAuth, srv))

	//直接转发e签宝平台
	router.GET("/api/v1/esign/default/*action", api.HandleReqEsignDirect)
	router.POST("/api/v1/esign/default/*action", api.HandleReqEsignDirect)
	router.PUT("/api/v1/esign/default/*action", api.HandleReqEsignDirect)

	http.ListenAndServe(Settings.HttpAddr, router)
}
