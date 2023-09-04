package config

type Conf struct {
	AHostname string
	AMongo    MongoConf `mapstructure:"AMongo"`
	BMongo    MongoConf `mapstructure:"BMongo"`
}

type MongoConf struct {
	Host string
}
