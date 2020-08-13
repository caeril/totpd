package config

import (
	"os"
	"strconv"
)

type Config struct {
	DataPath string
	HttpPort int
}

var _singleConfig Config

func Get() Config {
	return _singleConfig
}

func InitConfig() {

	_singleConfig = Config{DataPath: "/tmp", HttpPort: 8087}

	tmp := os.Getenv("TOTPD_DATA_PATH")
	if len(tmp) > 0 {
		_singleConfig.DataPath = tmp
	}

	tmp = os.Getenv("TOTPD_HTTP_PORT")
	if len(tmp) > 0 {
		var err error
		_singleConfig.HttpPort, err = strconv.Atoi(tmp)
		if err != nil {
			panic(err)
		}
	}
}
