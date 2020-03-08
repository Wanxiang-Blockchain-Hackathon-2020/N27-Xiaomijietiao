package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	Settings = new(Config)
)

func init() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./ioubackend configPath")
		os.Exit(-1)
	}

	// 加载配置
	filePath := os.Args[1]
	err := LoadConfig(filePath)
	if err != nil {
		fmt.Println("load config file failed!", err)
		os.Exit(-2)
	}
}

type Config struct {
	Cpunum      int    `json:"cpunum"`
	LogLevel    string `json:"loglevel"`    //日志级别
	LogFileName string `json:"logfilename"` //日志文件名
	HttpAddr    string `json:"httpaddr"`    //监听地址
	MongoIp     string `json:"mongoip"`     //MongoDB地址
	MongoPort   uint16 `json:"mongoport"`   //MongoDB端口
	MongoDbName string `json:"mongodbname"` //MongoDB名字
	Appid       string `josn:"appid"`       //e签宝的appid
	AppSecret   string `json:"appsecret"`   //e签宝的appsecret
	TemplateId  string `json:"templateid"`  //e签宝的合同模板ID
	NoticeUrl   string `json:"noticeurl"`   //e签宝统一回调地址
	TokenStore  string `json:"tokenstore"`  //token存储路径
	StatPrefix  string `json:"statprefix"` //属性统计前缀
	Version     string `json:"version"`     //版本号
	StatHost    string `json:"stathost"`   //统计监听server
}

func LoadConfig(filePath string) error {
	configuration, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	} else {
		err := json.Unmarshal(configuration, &Settings)
		return err
	}
}
