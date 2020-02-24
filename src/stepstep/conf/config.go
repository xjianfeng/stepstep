package conf

import (
	"github.com/go-ini/ini"
	"net/url"
	"time"
)

type server struct {
	AppName      string
	HTTPAddr     string
	RunMode      string
	LogPath      string
	DataPath     string
	LruCap       int32
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	ImageDomain  string
	PayServed    string
	PayCallBack  string
}

type dataBase struct {
	Type     string
	User     string
	Password string
	Host     string
	DBName   string
	TbPrefix string
}

type mongoDb struct {
	MongoUri string
	MongoDb  string
}

type redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	DefaultDb   int
}

type wechat struct {
	Appid        string
	Appsecret    string
	Mchid        string
	MchAppsecret string
	TemplateId   string
	Page         string
}

type cos struct {
	AppId     string
	SecretId  string
	SecretKey string
	Region    string
	ImagePath string
	SoundPath string
	Domain    string
	Bucket    string
}

var CfgServer = &server{}
var CfgRedis = &redis{}
var CfgDb = &dataBase{}
var CfgMongo = &mongoDb{}
var CfgWechat = &wechat{}
var CfgCos = &cos{}

func SetUp(configPath string) {
	var err error
	cfg, err := ini.Load(configPath)
	if err != nil {
		panic(err)
	}
	cfg.Section("server").MapTo(CfgServer)
	cfg.Section("redis").MapTo(CfgRedis)
	cfg.Section("database").MapTo(CfgDb)
	cfg.Section("mongodb").MapTo(CfgMongo)
	cfg.Section("wechat").MapTo(CfgWechat)
	cfg.Section("cos").MapTo(CfgCos)

	CfgRedis.Password, _ = url.QueryUnescape(CfgRedis.Password)
}

func IsRelease() bool {
	return CfgServer.RunMode == "release"
}
