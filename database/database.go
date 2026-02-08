package database

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB(conn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(conn)
	if err != nil {
		return nil, err
	}

	// 🔑 Paksa IPv4 (wajib di kasus kamu)
	cfg.ConnConfig.DialFunc = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, "tcp4", addr)
	}

	// 🔥 POOL CONFIG (setara dengan database/sql)
	cfg.MaxConns = 25              // = SetMaxOpenConns
	cfg.MinConns = 5               // = SetMaxIdleConns
	cfg.MaxConnLifetime = 30 * time.Minute
	cfg.MaxConnIdleTime = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	log.Println("Database connected & pool configured")
	return db, nil
}

func Migrate(db *pgxpool.Pool) error {
	ctx := context.Background()

	// Create categories table
	_, err := db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			description TEXT
		);
	`)
	if err != nil {
		return err
	}

	// Add category_id to products if not exists
	_, err = db.Exec(ctx, `
		ALTER TABLE products 
		ADD COLUMN IF NOT EXISTS category_id INT;
	`)
	if err != nil {
		return err
	}

	// Create carts table
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS carts (
			id SERIAL PRIMARY KEY,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Create cart_items table
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS cart_items (
			id SERIAL PRIMARY KEY,
			cart_id INT REFERENCES carts(id) ON DELETE CASCADE,
			product_id INT REFERENCES products(id),
			quantity INT NOT NULL CHECK (quantity > 0),
			UNIQUE(cart_id, product_id)
		);
	`)
	if err != nil {
		return err
	}

	// Create transactions table
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			invoice_number VARCHAR(50) UNIQUE,
			total_amount INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Add invoice_number column if not exists (for existing tables)
	_, err = db.Exec(ctx, `
		ALTER TABLE transactions 
		ADD COLUMN IF NOT EXISTS invoice_number VARCHAR(50) UNIQUE;
	`)
	if err != nil {
		return err
	}

	// Create transaction_items table
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS transaction_items (
			id SERIAL PRIMARY KEY,
			transaction_id INT REFERENCES transactions(id) ON DELETE CASCADE,
			product_id INT REFERENCES products(id) ON DELETE SET NULL,
			product_name VARCHAR(255),
			quantity INT NOT NULL,
			price INT NOT NULL,
			subtotal INT NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	// Add new columns to transaction_items if not exists
	_, err = db.Exec(ctx, `
		ALTER TABLE transaction_items 
		ADD COLUMN IF NOT EXISTS product_name VARCHAR(255),
		ADD COLUMN IF NOT EXISTS subtotal INT;
	`)
	if err != nil {
		return err
	}

	// Update foreign key to SET NULL on delete for persistence
	_, err = db.Exec(ctx, `
		ALTER TABLE transaction_items 
		DROP CONSTRAINT IF EXISTS transaction_items_product_id_fkey,
		ADD CONSTRAINT transaction_items_product_id_fkey 
		FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE SET NULL;
	`)
	if err != nil {
		return err
	}

	log.Println("✅ Database migration executed")
	return nil
}

