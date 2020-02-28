package model

import (
	"encoding/json"
	"io"
	"time"
)

// Session represents the user session data
type Session struct {
	ID             string    `json:"id"`
	Token          string    `json:"token"`
	UserID         string    `json:"user_id"`
	DeviceID       string    `json:"device_id"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
	LastActivityAt time.Time `json:"last_activity_at"`
}

// ToJSON converts session to json string
func (me *Session) ToJSON() string {
	b, _ := json.Marshal(me)
	return string(b)
}

// SessionFromJSON decodes the input and returns the Session
func SessionFromJSON(data io.Reader) (*Session, error) {
	var me *Session
	err := json.NewDecoder(data).Decode(&me)
	return me, err
}
