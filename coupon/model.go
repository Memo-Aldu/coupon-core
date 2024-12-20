package coupon

import (
	"log"
	"time"
	"encoding/json"
)

type DiscountType string


const (
	Percentage DiscountType = "percentage"
	Fixed      DiscountType = "fixed"	
)


func (dt DiscountType) IsValid() bool {
	return dt == Percentage || dt == Fixed
}


type CreateCouponRequest struct {
	Code                string    		`json:"code"`
	DiscountType        DiscountType    `json:"discount_type"`
	Value			    float64   		`json:"value"`
	MaxRedemptions      int       		`json:"max_redemptions"`
	ExpiryDate		    string    		`json:"expiry_date"`
	MinimumOrderValue   float64   		`json:"minimum_order_value"`
	ApplicableProducts  json.RawMessage `json:"applicable_products"`
	IsActive 		 	bool      		`json:"is_active"`
	UserSpecific 		bool      		`json:"user_specific"`
}


type UpdateCouponRequest struct {
	DiscountType        string    `json:"discount_type"`
	Value			    float64   `json:"value"`
	MaxRedemptions      int       `json:"max_redemptions"`
	ExpiryDate		    string 	  `json:"expiry_date"`
	MinimumOrderValue   float64   `json:"minimum_order_value"`
	ApplicableProducts  []string 	  `json:"applicable_products"`
	IsActive 		 	bool      `json:"is_active"`	
}


type Coupon struct {
    ID                  int       		`json:"id" db:"id"`
    Code                string    		`json:"code" db:"code"`
    DiscountType        DiscountType    `json:"discount_type" db:"discount_type"`
    Value			    float64   		`json:"value" db:"value"`    
	MinimumOrderValue   float64   		`json:"minimum_order_value" db:"minimum_order_value"`
	MaxRedemptions      int       		`json:"max_redemptions" db:"max_redemptions"`
	RedeemedCount 	    int       		`json:"redeemed_count" db:"redeemed_count"`
	ExpiryDate		    time.Time 		`json:"expiry_date" db:"expiry_date"`
	ApplicableProducts  []string		    `json:"applicable_products" db:"applicable_products"`
	CreatedAt 		 	time.Time 		`json:"created_at" db:"created_at"`
	UpdatedAt 		 	time.Time 		`json:"updated_at" db:"updated_at"`
	IsActive 		 	bool      		`json:"is_active" db:"is_active"`
	UserSpecific 		bool      		`json:"user_specific" db:"user_specific"`
}

type CouponUser struct {
	ID          int       `json:"id" db:"id"`
	ExternalId  int       `json:"external_id" db:"external_id"`
}

type CouponRedemption struct {
	ID           int       `json:"id" db:"id"`
	CouponId     int       `json:"coupon_id" db:"coupon_id"`
	CouponUserId int       `json:"coupon_user_id" db:"coupon_user_id"`
	OrderId      int       `json:"order_id" db:"order_id"`
	RedeemedAt   time.Time `json:"redeemed_at" db:"redeemed_at"`
}


func NewCoupon(code string, discountType DiscountType, value, minimumOrderValue float64, 
	maxRedemptions int, expiryDate string, applicableProducts []string, IsActive, isUserSpecific bool ) *Coupon {
	
	date, err := parseDate(expiryDate)
	if err != nil {
		log.Println("Error parsing date", err)
		return nil
	}

	return &Coupon{
		Code: code,
		DiscountType: discountType,
		Value: value,
		MinimumOrderValue: minimumOrderValue,
		MaxRedemptions: maxRedemptions,
		ExpiryDate: date,
		ApplicableProducts: applicableProducts,
		IsActive: IsActive,
		UserSpecific: isUserSpecific,
		CreatedAt: time.Now().UTC(),
	}
}