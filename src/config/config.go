package config

type dbConfig struct {
	Host            string
	Port            int
	Database        string
	User            string `default:"root"`
	Password        string
	ConnMaxLifetime int `default:"600"`
	MaxIdleConn     int `default:"32"`
}

type httpConfig struct {
	Port int `default:"80"`
}

type switchConfig struct {
	OpenMonitor bool `default:"false"`
	OpenHttp bool `default:"false"`
}

type logConfig struct {
	Name string
	Path string
}

type qcloudTool struct {
	Db  dbConfig
	Log logConfig
	Switch switchConfig
	Http httpConfig
}

var QcloudTool qcloudTool
