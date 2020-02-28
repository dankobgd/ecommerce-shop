package model

// NewBool returns pointer to bool
func NewBool(b bool) *bool { return &b }

// NewInt returns pointer to int
func NewInt(n int) *int { return &n }

// NewInt64 returns pointer to int64
func NewInt64(n int64) *int64 { return &n }

// NewString returns pointer to string
func NewString(s string) *string { return &s }
