package model

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Application error codes
const (
	ErrConflict        = "Conflict"        // action cannot be performed
	ErrInternal        = "Internal"        // internal error
	ErrBadRequest      = "Bad Request"     // bad request
	ErrInvalid         = "Invalid"         // validation failed
	ErrNotFound        = "Not Found"       // entity does not exist
	ErrUnauthorized    = "Unauthorized"    // permission denied
	ErrUnauthenticated = "Unauthenticated" // invalid token provided
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

// FieldError holds the key and the message of the field that was invalid
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a list of validation errors
type ValidationErrors []*FieldError

// Add appends the new error to the list
func (list *ValidationErrors) Add(err *FieldError) {
	*list = append(*list, err)
}

// IsZero returns true if there are no errors
func (list ValidationErrors) IsZero() bool {
	return len(list) == 0
}

// Invalid creates the input validation error
func Invalid(field string, l *i18n.Localizer, msg *i18n.Message) *FieldError {
	e := &FieldError{Field: field}
	e.Message = locale.LocalizeDefaultMessage(l, msg)
	return e
}

// NewValidationError builds the invalid user error
func NewValidationError(name string, msg *i18n.Message, userID string, errs ValidationErrors) *AppErr {
	details := map[string]interface{}{}
	if userID != "" {
		details["userID"] = userID
	}
	if !errs.IsZero() {
		details["validation"] = map[string]interface{}{"errors": errs}
	}

	errorName := name + ".Validate"
	return NewAppErr(errorName, ErrInvalid, locale.GetUserLocalizer("en"), msg, http.StatusUnprocessableEntity, details)
}
