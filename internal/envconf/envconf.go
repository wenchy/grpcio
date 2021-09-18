package envconf

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Mysql struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	Username string `yaml:"username"`
}
type Redis struct {
	Addrs    []string `yaml:"addrs"`
	Password string   `yaml:"password"`
}

type Node struct {
	ID             string `yaml:"id"`
	GRPCAddress    string `yaml:"grpc_address"`
	GatewayAddress string `yaml:"gateway_address"`
	HTTPAddress    string `yaml:"http_address"`
}

type Log struct {
	Level string `yaml:"level"`
	Dir   string `yaml:"dir"`
}

type envConf struct {
	Mysql Mysql             `yaml:"mysql"`
	Redis Redis             `yaml:"redis"`
	Nodes map[string][]Node `yaml:"nodes"`
}

var Conf *envConf

// Init initialize all env variables.
func Init(envPath string) {
	confBytes, err := ioutil.ReadFile(envPath)
	if err != nil {
		panic(err)
	}
	Conf = new(envConf)
	err = yaml.Unmarshal(confBytes, Conf)
	if err != nil {
		panic(err)
	}
}
