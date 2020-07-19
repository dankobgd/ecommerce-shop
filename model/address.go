package model

import (
	"encoding/json"
	"io"
	"time"
)

// AddrType is the contact address type
type AddrType int

// addr types
const (
	BillingAddress AddrType = iota
	ShippingAddress
	PhysicalAddress
)

func (addr AddrType) String() string {
	switch addr {
	case BillingAddress:
		return "billing"
	case ShippingAddress:
		return "shipping"
	case PhysicalAddress:
		return "physical"
	default:
		return "unknown"
	}
}

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

// AddressType is the type of the address
type AddressType struct {
	ID        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	AddressID int64  `json:"address_id" db:"address_id"`
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
