package atom

import (
	"fmt"

	"github.com/spf13/viper"
)

func (atom *Atom) InitViper(confFilename string) error {
	fmt.Printf("confFilename: %s\n", confFilename)
	// config
	viper := viper.New()
	// viper.SetConfigName("development.yaml") // name of config file (without extension)
	viper.SetConfigName(confFilename) // name of config file (without extension)
	viper.SetConfigType("yaml")       // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./conf/")

	viper.SetDefault("environment", "development")
	viper.SetDefault("server", map[string]string{"address": ":8080"})
	viper.SetDefault("redis", map[string]string{"address": "localhost:6379", "password": ""})
	viper.SetDefault("mongodb", map[string]string{"address": "mongodb://localhost:27017"})
	viper.SetDefault("elasticsearch", map[string]interface{}{"addresses": []string{"https://es.xxx.com"}, "username": "user", "password": "xxx"})
	viper.SetDefault("gin", map[string]string{"address": ":8080"})
	viper.SetDefault("zap", map[string]string{"level": "debug", "dir": "./log"})
	viper.SetDefault("mysql", map[string]string{"address": "localhost:3306", "database": "test", "username": "root", "password": "123456"})

	// dump demo.yaml
	viper.WriteConfigAs("./conf/demo.yaml")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	atom.Viper = viper

	return nil
}
