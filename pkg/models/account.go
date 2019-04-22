package models

import "time"

// Account type details
type Account struct {
	Num     string
	UserID  string
	Name    string
	OpenAt  time.Time
	Balance float64
}
