package types

import "time"

//PasswordResetCode...
type PasswordResetCode struct {
	Expires time.Time `json:"expires"`
	UserID  uint      `json:"userId"`
	Code    string    `json:"code"`
}
