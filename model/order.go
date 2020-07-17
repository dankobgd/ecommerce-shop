package model

import "time"

type orderStatus int

// order statuses
const (
	OrderStatusPending orderStatus = iota
	OrderStatusSuccess
	OrderStatusFailed
)

func (s orderStatus) String() string {
	switch s {
	case OrderStatusPending:
		return "pending"
	case OrderStatusSuccess:
		return "success"
	case OrderStatusFailed:
		return "fail"
	default:
		return "unknown"
	}
}

// Order represents the transaction
type Order struct {
	ID        int64      `json:"id" db:"id"`
	UserID    int64      `json:"user_id" db:"user_id"`
	Status    string     `json:"status" db:"status"`
	Total     int        `json:"total" db:"total"`
	ShippedAt *time.Time `json:"shipped_at" db:"shipped_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}

// PreSave fills the defaults
func (o *Order) PreSave() {
	o.CreatedAt = time.Now()
	if o.Status == "" {
		o.Status = OrderStatusPending.String()
	}
}
