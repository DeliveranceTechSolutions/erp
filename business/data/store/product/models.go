package product

import "time"

// user.Product is good form when naming
// and structuring core model types

type Category struct {
	Title string
	Type  string
}

// Product represents an individual user.
type Product struct {
	ID           string         `db:"product_id" json:"id"`
	Name         string         `db:"name" json:"name"`
	Quantity     string         `db:"quantity" json:"quantity"`
	OwnerID		 string			`db:"owner_id" json:"owner_id"`
	Reserve	 	 float32		`db:"reserve" json:"reserve"`
	StarterBid	 float32		`db:"starter_bid" json:"starter_bid"`
	Bids		 []string		`db:"bids" json:"bids"`
	Type 		 Category  		`db:"product_type" json:"product_type"`
	DateCreated  time.Time      `db:"date_created" json:"date_created"`
	DateUpdated  time.Time      `db:"date_updated" json:"date_updated"`
}

// using a New{CoreType} idiom allows you to circumvent
// clients inserting malicious code

// NewProduct contains information needed to create a new Product.
type NewProduct struct {
	Name         string   		`json:"name" validate:"required"`
	Quantity     string         `json:"quantity" validate:"quantity"`
	OwnerID		 string			`json:"user_id" validate:"owner_id"`
	Reserve		 float32		`json:"reserve" validate:"reserve"`
	StarterBid	 float32		`json:"starter_bid" validate:"starter_bid"`
	Bids		 []string		`json:"bids" validate:"bids"`
}

// UpdateProduct defines what information may be provided to modify an existing
// Product. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exceptions around
// marshalling/unmarshalling.
type UpdateProduct struct {
	Name         *string  		`json:"name"`
	Reserve		 *float32		`json:"reserve"`
	StarterBid	 *float32		`json:"starter_bid"`
}

// VENDORS CAN BUYITNOW WHILE NON-VENDORS NEED TO BID