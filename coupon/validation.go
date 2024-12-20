package coupon

import (
	"fmt"
	"strings"
	"time"
	"errors"
)


func ValidateCreateCouponRequest(coupon *Coupon) error {
	if coupon == nil {
		return fmt.Errorf("Coupon is nil")
	}

	if coupon.Code == "" {
		return fmt.Errorf("Coupon name is required")
	}

	if coupon.DiscountType == "" {
		return fmt.Errorf("Coupon discount type is required")
	}

	if coupon.Value == 0 {
		return fmt.Errorf("Coupon value is required")
	}

	return nil
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