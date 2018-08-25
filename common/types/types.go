package types

import "strings"

type Paging struct {
	Page  uint32 `json:"page" form:"page" query:"page"`
	Size  uint32 `json:"size" form:"size" query:"size"`
	Sort  string `json:"sort" form:"sort" query:"sort"`
	Order string `json:"order" form:"order" query:"order"`
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
