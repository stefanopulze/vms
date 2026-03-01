package influx

type Options struct {
	Enabled      bool   `yaml:"enabled" env:"INFLUX_ENABLED"`
	Host         string `yaml:"host" env:"INFLUX_HOST"`
	Token        string `yaml:"token" env:"INFLUX_TOKEN"`
	Database     string `yaml:"database" env:"INFLUX_DATABASE"`
	Organization string `yaml:"organization" env:"INFLUX_ORGANIZATION"`
}
