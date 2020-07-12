package model

// Cart is shopping cart model
type Cart struct {
	Products     map[int64]*Product
	Quantity     int
	Subtotal     int
	Total        int
	ChargeAmount int
}

// StripeAmount returns the stripe amount in cents
func (c *Cart) StripeAmount() int {
	return c.Total * 100
}

// AddItemToCart adds item to the cart
func (c *Cart) AddItemToCart(p *Product) {
	p, ok := c.Products[p.ID]
	if !ok {
		c.Products[p.ID] = p
	}

	c.Quantity++
}

// RemoveItemFromCart removes item from the cart
func (c *Cart) RemoveItemFromCart(pid int64) {
	delete(c.Products, pid)

	if _, ok := c.Products[pid]; ok {
		c.Quantity--
	}
}
