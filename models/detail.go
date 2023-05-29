package models

import "time"

type Details struct {
	Txname   string    `json:"tx_name"`
	Amount   int       `json:"amount"`
	Datetime time.Time `json:"datetime"`
}
