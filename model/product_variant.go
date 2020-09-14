package model

// ProductProperties contains the variant props
type ProductProperties []struct {
	Category string                 `json:"category"`
	Props    map[string]interface{} `json:"props"`
}
