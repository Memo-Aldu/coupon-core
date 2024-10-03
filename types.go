package main


import (
	"math/rand"
)


type Coupon struct {
    ID              int       `json:"id" db:"id"`
    Code            string    `json:"code" db:"code"`
    DiscountType    string    `json:"discount_type" db:"discount_type"`
    Value			float64   `json:"value" db:"value"`    
}


func NewCoupon(code, discountType string, value float64) *Coupon {
	return &Coupon{
		ID: rand.Intn(10000),
		Code: code,
		DiscountType: discountType,
		Value: value,
	}
}