package system

var Properties = SystemProperties{}
var Runtime = &Properties.Runtime
var Database = &Properties.Database
var OAuth2 = &Properties.Oauth2

type SystemProperties struct {
	Version string
	Runtime RuntimeProperties
	Database DatabaseProperties
	Oauth2 OAuth2Properties
}

type LogConfig struct {
	LogLevel string
	LogPath string
	LogFile string
	MaxSize int
	Rotate string
	Stdout bool
}

type RuntimeProperties struct {
	LogConfig
	WebPath string
	Debug bool
	Port string
}

type DatabaseProperties struct {
	Host string
	Port string
	Schema string
	User string
	Password string
}

type OAuth2Properties struct {
	Client []AuthClient
}

type AuthClient struct {
	Id string
	Secret string
	Domain string
}