package models

import "time"

// Account type details
type Account struct {
	Num     string    `json:"num"`
	UserID  string    `json:"userID"`
	Name    string    `json:"name"`
	OpenAt  time.Time `json:"openAt"`
	Balance int64     `json:"balance"`
}
