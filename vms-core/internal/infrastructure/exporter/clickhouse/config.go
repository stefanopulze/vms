package clickhouse

type Options struct {
	Addr     string `yaml:"addr" env:"CLICKHOUSE_ADDR"`
	Database string `yaml:"database" env:"CLICKHOUSE_DATABASE"`
	Username string `yaml:"username" env:"CLICKHOUSE_USERNAME"`
	Password string `yaml:"password" env:"CLICKHOUSE_PASSWORD"`
}
