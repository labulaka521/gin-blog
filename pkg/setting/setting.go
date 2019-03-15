package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

// 定义与配置文件中app相同的配置的结构体
type App struct {
	PageSize  int
	JwtSecret string

	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath  string
	ImageMaxSize   int
	ImageAllowExts []string

	LogSavePath string
	LogSavename string
	LogFileExt  string
	TimeFormat  string
}

var AppSetting = &App{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type        string
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	Host        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeOut time.Duration
}

var RedisSetting = &Redis{}

func Setup() {
	Cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	// 映射结构体
	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.Mapto AppSetting err: %v", err)
	}

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}

	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second

	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo Database err: %v", err)
	}

	err = Cfg.Section("redis").MapTo(RedisSetting)

	if err != nil {
		log.Fatalf("Cfg.MapTo Redis err: %v", err)
	}
	RedisSetting.IdleTimeOut = RedisSetting.IdleTimeOut * time.Second
}

//编写与配置项保持一致的结构体(APP SERVER Database)
// 使用MapTo将配置项映射到结构体上
// 对一些需特殊设置的配置项进行再赋值
