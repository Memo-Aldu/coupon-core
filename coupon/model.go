package coupon

import (
	"errors"
	"log"
	"time"
	"encoding/json"
	"strings"
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


func (req *CreateCouponRequest) Validate() error {
	var errors ValidationErrors
	if req.Code == "" {
		errors = append(errors, NewValidationError("code", "code is required"))
	}
	if !DiscountType(req.DiscountType).IsValid() {
		errors = append(errors, NewValidationError("discount_type", "invalid discount type"))
	}
	if req.Value <= 0 {
		errors = append(errors, NewValidationError("value", "value must be greater than zero"))
	}
	if req.MaxRedemptions < 0 {
		errors = append(errors, NewValidationError("max_redemptions", "max_redemptions cannot be negative"))
	}
	if req.MinimumOrderValue < 0 {
		errors = append(errors, NewValidationError("minimum_order_value", "minimum_order_value cannot be negative"))
	}

	if len(errors) > 0 {
        return errors
    }

	return nil
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


func (r *CouponRedemption) Validate() error {
	var errors ValidationErrors 
	if r.CouponId <= 0 {
		errors = append(errors, NewValidationError("coupon_id", "coupon_id must be greater than zero"))
	}

	if r.CouponUserId <= 0 {
		errors = append(errors, NewValidationError("coupon_user_id", "coupon_user_id must be greater than zero"))
	}

	if r.OrderId <= 0 {
		errors = append(errors, NewValidationError("order_id", "order_id must be greater than zero"))
	}

	if len(errors) > 0 {
		return errors
	}
	
	return nil
}


func NewCoupon(code string, discountType DiscountType, value, minimumOrderValue float64, maxRedemptions int, expiryDate string, 
	applicableProducts []string, IsActive, isUserSpecific bool ) *Coupon {
	
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


func parseDate(date string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02T15:04:05Z", date)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected '2006-01-02T15:04:05Z'")
	}
	return parsed, nil
}


type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}

func NewValidationError(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}


type ValidationErrors []error

func (v ValidationErrors) Error() string {
    var messages []string
    for _, err := range v {
        messages = append(messages, err.Error())
    }
    return strings.Join(messages, ", ")
}