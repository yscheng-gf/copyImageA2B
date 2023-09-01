package config

type Conf struct {
	AHostname string
	BMongo    MongoConf `mapstructure:"BMongo"`
}

type MongoConf struct {
	Host string
}
