package model

// ProductProperties holds valid product variant properties
type ProductProperties struct {
	Tshirts struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Theme []string `json:"theme"`
	} `json:"tshirts"`
	Shirts struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Theme []string `json:"theme"`
	} `json:"shirts"`
	Shoes struct {
		Color    []string `json:"color"`
		Size     []string `json:"size"`
		LaceSize []string `json:"lace_size"`
		Theme    []string `json:"theme"`
	} `json:"shoes"`
	Sneakers struct {
		Color    []string `json:"color"`
		Size     []string `json:"size"`
		LaceSize []string `json:"lace_size"`
		Theme    []string `json:"theme"`
	} `json:"sneakers"`
	Boots struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Theme []string `json:"theme"`
	} `json:"boots"`
	Jackets struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Theme []string `json:"theme"`
	} `json:"jackets"`
	Shorts struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Theme []string `json:"theme"`
	} `json:"shorts"`
	Jeans struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Theme []string `json:"theme"`
	} `json:"jeans"`
	Sweatpants struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Theme []string `json:"theme"`
	} `json:"sweatpants"`
	Balls struct {
		Color []string `json:"color"`
		Type  []string `json:"type"`
	} `json:"balls"`
	Backpacks struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Theme []string `json:"theme"`
	} `json:"backpacks"`
	Rackets struct {
		Size []string `json:"size"`
		Type []string `json:"type"`
	} `json:"rackets"`
	Bags struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Theme []string `json:"theme"`
	} `json:"bags"`
	Skateboards struct {
		Color []string `json:"color"`
	} `json:"skateboards"`
	Rollerblades struct {
		Color []string `json:"color"`
	} `json:"rollerblades"`
	Bycicles struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
		Gears []string `json:"gears"`
	} `json:"bycicles"`
	Hats struct {
		Color []string `json:"color"`
		Theme []string `json:"theme"`
	} `json:"hats"`
	Helmets struct {
		Color []string `json:"color"`
		Type  []string `json:"type"`
	} `json:"helmets"`
	Suitcases struct {
		Color []string `json:"color"`
		Size  []string `json:"size"`
	} `json:"suitcases"`
}
