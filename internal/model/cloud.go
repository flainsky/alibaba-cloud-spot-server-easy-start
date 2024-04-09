package model

type CloudConf struct {
	Version      float32    `yaml:"version"`
	AliYunConfig AliYunConf `yaml:"aliyun" mapstructure:"aliyun"`
}

type AliYunConf struct {
	AccessKey    string `yaml:"accessKey"`
	AccessSecret string `yaml:"accessSecret"`
	Ver          int32  `yaml:"ver"`
}
