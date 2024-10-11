package main

import (
	"log"

	"github.com/Memo-Aldu/coupon-core/coupon"
)

func main() {
	log.Println("Starting API Server")
	repo, err := coupon.NewPostgresRepository()

	if err != nil {
		log.Fatal(err)
	}

	err = repo.Init()
	if err != nil {
		log.Fatal(err)
	}

	svc := coupon.NewCouponService(repo)

	var server *coupon.APIServer = coupon.NewAPIServer(svc, ":3000", "v1", "/api")
	server.Start()
}