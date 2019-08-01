package master

import (
	"encoding/json"
	"io/ioutil"
)

// 创建单例,供外部使用
var (
	G_config Config
)

// 程序配置
type Config struct {
	ApiPort         int      `json:"api_port"`
	ApiReadTimeout  int      `json:"api_read_timeout"`
	ApiWriteTimeout int      `json:"api_write_timeout"`
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
	StaticDir       string   `json:"staticDir"`
}

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}
	G_config = conf
	return
}
