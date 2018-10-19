package common

const STATIC_DIR string = "../../webui"
const PROG_PATH string = "curl"

type LoggerConfig struct {
	Logfile  string
	LogLevel int
}

type DBServiceConfig struct {
	Server   string
	Port     string
	User     string
	Password string
	Driver   string
	DBName   string
}
