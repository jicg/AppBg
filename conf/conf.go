package conf

import (
	"io/ioutil"
	"encoding/json"
	"fmt"
)

var (
	conf *Conf
)

type Conf struct {
	Db      string `json:"db"`
	Logfile string `json:"logfile"`
	Logpath string `json:"logpath"`
	Port    string `json:"port"`
	Mode    string `json:"mode"`
}

func loadConf() (*Conf, error) {
	var (
		bs  []byte
		err error
	)
	bs, err = ioutil.ReadFile("conf.json")

	if err != nil {
		return nil, err
	}
	conf = new(Conf)
	if err = json.Unmarshal(bs, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

func init() {
	var err error
	conf, err = loadConf();
	if err != nil {
		fmt.Println("请配置好文件conf.json");
		return;
	}
	if (len(conf.Mode)) == 0 {
		conf.Mode = "debug"
	}
}

func GetConf() *Conf {
	return conf
}
