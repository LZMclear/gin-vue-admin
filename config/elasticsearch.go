package config

type Elasticsearch struct {
	Address    string `mapstructure:"address" json:"address" yaml:"address"`
	Timeout    int    `mapstructure:"timeout" json:"timeout" yaml:"timeout"`
	Index      string `mapstructure:"index" json:"index" yaml:"index"`
	BatchSize  int    `mapstructure:"batch-size" json:"batch-size" yaml:"batch-size"`
	MaxRetries int    `mapstructure:"max-retries" json:"max-retries" yaml:"max-retries"`
}
