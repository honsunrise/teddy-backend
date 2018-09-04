package config

type Server struct {
	Address string `json:"address"`
	Port    int    `json:"port"`
}

type Database struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
	AuthDB   string `json:"auth_db"`
}

type Config struct {
	Server    Server                `json:"server"`
	Databases map[string][]Database `json:"databases"`
	Casbin    string                `json:"casbin"`
	JWTPkcs8  string                `json:"jwt_pkcs8"`
}
