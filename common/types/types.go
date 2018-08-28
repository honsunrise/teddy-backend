package types

import (
	"errors"
	"strings"
)

var ErrOrderNotSupport = errors.New("order not support")

type Paging struct {
	Page  uint32 `json:"page"`
	Size  uint32 `json:"size"`
	Sort  string `json:"sort"`
	Order string `json:"order"`
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
