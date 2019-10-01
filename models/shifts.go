package models

import "time"

//Shifts data structure
//Used for processing
type Shifts struct {
	Id      uint64    `json:"id"`
	Worker  string    `json:"worker"`
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

//Struct that is sent to frontend, the StartAt and EndAt is in string,which is converted
type JsonShifts struct {
	Id      uint64 `json:"id"`
	Worker  string `json:"worker"`
	StartAt string `json:"start_at"`
	EndAt   string `json:"end_at"`
}
