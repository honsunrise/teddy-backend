package models

import "time"

type InBoxType uint32

const (
	SYSTEM InBoxType = iota
	AT
	REVIEW
	PRIVATE
	ALL
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
	Uid   string
	Items []InBoxItem
}
