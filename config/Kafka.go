package config

type Kafka struct {
	Topic        string   `mapstructure:"topic" json:"topic" yaml:"topic"`
	Brokers      []string `mapstructure:"brokers" json:"brokers" yaml:"brokers"`
	MysqlGroupId string   `mapstructure:"mysql-group-id" json:"mysql-group-id" yaml:"mysql-group-id"`
	EsGroupId    string   `mapstructure:"es-group-id" json:"es-group-id" yaml:"es-group-id"`
}
