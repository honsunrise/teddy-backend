package types

type Paging struct {
	Page uint32 `json:"page"`
	Size uint32 `json:"size"`
}

type Sort struct {
	Name string
	ASC  bool
}

// For config

type Server struct {
	Address string `json:"address" mapstructure:"address"`
	Port    int    `json:"port" mapstructure:"port"`
}

type Database struct {
	Address  string `json:"address" mapstructure:"address"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
	AuthDB   string `json:"auth_db" mapstructure:"auth_db"`
}

type Mail struct {
	Host     string `json:"host" mapstructure:"host"`
	Port     int    `json:"port" mapstructure:"port"`
	Username string `json:"username" mapstructure:"username"`
	Password string `json:"password" mapstructure:"password"`
}
