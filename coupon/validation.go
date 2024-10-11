package coupon

import (
	"fmt"
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