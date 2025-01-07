package coupon

import (
	"encoding/json"
	"fmt"
	"time"
)


type CouponService struct {
	repository Repository
}


func NewCouponService(repository Repository) *CouponService {
	return &CouponService{repository}
}


func (s *CouponService) GetCouponById(id int) (*Coupon, error) {
	return s.repository.GetCouponById(id)
}


func (s *CouponService) CreateCoupon(newCoupon CreateCouponRequest) (*Coupon, error) {
	var applicableProducts []string
	if err := json.Unmarshal(newCoupon.ApplicableProducts, &applicableProducts); err != nil {
		return nil, fmt.Errorf("invalid applicable_products JSON: %w", err)
	}

	coupon := NewCoupon(newCoupon.Code, newCoupon.DiscountType, newCoupon.Value,
		newCoupon.MinimumOrderValue, newCoupon.MaxRedemptions, newCoupon.ExpiryDate,
		applicableProducts, newCoupon.IsActive, newCoupon.UserSpecific)
	
	if err := ValidateCreateCouponRequest(coupon); err != nil {
		return nil, err
	}
	return s.repository.CreateCoupon(coupon)
}


func (s *CouponService) UpdateCoupon(id int, request UpdateCouponRequest) (*Coupon, error) {
	coupon, err := s.repository.GetCouponById(id)
	if err != nil {
		return  nil, fmt.Errorf("failed to get coupon: %w", err)
	}
	if request.DiscountType != "" {
		coupon.DiscountType = request.DiscountType
	}

	if request.Value >= 0 {
		coupon.Value = request.Value
	}

	if request.MaxRedemptions >= 0 {
		coupon.MaxRedemptions = request.MaxRedemptions
	}

	if request.ExpiryDate != "" {
		coupon.ExpiryDate, err = parseDate(request.ExpiryDate); if err != nil {
			return nil, fmt.Errorf("invalid expiry date: %w", err)
		}
	}

	if request.MinimumOrderValue >= 0 {
		coupon.MinimumOrderValue = request.MinimumOrderValue
	}

	if request.ApplicableProducts != nil {
		var applicableProducts []string
		if err := json.Unmarshal(request.ApplicableProducts, &applicableProducts); err != nil {
			return nil, fmt.Errorf("invalid applicable_products JSON: %w", err)
		}
		coupon.ApplicableProducts = applicableProducts
	}

	if request.IsActive != coupon.IsActive {
		coupon.IsActive = request.IsActive
	}

	coupon.UpdatedAt = time.Now()
	return s.repository.UpdateCoupon(coupon)
}

