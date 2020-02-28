package model

import (
	"encoding/json"
	"fmt"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Application error codes
const (
	ErrConflict        = "Conflict"         // action cannot be performed
	ErrInternal        = "Internal"         // internal error
	ErrInvalid         = "Invalid"          // validation failed
	ErrNotFound        = "Not Found"        // entity does not exist
	ErrPermision       = "Permision Denied" // permission denied
	ErrUnauthenticated = "Unauthenticated"  // unauthenticated access
)

// AppErr is the main app error
type AppErr struct {
	ID         string      `json:"id"`            // unique string which is the same as translation id
	Op         string      `json:"op"`            // operation where it failed (Struct.Func)
	Code       string      `json:"code"`          // machine readable error code
	StatusCode int         `json:"status_code"`   // http status code
	Message    string      `json:"message"`       // meaningful end user message
	Err        error       `json:"err,omitempty"` // embeded error
	Details    interface{} `json:"details,omitempty"`
}

// NewAppErr creates the new app error
func NewAppErr(op string, code string, l *i18n.Localizer, msg *i18n.Message, statusCode int, details interface{}) *AppErr {
	e := &AppErr{
		ID:         msg.ID,
		Op:         op,
		Code:       code,
		StatusCode: statusCode,
		Details:    details,
	}
	e.Message = locale.LocalizeDefaultMessage(l, msg)
	return e
}

func (e *AppErr) Error() string {
	return fmt.Sprintf("%v, %v: %v\n", e.Code, e.Op, e.Message)
}

// ToJSON converts AppErr to json string
func (e *AppErr) ToJSON() string {
	b, _ := json.Marshal(e)
	return string(b)
}

// ValidationErr holds the key and the message of the field that had errors
type ValidationErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a list of validation errors
type ValidationErrors []*ValidationErr

// Add appends the new error to the list
func (list *ValidationErrors) Add(verr *ValidationErr) {
	*list = append(*list, verr)
}

// IsZero returns true if there are no errors
func (list ValidationErrors) IsZero() bool {
	return len(list) == 0
}

// NewValidationErr creates the validationErr
func NewValidationErr(field string, l *i18n.Localizer, msg *i18n.Message) *ValidationErr {
	e := &ValidationErr{Field: field}
	e.Message = locale.LocalizeDefaultMessage(l, msg)
	return e
}
