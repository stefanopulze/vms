package database

type Config struct {
	Hostname string `yaml:"hostname" env:"DATABASE_HOSTNAME"`
	Username string `yaml:"username" env:"DATABASE_USERNAME"`
	Password string `yaml:"password" env:"DATABASE_PASSWORD"`
	Name     string `yaml:"name" env:"DATABASE_NAME"`
}
