package coupon


import (
	"database/sql"
	"log"
	"fmt"
	"encoding/json"
	// importing the pq only for the side effects
	_ "github.com/lib/pq"
)

type Repository interface {
	connect() error
	disconnect() error
	GetCouponById(id int) (*Coupon, error)
	CreateCoupon(coupon *Coupon) (*Coupon, error)
	UpdateCoupon(coupon *Coupon) (*Coupon, error)
	DeleteCoupon(id int) error
}

// PostgresDatabase struct that implements the Database interface
type PostgresRepository struct {
	db *sql.DB
}


func NewPostgresRepository() (*PostgresRepository, error) {
	db := &PostgresRepository{}
	if err := db.connect(); err != nil {
		return nil, err
	}

	return db, nil
}


func (db *PostgresRepository) Init() error {
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


func (db *PostgresRepository) connect() error {
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


func (db *PostgresRepository) disconnect() error {
	err := db.db.Close()
	if err != nil {
		return err
	}

	return nil
}


func (db *PostgresRepository) GetCouponById(id int) (*Coupon, error) {
    query := `SELECT id, code, discount_type, value, max_redemptions, redeemed_count, 
              expiry_date, minimum_order_value, applicable_products, is_active, 
              user_specific, created_at, updated_at FROM coupons WHERE id = $1`

	row := db.db.QueryRow(query, id)

    coupon := new(Coupon)
	var applicableProductsJSON []byte
    err := row.Scan(&coupon.ID, &coupon.Code, &coupon.DiscountType, &coupon.Value, 
		&coupon.MaxRedemptions, &coupon.RedeemedCount, &coupon.ExpiryDate, 
		&coupon.MinimumOrderValue, &applicableProductsJSON, &coupon.IsActive, 
		&coupon.UserSpecific, &coupon.CreatedAt, &coupon.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
				return nil, fmt.Errorf("Coupon not found")
		}
		return nil, err
	}

	if err := json.Unmarshal(applicableProductsJSON, &coupon.ApplicableProducts); err != nil {
        return nil, fmt.Errorf("failed to unmarshal applicable_products: %w", err)
    }

	return coupon, nil
}


func (db *PostgresRepository) CreateCoupon(coupon *Coupon) (*Coupon, error) {
    query := `INSERT INTO coupons (
        code, discount_type, value, max_redemptions, expiry_date, minimum_order_value, 
        applicable_products, is_active, user_specific)
        VALUES ($1, $2, $3, $4, $5, $6, $7::jsonb, $8, $9) RETURNING id`

    // Marshal applicable_products to JSON
    applicableProductsJSON, err := json.Marshal(coupon.ApplicableProducts)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal applicable_products: %w", err)
    }

    log.Printf("applicableProductsJSON: %s", string(applicableProductsJSON))

    var id int
    err = db.db.QueryRow(query, coupon.Code, coupon.DiscountType, coupon.Value,
        coupon.MaxRedemptions, coupon.ExpiryDate, coupon.MinimumOrderValue,
        applicableProductsJSON, coupon.IsActive, coupon.UserSpecific).Scan(&id)

    if err != nil {
        return nil, fmt.Errorf("failed to create coupon: %w", err)
    }
    coupon.ID = id

    log.Println("Coupon created successfully")
    return coupon, nil
}


func (db *PostgresRepository) UpdateCoupon(coupon *Coupon) (*Coupon, error) {
	query := `UPDATE coupons SET discount_type = $1, value = $2, 
	max_redemptions = $3, expiry_date = $4, minimum_order_value = $5, 
	applicable_products = $6::jsonb, is_active = $7, updated_at = $8 WHERE id = $9 RETURNING id`

	applicableProductsJSON, err := json.Marshal(coupon.ApplicableProducts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal applicable_products: %w", err)
	}

	var id int
	err = db.db.QueryRow(query, coupon.DiscountType, coupon.Value, coupon.MaxRedemptions, coupon.ExpiryDate, 
		coupon.MinimumOrderValue, applicableProductsJSON, coupon.IsActive, coupon.UpdatedAt, coupon.ID).Scan(&id)
	
	if err != nil {
		return nil, fmt.Errorf("failed to update coupon with id: %d: %w", coupon.ID, err)
	}

	coupon.ID = id
	return coupon, nil
}


func (db *PostgresRepository) DeleteCoupon(id int) error {
	query := `DELETE FROM coupons WHERE id = $1`
	_, err := db.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete coupon with id: %d: %w", id, err)
	}
	return nil
}


func (db *PostgresRepository) createCouponTable() error {
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
			applicable_products JSONB DEFAULT '[]',
			is_active BOOLEAN DEFAULT TRUE,
			user_specific BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW()
		);
	`

	_, err := db.db.Exec(query)

	return err
}


func (db *PostgresRepository) createCouponUserTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS coupon_users (
			id SERIAL PRIMARY KEY,
			external_id VARCHAR(255) NOT NULL
		)
	`
	_, err := db.db.Exec(query)

	return err
}


func (db *PostgresRepository) createCouponRedemptionTable() error {
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



