package main

import (
	"database/sql"
	"log"
	"fmt"
	"strings"
	"strconv"
	// importing the pq only for the side effects
	_ "github.com/lib/pq"
)

type Database interface {
	connect() error
	disconnect() error
	GetCouponById(id int) (*Coupon, error)
	CreateCoupon(coupon *Coupon) (*Coupon, error)
	UpdateCoupon(coupon *Coupon) (*Coupon, error)
	DeleteCoupon(id int) error
}

// PostgresDatabase struct that implements the Database interface
type PostgresDatabase struct {
	db *sql.DB
}


func NewPostgresDatabase() (*PostgresDatabase, error) {
	db := &PostgresDatabase{}
	if err := db.connect(); err != nil {
		return nil, err
	}

	return db, nil
}


func (db *PostgresDatabase) Init() error {
	err := db.createCouponTable()
	if err != nil {
		log.Fatal("Error creating coupon table", err)
		return err
	}

	err = db.createCouponUserTable()
	if err != nil {
		log.Fatal("Error creating coupon user table", err)
		return err
	}

	err = db.createCouponRedemptionTable()
	if err != nil {
		log.Fatal("Error creating coupon redemption table", err)
		return err
	}

	return nil
}


func (db *PostgresDatabase) connect() error {
	connectionString := "user=postgres dbname=postgres password=root sslmode=disable"
	database, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	err = database.Ping()
	if err != nil {
		return err
	}

	db.db = database

	return nil
}


func (db *PostgresDatabase) disconnect() error {
	err := db.db.Close()
	if err != nil {
		return err
	}

	return nil
}


func (db *PostgresDatabase) GetCouponById(id int) (*Coupon, error) {
    query := `SELECT id, code, discount_type, value, max_redemptions, redeemed_count, 
              expiry_date, minimum_order_value, applicable_products, is_active, 
              user_specific, created_at, updated_at FROM coupons WHERE id = $1`

	row := db.db.QueryRow(query, id)

    coupon := new(Coupon)
    var applicableProductsStr string
    err := row.Scan(&coupon.ID, &coupon.Code, &coupon.DiscountType, &coupon.Value, 
		&coupon.MaxRedemptions, &coupon.RedeemedCount, &coupon.ExpiryDate, 
		&coupon.MinimumOrderValue, &applicableProductsStr, &coupon.IsActive, 
		&coupon.UserSpecific, &coupon.CreatedAt, &coupon.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
				return nil, fmt.Errorf("Coupon not found")
		}
		return nil, err
	}

	productList, err := postgresArrayToIntSlice(applicableProductsStr )
	if err != nil {
		return nil, err
	}
	coupon.ApplicableProducts = productList

	return coupon, nil
}


func (db *PostgresDatabase) CreateCoupon(coupon *Coupon) (*Coupon, error) {
	applicableProductsStr := intSliceToPostgresArray(coupon.ApplicableProducts)
	query := `INSERT INTO coupons (
		code, discount_type, value, max_redemptions, expiry_date, minimum_order_value, 
		applicable_products, is_active, user_specific)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var id int
	err := db.db.QueryRow(query, coupon.Code, coupon.DiscountType, coupon.Value, 
		coupon.MaxRedemptions, coupon.ExpiryDate, coupon.MinimumOrderValue, 
		applicableProductsStr, coupon.IsActive, coupon.UserSpecific).Scan(&id)

	if err != nil {
		return nil, err
	}
	coupon.ID = id

	return coupon, nil
}


func (db *PostgresDatabase) UpdateCoupon(coupon *Coupon) (*Coupon, error) {
	return nil, nil
}


func (db *PostgresDatabase) DeleteCoupon(id int) error {
	return nil
}


func (db *PostgresDatabase) createCouponTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS coupons (
			id SERIAL PRIMARY KEY,
			code VARCHAR(50) UNIQUE NOT NULL,
			discount_type VARCHAR(20) NOT NULL,
			value DECIMAL(10, 2) NOT NULL,
			max_redemptions INT,
			redeemed_count INT DEFAULT 0,
			expiry_date TIMESTAMP NOT NULL,
			minimum_order_value DECIMAL(10, 2),
			applicable_products INT[] DEFAULT '{}',
			is_active BOOLEAN DEFAULT TRUE,
			user_specific BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`

	_, err := db.db.Exec(query)

	return err
}


func (db *PostgresDatabase) createCouponUserTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS coupon_users (
			id SERIAL PRIMARY KEY,
			external_id VARCHAR(255) NOT NULL
		)
	`
	_, err := db.db.Exec(query)

	return err
}


func (db *PostgresDatabase) createCouponRedemptionTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS coupon_redemptions (
			id SERIAL PRIMARY KEY,
			coupon_id INT NOT NULL REFERENCES coupons(id) ON DELETE CASCADE,
			coupon_user_id INT NOT NULL REFERENCES coupon_users(id),
			order_id INT NOT NULL,
			redeemed_at TIMESTAMP DEFAULT NOW()
		)
	`
	_, err := db.db.Exec(query)

	return err
}


func intSliceToPostgresArray(intSlice []int) string {
    strSlice := make([]string, len(intSlice))
    for i, val := range intSlice {
        strSlice[i] = fmt.Sprintf("%d", val)
    }
    return "{" + strings.Join(strSlice, ",") + "}"
}

func postgresArrayToIntSlice(pgArray string) ([]int, error) {
    pgArray = strings.Trim(pgArray, "{}")
    strSlice := strings.Split(pgArray, ",")
    
    intSlice := make([]int, len(strSlice))
    for i, s := range strSlice {
        num, err := strconv.Atoi(s)
        if err != nil {
            return nil, err
        }
        intSlice[i] = num
    }
    return intSlice, nil
}




