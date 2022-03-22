package config

type Server struct {
	Host string
	Port int
	Name string
}

type Database struct {
	Host string
	Username string
	Passwd string
	Port int32
}