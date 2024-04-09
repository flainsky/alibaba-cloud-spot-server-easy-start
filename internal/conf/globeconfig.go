package conf

import (
	"CloudBuild/internal/model"
	"github.com/spf13/viper"
)

var cloudConf model.CloudConf

func GetCloudConfig() (model.CloudConf, error) {

	vip := viper.New()
	vip.SetConfigName("cloud")
	vip.SetConfigType("yaml")
	vip.AddConfigPath("./conf")

	err := vip.ReadInConfig()
	if err != nil {
		return cloudConf, err
	}

	err = vip.Unmarshal(&cloudConf)
	if err != nil {
		return cloudConf, err
	}
	return cloudConf, nil

}
