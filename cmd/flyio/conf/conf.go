package conf

import (
	"io/ioutil"

	"github.com/wenchy/grpcio/internal/envconf"
	"gopkg.in/yaml.v3"
)

type serverConf struct {
	Server envconf.Node `yaml:"server"`
	Log    envconf.Log  `yaml:"log"`
}

var Conf serverConf

func InitConf(path string) error {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(d, &Conf)
	if err != nil {
		panic(err)
	}

	//envconf.Init("../../conf/common/env_conf.yaml")

	return nil
}
