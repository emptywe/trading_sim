package model

import "time"

type Session struct {
	Sid         string
	Name        string
	Value       string
	Valid       bool
	Established time.Time
}

type CleanSession struct {
	Sid         string
	Established time.Time
}
