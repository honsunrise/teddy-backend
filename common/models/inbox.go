package models

import "time"

type InBoxType uint32

const (
	SYSTEM InBoxType = iota
	AT
	REVIEW
	PRIVATE
)

type InBoxItem struct {
	ID       string
	From     string
	Type     InBoxType
	Subject  string
	Content  string
	Unread   bool
	SendTime time.Time
	ReadTime time.Time
}

type InBox struct {
	Uid         string
	Items       []InBoxItem
	UnReadCount uint64 `json:"unread_count"`
}
