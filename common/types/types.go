package types

import (
	"errors"
	"strings"
)

var ErrOrderNotSupport = errors.New("order not support")

type Paging struct {
	Page uint32 `json:"page"`
	Size uint32 `json:"size"`
}

type Order uint8

const (
	ASC Order = 1 << iota
	DESC
)

func (o Order) String() string {
	switch o {
	case ASC:
		return "ASC"
	case DESC:
		return "DESC"
	default:
		return "ASC"
	}
}

func ParseOrder(order string) (Order, error) {
	switch strings.ToUpper(order) {
	case "ASC":
		return ASC, nil
	case "DESC":
		return DESC, nil
	default:
		return 0, ErrOrderNotSupport
	}
}

type Sort struct {
	Name  string
	Order Order
}

// For config

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

type Mail struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Config struct {
	Server    Server                `json:"server"`
	Databases map[string][]Database `json:"databases"`
	JWTPkcs8  string                `json:"jwt_pkcs8"`
	Mail      Mail                  `json:"mail"`
}
