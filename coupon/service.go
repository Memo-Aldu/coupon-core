package coupon


type CouponService struct {
	repository Repository
}


func NewCouponService(repository Repository) *CouponService {
	return &CouponService{repository}
}


func (s *CouponService) GetCouponById(id int) (*Coupon, error) {
	return s.repository.GetCouponById(id)
}


func (s *CouponService) CreateCoupon(coupon *Coupon) (*Coupon, error) {
	return s.repository.CreateCoupon(coupon)
}


func (s *CouponService) UpdateCoupon(coupon *Coupon) (*Coupon, error) {
	return s.repository.UpdateCoupon(coupon)
}

