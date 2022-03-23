package config

type Server struct {
	Host string
	Port int
	Name string
	Database *Database
}

type Database struct {
	Host string `yaml:"host"`
	Username string `yaml:"username"`
	Passwd string `yaml:"passwd"`
	Port int`yaml:"port"`
}