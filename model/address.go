package model

import (
	"encoding/json"
	"io"
	"time"

	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgInvalidAddress    = &i18n.Message{ID: "model.address.validate.app_error", Other: "Invalid address data"}
	msgValidateAddressID = &i18n.Message{ID: "model.address.validate.address_id.app_error", Other: "Invalid address id"}
	msgValidateLine1     = &i18n.Message{ID: "model.address.validate.line_1.app_error", Other: "Invalid address line 1"}
	msgValidateCity      = &i18n.Message{ID: "model.address.validate.city.app_error", Other: "Invalid address city"}
	msgValidateCountry   = &i18n.Message{ID: "model.address.validate.country.app_error", Other: "Invalid address country"}
)

// Address holds the contact info
type Address struct {
	ID        int64      `json:"id" db:"id"`
	Line1     string     `json:"line_1" db:"line_1"`
	Line2     *string    `json:"line_2" db:"line_2"`
	City      string     `json:"city" db:"city"`
	Country   string     `json:"country" db:"country"`
	State     *string    `json:"state" db:"state"`
	ZIP       *string    `json:"zip" db:"zip"`
	Latitude  *float64   `json:"latitude" db:"latitude"`
	Longitude *float64   `json:"longitude" db:"longitude"`
	Phone     *string    `json:"phone" db:"phone"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}

// GeocodingResultList is a list of geocoding results
type GeocodingResultList []*GeocodingResult

// GeocodingResult is a geocoding item result
type GeocodingResult struct {
	PlaceID     string   `json:"place_id,omitempty"`
	Licence     string   `json:"licence,omitempty"`
	OsmType     string   `json:"osm_type,omitempty"`
	OsmID       string   `json:"osm_id,omitempty"`
	Boundingbox []string `json:"boundingbox,omitempty"`
	Lat         string   `json:"lat,omitempty"`
	Lon         string   `json:"lon,omitempty"`
	DisplayName string   `json:"display_name,omitempty"`
	Class       string   `json:"class,omitempty"`
	Type        string   `json:"type,omitempty"`
	Importance  float64  `json:"importance,omitempty"`
}

// PreSave will set missing defaults and fill CreatedAt and UpdatedAt times
func (addr *Address) PreSave() {
	addr.CreatedAt = time.Now()
	addr.UpdatedAt = addr.CreatedAt
}

// PreUpdate sets the update timestamp
func (addr *Address) PreUpdate() {
	addr.UpdatedAt = time.Now()
}

// Validate validates the address and returns an error if it doesn't pass criteria
func (addr *Address) Validate() *AppErr {
	var errs ValidationErrors
	l := locale.GetUserLocalizer("en")

	if addr.ID != 0 {
		errs.Add(Invalid("address_id", l, msgValidateAddressID))
	}
	if addr.Line1 == "" {
		errs.Add(Invalid("line_1", l, msgValidateLine1))
	}
	if addr.City == "" {
		errs.Add(Invalid("city", l, msgValidateCity))
	}
	if addr.Country == "" {
		errs.Add(Invalid("country", l, msgValidateCountry))
	}

	if !errs.IsZero() {
		return NewValidationError("Address", msgInvalidAddress, "", errs)
	}
	return nil
}

// AddressPatch is the patch
type AddressPatch struct {
	Line1     *string  `json:"line_1"`
	Line2     *string  `json:"line_2"`
	City      *string  `json:"city"`
	Country   *string  `json:"country"`
	State     *string  `json:"state"`
	ZIP       *string  `json:"zip"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Phone     *string  `json:"phone"`
}

// Patch patches the address fields that are provided
func (addr *Address) Patch(patch *AddressPatch) {
	if patch.Line1 != nil {
		addr.Line1 = *patch.Line1
	}
	if patch.Line2 != nil {
		addr.Line2 = patch.Line2
	}
	if patch.City != nil {
		addr.City = *patch.City
	}
	if patch.Country != nil {
		addr.Country = *patch.Country
	}
	if patch.State != nil {
		addr.State = patch.State
	}
	if patch.ZIP != nil {
		addr.ZIP = patch.ZIP
	}
	if patch.Latitude != nil {
		addr.Latitude = patch.Latitude
	}
	if patch.Longitude != nil {
		addr.Longitude = patch.Longitude
	}
	if patch.Phone != nil {
		addr.Phone = patch.Phone
	}
}

// AddressPatchFromJSON decodes the input and returns the AddressPatch
func AddressPatchFromJSON(data io.Reader) (*AddressPatch, error) {
	var patch *AddressPatch
	err := json.NewDecoder(data).Decode(&patch)
	return patch, err
}

// AddressFromJSON decodes the input and return the Address
func AddressFromJSON(data io.Reader) (*Address, error) {
	var addr *Address
	err := json.NewDecoder(data).Decode(&addr)
	return addr, err
}

// ToJSON converts Address to json string
func (addr *Address) ToJSON() string {
	b, _ := json.Marshal(addr)
	return string(b)
}
