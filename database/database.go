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

	log.Println("✅ Database migration executed")
	return nil
}

